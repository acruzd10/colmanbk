package dyno

import (
	"colmanback/db"
	"colmanback/objects"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var Conn *dynamodb.DynamoDB

var cacheMap map[string][]db.CacheMapElement = make(map[string][]db.CacheMapElement)

type Dyno[K objects.Object] struct {
	tableName   string
	codeName    string
	sortName    string
	keepCache   bool
	constructor func() K
	cacheMap    func([]K) []db.CacheMapElement

	cache map[string]K
}

type CacheMapEntry struct {
	Key      string               `json:"key"`
	Elements []db.CacheMapElement `json:"elements"`
}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) getObjectListFromDB() ([]K, error) {
	var objectList []K

	params := &dynamodb.ScanInput{
		TableName: aws.String(dynoInst.tableName),
	}

	result, err := Conn.Scan(params)
	if err != nil {
		return objectList, fmt.Errorf("query API call on table %s failed. Err: %s", dynoInst.tableName, err)
	}

	for _, i := range result.Items {
		objectInst := dynoInst.constructor()

		err = dynamodbattribute.UnmarshalMap(i, &objectInst)
		if err != nil {
			return objectList, fmt.Errorf("cannot unmarshall. Err: %s", err)
		}

		if dynoInst.sortName == "" || objectInst.CodeValue() == objectInst.SortValue() {
			objectList = append(objectList, objectInst)
		}
	}

	return objectList, nil
}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) getObjectByCodeFromDB(codeValue string) K {
	var input *dynamodb.GetItemInput

	if dynoInst.sortName == "" {
		input = &dynamodb.GetItemInput{
			TableName: aws.String(dynoInst.tableName),
			Key: map[string]*dynamodb.AttributeValue{
				dynoInst.codeName: {
					S: aws.String(codeValue),
				},
			},
		}
	} else {
		input = &dynamodb.GetItemInput{
			TableName: aws.String(dynoInst.tableName),
			Key: map[string]*dynamodb.AttributeValue{
				dynoInst.codeName: {
					S: aws.String(codeValue),
				},
				dynoInst.sortName: {
					S: aws.String(codeValue), //tables with a sort index store the key as sort value for the actual objects.
				},
			},
		}
	}

	result, err := Conn.GetItem(input)

	objectInst := dynoInst.constructor()
	if err != nil {
		log.Fatalf("Error attempting to retrieve object with %s = %s. Error: %s", dynoInst.codeName, codeValue, err)
	} else if result.Item == nil {
		log.Printf("Object with %s = %s could not be found in %s.", dynoInst.codeName, codeValue, dynoInst.tableName)
	} else {
		err = dynamodbattribute.UnmarshalMap(result.Item, &objectInst)
		if err != nil {
			log.Fatalf("Error unmarshalling result for object with %s = %s. Error: %s", dynoInst.codeName, codeValue, err)
		}
	}

	return objectInst
}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) initCache() {
	dynoInst.cache = make(map[string]K)
	cacheList, err := dynoInst.getObjectListFromDB()

	if err != nil {
		log.Fatalf("object cache cannot be initialised. Error: %v", err)
	}

	for _, objectInst := range cacheList {
		dynoInst.cache[objectInst.CodeValue()] = objectInst
	}

	if dynoInst.cacheMap != nil {
		cacheArray := dynoInst.cacheMap(cacheList)

		for _, cacheMapElement := range cacheArray {
			var cacheElementArray []db.CacheMapElement

			if currentElement, hasKey := cacheMap[cacheMapElement.Tag]; hasKey {
				cacheElementArray = currentElement
			}

			cacheMapElement.Type = dynoInst.tableName
			cacheMap[cacheMapElement.Tag] = append(cacheElementArray, cacheMapElement)
		}
	}
}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) Config(tableName string, codeName string, keepCache bool, constructor func() K, getCacheMap func([]K) []db.CacheMapElement) {
	dynoInst.tableName = tableName
	dynoInst.codeName = codeName
	dynoInst.keepCache = keepCache
	dynoInst.constructor = constructor
	dynoInst.cacheMap = getCacheMap

	if keepCache {
		dynoInst.initCache()
	}
}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) SetSortName(sortName string) {
	if dynoInst.keepCache {
		log.Fatalf("A sort name for adapter %s cannot be set as this adapter has been already configured to keep a cache. Set the sort name before calling Config.", dynoInst.tableName)
	} else {
		dynoInst.sortName = sortName
	}
}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) DeleteObjectByCode(codeValue string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			dynoInst.codeName: {
				S: aws.String(codeValue),
			},
		},
		TableName: aws.String(dynoInst.tableName),
	}

	_, err := Conn.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("cannot delete object with %s %s. Err: %s", dynoInst.codeName, codeValue, err)
	} else {
		if dynoInst.keepCache {
			delete(dynoInst.cache, codeValue)
		}
	}

	return nil
}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) DeleteObject(objectInst K) error {
	return dynoInst.DeleteObjectByCode(objectInst.CodeValue())
}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) GetObjectList() ([]K, error) {
	var objectList []K
	var err error

	if dynoInst.keepCache {
		objectList = []K{}
		for _, objectInst := range dynoInst.cache {
			objectList = append(objectList, objectInst)
		}
	} else {
		objectList, err = dynoInst.getObjectListFromDB()
	}

	return objectList, err
}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) GetObjectListJSON() ([]byte, error) {
	objectList, errList := dynoInst.GetObjectList()
	if errList != nil {
		return nil, errList
	}

	out, err := json.MarshalIndent(objectList, db.JSON_PREFIX, db.JSON_INDENT)
	if err != nil {
		return nil, fmt.Errorf("got error when trying to return object list from table %s as API response. Error: %s", dynoInst.tableName, err)
	}

	return out, nil
}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) GetObjectByCodeJSON(codeValue string) ([]byte, error) {
	objectInst, getErr := dynoInst.GetObjectByCode(codeValue)

	if getErr == nil {
		out, err := json.MarshalIndent(objectInst, db.JSON_PREFIX, db.JSON_INDENT)

		if err != nil {
			marshallErr := fmt.Errorf("got error when trying to return object list from table %s as API response. Error: %s", dynoInst.tableName, err)
			return nil, marshallErr
		} else {
			return out, nil
		}
	} else {
		return nil, getErr
	}

}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) GetObjectByCode(codeValue string) (K, error) {
	var objectInst K
	var notFoundError error
	var isFound bool = false

	if dynoInst.keepCache {
		if objectInstVal, ok := dynoInst.cache[codeValue]; ok {
			objectInst = objectInstVal
			isFound = true
		}
	}

	if !isFound {
		objectInst = dynoInst.getObjectByCodeFromDB(codeValue)
		if dynoInst.keepCache && len(objectInst.CodeValue()) > 0 {
			dynoInst.cache[objectInst.CodeValue()] = objectInst
			isFound = true
		}
	}

	if !isFound {
		notFoundError = fmt.Errorf("object with key %s not found", codeValue)
	}

	return objectInst, notFoundError
}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) PutObject(objectInst K) error {
	objectMarshalled, err := dynamodbattribute.MarshalMap(objectInst)

	if err != nil {
		return fmt.Errorf("got error marshalling map for object with key = %s. Error: %s", objectInst.CodeValue(), err)
	}

	input := &dynamodb.PutItemInput{
		Item:      objectMarshalled,
		TableName: aws.String(dynoInst.tableName),
	}

	_, err = Conn.PutItem(input)
	if err != nil {
		return fmt.Errorf("got error calling PutItem for object with key = %s into %s. Error: %s", objectInst.CodeValue(), dynoInst.tableName, err)
	}

	if dynoInst.keepCache {
		dynoInst.cache[objectInst.CodeValue()] = objectInst
	}

	return nil
}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) PutObjectList(objectList []K) {
	for _, objectInst := range objectList {
		dynoInst.PutObject(objectInst)
	}
}

//----------------------------------------------------------------------------------------
func (dynoInst *Dyno[K]) ResetCache() {
	if !dynoInst.keepCache {
		return
	}

	dynoInst.initCache()
}

//----------------------------------------------------------------------------------------
func PrintCacheMap() {
	for key, elementArray := range cacheMap {
		fmt.Printf("Key %s\n", key)
		for position, element := range elementArray {
			fmt.Printf(" -- Position %d: {Code: %s, Tag: %s, Type: %s}\n", position, element.Code, element.Tag, element.Type)
		}
	}
}

//----------------------------------------------------------------------------------------
func SearchCacheMap(searchKey string) []db.CacheMapElement {
	var exactMatches []db.CacheMapElement
	var approxMatches []db.CacheMapElement
	var allMatches []db.CacheMapElement
	var resultSlice []db.CacheMapElement
	var resultMap map[string]string = make(map[string]string)
	var tempKey string
	searchKeyLower := strings.ToLower(searchKey)

	if elementArray, hasElementArray := cacheMap[searchKeyLower]; hasElementArray {
		exactMatches = elementArray
	}

	for key, cacheMapElement := range cacheMap {
		if key != searchKeyLower && strings.Contains(key, searchKeyLower) {
			approxMatches = append(approxMatches, cacheMapElement...)
		}
	}

	allMatches = append(allMatches, exactMatches...)
	allMatches = append(allMatches, approxMatches...)

	for _, matchElement := range allMatches {
		tempKey = matchElement.Code + "|" + matchElement.Type
		if _, hasKey := resultMap[tempKey]; !hasKey {
			resultMap[tempKey] = tempKey
			resultSlice = append(resultSlice, matchElement)
		}
	}

	return resultSlice
}

//----------------------------------------------------------------------------------------
func SearchCacheMapJSON(searchKey string) []byte {
	allMatches := SearchCacheMap(searchKey)
	out, err := json.MarshalIndent(allMatches, db.JSON_PREFIX, db.JSON_INDENT)

	if err != nil {
		log.Fatalf("Got error when trying to return cache map entries for search key: %s. Error: %s", searchKey, err)
		return nil
	}

	return out
}

//----------------------------------------------------------------------------------------
func CacheMapJSON() []byte {
	var allMatches []CacheMapEntry

	for key, arrayElement := range cacheMap {
		var entry CacheMapEntry
		entry.Key = key
		entry.Elements = arrayElement

		allMatches = append(allMatches, entry)
	}

	out, err := json.MarshalIndent(allMatches, db.JSON_PREFIX, db.JSON_INDENT)

	if err != nil {
		log.Fatalf("Got error when trying to return full cache map %s", err)
		return nil
	}

	return out
}

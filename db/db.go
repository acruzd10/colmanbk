package db

import (
	"colmanback/objects"
	"encoding/json"
	"log"
	"strings"
)

const (
	JSON_PREFIX = ""
	JSON_INDENT = "    "
)

type CacheMapElement struct {
	Code string `json:"code"`
	Tag  string `json:"tag"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type Adapter[K objects.Object] interface {
	//Config
	Config(tableName string, codeName string, keepCache bool, constructor func() K, getCacheMap func([]K) []CacheMapElement)

	//Raw Operations
	DeleteObjectByCode(codeValue string) error
	DeleteObjectByCodeAndSort(codeValue string, sortValue string) error
	DeleteObject(objectInst K) error
	GetObjectList() ([]K, error)
	GetObjectByCode(codeValue string) (K, error)
	GetObjectByCodeJSON(codeValue string) ([]byte, error)
	GetObjectListJSON() ([]byte, error)
	GetSortKeyList(codeValue string) ([]string, error)
	PutObject(objectInst K) error
	PutObjectList(objectList []K)
	ResetCache()
}

//----------------------------------------------------------------------------------------
func FromJson(objectInst objects.Object, jsonInst []byte) {
	err := json.Unmarshal(jsonInst, objectInst)
	if err != nil {
		log.Fatalf("Cannot unmarshall %s into object %s", string(jsonInst), objectInst.ToString())
	}
}

//----------------------------------------------------------------------------------------
func ToJson(objectInst objects.Object) []byte {
	out, err := json.MarshalIndent(objectInst, JSON_PREFIX, JSON_INDENT)

	if err != nil {
		log.Fatalf("Cannot marshal object:\n %s", objectInst.ToString())
		return nil
	}

	return out
}

//----------------------------------------------------------------------------------------
func AddToCacheMap(cacheMap []CacheMapElement, tag string, code string, name string) []CacheMapElement {
	if len(tag) > 0 {
		var cacheMapElement CacheMapElement

		cacheMapElement.Tag = strings.ToLower(tag)
		cacheMapElement.Code = code
		cacheMapElement.Name = name

		cacheMap = append(cacheMap, cacheMapElement)
	}

	return cacheMap
}

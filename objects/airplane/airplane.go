package airplane

import (
	"colmanback/api_util"
	"colmanback/db"
	"colmanback/db/dyno"
	"colmanback/objects/airplanemake"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Airplane struct {
	Code     string                     `json:"code"`
	Name     string                     `json:"name"`
	Iata     string                     `json:"iata"`
	Icao     string                     `json:"icao"`
	Make     string                     `json:"make"`
	MakeInst *airplanemake.AirplaneMake `json:"makeDetails,omitempty"`
}

var AdapterInst db.Adapter[*Airplane]

//----------------------------------------------------------------------------------------
func (airplaneInst *Airplane) CodeValue() string {
	return airplaneInst.Code
}

//----------------------------------------------------------------------------------------
func (airplaneInst *Airplane) SortValue() string {
	return ""
}

//----------------------------------------------------------------------------------------
func (airplaneInst *Airplane) ToString() string {
	str := fmt.Sprintf(`
	----------------------
	  Code ......: %s
	  Name ......: %s	
	  Iata ......: %s
	  Icao ......: %s
	  Make ......: %s`,
		airplaneInst.Code,
		airplaneInst.Name,
		airplaneInst.Iata,
		airplaneInst.Icao,
		airplaneInst.Make)

	return str
}

//----------------------------------------------------------------------------------------
func (airplaneInst *Airplane) FromJson(jsonInst []byte) {
	db.FromJson(airplaneInst, jsonInst)
}

//----------------------------------------------------------------------------------------
func (airplaneInst *Airplane) Print() {
	fmt.Println(airplaneInst.ToString())
}

//----------------------------------------------------------------------------------------
func GetCacheMap(airplaneList []*Airplane) []db.CacheMapElement {
	var cacheMap []db.CacheMapElement

	for _, airplane := range airplaneList {
		cacheMap = db.AddToCacheMap(cacheMap, airplane.Iata, airplane.Code, airplane.Name)
		cacheMap = db.AddToCacheMap(cacheMap, airplane.Icao, airplane.Code, airplane.Name)
		cacheMap = db.AddToCacheMap(cacheMap, airplane.Name, airplane.Code, airplane.Name)
	}

	return cacheMap
}

//----------------------------------------------------------------------------------------
func (airplaneInst *Airplane) Put() {
	makeInst := airplaneInst.MakeInst

	airplaneInst.MakeInst = nil
	AdapterInst.PutObject(airplaneInst)

	airplaneInst.MakeInst = makeInst
}

//----------------------------------------------------------------------------------------
func (airplaneInst *Airplane) Delete() {
	AdapterInst.DeleteObject(airplaneInst)
}

//----------------------------------------------------------------------------------------
func (airplaneInst *Airplane) WriteObject(writer http.ResponseWriter, request *http.Request) {
	api_util.WriteObject(airplaneInst, writer, request)
}

//----------------------------------------------------------------------------------------
func (airplaneInst *Airplane) InitRefObjs() {
	var err error

	if len(airplaneInst.Make) != 0 {
		airplaneInst.MakeInst, err = airplanemake.GetByCode(airplaneInst.Make)
		if err != nil {
			panic(fmt.Sprintf("Error initialising airplane make with code %s for airplane %s", airplaneInst.Make, airplaneInst.Code))
		}
	}
}

//----------------------------------------------------------------------------------------
func GetList() ([]*Airplane, error) {
	return AdapterInst.GetObjectList()
}

//----------------------------------------------------------------------------------------
func GetByCode(code string) (*Airplane, error) {
	airplaneInst, err := AdapterInst.GetObjectByCode(code)

	if airplaneInst != nil {
		airplaneInst.InitRefObjs()
	}

	return airplaneInst, err
}

//----------------------------------------------------------------------------------------
func LoadAirplaneList(svc *dynamodb.DynamoDB, airplaneList []Airplane) {
	var airplanePointerList []*Airplane = []*Airplane{}
	var airplaneMakeMap map[string]string = make(map[string]string)

	airplaneMakeList, err := airplanemake.GetList()
	if err != nil {
		panic("Error getting airplaneMakeList")
	}

	for _, airplaneMake := range airplaneMakeList {
		airplaneMakeMap[airplaneMake.Code] = airplaneMake.Name
	}

	for _, airplaneInst := range airplaneList {
		var newAirplaneInst Airplane = airplaneInst
		var name string = strings.ToLower(newAirplaneInst.Name)

		for key := range airplaneMakeMap {
			if strings.Contains(name, key) {
				newAirplaneInst.Make = key
				break
			}
		}

		airplanePointerList = append(airplanePointerList, &newAirplaneInst)
	}

	AdapterInst.PutObjectList(airplanePointerList)
}

//----------------------------------------------------------------------------------------
func ObjectFactory() *Airplane {
	var airplaneInst Airplane = Airplane{}

	return &airplaneInst
}

//----------------------------------------------------------------------------------------
func InitConn() {
	dynoInstAirline := &dyno.Dyno[*Airplane]{}
	dynoInstAirline.Config("airplane", "code", true, ObjectFactory, GetCacheMap)
	AdapterInst = dynoInstAirline
}

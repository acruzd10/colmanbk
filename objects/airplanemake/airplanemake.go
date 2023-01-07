package airplanemake

import (
	"colmanback/api_util"
	"colmanback/db"
	"colmanback/db/dyno"
	"fmt"
	"net/http"
)

type AirplaneMake struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	Abbreviation string `json:"abbreviation"`
}

var AdapterInst db.Adapter[*AirplaneMake]

//----------------------------------------------------------------------------------------
func (airplaneMakeInst *AirplaneMake) CodeValue() string {
	return airplaneMakeInst.Code
}

//----------------------------------------------------------------------------------------
func (airplaneMakeInst *AirplaneMake) SortValue() string {
	return ""
}

//----------------------------------------------------------------------------------------
func (airplaneMakeInst *AirplaneMake) FromJson(jsonInst []byte) {
	db.FromJson(airplaneMakeInst, jsonInst)
}

//----------------------------------------------------------------------------------------
func (airplaneMakeInst *AirplaneMake) ToString() string {
	str := fmt.Sprintf(`
	----------------------
	  Code ......: %s
	  Name ......: %s	
	  Country ...: %s
	  Abbrev. ...: %s`,
		airplaneMakeInst.Code,
		airplaneMakeInst.Name,
		airplaneMakeInst.Country,
		airplaneMakeInst.Abbreviation)

	return str
}

//----------------------------------------------------------------------------------------
func (airplaneMakeInst *AirplaneMake) Print() {
	fmt.Println(airplaneMakeInst.ToString())
}

//----------------------------------------------------------------------------------------
func (airplaneMakeInst *AirplaneMake) WriteObject(writer http.ResponseWriter, request *http.Request) {
	api_util.WriteObject(airplaneMakeInst, writer, request)
}

//----------------------------------------------------------------------------------------
func (airplaneMakeInst *AirplaneMake) Put() {
	AdapterInst.PutObject(airplaneMakeInst)
}

//----------------------------------------------------------------------------------------
func ObjectFactory() *AirplaneMake {
	var airplaneMakeInst AirplaneMake = AirplaneMake{}

	return &airplaneMakeInst
}

//----------------------------------------------------------------------------------------
func GetCacheMap(airplaneMakeList []*AirplaneMake) []db.CacheMapElement {
	var cacheMap []db.CacheMapElement

	for _, airplaneMake := range airplaneMakeList {
		cacheMap = db.AddToCacheMap(cacheMap, airplaneMake.Abbreviation, airplaneMake.Code, airplaneMake.Name)
		cacheMap = db.AddToCacheMap(cacheMap, airplaneMake.Name, airplaneMake.Code, airplaneMake.Name)
	}

	return cacheMap
}

//----------------------------------------------------------------------------------------
func (airplaneMakeInst *AirplaneMake) Delete() {
	AdapterInst.DeleteObject(airplaneMakeInst)
}

//----------------------------------------------------------------------------------------
func GetList() ([]*AirplaneMake, error) {
	return AdapterInst.GetObjectList()
}

//----------------------------------------------------------------------------------------
func GetByCode(code string) (*AirplaneMake, error) {
	return AdapterInst.GetObjectByCode(code)
}

//----------------------------------------------------------------------------------------
func LoadList(countryList []AirplaneMake) {
	var makePointerList []*AirplaneMake = []*AirplaneMake{}

	for _, airplaneMakeInst := range countryList {
		var newAirplaneMakeInst AirplaneMake = airplaneMakeInst
		makePointerList = append(makePointerList, &newAirplaneMakeInst)
	}

	AdapterInst.PutObjectList(makePointerList)
}

//----------------------------------------------------------------------------------------
func InitConn() {
	dynoInstAirplaneMake := &dyno.Dyno[*AirplaneMake]{}
	dynoInstAirplaneMake.Config("airplanemake", "code", true, ObjectFactory, GetCacheMap)
	AdapterInst = dynoInstAirplaneMake
}

package airplanemake

import (
	"colmanback/api_util"
	"colmanback/db"
	"colmanback/db/dyno"
	"colmanback/objects/country"
	"fmt"
	"net/http"
)

type AirplaneMake struct {
	Code         string           `json:"code"`
	Name         string           `json:"name"`
	Abbreviation string           `json:"abbreviation"`
	Country      string           `json:"country"`
	CountryInst  *country.Country `json:"countryDetails,omitempty"`
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
	countryInst := airplaneMakeInst.CountryInst

	airplaneMakeInst.CountryInst = nil
	AdapterInst.PutObject(airplaneMakeInst)

	airplaneMakeInst.CountryInst = countryInst
}

//----------------------------------------------------------------------------------------
func (airplaneMakeInst *AirplaneMake) InitRefObjs() {
	var err error

	if len(airplaneMakeInst.Country) != 0 {
		airplaneMakeInst.CountryInst, err = country.GetCountryByISO(airplaneMakeInst.Country)
		if err != nil {
			panic(fmt.Sprintf("Error initialising country with code %s for model %s", airplaneMakeInst.Country, airplaneMakeInst.Code))
		}
	}
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
	airplaneMakeInst, err := AdapterInst.GetObjectByCode(code)

	if airplaneMakeInst != nil {
		airplaneMakeInst.InitRefObjs()
	}

	return airplaneMakeInst, err
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

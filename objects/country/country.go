package country

import (
	"colmanback/api_util"
	"colmanback/db"
	"colmanback/db/dyno"
	"fmt"
	"net/http"
)

type Country struct {
	Code      string `json:"code"`
	Continent string `json:"continent"`
	Name      string `json:"name"`
}

var AdapterInst db.Adapter[*Country]

//----------------------------------------------------------------------------------------
func (countryInst *Country) CodeValue() string {
	return countryInst.Code
}

//----------------------------------------------------------------------------------------
func (countryInst *Country) SortValue() string {
	return ""
}

//----------------------------------------------------------------------------------------
func (countryInst *Country) FromJson(jsonInst []byte) {
	db.FromJson(countryInst, jsonInst)
}

//----------------------------------------------------------------------------------------
func (countryInst *Country) ToString() string {
	str := fmt.Sprintf(`
	----------------------
	  Code ......: %s
	  Name ......: %s
	  Continent .: %s`,
		countryInst.Code,
		countryInst.Name,
		countryInst.Continent)

	return str
}

//----------------------------------------------------------------------------------------
func (countryInst *Country) Print() {
	fmt.Println(countryInst.ToString())
}

//----------------------------------------------------------------------------------------
func (countryInst *Country) WriteObject(writer http.ResponseWriter, request *http.Request) {
	api_util.WriteObject(countryInst, writer, request)
}

//----------------------------------------------------------------------------------------
func CountryFactory() *Country {
	var countryInst Country = Country{}

	return &countryInst
}

//----------------------------------------------------------------------------------------
func GetCacheMap(countryList []*Country) []db.CacheMapElement {
	var cacheMap []db.CacheMapElement

	for _, country := range countryList {
		cacheMap = db.AddToCacheMap(cacheMap, country.Code, country.Code, country.Name)
		cacheMap = db.AddToCacheMap(cacheMap, country.Name, country.Code, country.Name)
	}

	return cacheMap
}

//----------------------------------------------------------------------------------------
func (countryInst *Country) Delete() {
	AdapterInst.DeleteObject(countryInst)
}

//----------------------------------------------------------------------------------------
func (countryInst *Country) Put() {
	panic("Method not implemented.")
}

//----------------------------------------------------------------------------------------
func GetCountryList() ([]*Country, error) {
	return AdapterInst.GetObjectList()
}

//----------------------------------------------------------------------------------------
func GetCountryByISO(code string) (*Country, error) {
	return AdapterInst.GetObjectByCode(code)
}

//----------------------------------------------------------------------------------------
func LoadCountryList(countryList []Country) {
	var countryPointerList []*Country = []*Country{}

	for _, countryInst := range countryList {
		var newCountryInst Country = countryInst
		countryPointerList = append(countryPointerList, &newCountryInst)
	}

	AdapterInst.PutObjectList(countryPointerList)
}

//----------------------------------------------------------------------------------------
func InitConn() {
	dynoInstCountry := &dyno.Dyno[*Country]{}
	dynoInstCountry.Config("country", "code", true, CountryFactory, GetCacheMap)
	AdapterInst = dynoInstCountry
}

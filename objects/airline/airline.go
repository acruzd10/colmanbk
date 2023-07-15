package airline

import (
	"colmanback/api_util"
	"colmanback/db"
	"colmanback/db/dyno"
	"colmanback/objects/country"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	ICAO_PREFIX = "icao:"
	IATA_PREFIX = "iata:"
)

type Airline struct {
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Iata        string           `json:"iata"`
	Icao        string           `json:"icao"`
	Callsign    string           `json:"callsign"`
	Country     string           `json:"country"`
	CountryInst *country.Country `json:"countryDetails,omitempty"`
}

var AdapterInst db.Adapter[*Airline]

//----------------------------------------------------------------------------------------
func (airlineInst *Airline) makeCode() {
	if airlineInst.Iata != "" {
		airlineInst.Code = IATA_PREFIX + airlineInst.Iata
	} else if airlineInst.Icao != "" {
		airlineInst.Code = ICAO_PREFIX + airlineInst.Icao
	} else {
		log.Panicf("Cannot make code for the following airline instance: %s", airlineInst.ToString())
	}
}

//----------------------------------------------------------------------------------------
func (airlineInst *Airline) CodeValue() string {
	return airlineInst.Code
}

//----------------------------------------------------------------------------------------
func (airlineInst *Airline) SortValue() string {
	return ""
}

//----------------------------------------------------------------------------------------
func (airlineInst *Airline) FromJson(jsonInst []byte) {
	db.FromJson(airlineInst, jsonInst)
}

//----------------------------------------------------------------------------------------
func (airlineInst *Airline) ToString() string {
	str := fmt.Sprintf(`
	----------------------
	  Code ......: %s
	  Name ......: %s
	  Iata ......: %s
	  Icao ......: %s
	  Callsign ..: %s
	  Country ...: %s`,
		airlineInst.Code,
		airlineInst.Name,
		airlineInst.Iata,
		airlineInst.Icao,
		airlineInst.Callsign,
		airlineInst.Country)

	return str
}

//----------------------------------------------------------------------------------------
func (airlineInst *Airline) Print() {
	fmt.Println(airlineInst.ToString())
}

//----------------------------------------------------------------------------------------
func (airlineInst *Airline) WriteObject(writer http.ResponseWriter, request *http.Request) {
	api_util.WriteObject(airlineInst, writer, request)
}

//----------------------------------------------------------------------------------------
func (airlineInst *Airline) Delete() {
	AdapterInst.DeleteObject(airlineInst)
}

//----------------------------------------------------------------------------------------
func (airlineInst *Airline) Put() {
	countryInst := airlineInst.CountryInst

	if airlineInst.Code == "" {
		airlineInst.makeCode()
	}

	airlineInst.CountryInst = nil
	AdapterInst.PutObject(airlineInst)

	airlineInst.CountryInst = countryInst
}

//----------------------------------------------------------------------------------------
func ObjectFactory() *Airline {
	var airlineInst Airline = Airline{}

	return &airlineInst
}

//----------------------------------------------------------------------------------------
func (airlineInst *Airline) InitRefObjs() {
	var err error

	if len(airlineInst.Country) != 0 {
		airlineInst.CountryInst, err = country.GetCountryByISO(airlineInst.Country)
		if err != nil {
			panic(fmt.Sprintf("Error initialising country with code %s for model %s", airlineInst.Country, airlineInst.Code))
		}
	}
}

//----------------------------------------------------------------------------------------
func GetCacheMap(airlineList []*Airline) []db.CacheMapElement {
	var cacheMap []db.CacheMapElement

	for _, airline := range airlineList {
		cacheMap = db.AddToCacheMap(cacheMap, airline.Iata, airline.Code, airline.Name)
		cacheMap = db.AddToCacheMap(cacheMap, airline.Icao, airline.Code, airline.Name)
		cacheMap = db.AddToCacheMap(cacheMap, airline.Name, airline.Code, airline.Name)
	}

	return cacheMap
}

//----------------------------------------------------------------------------------------
func GetList() ([]*Airline, error) {
	return AdapterInst.GetObjectList()
}

//----------------------------------------------------------------------------------------
func getAirlineByCodeIntl(code string) (*Airline, error) {
	return AdapterInst.GetObjectByCode(code)
}

//----------------------------------------------------------------------------------------
func GetByCode(code string) (*Airline, error) {
	var searchErr error = nil
	airlineInst, apiErr := getAirlineByCodeIntl(code)

	if apiErr != nil {
		if !strings.HasPrefix(code, IATA_PREFIX) &&
			!strings.HasPrefix(code, ICAO_PREFIX) {

			iataCode := IATA_PREFIX + code
			airlineInst, apiErr = getAirlineByCodeIntl(iataCode)

			if apiErr != nil {
				icaoCode := ICAO_PREFIX + code
				airlineInst, apiErr = getAirlineByCodeIntl(icaoCode)

				if apiErr != nil {
					searchErr = fmt.Errorf("airline with code %s could not be found", code)
				}
			}
		} else {
			searchErr = apiErr
		}
	}

	if airlineInst != nil {
		airlineInst.InitRefObjs()
	}

	return airlineInst, searchErr
}

//----------------------------------------------------------------------------------------
func LoadAirlineList(airlineList []Airline) {
	var airlinePointerList []*Airline = []*Airline{}
	var countryList []*country.Country
	var countryMap map[string]string = make(map[string]string)
	var err error

	countryList, err = country.GetCountryList()
	if err != nil {
		log.Fatalf("cannot load the list of countries whilst loading airline list. Error: %v", err)
	}

	for _, countryInst := range countryList {
		countryMap[strings.ToLower(countryInst.Name)] = countryInst.Code
	}

	for _, airlineInst := range airlineList {
		var airlineNewInst *Airline = &Airline{}

		airlineNewInst.Code = strings.ToLower(airlineInst.Code)
		airlineNewInst.Name = strings.ToTitle(airlineInst.Name)
		airlineNewInst.Iata = strings.ToLower(airlineInst.Iata)
		airlineNewInst.Icao = strings.ToLower(airlineInst.Icao)
		airlineNewInst.Callsign = strings.ToLower(airlineInst.Callsign)

		if iso, found := countryMap[strings.ToLower(airlineInst.Country)]; found {
			airlineNewInst.Country = iso
		} else {
			airlineNewInst.Country = ""
		}

		airlinePointerList = append(airlinePointerList, airlineNewInst)
	}

	AdapterInst.PutObjectList(airlinePointerList)
}

//----------------------------------------------------------------------------------------
func InitConn() {
	dynoInstAirline := &dyno.Dyno[*Airline]{}
	dynoInstAirline.Config("airline", "code", true, ObjectFactory, GetCacheMap)
	AdapterInst = dynoInstAirline
}

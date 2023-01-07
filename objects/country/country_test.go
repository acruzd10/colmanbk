package country

import (
	"colmanback/db/dyno"
	"colmanback/test_util"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	countryCount     = 250
	countryISO       = "cn"
	countryContinent = "Asia"
	countryName      = "China"
)

func initDyno() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dyno.Conn = dynamodb.New(sess)

	//Airline adapter
	dynoInstAirline := &dyno.Dyno[*Country]{}
	dynoInstAirline.Config("country", "code", true, CountryFactory, GetCacheMap)
	AdapterInst = dynoInstAirline
}

func chkList(t *testing.T) {
	countryList, err := GetCountryList()

	if err != nil {
		t.Errorf("Error when retrieving list of countries. Error: %v", err)
	} else {
		countryNumber := len(countryList)
		if countryNumber != countryCount {
			t.Errorf("Not all the expected countries were found. Expected %d but got %d", countryCount, countryNumber)
		}
	}
}

func chkSingleCountry(t *testing.T) {
	countryObj, err := GetCountryByISO("cn")

	if err != nil {
		t.Errorf("Error when retrieving single country. Error: %v", err)
	} else {
		test_util.CheckField(t, "code", countryISO, countryObj.Code)
		test_util.CheckField(t, "name", countryName, countryObj.Name)
		test_util.CheckField(t, "continent", countryContinent, countryObj.Continent)
	}
}

func TestAirplane(t *testing.T) {
	initDyno()

	t.Log("Ensure that the list has the expected number of countries.")
	chkList(t)

	t.Log("Ensure that a country can be retrieved by code.")
	chkSingleCountry(t)
}

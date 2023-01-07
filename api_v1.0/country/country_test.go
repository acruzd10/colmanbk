package country

import (
	"colmanback/db/dyno"
	"colmanback/objects/country"
	"colmanback/test_util"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
)

const (
	countryCount     = 250
	countryISO       = "cn"
	countryContinent = "Asia"
	countryName      = "China"
)

func chkSingleCountry(t *testing.T, router *mux.Router) {
	req, err := http.NewRequest(http.MethodGet, ApiURL+strings.Replace(ResourceURL, "{"+ObjectID+"}", countryISO, 1), nil)
	var countryInst country.Country

	if err != nil {
		t.Errorf("Error preparing single country request. Error: %v", err)
	} else {
		resp := test_util.ExecuteRequest(router, req)
		if resp.Code != http.StatusOK {
			t.Errorf("An unexpected response code has been returned: %d", resp.Code)
		} else {
			unmarshallErr := json.Unmarshal(resp.Body.Bytes(), &countryInst)
			if unmarshallErr != nil {
				t.Errorf("An error has occurred whilst unmarshalling country: %v", unmarshallErr)
			} else {
				test_util.CheckField(t, "code", countryISO, countryInst.Code)
				test_util.CheckField(t, "name", countryName, countryInst.Name)
				test_util.CheckField(t, "continent", countryContinent, countryInst.Continent)
			}
		}
	}
}

func chkListCountry(t *testing.T, router *mux.Router) {
	req, err := http.NewRequest(http.MethodGet, ApiURL+BaseURL, nil)
	var countryList []country.Country

	if err != nil {
		t.Errorf("Error preparing the full country list get. Error: %v", err)
	} else {
		resp := test_util.ExecuteRequest(router, req)
		if resp.Code != http.StatusOK {
			t.Errorf("An unexpected response code has been returned: %d", resp.Code)
		} else {
			unmarshallErr := json.Unmarshal(resp.Body.Bytes(), &countryList)
			if unmarshallErr != nil {
				t.Errorf("An error has occurred whilst unmarshalling the list of countries: %v", unmarshallErr)
			} else {
				countryNumber := len(countryList)
				if countryNumber != countryCount {
					t.Errorf("Expected %d countries in the list but got %d", countryCount, countryNumber)
				}
			}
		}
	}
}

func TestCountry(t *testing.T) {
	router := mux.NewRouter()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dyno.Conn = dynamodb.New(sess)
	country.InitConn()

	InitRouter(router)

	t.Log("Ensure a single country can be retrieved")
	chkSingleCountry(t, router)

	t.Log("Ensure the list of countries can he retrieved")
	chkListCountry(t, router)
}

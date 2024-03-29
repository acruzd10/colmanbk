package model

import (
	"bytes"
	airlineAPI "colmanback/api_v1.0/airline"
	airplaneAPI "colmanback/api_v1.0/airplane"
	airplaneMakeAPI "colmanback/api_v1.0/airplanemake"
	modelMakeAPI "colmanback/api_v1.0/modelmake"
	"colmanback/db"
	"colmanback/db/dyno"
	"colmanback/objects"
	airlineObject "colmanback/objects/airline"
	airplaneObject "colmanback/objects/airplane"
	airplaneMakeObject "colmanback/objects/airplanemake"
	countryObject "colmanback/objects/country"
	modelObject "colmanback/objects/model"
	modelMakeObject "colmanback/objects/modelmake"
	"colmanback/test_util"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
)

const (
	modelMakeCode = "modelTestModelMakeCode"
	modelMakeName = modelMakeCode + "Name"

	airlineCode    = "modelTestAirline"
	airlineName    = airlineCode + "Name"
	airlineIATA    = "modelTestAA"
	airlineCountry = "gb"

	airplaneMakeCode    = "modelTestBus"
	airplaneMakeName    = airplaneMakeCode + "Name"
	airplaneMakeCountry = "fr"

	airplaneCode = "modelTestAirplane"
	airplaneName = airplaneCode + "Name"
	airplaneIATA = "modelTest77"

	modelIsCargo         = true
	modelIsOldLivery     = true
	modelIsSpecialLivery = true
	modelNotes           = "modelTest NOTES"
	modelReg             = "TE-STY"
	modelScale           = objects.Scale1400
)

var router *mux.Router
var modelCode string

func boolToString(value bool) string {
	if value {
		return "true"
	}

	return "false"
}

func makeModelMakeInstance() modelMakeObject.ModelMake {
	var objectInst modelMakeObject.ModelMake

	objectInst.Code = modelMakeCode
	objectInst.Name = modelMakeName

	return objectInst
}

func makeAirlineInstance() airlineObject.Airline {
	var objectInst airlineObject.Airline

	objectInst.Code = airlineCode
	objectInst.Name = airlineName
	objectInst.Iata = airlineIATA
	objectInst.Country = airlineCountry

	return objectInst
}

func makeAirplaneMakeInstance() airplaneMakeObject.AirplaneMake {
	var objectInst airplaneMakeObject.AirplaneMake

	objectInst.Code = airplaneMakeCode
	objectInst.Name = airplaneMakeName
	objectInst.Country = airplaneMakeCountry

	return objectInst
}

func makeAirplaneInstance() airplaneObject.Airplane {
	var objectInst airplaneObject.Airplane

	objectInst.Code = airplaneCode
	objectInst.Make = airplaneMakeCode
	objectInst.Name = airplaneName
	objectInst.Iata = airplaneIATA

	return objectInst
}

func makeModelInstance() modelObject.Model {
	var objectInst modelObject.Model

	objectInst.Airline = airlineCode
	objectInst.Airplane = airplaneCode
	objectInst.IsCargo = modelIsCargo
	objectInst.IsOldLivery = modelIsOldLivery
	objectInst.IsSpecialLivery = modelIsSpecialLivery
	objectInst.ModelMake = modelMakeCode
	objectInst.Reg = modelReg
	objectInst.Scale = modelScale

	objectInst.Notes = modelNotes

	return objectInst
}

func createModelMake(t *testing.T) {
	objectInst := makeModelMakeInstance()

	modelMakeObject.InitConn()
	modelMakeAPI.InitRouter(router)

	jsonString := string(db.ToJson(&objectInst))
	test_util.CheckPut(t, router, jsonString, modelMakeAPI.ApiURL+modelMakeAPI.BaseURL)
}

func createAirline(t *testing.T) {
	objectInst := makeAirlineInstance()

	airlineObject.InitConn()
	airlineAPI.InitRouter(router)

	jsonString := string(db.ToJson(&objectInst))
	test_util.CheckPut(t, router, jsonString, airlineAPI.ApiURL+airlineAPI.BaseURL)
}

func createAirplaneMake(t *testing.T) {
	objectInst := makeAirplaneMakeInstance()

	airplaneMakeObject.InitConn()
	airplaneMakeAPI.InitRouter(router)

	jsonString := string(db.ToJson(&objectInst))
	test_util.CheckPut(t, router, jsonString, airplaneMakeAPI.ApiURL+airplaneMakeAPI.BaseURL)
}

func createAirplane(t *testing.T) {
	objectInst := makeAirplaneInstance()

	airplaneObject.InitConn()
	airplaneAPI.InitRouter(router)

	jsonString := string(db.ToJson(&objectInst))
	test_util.CheckPut(t, router, jsonString, airplaneAPI.ApiURL+airplaneAPI.BaseURL)
}

func createModel(t *testing.T) {
	objectInst := makeModelInstance()
	var newObjectInst modelObject.Model

	jsonString := string(db.ToJson(&objectInst))
	req, err := http.NewRequest(http.MethodPut, ApiURL+BaseURL, bytes.NewBuffer([]byte(jsonString)))

	if err != nil {
		t.Errorf("An error has been reported when preparing the put request for %s. Error: %v\n", jsonString, err)
	} else {
		resp := test_util.ExecuteRequest(router, req)
		if resp.Code != http.StatusOK {
			t.Errorf("Status code not as expected after put.")
		} else {
			fmt.Println(resp.Body.String())
			(&newObjectInst).FromJson(resp.Body.Bytes())
			modelCode = newObjectInst.Code
		}
	}
}

func chkDefaultModelInst(t *testing.T, objectInst modelObject.Model, modelCodeLocal string) {
	t.Logf("Checking model: %s", objectInst.ToString())

	test_util.CheckField(t, "Airline", airlineCode, objectInst.Airline)
	test_util.CheckField(t, "Airplane", airplaneCode, objectInst.Airplane)
	test_util.CheckField(t, "isCargo", boolToString(modelIsCargo), boolToString(objectInst.IsCargo))
	test_util.CheckField(t, "isOldLivery", boolToString(modelIsOldLivery), boolToString(objectInst.IsOldLivery))
	test_util.CheckField(t, "isSpecialLivery", boolToString(modelIsSpecialLivery), boolToString(objectInst.IsSpecialLivery))
	test_util.CheckField(t, "modelMake", modelMakeCode, objectInst.ModelMake)
	test_util.CheckField(t, "reg", modelReg, objectInst.Reg)
	test_util.CheckField(t, "scale", string(modelScale), string(objectInst.Scale))
	test_util.CheckField(t, "code", modelCodeLocal, objectInst.Code)
}

func retrieveModel(t *testing.T) {
	getURL := ApiURL + strings.Replace(ResourceURL, "{"+ObjectID+"}", url.QueryEscape(modelCode), 1)
	var objectInst modelObject.Model

	t.Logf("Checking URL: %s", getURL)
	test_util.CheckExists(t, router, getURL, true)
	req, err := http.NewRequest(http.MethodGet, getURL, nil)
	if err == nil {
		resp := test_util.ExecuteRequest(router, req)
		unmarshallErr := json.Unmarshal(resp.Body.Bytes(), &objectInst)
		if unmarshallErr == nil {
			chkDefaultModelInst(t, objectInst, modelCode)
		} else {
			t.Errorf("An error has occurred whilst unmarshalling object. Error: %v\n", unmarshallErr)
		}
	} else {
		t.Errorf("An error has occurred whist building get request. Error: %v\n", err)
	}
}

func testSetup(t *testing.T) {
	router = mux.NewRouter().SkipClean(true).UseEncodedPath()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dyno.Conn = dynamodb.New(sess)

	modelObject.InitConn()
	countryObject.InitConn() //Needs explicit init because no country is being created.

	InitRouter(router)

	createModelMake(t)
	createAirline(t)
	createAirplaneMake(t)
	createAirplane(t)

	createModel(t)
	retrieveModel(t)
}

func testTearDown(t *testing.T) {
	test_util.CheckDelete(t, router, modelMakeAPI.ApiURL+strings.Replace(modelMakeAPI.ResourceURL, "{"+modelMakeAPI.ObjectID+"}", modelMakeCode, 1), true)
	test_util.CheckDelete(t, router, airlineAPI.ApiURL+strings.Replace(airlineAPI.ResourceURL, "{"+airlineAPI.ObjectID+"}", airlineCode, 1), true)
	test_util.CheckDelete(t, router, airplaneMakeAPI.ApiURL+strings.Replace(airplaneMakeAPI.ResourceURL, "{"+airplaneMakeAPI.ObjectID+"}", airplaneMakeCode, 1), true)
	test_util.CheckDelete(t, router, airplaneAPI.ApiURL+strings.Replace(airplaneAPI.ResourceURL, "{"+airplaneAPI.ObjectID+"}", airplaneCode, 1), true)
	test_util.CheckDelete(t, router, ApiURL+strings.Replace(ResourceURL, "{"+ObjectID+"}", modelCode, 1), true)
}

func TestModel(t *testing.T) {
	testSetup(t)

	testTearDown(t)
}

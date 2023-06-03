package model

import (
	airlineAPI "colmanback/api_v1.0/airline"
	airplaneAPI "colmanback/api_v1.0/airplane"
	airplaneMakeAPI "colmanback/api_v1.0/airplanemake"
	modelMakeAPI "colmanback/api_v1.0/modelmake"
	"colmanback/db"
	"colmanback/db/dyno"
	airlineObject "colmanback/objects/airline"
	airplaneObject "colmanback/objects/airplane"
	airplaneMakeObject "colmanback/objects/airplanemake"
	modelMakeObject "colmanback/objects/modelmake"
	"colmanback/test_util"
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
)

var router *mux.Router

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

func makeAirplane() airplaneObject.Airplane {
	var objectInst airplaneObject.Airplane

	objectInst.Code = airplaneCode
	objectInst.Make = airplaneMakeCode
	objectInst.Name = airplaneName
	objectInst.Iata = airplaneIATA

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
	objectInst := makeAirplane()

	airplaneObject.InitConn()
	airplaneAPI.InitRouter(router)

	jsonString := string(db.ToJson(&objectInst))
	test_util.CheckPut(t, router, jsonString, airplaneAPI.ApiURL+airplaneAPI.BaseURL)
}

func testSetup(t *testing.T) {
	router = mux.NewRouter()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dyno.Conn = dynamodb.New(sess)

	createModelMake(t)
	createAirline(t)
	createAirplaneMake(t)
	createAirplane(t)
}

func testTearDown(t *testing.T) {
	test_util.CheckDelete(t, router, modelMakeAPI.ApiURL+strings.Replace(modelMakeAPI.ResourceURL, "{"+modelMakeAPI.ObjectID+"}", modelMakeCode, 1), true)
	test_util.CheckDelete(t, router, airlineAPI.ApiURL+strings.Replace(airlineAPI.ResourceURL, "{"+airlineAPI.ObjectID+"}", airlineCode, 1), true)
	test_util.CheckDelete(t, router, airplaneMakeAPI.ApiURL+strings.Replace(airplaneMakeAPI.ResourceURL, "{"+airplaneMakeAPI.ObjectID+"}", airplaneMakeCode, 1), true)
	test_util.CheckDelete(t, router, airplaneAPI.ApiURL+strings.Replace(airplaneAPI.ResourceURL, "{"+airplaneAPI.ObjectID+"}", airplaneCode, 1), true)
}

/*
func getTestInstance() *model.Model {
	var objectInst *model.Model = &model.Model{
		Code:      "apiTest",
		ModelMake: "",
	}

	return objectInst
}
*/

func TestModel(t *testing.T) {
	testSetup(t)

	testTearDown(t)
}

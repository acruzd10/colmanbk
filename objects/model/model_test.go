package model

import (
	"colmanback/db/dyno"
	"colmanback/objects"
	"colmanback/objects/airline"
	"colmanback/objects/airplane"
	"colmanback/objects/airplanemake"
	"colmanback/objects/modelmake"
	"colmanback/test_util"
	"encoding/json"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	//Foreign Keys
	modelMakeConst          = "model_model_make_test"
	modelMakeNameConst      = modelMakeConst + "_name"
	airlineConst            = "iata:model_airline_test"
	airlineNameConst        = airlineConst + "_name"
	airlineIataConst        = airlineConst + "_iata"
	airlineCallsignConst    = airlineConst + "_callsign"
	airplaneConst           = "model_airplane_test"
	airplaneIataConst       = airplaneConst + "_iata"
	airplaneNameConst       = airplaneConst + "_name"
	airplaneMakeConst       = "model_airplane_make_test"
	airplaneMakeNameConst   = airplaneMakeConst + "_name"
	airplaneMakeAbbrevConst = airplaneMakeConst + "_abbrev"
	countryConst            = "gb"

	//Properties
	scaleConst           = objects.Scale1400
	regConst             = "modeltest001"
	regNewConst          = regConst + "N"
	notesConst           = "model_test_notes!"
	isCargoConst         = true
	isOldLiveryConst     = false
	isSpecialLiveryConst = true
)

var airlineInst airline.Airline
var airplaneInst airplane.Airplane
var airplaneMakeInst airplanemake.AirplaneMake
var modelMakeInst modelmake.ModelMake

func initDyno() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dyno.Conn = dynamodb.New(sess)

	InitConn()
}

func boolToString(boolValue bool) string {
	if boolValue {
		return "true"
	}

	return "false"
}

func chkModel(t *testing.T, objectInst *Model, emptyCode bool) {
	test_util.CheckField(t, "modelMake", modelMakeConst, objectInst.ModelMake)
	test_util.CheckField(t, "airline", airlineConst, objectInst.Airline)
	test_util.CheckField(t, "airplane", airplaneConst, objectInst.Airplane)
	test_util.CheckField(t, "scale", string(scaleConst), string(objectInst.Scale))
	test_util.CheckField(t, "reg", strings.ToUpper(regConst), strings.ToUpper(objectInst.Reg))
	test_util.CheckField(t, "notes", notesConst, objectInst.Notes)
	test_util.CheckField(t, "isCargo", boolToString(isCargoConst), boolToString(objectInst.IsCargo))
	test_util.CheckField(t, "isOldLivery", boolToString(isOldLiveryConst), boolToString(objectInst.IsOldLivery))
	test_util.CheckField(t, "isSpecialLivery", boolToString(isSpecialLiveryConst), boolToString(objectInst.IsSpecialLivery))

	if emptyCode != (len(objectInst.Code) == 0) {
		t.Errorf("The value of the code %s is not as expected.", objectInst.Code)
	}
}

func testSetup(t *testing.T) {
	initDyno()

	modelmake.InitConn()
	airline.InitConn()
	airplane.InitConn()
	airplanemake.InitConn()

	modelMakeInst = modelmake.ModelMake{}
	modelMakeInst.Code = modelMakeConst
	modelMakeInst.Name = modelMakeNameConst
	(&modelMakeInst).Put()

	airlineInst = airline.Airline{}
	airlineInst.Code = airlineConst
	airlineInst.Callsign = airlineCallsignConst
	airlineInst.Country = countryConst
	airlineInst.Iata = airlineIataConst
	airlineInst.Name = airlineNameConst
	(&airlineInst).Put()

	airplaneMakeInst = airplanemake.AirplaneMake{}
	airplaneMakeInst.Code = airplaneMakeConst
	airplaneMakeInst.Name = airplaneMakeNameConst
	airplaneMakeInst.Abbreviation = airplaneMakeAbbrevConst
	airplaneMakeInst.Country = countryConst
	(&airplaneMakeInst).Put()

	airplaneInst = airplane.Airplane{}
	airplaneInst.Code = airplaneConst
	airplaneInst.Name = airplaneNameConst
	airplaneInst.Iata = airplaneIataConst
	airplaneInst.Make = airplaneMakeInst.Code
	(&airplaneInst).Put()

	t.Log("Test setup completed.")
}

func testTearDown(t *testing.T) {
	modelMakeInst.Delete()
	airlineInst.Delete()
	airplaneMakeInst.Delete()
	airplaneInst.Delete()

	t.Log("Teardown completed.")
}

func TestModel(t *testing.T) {
	if AdapterInst == nil {
		testSetup(t)
	}

	//Create airline from JSON
	objectInst := ObjectFactory()
	objectInstLoad := ObjectFactory()

	objectInst.Airline = airlineConst
	objectInst.Airplane = airplaneConst
	objectInst.ModelMake = modelMakeConst
	objectInst.Scale = scaleConst
	objectInst.Reg = regConst
	objectInst.Notes = notesConst
	objectInst.IsCargo = isCargoConst
	objectInst.IsOldLivery = isOldLiveryConst
	objectInst.IsSpecialLivery = isSpecialLiveryConst

	jsonString, _ := json.MarshalIndent(&objectInst, " ", " ")
	objectInstLoad.FromJson([]byte(jsonString))

	t.Log("Initial check")
	chkModel(t, objectInstLoad, false)

	t.Log("Put check")
	objectInstLoad.Put()
	chkModel(t, objectInstLoad, false)

	modelRetrInst, getRetrErr := GetByCode(objectInstLoad.Code)
	if getRetrErr != nil {
		t.Errorf("Error in get after putting model. %v\n", getRetrErr)
	}
	chkModel(t, modelRetrInst, false)

	t.Log("Update, put again and retrieve")
	modelRetrInst.Reg = regNewConst
	modelRetrInst.Put()
	modelUpdtInst, getUpdtErr := GetByCode(objectInstLoad.Code)
	if getUpdtErr != nil {
		t.Errorf("Error in get after updating airline\n")
	}
	test_util.CheckField(t, "reg", regNewConst, modelUpdtInst.Reg)

	t.Log("Delete and check it's gone!")
	modelUpdtInst.Delete()
	modelEmptyInst, getEmptyErr := GetByCode(objectInstLoad.Code)
	if getEmptyErr == nil {
		t.Errorf("Error for unexistent object not produced when expected. Perhaps the object still exists?\n")
	}
	test_util.CheckField(t, "reg", "", modelEmptyInst.Reg)

	testTearDown(t)
	t.Log("Test for model has finished.")
}

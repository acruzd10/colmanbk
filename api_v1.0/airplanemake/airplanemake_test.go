package airplanemake

import (
	"colmanback/db"
	"colmanback/db/dyno"
	"colmanback/objects/airplanemake"
	"colmanback/test_util"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
)

const (
	codeConst         = "dummyx"
	nameConst         = "Namex dummyx"
	nameNewConst      = "New Namex dummyx"
	countryConst      = "ru"
	abbreviationConst = "dmx"
)

func getTestInstance() airplanemake.AirplaneMake {
	var airplaneMakeInst airplanemake.AirplaneMake

	airplaneMakeInst.Code = codeConst
	airplaneMakeInst.Name = nameConst
	airplaneMakeInst.Country = countryConst
	airplaneMakeInst.Abbreviation = abbreviationConst

	return airplaneMakeInst
}

func compareFields(t *testing.T, objectInst *airplanemake.AirplaneMake, newObjectInst *airplanemake.AirplaneMake) {

	test_util.CheckField(t, "code", objectInst.Code, newObjectInst.Code)
	test_util.CheckField(t, "callsign", objectInst.Abbreviation, newObjectInst.Abbreviation)
	test_util.CheckField(t, "country", objectInst.Country, newObjectInst.Country)
	test_util.CheckField(t, "name", objectInst.Name, newObjectInst.Name)
}

func chkExists(t *testing.T, router *mux.Router, expectExists bool) {
	test_util.CheckExists(t, router, ApiURL+strings.Replace(ResourceURL, "{"+ObjectID+"}", codeConst, 1), expectExists)
}

func chkPut(t *testing.T, router *mux.Router, objectInst airplanemake.AirplaneMake) {
	jsonString := string(db.ToJson(&objectInst))
	test_util.CheckPut(t, router, jsonString, ApiURL+BaseURL)
}

func chkList(t *testing.T, router *mux.Router) {
	var objectList []*airplanemake.AirplaneMake

	test_util.CheckList(t, router, ApiURL+BaseURL, &objectList)
}

func chkFields(t *testing.T, router *mux.Router, expectedObjectInst airplanemake.AirplaneMake) {
	var newObjectInstMem airplanemake.AirplaneMake

	test_util.CheckFields(t, router, &expectedObjectInst, ApiURL+strings.Replace(ResourceURL, "{"+ObjectID+"}", codeConst, 1), &newObjectInstMem, compareFields)
}

func chkDelete(t *testing.T, router *mux.Router, expectOK bool) {
	test_util.CheckDelete(t, router, ApiURL+strings.Replace(ResourceURL, "{"+ObjectID+"}", codeConst, 1), expectOK)
}

func TestAirplaneMake(t *testing.T) {
	var origObjectInst airplanemake.AirplaneMake = getTestInstance()
	var newObjectInst airplanemake.AirplaneMake = getTestInstance()

	newObjectInst.Name = nameNewConst

	router := mux.NewRouter()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dyno.Conn = dynamodb.New(sess)
	airplanemake.InitConn()

	InitRouter(router)

	t.Log("Ensure the object does not exist to beging with")
	chkExists(t, router, false)

	t.Log("Check put")
	chkPut(t, router, origObjectInst)

	t.Log("Ensure that the object exists after put")
	chkExists(t, router, true)

	t.Log("Ensure the object has the expected field values")
	chkFields(t, router, origObjectInst)

	t.Log("Check put (update)")
	chkPut(t, router, newObjectInst)

	t.Log("Ensure the object has the expected (updated) field values")
	chkFields(t, router, newObjectInst)

	t.Log("Ensure the full list has at least one element")
	chkList(t, router)

	t.Log("Check delete")
	chkDelete(t, router, true)

	t.Log("Ensure that the object does NOT exist after delete")
	chkExists(t, router, false)

	t.Log("Ensure that an attempt to delete non-existent object also succeeds")
	chkDelete(t, router, true)
}

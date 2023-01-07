package airline

import (
	"colmanback/db/dyno"
	"colmanback/test_util"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	codeConst     = "iata:test_code"
	nameConst     = "test_name"
	nameNewConst  = "test_name_2"
	iataConst     = "test_iata"
	icaoConst     = "test_icao"
	callsignConst = "test_callsign"
	countryConst  = "gb"
)

func initDyno() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dyno.Conn = dynamodb.New(sess)

	InitConn()
}

func chkAirline(t *testing.T, airlineInst *Airline) {
	test_util.CheckField(t, "code", codeConst, airlineInst.Code)
	test_util.CheckField(t, "name", nameConst, airlineInst.Name)
	test_util.CheckField(t, "iata", iataConst, airlineInst.Iata)
	test_util.CheckField(t, "icao", icaoConst, airlineInst.Icao)
	test_util.CheckField(t, "callsign", callsignConst, airlineInst.Callsign)
	test_util.CheckField(t, "country", countryConst, airlineInst.Country)
}

func testSetup(t *testing.T) {
	initDyno()

	objectInst, _ := GetAirlineByCode(codeConst)

	if len(objectInst.Code) > 0 {
		objectInst.Delete()
	}

	t.Log("Test setup completed.")
}

func TestAirplane(t *testing.T) {
	if AdapterInst == nil {
		testSetup(t)
	}

	//Create airline from JSON
	airlineInst := AirlineFactory()
	jsonString := fmt.Sprintf("{\"code\":\"%s\", \"name\":\"%s\", \"iata\":\"%s\", \"icao\":\"%s\", \"callsign\":\"%s\", \"country\":\"%s\"}",
		codeConst,
		nameConst,
		iataConst,
		icaoConst,
		callsignConst,
		countryConst)

	airlineInst.FromJson([]byte(jsonString))

	t.Log("Initial check")
	chkAirline(t, airlineInst)

	t.Log("Put check")
	airlineInst.Put()
	airlineRetrInst, getRetrErr := GetAirlineByCode(codeConst)
	if getRetrErr != nil {
		t.Errorf("Error in get after putting airline\n")
	}
	chkAirline(t, airlineRetrInst)

	t.Log("Update, put again and retrieve")
	airlineRetrInst.Name = nameNewConst
	airlineRetrInst.Put()
	airlineUpdtInst, getUpdtErr := GetAirlineByCode(codeConst)
	if getUpdtErr != nil {
		t.Errorf("Error in get after updating airline\n")
	}
	test_util.CheckField(t, "name", nameNewConst, airlineUpdtInst.Name)

	t.Log("Delete and check it's gone!")
	airlineUpdtInst.Delete()
	airlineEmptyInst, getEmptyErr := GetAirlineByCode(codeConst)
	if getEmptyErr == nil {
		t.Errorf("Error for unexistent object not produced when expected. Perhaps the object still exists?\n")
	}
	test_util.CheckField(t, "name", "", airlineEmptyInst.Name)

	t.Log("Test for airline has finished.")
}

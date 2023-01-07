package airplanemake

import (
	"colmanback/db/dyno"
	"colmanback/test_util"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	codeConst         = "dummyx"
	nameConst         = "Namex dummyx"
	nameNewConst      = "New Namex dummyx"
	countryConst      = "ru"
	abbreviationConst = "dmx"
)

func chkObject(t *testing.T, objectInst *AirplaneMake) {
	test_util.CheckField(t, "code", codeConst, objectInst.Code)
	test_util.CheckField(t, "name", nameConst, objectInst.Name)
	test_util.CheckField(t, "country", countryConst, objectInst.Country)
	test_util.CheckField(t, "abbreviation", abbreviationConst, objectInst.Abbreviation)
}

func initDyno() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dyno.Conn = dynamodb.New(sess)

	InitConn()
}

func testSetup(t *testing.T) {
	initDyno()

	objectInst, _ := GetByCode(codeConst)

	if len(objectInst.Code) > 0 {
		objectInst.Delete()
	}

	t.Log("Test setup completed.")
}

func TestAirplaneMake(t *testing.T) {
	testSetup(t)

	//Create object from JSON
	objectInst := ObjectFactory()
	jsonString := fmt.Sprintf("{\"code\":\"%s\", \"name\":\"%s\", \"country\":\"%s\", \"abbreviation\":\"%s\"}",
		codeConst,
		nameConst,
		countryConst,
		abbreviationConst)

	objectInst.FromJson([]byte(jsonString))

	t.Log("Initial check")
	chkObject(t, objectInst)

	t.Log("Put check")
	objectInst.Put()
	objectRetrInst, getRetrErr := GetByCode(codeConst)
	if getRetrErr != nil {
		t.Errorf("Error in get after putting object\n")
	}
	chkObject(t, objectRetrInst)

	t.Log("Update, put again and retrieve")
	objectRetrInst.Name = nameNewConst
	objectRetrInst.Put()
	objectUpdtInst, getUpdtErr := GetByCode(codeConst)
	if getUpdtErr != nil {
		t.Errorf("Error in get after updating object\n")
	}
	test_util.CheckField(t, "name", nameNewConst, objectUpdtInst.Name)

	t.Log("Delete and check it's gone!")
	objectUpdtInst.Delete()
	objectEmptyInst, getEmptyErr := GetByCode(codeConst)
	if getEmptyErr == nil {
		t.Errorf("Error for unexistent object not produced when expected. Perhaps the object still exists?\n")
	}
	test_util.CheckField(t, "name", "", objectEmptyInst.Name)

	t.Log("Test for airplanemake has finished.")
}

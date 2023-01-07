package modelmake

import (
	"colmanback/db/dyno"
	"colmanback/test_util"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	codeConst    = "dummyx"
	nameConst    = "Namex dummyx"
	nameNewConst = "New Namex dummyx"
)

func chkObject(t *testing.T, objectInst *ModelMake) {
	test_util.CheckField(t, "code", codeConst, objectInst.Code)
	test_util.CheckField(t, "name", nameConst, objectInst.Name)
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

func TestModelMake(t *testing.T) {
	testSetup(t)

	//Create object from JSON
	objectInst := ObjectFactory()
	jsonString := fmt.Sprintf("{\"code\":\"%s\", \"name\":\"%s\"}",
		codeConst,
		nameConst)

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

package api_util

import (
	"colmanback/db"
	"colmanback/objects"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	ContentType          = "Content-Type"
	ContentTypeAppJSON   = "application/json"
	ResponseMessageField = "message"
)

type GenAPI[K objects.Object] struct {
	ObjectID string
	ApiURL   string
	BaseURL  string

	Constructor        func() K
	GetObjectByCode    func(objectID string) (K, error)
	GetObjectList      func() ([]K, error)
	DeleteObjectByCode func(objectID string) error
}

//----------------------------------------------------------------------------------------
func (apiInst *GenAPI[K]) Get(writer http.ResponseWriter, request *http.Request) {
	var concreteObjectList []objects.Object

	SetupCORSResponse(&writer)

	pathParams := mux.Vars(request)
	if objectID, ok := pathParams[apiInst.ObjectID]; ok {
		objectInst, getErr := apiInst.GetObjectByCode(objectID)
		if getErr == nil && objectInst.CodeValue() != "" {
			objectInst.WriteObject(writer, request)
		} else {
			WriteMsg(&writer, http.StatusNotFound, fmt.Sprintf("%s with code %s not found", apiInst.ObjectID, objectID))
		}
	} else {
		objectListInst, getListErr := apiInst.GetObjectList()

		if getListErr == nil {
			for _, objectInst := range objectListInst {
				concreteObjectList = append(concreteObjectList, objectInst)
			}
		} else {
			WriteMsg(&writer, http.StatusInternalServerError, fmt.Sprintf("internal server error: %v", getListErr))
		}

		WriteObjectList(concreteObjectList, writer, request)
	}
}

//----------------------------------------------------------------------------------------
func (apiInst *GenAPI[K]) Put(writer http.ResponseWriter, request *http.Request) {
	objectInst := apiInst.Constructor()

	SetupCORSResponse(&writer)
	err := json.NewDecoder(request.Body).Decode(objectInst)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}

	objectInst.Put()
	objectInst.WriteObject(writer, request)
}

//----------------------------------------------------------------------------------------
func (apiInst *GenAPI[K]) Delete(writer http.ResponseWriter, request *http.Request) {
	pathParams := mux.Vars(request)
	if objectID, ok := pathParams[apiInst.ObjectID]; ok {
		apiInst.DeleteObjectByCode(objectID)
		WriteMsg(&writer, http.StatusOK, "Object with code "+objectID+" deleted")
	} else {
		WriteMsg(&writer, http.StatusNotFound, "Object resource not found / an error has been produced.")
	}
}

//----------------------------------------------------------------------------------------
func SetupCORSResponse(writer *http.ResponseWriter) {
	(*writer).Header().Set("Access-Control-Allow-Origin", "*")
	(*writer).Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	(*writer).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
}

//----------------------------------------------------------------------------------------
func WriteMsg(write *http.ResponseWriter, statusCode int, message string) {
	(*write).WriteHeader(statusCode)
	(*write).Header().Set(ContentType, ContentTypeAppJSON)

	responseMap := make(map[string]string)
	responseMap[ResponseMessageField] = message
	jsonResponseMap, _ := json.Marshal(responseMap)
	(*write).Write(jsonResponseMap)
}

//----------------------------------------------------------------------------------------
func WriteObject(objectInst objects.Object, writer http.ResponseWriter, request *http.Request) {
	out, err := json.MarshalIndent(objectInst, db.JSON_PREFIX, db.JSON_INDENT)

	if err != nil {
		log.Panicf("Write Object failed for airline: %s\n", objectInst.ToString())
	}

	SetupCORSResponse(&writer)
	if out != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(out)
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

//----------------------------------------------------------------------------------------
func WriteObjectList(objectInstList []objects.Object, writer http.ResponseWriter, request *http.Request) {
	out, err := json.MarshalIndent(objectInstList, db.JSON_PREFIX, db.JSON_INDENT)

	if err != nil {
		log.Fatalf("Got error when trying to return object list %s", err)
	}

	SetupCORSResponse(&writer)
	if out != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(out)
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

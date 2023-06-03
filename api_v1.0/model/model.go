package model

import (
	"colmanback/api_util"
	"colmanback/objects"
	"colmanback/objects/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	// Standard model adapter constants
	ObjectID    = "modelID"
	ApiURL      = "/api/v1/model"
	BaseURL     = ""
	ResourceURL = "/{" + ObjectID + "}"

	// Model-specific constants
	PictureID     = "pictureID"
	PutPicture    = "/picture/add"
	ListByPicture = "/picture/{" + PictureID + "}/list-models"
	TagPicture    = "/picture/{" + PictureID + "}/tag"
	UntagPicture  = "/picture/{" + PictureID + "}/untag"
	DeletePicture = "/picture/{" + PictureID + "}"
)

var apiInst api_util.GenAPI[*model.Model]
var apiInstListByPicture api_util.GenAPI[*model.Model]

//----------------------------------------------------------------------------------------
func handlePictureCommand(writer http.ResponseWriter, request *http.Request, handler func(string, string) (*model.Model, error)) {
	var objectInst *model.Model = &model.Model{}

	api_util.SetupCORSResponse(&writer)
	err := json.NewDecoder(request.Body).Decode(objectInst)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}

	requestVars := mux.Vars(request)

	if picture, ok := requestVars[PictureID]; ok {
		modelInst, tagErr := handler(picture, objectInst.Code)
		if tagErr != nil {
			api_util.WriteMsg(&writer, http.StatusInternalServerError, fmt.Sprintf("An internal error has occurred. Error: %v", tagErr))
		} else {
			modelInst.WriteObject(writer, request)
		}
	} else {
		api_util.WriteMsg(&writer, http.StatusBadRequest, "The picture ID to be processed could not be found in the request.")
	}
}

//----------------------------------------------------------------------------------------
func modelListToObjectList(modelList []*model.Model) []objects.Object {
	var objectInstList []objects.Object

	for _, modelInst := range modelList {
		objectInstList = append(objectInstList, modelInst)
	}

	return objectInstList
}

//----------------------------------------------------------------------------------------
func handleAddModelPicture(writer http.ResponseWriter, request *http.Request) {
	modelCodeListStr := request.PostFormValue("modelList")
	modelPicture, _, formErr := request.FormFile("picture")

	api_util.SetupCORSResponse(&writer)
	if formErr == nil && len(modelCodeListStr) > 0 {
		modelCodeList := strings.Split(modelCodeListStr, ",")
		modelList, addErr := model.AddModelPicture(modelPicture, modelCodeList)

		if addErr == nil {
			objectInstList := modelListToObjectList(modelList)
			defer modelPicture.Close()
			api_util.WriteObjectList(objectInstList, writer, request)
		} else {
			api_util.WriteMsg(&writer, http.StatusInternalServerError, fmt.Sprintf("An internal error has occurred. Error: %v", addErr))
		}
	} else {
		api_util.WriteMsg(&writer, http.StatusBadRequest, fmt.Sprintf("The request was malformed. Error: %v", formErr))
	}
}

//----------------------------------------------------------------------------------------
func handleTagModelPicture(writer http.ResponseWriter, request *http.Request) {
	handlePictureCommand(writer, request, model.TagModelPicture)
}

//----------------------------------------------------------------------------------------
func handleUntagModelPicture(writer http.ResponseWriter, request *http.Request) {
	handlePictureCommand(writer, request, model.RemoveModelPicture)
}

//----------------------------------------------------------------------------------------
func handleDeleteModelPicture(writer http.ResponseWriter, request *http.Request) {
	api_util.SetupCORSResponse(&writer)

	requestVars := mux.Vars(request)

	if picture, ok := requestVars[PictureID]; ok {
		modelList, tagErr := model.DeleteModelPicture(picture)
		if tagErr != nil {
			api_util.WriteMsg(&writer, http.StatusInternalServerError, fmt.Sprintf("An internal error has occurred. Error: %v", tagErr))
		} else {
			objectInstList := modelListToObjectList(modelList)
			api_util.WriteObjectList(objectInstList, writer, request)
		}
	} else {
		api_util.WriteMsg(&writer, http.StatusBadRequest, "The picture ID to be processed could not be found in the request.")
	}
}

//----------------------------------------------------------------------------------------
func initBase(subRouter *mux.Router) {
	subRouter.HandleFunc(BaseURL, apiInst.GetList).Methods(http.MethodGet)
	subRouter.HandleFunc(BaseURL, apiInst.Put).Methods(http.MethodPut)
	subRouter.HandleFunc(ResourceURL, apiInst.Get).Methods(http.MethodGet)
	subRouter.HandleFunc(ResourceURL, apiInst.Delete).Methods(http.MethodDelete)

	apiInst.ApiURL = ApiURL
	apiInst.BaseURL = BaseURL
	apiInst.ObjectID = ObjectID

	apiInst.Constructor = model.ObjectFactory
	apiInst.GetObjectByCode = model.GetByCode
	apiInst.GetObjectList = model.GetList
	apiInst.DeleteObjectByCode = model.AdapterInst.DeleteObjectByCode
}

//----------------------------------------------------------------------------------------
func initListByPicture(subRouter *mux.Router) {
	subRouter.HandleFunc(BaseURL+ListByPicture, apiInstListByPicture.GetList).Methods(http.MethodGet)

	apiInstListByPicture.ApiURL = ApiURL
	apiInstListByPicture.BaseURL = BaseURL
	apiInstListByPicture.ObjectID = PictureID

	apiInstListByPicture.Constructor = model.ObjectFactory
	apiInstListByPicture.GetObjectListByCode = model.GetModelByPicture
}

//----------------------------------------------------------------------------------------
func initAddModelPicture(subRouter *mux.Router) {
	subRouter.HandleFunc(PutPicture, handleAddModelPicture).Methods(http.MethodPost)
}

//----------------------------------------------------------------------------------------
func initTagModelPicture(subRouter *mux.Router) {
	subRouter.HandleFunc(TagPicture, handleTagModelPicture).Methods(http.MethodPut)
}

//----------------------------------------------------------------------------------------
func initUntagModelPicture(subRouter *mux.Router) {
	subRouter.HandleFunc(TagPicture, handleUntagModelPicture).Methods(http.MethodPut)
}

//----------------------------------------------------------------------------------------
func initDeleteModelPicture(subRouter *mux.Router) {
	subRouter.HandleFunc(DeletePicture, handleDeleteModelPicture).Methods(http.MethodDelete)
}

//----------------------------------------------------------------------------------------
func InitRouter(router *mux.Router) {
	subRouter := router.PathPrefix(ApiURL).Subrouter()

	//Generic model API.
	initBase(subRouter)

	//Model-specific APIs.
	initListByPicture(subRouter)
	initAddModelPicture(subRouter)
	initTagModelPicture(subRouter)
	initUntagModelPicture(subRouter)
	initDeleteModelPicture(subRouter)

	//TODO: add routes for the remaining methods.
	//TODO: implement the damn test for MODEL!
}

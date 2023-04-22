package model

import (
	"colmanback/api_util"
	"colmanback/objects"
	"colmanback/objects/model"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	ObjectID    = "modelID"
	ApiURL      = "/api/v1/model"
	BaseURL     = ""
	ResourceURL = "/{" + ObjectID + "}"

	// Model-specific constants
	PictureID     = "pictureID"
	ListByPicture = "/list-by-picture/{" + PictureID + "}"
	PutPicture    = "/add-picture"
)

var apiInst api_util.GenAPI[*model.Model]
var apiInstListByPicture api_util.GenAPI[*model.Model]

//----------------------------------------------------------------------------------------
func handleAddModelPicture(writer http.ResponseWriter, request *http.Request) {
	var objectInstList []objects.Object

	modelCodeListStr := request.PostFormValue("modelList")
	modelPicture, _, formErr := request.FormFile("picture")

	api_util.SetupCORSResponse(&writer)
	if formErr == nil && len(modelCodeListStr) > 0 {
		modelCodeList := strings.Split(modelCodeListStr, ",")
		modelList, addErr := model.AddModelPicture(modelPicture, modelCodeList)
		if addErr == nil {
			for _, modelInst := range modelList {
				objectInstList = append(objectInstList, modelInst)
			}

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
func InitRouter(router *mux.Router) {
	subRouter := router.PathPrefix(ApiURL).Subrouter()

	initBase(subRouter)

	//Model-specific APIs
	initListByPicture(subRouter)
	initAddModelPicture(subRouter)

	//TODO: add routes for the remaining methods.
	//TODO: implement the damn test for MODEL!
}

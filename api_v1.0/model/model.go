package model

import (
	"colmanback/api_util"
	"colmanback/objects/model"
	"net/http"

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
)

var apiInst api_util.GenAPI[*model.Model]
var apiInstListByPicture api_util.GenAPI[*model.Model]

//----------------------------------------------------------------------------------------
func initListByPicture(router *mux.Router, subRouter *mux.Router) {
	subRouter.HandleFunc(ListByPicture, apiInstListByPicture.GetList).Methods(http.MethodGet)

	apiInstListByPicture.ApiURL = ApiURL
	apiInstListByPicture.BaseURL = BaseURL
	apiInstListByPicture.ObjectID = PictureID

	apiInstListByPicture.Constructor = model.ObjectFactory
	apiInstListByPicture.GetObjectListByCode = model.GetModelByPicture
}

//----------------------------------------------------------------------------------------
func InitRouter(router *mux.Router) {
	subRouter := router.PathPrefix(ApiURL).Subrouter()

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

	//Model-specific APIs
	initListByPicture(router, subRouter)

	//TODO: add routes for the remaining methods.
}

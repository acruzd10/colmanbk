package modelmake

import (
	"colmanback/api_util"
	"colmanback/objects/modelmake"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	ObjectID    = "modelMakeID"
	ApiURL      = "/api/v1/modelmake"
	BaseURL     = ""
	ResourceURL = "/{" + ObjectID + "}"
)

var apiInst api_util.GenAPI[*modelmake.ModelMake]

//----------------------------------------------------------------------------------------
func InitRouter(router *mux.Router) {
	subRouter := router.PathPrefix(ApiURL).Subrouter()

	subRouter.HandleFunc(BaseURL, apiInst.Get).Methods(http.MethodGet)
	subRouter.HandleFunc(BaseURL, apiInst.Put).Methods(http.MethodPut)
	subRouter.HandleFunc(ResourceURL, apiInst.Get).Methods(http.MethodGet)
	subRouter.HandleFunc(ResourceURL, apiInst.Delete).Methods(http.MethodDelete)

	apiInst.ApiURL = ApiURL
	apiInst.BaseURL = BaseURL
	apiInst.ObjectID = ObjectID

	apiInst.Constructor = modelmake.ObjectFactory
	apiInst.GetObjectByCode = modelmake.GetByCode
	apiInst.GetObjectList = modelmake.GetList
	apiInst.DeleteObjectByCode = modelmake.AdapterInst.DeleteObjectByCode
}

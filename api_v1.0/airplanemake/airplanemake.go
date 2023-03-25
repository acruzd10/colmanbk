package airplanemake

import (
	"colmanback/api_util"
	"colmanback/objects/airplanemake"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	ObjectID    = "objectID"
	ApiURL      = "/api/v1/airplanemake"
	BaseURL     = ""
	ResourceURL = "/{" + ObjectID + "}"
)

var apiInst api_util.GenAPI[*airplanemake.AirplaneMake]

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

	apiInst.Constructor = airplanemake.ObjectFactory
	apiInst.GetObjectByCode = airplanemake.GetByCode
	apiInst.GetObjectList = airplanemake.GetList
	apiInst.DeleteObjectByCode = airplanemake.AdapterInst.DeleteObjectByCode
}

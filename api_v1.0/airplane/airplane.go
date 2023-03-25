package airplane

import (
	"colmanback/api_util"
	"colmanback/objects/airplane"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	ObjectID    = "airplaneID"
	ApiURL      = "/api/v1/airplane"
	BaseURL     = ""
	ResourceURL = "/{" + ObjectID + "}"
)

var apiInst api_util.GenAPI[*airplane.Airplane]

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

	apiInst.Constructor = airplane.ObjectFactory
	apiInst.GetObjectByCode = airplane.GetByCode
	apiInst.GetObjectList = airplane.GetList
	apiInst.DeleteObjectByCode = airplane.AdapterInst.DeleteObjectByCode
}

package airline

import (
	"colmanback/api_util"
	"colmanback/objects/airline"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	ObjectID    = "airlineID"
	ApiURL      = "/api/v1/airline"
	BaseURL     = ""
	ResourceURL = "/{" + ObjectID + "}"
)

var apiInst api_util.GenAPI[*airline.Airline]

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

	apiInst.Constructor = airline.AirlineFactory
	apiInst.GetObjectByCode = airline.GetAirlineByCode
	apiInst.GetObjectList = airline.GetAirlineList
	apiInst.DeleteObjectByCode = airline.AdapterInst.DeleteObjectByCode
}

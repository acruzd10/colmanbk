package country

import (
	"colmanback/api_util"
	"colmanback/objects/country"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	ObjectID    = "countryID"
	ApiURL      = "/api/v1/country"
	BaseURL     = ""
	ResourceURL = "/{" + ObjectID + "}"
)

var apiInst api_util.GenAPI[*country.Country]

//----------------------------------------------------------------------------------------
func InitRouter(router *mux.Router) {
	subRouter := router.PathPrefix(ApiURL).Subrouter()

	subRouter.HandleFunc(BaseURL, apiInst.GetList).Methods(http.MethodGet)
	subRouter.HandleFunc(ResourceURL, apiInst.Get).Methods(http.MethodGet)

	apiInst.ApiURL = ApiURL
	apiInst.BaseURL = BaseURL
	apiInst.ObjectID = ObjectID

	apiInst.GetObjectByCode = country.GetCountryByISO
	apiInst.GetObjectList = country.GetCountryList
}

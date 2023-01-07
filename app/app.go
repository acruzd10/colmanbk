package app

import (
	airlineapi "colmanback/api_v1.0/airline"
	airplaneapi "colmanback/api_v1.0/airplane"
	airplanemakeapi "colmanback/api_v1.0/airplanemake"
	countryapi "colmanback/api_v1.0/country"
	"colmanback/db/dyno"
	airlineobject "colmanback/objects/airline"
	airplaneobject "colmanback/objects/airplane"
	airplanemakeobject "colmanback/objects/airplanemake"
	countryobject "colmanback/objects/country"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
)

type App struct {
	Sess *session.Session
	Port string
}

//----------------------------------------------------------------------------------------
func (appInst *App) initConn() {
	appInst.Sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dyno.Conn = dynamodb.New(appInst.Sess)

	airlineobject.InitConn()
	airplanemakeobject.InitConn()
	airplaneobject.InitConn()
	countryobject.InitConn()
}

//----------------------------------------------------------------------------------------
func (appInst *App) initRoutes() *mux.Router {
	router := mux.NewRouter()

	airlineapi.InitRouter(router)
	airplanemakeapi.InitRouter(router)
	airplaneapi.InitRouter(router)
	countryapi.InitRouter(router)

	return router
}

//----------------------------------------------------------------------------------------
func (appInst *App) Serve() {
	appInst.initConn()
	router := appInst.initRoutes()

	log.Printf("Staring web server on port %s\n", appInst.Port)
	http.ListenAndServe(appInst.Port, router)
}

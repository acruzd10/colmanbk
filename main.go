package main

import "colmanback/app"

//----------------------------------------------------------------------------------------
func main() {
	appInst := app.App{}

	appInst.Port = ":8081"
	appInst.Serve()
}

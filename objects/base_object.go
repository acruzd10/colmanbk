package objects

import "net/http"

type Object interface {
	CodeValue() string
	SortValue() string
	ToString() string
	FromJson(jsonInst []byte)
	Print()
	Put()
	Delete()
	WriteObject(writer http.ResponseWriter, request *http.Request)
}

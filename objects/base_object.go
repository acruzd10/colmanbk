package objects

import "net/http"

type ModelScale string

const (
	Scale1200 ModelScale = "1/200"
	Scale1400 ModelScale = "1/400"
)

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

package modelmake

import (
	"colmanback/api_util"
	"colmanback/db"
	"colmanback/db/dyno"
	"fmt"
	"net/http"
)

type ModelMake struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

var AdapterInst db.Adapter[*ModelMake]

//----------------------------------------------------------------------------------------
func (modelMakeInst *ModelMake) CodeValue() string {
	return modelMakeInst.Code
}

//----------------------------------------------------------------------------------------
func (modelMakeInst *ModelMake) SortValue() string {
	return ""
}

//----------------------------------------------------------------------------------------
func (modelMakeInst *ModelMake) FromJson(jsonInst []byte) {
	db.FromJson(modelMakeInst, jsonInst)
}

//----------------------------------------------------------------------------------------
func (modelMakeInst *ModelMake) ToString() string {
	str := fmt.Sprintf(`
	----------------------
	  Code ......: %s
	  Name ......: %s`,
		modelMakeInst.Code,
		modelMakeInst.Name)

	return str
}

//----------------------------------------------------------------------------------------
func (modelMakeInst *ModelMake) Print() {
	fmt.Println(modelMakeInst.ToString())
}

//----------------------------------------------------------------------------------------
func (modelMakeInst *ModelMake) WriteObject(writer http.ResponseWriter, request *http.Request) {
	api_util.WriteObject(modelMakeInst, writer, request)
}

//----------------------------------------------------------------------------------------
func ObjectFactory() *ModelMake {
	var modelMakeInst ModelMake = ModelMake{}

	return &modelMakeInst
}

//----------------------------------------------------------------------------------------
func GetCacheMap(modelMakeList []*ModelMake) []db.CacheMapElement {
	var cacheMap []db.CacheMapElement

	for _, modelMake := range modelMakeList {
		cacheMap = db.AddToCacheMap(cacheMap, modelMake.Code, modelMake.Code, modelMake.Name)
		cacheMap = db.AddToCacheMap(cacheMap, modelMake.Name, modelMake.Code, modelMake.Name)
	}

	return cacheMap
}

//----------------------------------------------------------------------------------------
func (modelMakeInst *ModelMake) Delete() {
	AdapterInst.DeleteObject(modelMakeInst)
}

//----------------------------------------------------------------------------------------
func (modelMakeInst *ModelMake) Put() {
	AdapterInst.PutObject(modelMakeInst)
}

//----------------------------------------------------------------------------------------
func GetList() ([]*ModelMake, error) {
	return AdapterInst.GetObjectList()
}

//----------------------------------------------------------------------------------------
func GetByCode(code string) (*ModelMake, error) {
	return AdapterInst.GetObjectByCode(code)
}

//----------------------------------------------------------------------------------------
func LoadModelMakeList(ModelMakeList []ModelMake) {
	var modelMakePointerList []*ModelMake = []*ModelMake{}

	for _, modelMakeInst := range ModelMakeList {
		var newModelMakeInst ModelMake = modelMakeInst
		modelMakePointerList = append(modelMakePointerList, &newModelMakeInst)
	}

	AdapterInst.PutObjectList(modelMakePointerList)
}

//----------------------------------------------------------------------------------------
func InitConn() {
	dynoInstAirline := &dyno.Dyno[*ModelMake]{}
	dynoInstAirline.Config("modelmake", "code", true, ObjectFactory, GetCacheMap)
	AdapterInst = dynoInstAirline
}

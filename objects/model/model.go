package model

import (
	"colmanback/api_util"
	"colmanback/db"
	"colmanback/db/dyno"
	"colmanback/db/s3"
	"colmanback/objects"
	"colmanback/objects/airline"
	"colmanback/objects/airplane"
	"colmanback/objects/modelmake"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	SEP = "#"
)

type Model struct {
	Code    string `json:"code"`
	Picture string `json:"picture",omitempty`

	//Foreign Keys
	ModelMake string `json:"modelMake"`
	Airline   string `json:"airline"`
	Airplane  string `json:"airplane"`

	//Properties
	Scale           objects.ModelScale `json:"scale"`
	Reg             string             `json:"reg"`
	Notes           string             `json:"notes"`
	IsCargo         bool               `json:"isCargo"`
	IsOldLivery     bool               `json:"isOldLivery"`
	IsSpecialLivery bool               `json:"isSpecialLivery"`
	PictureList     []string           `json:"pictureList,omitempty"`

	//Reference Instances
	ModelMakeInst *modelmake.ModelMake `json:"modelMakeDetails,omitempty"`
	AirlineInst   *airline.Airline     `json:"airlineDetails,omitempty"`
	AirplaneInst  *airplane.Airplane   `json:"airplaneDetails,omitempty"`
}

var AdapterInst db.Adapter[*Model]
var FileInst *s3.S3Adapter

//----------------------------------------------------------------------------------------
func (modelInst *Model) makeCode() {
	if len(modelInst.Code) > 0 {
		return
	}

	code := modelInst.ModelMake + SEP + string(modelInst.Scale) + SEP + modelInst.Reg
	modelInst.Code = strings.ToLower(code)
}

//----------------------------------------------------------------------------------------
func (modelInst *Model) CodeValue() string {
	return modelInst.Code
}

//----------------------------------------------------------------------------------------
func (modelInst *Model) SortValue() string {
	if len(modelInst.Picture) > 0 {
		return modelInst.Picture
	} else {
		return modelInst.Code
	}
}

//----------------------------------------------------------------------------------------
func (modelInst *Model) FromJson(jsonInst []byte) {
	db.FromJson(modelInst, jsonInst)

	if len(modelInst.Code) == 0 {
		modelInst.makeCode()
	}

	modelInst.Reg = strings.ToUpper(modelInst.Reg)
}

//----------------------------------------------------------------------------------------
func (modelInst *Model) ToString() string {
	str := fmt.Sprintf(`
	----------------------
	  Code ......: %s
	  ModelMake .: %s
	  Airline ...: %s
	  Airplane ..: %s
	  Scale .....: %s
	  Reg. ......: %s
	  Notes .....: %s
	  Is Cargo ..: %t
	  Is Old Liv.: %t
	  Is Spc Liv.: %t`,
		modelInst.Code,
		modelInst.ModelMake,
		modelInst.Airline,
		modelInst.Airplane,
		modelInst.Scale,
		modelInst.Reg,
		modelInst.Notes,
		modelInst.IsCargo,
		modelInst.IsOldLivery,
		modelInst.IsSpecialLivery)

	return str
}

func (modelInst *Model) Print() {
	fmt.Println(modelInst.ToString())
}

//----------------------------------------------------------------------------------------
func ObjectFactory() *Model {
	var modelInst Model = Model{}

	return &modelInst
}

//----------------------------------------------------------------------------------------
func (modelInst *Model) InitRefObjs() {
	var err error

	if len(modelInst.Airline) != 0 {
		modelInst.AirlineInst, err = airline.GetByCode(modelInst.Airline)
		if err != nil {
			panic(fmt.Sprintf("Error initialising airline with code %s for model %s", modelInst.Airline, modelInst.Code))
		}
	}

	if len(modelInst.Airplane) != 0 {
		modelInst.AirplaneInst, err = airplane.GetByCode(modelInst.Airplane)
		if err != nil {
			panic(fmt.Sprintf("Error initialising airplane with code %s for model %s", modelInst.Airline, modelInst.Code))
		}
	}

	if len(modelInst.ModelMake) != 0 {
		modelInst.ModelMakeInst, err = modelmake.GetByCode(modelInst.ModelMake)
		if err != nil {
			panic(fmt.Sprintf("Error initialising modelmake with code %s for model %s", modelInst.ModelMake, modelInst.Code))
		}
	}
}

//----------------------------------------------------------------------------------------
func (modelInst *Model) LoadPictures() {
	pictureList, err := AdapterInst.GetSortKeyList(modelInst.Code)

	modelInst.PictureList = pictureList
	if err != nil {
		log.Fatalf("An error has occurred while retrieving the pictures for a model with code %s. Error: %v", modelInst.Code, err)
	}
}

//----------------------------------------------------------------------------------------
func (modelInst *Model) Delete() {
	var err error

	if len(modelInst.Picture) == 0 {
		err = AdapterInst.DeleteObject(modelInst)
	} else {
		err = AdapterInst.DeleteObjectByCodeAndSort(modelInst.Code, modelInst.Picture)
	}

	if err != nil {
		log.Fatalf("An error has occurred while deleting model with code %s. Error: %v\n", modelInst.Code, err)
	}
}

//----------------------------------------------------------------------------------------
func (modelInst *Model) Put() {
	var modelInstPut Model = *modelInst

	//Ensure that redundant ref info is not persisted.
	modelInstPut.AirlineInst = nil
	modelInstPut.AirplaneInst = nil
	modelInstPut.ModelMakeInst = nil

	if len(modelInstPut.Code) == 0 {
		modelInstPut.makeCode()
		modelInst.Code = modelInstPut.Code
	}

	err := AdapterInst.PutObject(&modelInstPut)
	if err != nil {
		log.Fatalf("An error has occurred while putting model with code %s. Error: %v\n", modelInst.Code, err)
	}
}

//----------------------------------------------------------------------------------------
func (airlineInst *Model) WriteObject(writer http.ResponseWriter, request *http.Request) {
	api_util.WriteObject(airlineInst, writer, request)
}

//----------------------------------------------------------------------------------------
func GetList() ([]*Model, error) {
	objectList, err := AdapterInst.GetObjectList()

	if err == nil {
		for index, objectInst := range objectList {
			objectInst.InitRefObjs()
			objectList[index] = objectInst
		}
	}

	return objectList, err
}

//----------------------------------------------------------------------------------------
func GetByCode(code string) (*Model, error) {
	objectInst, err := AdapterInst.GetObjectByCode(code)
	if err == nil {
		objectInst.InitRefObjs()
	}

	return objectInst, err
}

//----------------------------------------------------------------------------------------
func InitConn() {
	dynoInstModel := &dyno.Dyno[*Model]{}
	dynoInstModel.SetSortName("picture")
	dynoInstModel.Config("model", "code", true, ObjectFactory, nil)
	AdapterInst = dynoInstModel

	fileInstModel := &s3.S3Adapter{}
	fileInstModel.Config("colman-pics", 1000)
	FileInst = fileInstModel
}

//----------------------------------------------------------------------------------------
func AddModelPicture(file multipart.File, modelCodeList []string) ([]*Model, error) {
	nowTime := time.Now().Format(time.RFC3339)
	uuidName := uuid.New().String()
	fileName := nowTime + "#" + uuidName

	var modelInst *Model
	var intlErr error
	modelList := []*Model{}

	response, err := FileInst.AddFile(fileName, file)
	if err == nil {
		if len(response.FileLocation) != 0 {
			for _, code := range modelCodeList {
				modelInst = &Model{}
				modelInst.Code = code
				modelInst.Picture = response.FileLocation
				modelInst.Put()

				modelInst, intlErr = GetByCode(code)
				if intlErr == nil {
					modelInst.PictureList = append(modelInst.PictureList, response.FileLocation)
					modelList = append(modelList, modelInst)
				}
			}
		}
	} else {
		log.Printf("Error whilst saving image for models %v. Error: %v", modelCodeList, err)
	}

	return modelList, err
}

/*
//----------------------------------------------------------------------------------------
func DeleteModelPicture(fileName string) ([]*Model, error) {

}
*/

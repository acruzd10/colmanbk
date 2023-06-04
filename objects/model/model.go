package model

import (
	"colmanback/api_util"
	"colmanback/db"
	"colmanback/db/s3"
	"colmanback/objects"
	"colmanback/objects/airline"
	"colmanback/objects/airplane"
	"colmanback/objects/modelmake"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	SEP = "#"
)

type Model struct {
	Code     string   `json:"code"`
	Picture  string   `json:"picture,omitempty"` //Used by the picture index stub only
	CodeList []string `json:"codeList,omitempty"`

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
	PictureList     []string           `json:"pictureList,omitempty"` //Used by the actual model instances

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
	  Is Spc Liv.: %t
	  Picture ...: %s
	  PictureList: %v`,
		modelInst.Code,
		modelInst.ModelMake,
		modelInst.Airline,
		modelInst.Airplane,
		modelInst.Scale,
		modelInst.Reg,
		modelInst.Notes,
		modelInst.IsCargo,
		modelInst.IsOldLivery,
		modelInst.IsSpecialLivery,
		modelInst.Picture,
		modelInst.PictureList,
	)

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
		if modelInst.PictureList != nil && len(modelInst.PictureList) > 0 {
			for _, picture := range modelInst.PictureList {
				AdapterInst.DeleteObjectByCodeAndSort(modelInst.CodeValue(), picture)
			}
		}
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
	var getErr error

	airlineInst := modelInst.AirlineInst
	airplaneInst := modelInst.AirplaneInst
	modelMakeInst := modelInst.ModelMakeInst

	//Ensure that redundant ref info is not persisted.
	modelInst.AirlineInst = nil
	modelInst.AirplaneInst = nil
	modelInst.ModelMakeInst = nil

	if len(modelInst.Code) == 0 {
		modelInst.makeCode()
	}

	err := AdapterInst.PutObject(modelInst)
	if err != nil {
		log.Fatalf("An error has occurred while putting model with code %s. Error: %v\n", modelInst.Code, err)
	}

	if airlineInst == nil && len(modelInst.Airline) > 0 {
		airlineInst, getErr = airline.GetByCode(modelInst.Airline)
		if getErr != nil {
			log.Printf("Model has airline with code %s, but an attempt to retrieve an airline with that code produced error %v\n", modelInst.Airline, getErr)
		}
	}

	if airplaneInst == nil && len(modelInst.Airplane) > 0 {
		airplaneInst, getErr = airplane.GetByCode(modelInst.Airplane)
		if getErr != nil {
			log.Printf("Model has airplane with code %s, but an attempt to retrieve an airplane with that code produced error %v\n", modelInst.Airplane, getErr)
		}
	}

	if modelMakeInst == nil && len(modelInst.ModelMake) > 0 {
		modelMakeInst, getErr = modelmake.GetByCode(modelInst.ModelMake)
		if getErr != nil {
			log.Printf("Model has model make with code %s, but an attempt to retrieve an model make with that code produced error %v\n", modelInst.ModelMake, getErr)
		}
	}

	modelInst.AirlineInst = airlineInst
	modelInst.AirplaneInst = airplaneInst
	modelInst.ModelMakeInst = modelMakeInst
}

//----------------------------------------------------------------------------------------
func (modelInst *Model) WriteObject(writer http.ResponseWriter, request *http.Request) {
	api_util.WriteObject(modelInst, writer, request)
}

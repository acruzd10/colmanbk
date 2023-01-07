package model

import (
	"colmanback/api_util"
	"colmanback/db"
	"colmanback/objects"
	"colmanback/objects/airline"
	"colmanback/objects/airplane"
	"colmanback/objects/country"
	"colmanback/objects/modelmake"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	SEP = "#"
)

type Model struct {
	Code string `json:"code"`

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
	PictureList     []string           `json:"pictureList"`

	//Reference Instances
	ModelMakeInst *modelmake.ModelMake
	AirlineInst   *airline.Airline
	AirplaneInst  *airplane.Airplane
}

type ModelIntl struct {
	Code               string   `json:"code"`
	ModelMake          string   `json:"modelMake"`
	ModelMakeName      string   `json:"modelMakeName"`
	Airline            string   `json:"airline"`
	AirlineName        string   `json:"airlineName"`
	AirlineCountry     string   `json:"country"`
	AirlineCountryName string   `json:"countryName"`
	Airplane           string   `json:"airplane"`
	AirplaneName       string   `json:"airplaneName"`
	Notes              string   `json:"notes"`
	Scale              string   `json:"scale"`
	Reg                string   `json:"reg"`
	IsCargo            bool     `json:"isCargo"`
	IsOldLivery        bool     `json:"isOldLivery"`
	IsSpecialLivery    bool     `json:"isSpecialLivery"`
	PictureList        []string `json:"pictureList"`
}

var AdapterInst db.Adapter[*Model]

//----------------------------------------------------------------------------------------
func (modelInst *Model) getModelIntl() *ModelIntl {
	modelIntlInst := ModelIntl{}

	modelInst.InitRefObjs()
	modelIntlInst.Code = modelInst.Code
	modelIntlInst.ModelMake = modelInst.ModelMake
	modelIntlInst.Airline = modelInst.Airline
	modelIntlInst.Airplane = modelInst.Airplane
	modelIntlInst.Scale = string(modelInst.Scale)
	modelIntlInst.Reg = modelInst.Reg
	modelIntlInst.Notes = modelInst.Notes
	modelIntlInst.IsCargo = modelInst.IsCargo
	modelIntlInst.IsOldLivery = modelInst.IsOldLivery
	modelIntlInst.IsSpecialLivery = modelInst.IsSpecialLivery
	modelIntlInst.PictureList = modelInst.PictureList

	if len(modelIntlInst.ModelMake) > 0 {
		modelIntlInst.ModelMakeName = modelInst.ModelMakeInst.Name
	}

	if len(modelIntlInst.Airline) > 0 {
		modelIntlInst.AirlineName = modelInst.AirlineInst.Name
		modelIntlInst.AirlineCountry = modelInst.AirlineInst.Country
		if len(modelIntlInst.AirlineCountry) > 0 {
			countryInst, err := country.GetCountryByISO(modelIntlInst.AirlineCountry)
			if err == nil {
				modelIntlInst.AirlineCountryName = countryInst.Name
			}
		}
	}

	if len(modelIntlInst.Airplane) > 0 {
		modelIntlInst.AirplaneName = modelInst.AirplaneInst.Name
	}

	return &modelIntlInst
}

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
	return ""
}

//----------------------------------------------------------------------------------------
func (modelInst *Model) FromJson(jsonInst []byte) {
	db.FromJson(modelInst, jsonInst)

	if len(modelInst.Code) == 0 {
		modelInst.makeCode()
	}

	modelInst.Reg = strings.ToUpper(modelInst.Reg)
}

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
func (modelInst *Model) Delete() {
	AdapterInst.DeleteObject(modelInst)
}

//----------------------------------------------------------------------------------------
func (modelInst *Model) Put() {
	AdapterInst.PutObject(modelInst)
}

//----------------------------------------------------------------------------------------
func (airlineInst *Model) WriteObject(writer http.ResponseWriter, request *http.Request) {
	api_util.WriteObject(airlineInst, writer, request)
}

//----------------------------------------------------------------------------------------
func (modelInst *Model) ToJSON() []byte {
	out, err := json.MarshalIndent(modelInst, db.JSON_PREFIX, db.JSON_INDENT)

	if err != nil {
		log.Fatalf("Error marshalling model with code %s. Error:%s", modelInst.Code, err)
		return nil
	}

	return out
}

//----------------------------------------------------------------------------------------
func (modelInst *Model) PutModel() []byte {
	modelInst.makeCode()
	AdapterInst.PutObject(modelInst)

	modelInstIntl := modelInst.getModelIntl()
	out, err := json.MarshalIndent(modelInstIntl, db.JSON_PREFIX, db.JSON_INDENT)

	if err != nil {
		log.Fatalf("Error marshalling model with code %s. Error:%s", modelInst.Code, err)
		return nil
	}

	return out
}

//----------------------------------------------------------------------------------------
func GetList() ([]*Model, error) {
	return AdapterInst.GetObjectList()
}

//----------------------------------------------------------------------------------------
func GetModelByCode(code string) (*Model, error) {
	return AdapterInst.GetObjectByCode(code)
}

//----------------------------------------------------------------------------------------
func GetModelListExtendedAPI(writer http.ResponseWriter) {
	var modelIntlList []*ModelIntl = []*ModelIntl{}
	modelList, _ := GetList()

	for _, model := range modelList {
		modelIntlList = append(modelIntlList, model.getModelIntl())
	}

	out, err := json.MarshalIndent(modelIntlList, db.JSON_PREFIX, db.JSON_INDENT)
	if err == nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(out)
	}
}

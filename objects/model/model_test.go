package model

import (
	"colmanback/db/dyno"
	"colmanback/objects"
	"colmanback/objects/airline"
	"colmanback/objects/airplane"
	"colmanback/objects/airplanemake"
	"colmanback/objects/modelmake"
	"colmanback/test_util"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	//Foreign Keys
	modelMakeConst          = "model_model_make_test"
	modelMakeNameConst      = modelMakeConst + "_name"
	airlineConst            = "iata:model_airline_test"
	airlineNameConst        = airlineConst + "_name"
	airlineIataConst        = airlineConst + "_iata"
	airlineCallsignConst    = airlineConst + "_callsign"
	airplaneConst           = "model_airplane_test"
	airplaneIataConst       = airplaneConst + "_iata"
	airplaneNameConst       = airplaneConst + "_name"
	airplaneMakeConst       = "model_airplane_make_test"
	airplaneMakeNameConst   = airplaneMakeConst + "_name"
	airplaneMakeAbbrevConst = airplaneMakeConst + "_abbrev"
	countryConst            = "gb"

	//Properties
	scaleConst           = objects.Scale1400
	regConst             = "modeltest001"
	regNewConst          = regConst + "N"
	notesConst           = "model_test_notes!"
	isCargoConst         = true
	isOldLiveryConst     = false
	isSpecialLiveryConst = true
)

var airlineInst airline.Airline
var airplaneInst airplane.Airplane
var airplaneMakeInst airplanemake.AirplaneMake
var modelMakeInst modelmake.ModelMake

func initDyno() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dyno.Conn = dynamodb.New(sess)

	InitConn()
}

func boolToString(boolValue bool) string {
	if boolValue {
		return "true"
	}

	return "false"
}

func chkModel(t *testing.T, objectInst *Model, emptyCode bool) {
	test_util.CheckField(t, "modelMake", modelMakeConst, objectInst.ModelMake)
	test_util.CheckField(t, "airline", airlineConst, objectInst.Airline)
	test_util.CheckField(t, "airplane", airplaneConst, objectInst.Airplane)
	test_util.CheckField(t, "scale", string(scaleConst), string(objectInst.Scale))
	test_util.CheckField(t, "reg", strings.ToUpper(regConst), strings.ToUpper(objectInst.Reg))
	test_util.CheckField(t, "notes", notesConst, objectInst.Notes)
	test_util.CheckField(t, "isCargo", boolToString(isCargoConst), boolToString(objectInst.IsCargo))
	test_util.CheckField(t, "isOldLivery", boolToString(isOldLiveryConst), boolToString(objectInst.IsOldLivery))
	test_util.CheckField(t, "isSpecialLivery", boolToString(isSpecialLiveryConst), boolToString(objectInst.IsSpecialLivery))

	if emptyCode != (len(objectInst.Code) == 0) {
		t.Errorf("The value of the code %s is not as expected.", objectInst.Code)
	}
}

func testSetup(t *testing.T) {
	initDyno()

	modelmake.InitConn()
	airline.InitConn()
	airplane.InitConn()
	airplanemake.InitConn()

	modelMakeInst = modelmake.ModelMake{}
	modelMakeInst.Code = modelMakeConst
	modelMakeInst.Name = modelMakeNameConst
	(&modelMakeInst).Put()

	airlineInst = airline.Airline{}
	airlineInst.Code = airlineConst
	airlineInst.Callsign = airlineCallsignConst
	airlineInst.Country = countryConst
	airlineInst.Iata = airlineIataConst
	airlineInst.Name = airlineNameConst
	(&airlineInst).Put()

	airplaneMakeInst = airplanemake.AirplaneMake{}
	airplaneMakeInst.Code = airplaneMakeConst
	airplaneMakeInst.Name = airplaneMakeNameConst
	airplaneMakeInst.Abbreviation = airplaneMakeAbbrevConst
	airplaneMakeInst.Country = countryConst
	(&airplaneMakeInst).Put()

	airplaneInst = airplane.Airplane{}
	airplaneInst.Code = airplaneConst
	airplaneInst.Name = airplaneNameConst
	airplaneInst.Iata = airplaneIataConst
	airplaneInst.Make = airplaneMakeInst.Code
	(&airplaneInst).Put()

	t.Log("Test setup completed.")
}

func testTearDown(t *testing.T) {
	modelMakeInst.Delete()
	airlineInst.Delete()
	airplaneMakeInst.Delete()
	airplaneInst.Delete()

	t.Log("Teardown completed.")
}

func createImage() image.Image {
	width := 200
	height := 100

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{100, 200, 200, 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}

	return img
}

func checkImageInModel(t *testing.T, modelInstList []*Model, expectedListCount int, expectedImageCount int) {
	if len(modelInstList) != expectedListCount {
		t.Errorf("The list of models for which the picture was added has not been successfully retrieved.")
	} else {
		for _, modelInst := range modelInstList {
			if modelInst.PictureList == nil || len(modelInst.PictureList) != expectedImageCount {
				t.Errorf("The number of images for model with code %s is not correct. Expected %d but got %d", modelInst.CodeValue(), expectedImageCount, len(modelInst.PictureList))
			}
		}
	}
}

func testFile(t *testing.T, codeArr []string) string {
	pipeRead, pipeWrite := io.Pipe()
	writer := multipart.NewWriter(pipeWrite)
	var filename string

	go func() {
		defer writer.Close()

		part, _ := writer.CreateFormFile("image", "imagename.png")
		img := createImage()

		err := png.Encode(part, img)
		if err != nil {
			t.Errorf("An error occurred encoding the image: %v", err)
		}
	}()

	request := httptest.NewRequest("POST", "/", pipeRead)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	imageFile, _, _ := request.FormFile("image")

	if imageFile == nil {
		t.Errorf("The image file is NIL!")
	} else {
		modelInstList, modelErr := AddModelPicture(imageFile, codeArr)
		if modelErr != nil {
			t.Errorf("An error has been returned by the AddMOdelPicture method: %v", modelErr)
		} else if modelInstList == nil {
			t.Errorf("The model list returned by AddModelPicture has come back empty.")
		} else {
			checkImageInModel(t, modelInstList, 2, 1)

			// Add again to simulate a different image.
			modelInstList, modelErr = AddModelPicture(imageFile, codeArr)
			if modelErr != nil {
				t.Errorf("The AddModelPicture returned an error: %v", modelErr)
			} else {
				lenList := len(modelInstList)
				if lenList != len(codeArr) {
					t.Errorf("The length of the modelInstList is not as expected: %d", lenList)
				} else {
					log.Printf("Model Retrieved: \n%v", modelInstList)
					checkImageInModel(t, modelInstList, 2, 2)

					if len(modelInstList[0].PictureList) > 0 {
						filename = modelInstList[0].PictureList[0]
					}
				}
			}
		}
	}

	return filename
}

func CreateObjectInst(reg string) *Model {
	objectInst := ObjectFactory()

	objectInst.Airline = airlineConst
	objectInst.Airplane = airplaneConst
	objectInst.ModelMake = modelMakeConst
	objectInst.Scale = scaleConst
	objectInst.Reg = reg
	objectInst.Notes = notesConst
	objectInst.IsCargo = isCargoConst
	objectInst.IsOldLivery = isOldLiveryConst
	objectInst.IsSpecialLivery = isSpecialLiveryConst

	return objectInst
}

func TestModel(t *testing.T) {
	if AdapterInst == nil {
		testSetup(t)
	}

	//Create airline from JSON
	objectInst := CreateObjectInst(regConst)
	objectInstLoad := ObjectFactory()

	jsonString, _ := json.MarshalIndent(&objectInst, " ", " ")
	objectInstLoad.FromJson([]byte(jsonString))

	t.Log("Initial check")
	chkModel(t, objectInstLoad, false)

	t.Log("Put check")
	objectInstLoad.Put()
	chkModel(t, objectInstLoad, false)

	modelRetrInst, getRetrErr := GetByCode(objectInstLoad.Code)
	if getRetrErr != nil {
		t.Errorf("Error in get after putting model. %v\n", getRetrErr)
	}
	chkModel(t, modelRetrInst, false)

	t.Log("Update, put again and retrieve")
	modelRetrInst.Reg = regNewConst
	modelRetrInst.Put()
	modelUpdtInst, getUpdtErr := GetByCode(objectInstLoad.Code)
	if getUpdtErr != nil {
		t.Errorf("Error in get after updating airline\n")
	}
	test_util.CheckField(t, "reg", regNewConst, modelUpdtInst.Reg)

	t.Log("Testing pic upload")
	objectInstSnd := CreateObjectInst(regConst + "2")
	objectInstSnd.Put()
	objectInstSndCode := objectInstSnd.Code
	filename := testFile(t, []string{objectInstLoad.Code, objectInstSnd.Code})

	t.Log("Check that the picture is linked to the two models.")
	objectInstListFromFile, fromFileErr := GetModelByPicture(filename)
	if fromFileErr != nil {
		t.Errorf("The model list cannot be retrieved for image %s. Error: %v", filename, fromFileErr)
	} else if objectInstListFromFile == nil {
		t.Errorf("The model list is empty for image %s. Not as expected...", filename)
	} else if len(objectInstListFromFile) != 2 {
		t.Errorf("The model list has the wrong number of items. Expected 2 but got %d", len(objectInstListFromFile))
	}

	if len(objectInstListFromFile) > 0 {
		t.Log("Step: Check that a picture can be removed from a single model")
		singlePicModelInst, removeErr := RemoveModelPicture(filename, objectInstListFromFile[0].CodeValue())
		if removeErr == nil {
			if len(singlePicModelInst.PictureList) != 1 {
				t.Errorf("The model with code %s has an unexpected number of pictures in its list, which is %v", singlePicModelInst.CodeValue(), singlePicModelInst.PictureList)
			} else {
				objectInstListFromFile, fromFileErr = GetModelByPicture(filename)
				if fromFileErr != nil {
					t.Errorf("The model list cannot be retrieved for image %s after removing pic from %s. Error: %v", filename, singlePicModelInst.CodeValue(), fromFileErr)
				} else {
					if len(objectInstListFromFile) != 1 {
						t.Errorf("The list of models has the wrong number of items. List: %v", objectInstListFromFile)
					} else if objectInstListFromFile[0].CodeValue() == singlePicModelInst.CodeValue() {
						t.Errorf("The retrieved model is wrong. Expected %s but got %s", singlePicModelInst.CodeValue(), objectInstListFromFile[0].CodeValue())
					}
				}
			}
		} else {
			t.Errorf("There was an error whilst trying to delete picture %s from model with code %s. Error: %v", filename, singlePicModelInst.CodeValue(), removeErr)
		}
	}

	t.Log("Step: Delete and check it's gone!")
	t.Logf("For model with code %s, picture %s, pictureList %v\n", modelUpdtInst.Code, modelUpdtInst.Picture, modelUpdtInst.PictureList)
	modelUpdtInst.Delete()
	modelEmptyInst, getEmptyErr := GetByCode(objectInstLoad.Code)
	if getEmptyErr == nil {
		t.Errorf("Error for unexistent object not produced when expected. Perhaps the object still exists?\n")
	}
	test_util.CheckField(t, "reg", "", modelEmptyInst.Reg)

	t.Logf("For model with code %s, picture %s, pictureList %v\n", objectInstSnd.Code, objectInstSnd.Picture, objectInstSnd.PictureList)
	objectInstSnd.Delete()
	modelEmptyInst, getEmptyErr = GetByCode(objectInstSndCode)
	if getEmptyErr == nil {
		t.Errorf("Error for unexistent object not produced when expected. Perhaps the object still exists?\n")
	}
	test_util.CheckField(t, "reg", "", modelEmptyInst.Reg)

	testTearDown(t)
	t.Log("Test for model has finished.")
}

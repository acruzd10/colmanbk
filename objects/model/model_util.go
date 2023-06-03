package model

import (
	"colmanback/db/dyno"
	"colmanback/db/s3"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
)

//----------------------------------------------------------------------------------------
func tagModelPicture(filename string, modelCodeList []string) ([]*Model, error) {
	var modelInst *Model
	var intlErr error
	modelList := []*Model{}

	for _, code := range modelCodeList {
		// Save model stub for the index.
		modelInst = &Model{}
		modelInst.Code = code
		modelInst.Picture = filename
		modelInst.Put()

		//Append the image to the actual model objects (in memory only)
		modelInst, intlErr = GetByCode(code)
		if intlErr == nil {
			modelInst.PictureList = append(modelInst.PictureList, filename)
			modelList = append(modelList, modelInst)
		} else {
			log.Printf("An error has occurred while tagging model with code %s for picture with filename %s. Error: %v", code, filename, intlErr)
			return nil, intlErr
		}
	}

	return modelList, nil
}

//----------------------------------------------------------------------------------------
func chkModelCode(modelCodeList []string) []string {
	var err error
	var validModelCodeList []string

	for _, modelCode := range modelCodeList {
		_, err = GetByCode(modelCode)
		if err == nil {
			validModelCodeList = append(validModelCodeList, modelCode)
		}
	}

	return validModelCodeList
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
	dynoInstModel.SetSortGSIName("picture-code-index")
	dynoInstModel.Config("model", "code", true, ObjectFactory, nil)
	AdapterInst = dynoInstModel

	fileInstModel := &s3.S3Adapter{}
	fileInstModel.Config("colman-pics", 1000)
	FileInst = fileInstModel
}

//----------------------------------------------------------------------------------------
func AddModelPicture(file multipart.File, modelCodeList []string) ([]*Model, error) {
	var addErr error
	var modelList []*Model

	nowTime := time.Now().Format(time.RFC3339)
	uuidName := uuid.New().String()
	filename := strings.Replace(nowTime, ":", "_", -1) + "-" + uuidName

	validModelCodeList := chkModelCode(modelCodeList)
	if len(validModelCodeList) > 0 {
		response, err := FileInst.AddFile(filename, file)
		if err == nil {
			if len(response.FileLocation) != 0 {
				modelList, addErr = tagModelPicture(filename, validModelCodeList)
			}
		} else {
			addErr = err
		}
	}

	return modelList, addErr
}

//----------------------------------------------------------------------------------------
func TagModelPicture(filename string, modelCode string) (*Model, error) {
	var modelCodeList []string = []string{modelCode}

	modelList, tagErr := tagModelPicture(filename, modelCodeList)
	if tagErr == nil && len(modelList) > 0 {
		return modelList[0], nil
	} else {
		return nil, tagErr
	}
}

//----------------------------------------------------------------------------------------
func GetModelByPicture(filename string) ([]*Model, error) {
	objectList, err := AdapterInst.GetObjectListBySort(filename)

	if err != nil {
		log.Printf("Error whilst retrieving list of models for picture %s\n. Error %v", filename, err)
	}

	return objectList, err
}

//----------------------------------------------------------------------------------------
func deleteModelPicture(filename string, isDeletingFromDB bool) ([]*Model, error) {
	var returnErr error
	var objectErr error
	var objectList []*Model

	storageErr := FileInst.DeleteFile(filename)
	if storageErr == nil {
		objectList, objectErr = GetModelByPicture(filename)
		if objectErr == nil {
			if !isDeletingFromDB {
				for _, objectInst := range objectList {
					removeModelPicture(objectInst, filename, true)
				}
			}
		} else {
			log.Printf("A database-related error has occurred while attempting to delete picture %s. Error %v", filename, objectErr)
			returnErr = objectErr
		}
	} else {
		log.Printf("A storage-related error has occurred while attempting to delete picture %s. Error %v", filename, storageErr)
		returnErr = storageErr
	}

	return objectList, returnErr
}

//----------------------------------------------------------------------------------------
func DeleteModelPicture(filename string) ([]*Model, error) {
	return deleteModelPicture(filename, false)
}

//----------------------------------------------------------------------------------------
func removeModelPicture(objectInst *Model, filename string, isDeletingFromFileStorage bool) (*Model, error) {
	var objectInstSub *Model
	var pictureList []string
	var retErr error

	objectInstSub = ObjectFactory()
	objectInstSub.Code = objectInst.Code
	objectInstSub.Picture = filename
	objectInstSub.Delete()

	for _, pictureInst := range objectInst.PictureList {
		if pictureInst != filename {
			pictureList = append(pictureList, pictureInst)
		}
	}

	if !isDeletingFromFileStorage {
		otherModelsList, otherModelsErr := GetModelByPicture(filename)
		if otherModelsErr == nil {
			if len(otherModelsList) == 0 {
				_, retErr = deleteModelPicture(filename, true)
			}
		} else {
			retErr = otherModelsErr
		}
	}

	objectInst.PictureList = pictureList

	return objectInst, retErr
}

//----------------------------------------------------------------------------------------
func RemoveModelPicture(filename string, modelCode string) (*Model, error) {
	objectInst, objectErr := GetByCode(modelCode)

	log.Printf("Trying to delete for code %s, model code %s, file %s, models pic list %v", modelCode, objectInst.Code, filename, objectInst.PictureList)

	if objectErr == nil {
		objectInst, objectErr = removeModelPicture(objectInst, filename, false)
	}

	if objectErr != nil {
		log.Printf("An error has occurred while trying to remove picture %s from model with code %s. Error: %v", filename, modelCode, objectErr)
	}

	return objectInst, objectErr
}

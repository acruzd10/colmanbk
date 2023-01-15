package s3

import (
	"log"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Adapter struct {
	sess          *session.Session
	uploader      *s3manager.Uploader
	s3svc         *s3.S3
	bucketName    string
	maxGetEntries int
}

type S3Response struct {
	FileLocation string `json:"fileLocation"`
}

//----------------------------------------------------------------------------------------
func (s3Adapter *S3Adapter) Config(bucketName string, maxGetEntries int) {
	s3Adapter.bucketName = bucketName
	s3Adapter.maxGetEntries = maxGetEntries
	s3Adapter.sess = session.Must(session.NewSession())
	s3Adapter.uploader = s3manager.NewUploader(s3Adapter.sess)
	s3Adapter.s3svc = s3.New(s3Adapter.sess)
}

//----------------------------------------------------------------------------------------
func (s3Adapter *S3Adapter) AddFile(fileName string, file multipart.File) (S3Response, error) {
	var response S3Response

	result, err := s3Adapter.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Adapter.bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		log.Printf("S3 failed to upload file %s into bucket %s. Error: %v", fileName, s3Adapter.bucketName, err)
	} else {
		response.FileLocation = result.Location
	}

	return response, err
}

//----------------------------------------------------------------------------------------
func (s3Adapter *S3Adapter) DeleteFiles(fileNameArr []string) error {
	var identifiersArr []*s3.ObjectIdentifier

	for _, fileName := range fileNameArr {
		var identifier s3.ObjectIdentifier

		identifier.Key = aws.String(strings.Trim(fileName, " "))
		identifiersArr = append(identifiersArr, &identifier)
	}

	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(s3Adapter.bucketName),
		Delete: &s3.Delete{
			Objects: identifiersArr,
			Quiet:   aws.Bool(false),
		},
	}

	_, err := s3Adapter.s3svc.DeleteObjects(input)
	if err != nil {
		log.Fatalf("S3 failed to delete files %v from bucket %s. Error: %v", fileNameArr, s3Adapter.bucketName, err)
	}

	return err
}

//----------------------------------------------------------------------------------------
func (s3Adapter *S3Adapter) DeleteFile(fileName string) error {
	fileNameArr := []string{fileName}

	return s3Adapter.DeleteFiles(fileNameArr)
}

//----------------------------------------------------------------------------------------
/*
func ListModelPictures(Code string) []byte {
	var fileList []string = []string{}
	var arrayString string = ""
	var arrayStringInit bool = false

	if sess == nil {
		sess = session.Must(session.NewSession())
	}

	svc := s3.New(sess)
	input := &s3.ListObjectsInput{
		Bucket:  aws.String(BUCKET_NAME),
		Prefix:  aws.String(Code + "/"),
		MaxKeys: aws.Int64(MAX_GET_ENTRIES),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		log.Printf("An error has occurred. %v", err)
	} else {
		for _, entry := range result.Contents {
			fileList = append(fileList, aws.StringValue(entry.Key))
		}
	}

	arrayString = "["
	for _, entryName := range fileList {
		if arrayStringInit {
			arrayString = arrayString + ","
		}
		arrayString = arrayString + "{\"entry\":" + "\"" + entryName + "\"}"
		arrayStringInit = true
	}

	arrayString = arrayString + "]"

	return []byte(arrayString)
}
*/

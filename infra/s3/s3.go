package s3

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"olx-clone/conf"
	"olx-clone/functions/logger"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/getsentry/sentry-go"
)

var log = logger.Log

var sessionObj, _ = session.NewSession(&aws.Config{
	Region: aws.String(conf.REGION)},
)

/*
Uploads a file to S3
Accepts filePath of local file and objectKey indicating the file name in S3
*/
func UploadFileS3(filePath string, objectKey string) (string, bool) {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		log.Errorln("Failed to open file:", filePath, err)
		return "", false
	}
	defer func() {
		if err := file.Close(); err != nil {
			sentry.CaptureException(err)
		}
	}()
	svc := s3manager.NewUploader(sessionObj)
	data, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(conf.S3Bucket),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		log.Errorln("error in file upload:", err)
		return "", false
	}
	log.Printf("Successfully uploaded %s to %s\n", filePath, objectKey)
	return data.Location, true
}

/*
Uploads a file to S3 to specific bucket
Accepts filePath of local file and objectKey indicating the file name in S3
*/
func UploadFileToS3Bucket(filePath, objectKey, bucketName string) (string, bool) {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		log.Errorln("Failed to open file:", filePath, err)
		return "", false
	}
	defer func() {
		if err := file.Close(); err != nil {
			sentry.CaptureException(err)
		}
	}()
	svc := s3manager.NewUploader(sessionObj)
	data, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		log.Errorln("error in file upload:", err)
		return "", false
	}
	log.Printf("Successfully uploaded %s to %s\n", filePath, objectKey)
	return data.Location, true
}

/*
Generates a presigned URL for S3 Object
objectKey is the file name / path for the file in S3
and timeInMin indicates time in minutes for which link will be active
*/
func GetPresignedURLS3(objectKey string, timeInMin int) string {
	if objectKey == "" {
		return ""
	}
	svc := s3.New(sessionObj)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(conf.S3Bucket),
		Key:    aws.String(objectKey),
	})
	urlStr, err := req.Presign(time.Duration(timeInMin) * time.Minute)
	if err != nil {
		log.Errorln("Failed to sign request", err)
		return ""
	}
	return urlStr
}

/*
Generates a presigned URL for S3 Object from given bucket
objectKey is the file name / path for the file in S3
and timeInMin indicates time in minutes for which link will be active
*/
func GetPresignedURLFromS3Bucket(objectKey string, bucketName string, timeInMin int) string {
	if objectKey == "" {
		return ""
	}
	svc := s3.New(sessionObj)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	urlStr, err := req.Presign(time.Duration(timeInMin) * time.Minute)
	if err != nil {
		log.Errorln("Failed to sign request", err)
		return ""
	}
	return urlStr
}

/*
GetPresignedURLForHTML generates a presigned url for HTML File
objectKey is the file name / path for the file in S3
and timeInMin indicates time in minutes for which link will be active
*/
func GetPresignedURLForHTML(objectKey string, timeInMin int) string {
	svc := s3.New(sessionObj)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(conf.S3Bucket),
		Key:                        aws.String(objectKey),
		ResponseContentType:        aws.String("text/html"),
		ResponseContentDisposition: aws.String("inline"),
	})
	urlStr, err := req.Presign(time.Duration(timeInMin) * time.Minute)
	if err != nil {
		log.Errorln("Failed to sign request", err)
		return ""
	}
	return urlStr
}

/*
GetLocalFilePath downloads the file locally and shares the file path
*/
func GetLocalFilePath(objectKey string) (string, error) {
	var fileName = objectKey
	if strings.Contains(objectKey, "/") {
		fileName = strings.Split(objectKey, "/")[1]
	}

	filePath, err := DownloadFile(objectKey, fileName)
	if err != nil {
		log.Errorln(err)
		return "", err
	}

	return filePath, nil
}

/*
GetPresignedURLS3CustomName generates a presigned URL for S3 Object with custom file name
objectKey is the file name / path for the file in S3
timeInMin indicates time in minutes for which link will be active
fileName indicates how the file will look on downloading
*/
func GetPresignedURLS3CustomName(objectKey string, timeInMin int, fileName string) string {
	svc := s3.New(sessionObj)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(conf.S3Bucket),
		Key:                        aws.String(objectKey),
		ResponseContentDisposition: aws.String(fmt.Sprintf("attachment;filename=%s", fileName)),
	})
	urlStr, err := req.Presign(time.Duration(timeInMin) * time.Minute)
	if err != nil {
		log.Errorln("Failed to sign request", err)
		return ""
	}
	return urlStr
}

/*
Uploads a file to S3
Accepts http req objec and attribute key name and objectKey indicating the file name in S3
*/
func ReadFromReqAndUploadFileS3(r *http.Request, fileName string, objectKey string) (string, bool) {

	file, _, err := r.FormFile(fileName)
	if err != nil {
		log.Errorln("Failed to open file:", fileName, err)
		return "", false
	}
	defer file.Close()

	svc := s3manager.NewUploader(sessionObj)
	data, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(conf.S3Bucket),
		Key:    aws.String(objectKey),
		Body:   file,
	})

	if err != nil {
		log.Errorln("error in file upload:", err)
		return "", false
	}
	log.Printf("Successfully uploaded to %s\n", data.Location)
	return data.Location, true
}

/*
Uploads a file to S3
Accepts file object, objectKey indicating the file name in S3 and bucket name
*/
func ReadFromMultipartFileAndUploadFileS3ToBucket(file multipart.File, objectKey, bucketName string) (string, bool) {
	defer file.Close()

	svc := s3manager.NewUploader(sessionObj)
	data, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})

	if err != nil {
		log.Errorln("error in file upload:", err)
		return "", false
	}
	log.Printf("Successfully uploaded to %s\n", data.Location)
	return data.Location, true
}

/*
Uploads a file to S3
Accepts file object and objectKey indicating the file name in S3
*/
func ReadFromMultipartFileAndUploadFileS3(file multipart.File, objectKey string) (string, bool) {
	defer file.Close()

	svc := s3manager.NewUploader(sessionObj)
	data, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(conf.S3Bucket),
		Key:    aws.String(objectKey),
		Body:   file,
	})

	if err != nil {
		log.Errorln("error in file upload:", err)
		return "", false
	}
	log.Printf("Successfully uploaded to %s\n", data.Location)
	return data.Location, true
}

/*
Uploads file to S3 based on io Reader
*/
func UploadRawFileS3(dataInput io.Reader, objectKey string) (string, bool) {
	svc := s3manager.NewUploader(sessionObj)
	data, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(conf.S3Bucket),
		Key:    aws.String(objectKey),
		Body:   dataInput,
	})

	if err != nil {
		log.Errorln("error in file upload:", err)
		return "", false
	}
	log.Printf("Successfully uploaded to %s\n", data.Location)
	return data.Location, true
}

/*
Returns result pointer for an s3 object
*/
func GetFileStream(objectKey string) *io.ReadCloser {
	result, err := s3.New(sessionObj).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(conf.S3Bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Errorln("error in getobject S3. Object Key: ", objectKey, "error: ", err)
		return nil
	}
	return &result.Body
}

/*
DownloadFile from s3 and return error (if any) and local file path
*/
func DownloadFile(objectKey string, fileName string) (string, error) {

	downloader := s3manager.NewDownloader(sessionObj)

	path := "/tmp/" + fileName
	file, err := os.Create(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(conf.S3Bucket),
			Key:    aws.String(objectKey),
		})

	if err != nil {
		return "", err
	}

	return path, nil
}

// ReadFromURLAndUploadFileS3 reads the file from a URL and uploads it to specified objectKey
func ReadFromURLAndUploadFileS3(url string, objectKey string) (string, bool) {

	response, err := http.Get(url)
	if err != nil {
		log.Errorln("error getting file", err)
		return "", false
	}
	defer response.Body.Close()

	svc := s3manager.NewUploader(sessionObj)
	data, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(conf.S3Bucket),
		Key:    aws.String(objectKey),
		Body:   response.Body,
	})

	if err != nil {
		log.Errorln("error in file upload:", err)
		return "", false
	}
	log.Printf("Successfully uploaded to %s\n", data.Location)
	return data.Location, true
}

// GetBase64FromS3 gets the base64 encoded value of a file on s3
func GetBase64FromS3(objectKey string) string {
	ioRPointer := GetFileStream(objectKey)
	if ioRPointer == nil {
		return ""
	}
	ioR := *ioRPointer
	buf := new(bytes.Buffer)
	buf.ReadFrom(ioR)
	newStr := buf.String()
	bstr := base64.StdEncoding.EncodeToString([]byte(newStr))
	return bstr
}

// UploadFromPublicURL gets the file from URL and uploads in s3 to specified key
func UploadFromPublicURL(fileURL string, objectKey string) bool {
	resp, err := http.Get(fileURL)
	if err != nil {
		log.Errorln(err)
		return false
	}
	defer resp.Body.Close()
	_, done := UploadRawFileS3(resp.Body, objectKey)
	return done
}

// CopyBetweenBuckets copies file between buckets
func CopyBetweenBuckets(objectKey string, sourceBucketName string, targetBucketName string) error {
	svc := s3.New(sessionObj)
	source := sourceBucketName + "/" + objectKey
	// copy the item
	_, err := svc.CopyObject(&s3.CopyObjectInput{Bucket: aws.String(targetBucketName), CopySource: aws.String(url.PathEscape(source)), Key: aws.String(objectKey)})
	if err != nil {
		log.Errorln(err)
		return err
	}
	// Wait to see if the item got copied
	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{Bucket: aws.String(targetBucketName), Key: aws.String(objectKey)})
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

// DeleteObject deletes an object from s3
func DeleteObject(objectKey string, bucketName string) error {
	svc := s3.New(sessionObj)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}
	_, err := svc.DeleteObject(input)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}

/*
Checks if file is in S3. Reads the Meta data of the object to check if the file exists
*/
func IsObjectInS3(objectKey string) bool {
	if objectKey == "" {
		return false
	}
	svc := s3.New(sessionObj)
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(conf.S3Bucket),
		Key:    aws.String(objectKey),
	})
	return err == nil
}

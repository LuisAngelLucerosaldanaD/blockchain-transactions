package files_s3

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"blockchain-transactions/internal/logger"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	AWS_BUCKET = "bjungle-files"
	AWS_REGION = "us-east-1"
)

func UploadFile(file *bytes.Reader, fullPath, bucket string) error {

	sess, err := sessionAws()
	if err != nil {
		return err
	}
	svc := s3.New(sess)
	input := &s3.PutObjectInput{
		ACL:         aws.String("public-read"),
		Body:        file,
		Bucket:      aws.String(bucket),
		Key:         aws.String(fullPath),
		ContentType: aws.String("application/pdf"),
	}
	_, err = svc.PutObject(input)
	if err != nil {
		if aErr, ok := err.(awserr.Error); ok {
			switch aErr.Code() {
			default:
				return aErr
			}
		} else {
			return aErr
		}
	}
	return nil
}

func GetObjectS3(bucket, path, fileName string) (string, error) {
	sess, err := sessionAws()
	if err != nil {
		return "", err
	}

	pathDest := "./download/" + fileName
	pathS3 := path + fileName

	file, err := os.Create(pathDest)
	if err != nil {
		return "", err
	}

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(pathS3),
		})
	if err != nil {
		logger.Error.Printf("couldn't get file. %v", err)
		return "", err
	}
	reader := bufio.NewReader(file)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		logger.Error.Printf("couldn't get file. %v", err)
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(content)

	if err := file.Close(); err != nil {
		logger.Error.Printf("couldn't close *file. %v", err)
		return "", err
	}
	err = os.Remove(pathDest)
	if err != nil {
		logger.Error.Printf("couldn't delete *file. %v", err)
		return "", err
	}

	return encoded, nil
}

func GetFileLink(bucket, _path string, fileNameS3 string) (string, error) {

	sess, err := sessionAws()
	svc := s3.New(sess)

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket + _path),
		Key:    aws.String(fileNameS3),
	}

	req, _ := svc.GetObjectRequest(params)

	url, err := req.Presign(15 * time.Minute) // Set link expiration time
	if err != nil {
		fmt.Println("Error en la peticion del URL", params, err)
	}
	return url, nil
}

func deleteFile() error {
	sess, err := sessionAws()
	if err != nil {
		return err
	}
	svc := s3.New(sess)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(AWS_BUCKET),
		Key:    aws.String("/01/01.pdf"),
	}

	_, err = svc.DeleteObject(input)
	if err != nil {
		aErr, ok := err.(awserr.Error)
		if ok {
			switch aErr.Code() {
			default:
				logger.Error.Printf("Error eliminando el archivo: %v", err)
				return err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			logger.Error.Printf("Error eliminando el archivo: %v", aErr.Error())
			return err
		}
		return nil
	}
	return nil
}

func sessionAws() (client.ConfigProvider, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_REGION),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		//logger.Error.Printf("iniciando sesión con aws: #{err}")
		fmt.Println("error iniciando sesión con aws: #{err}")
		return sess, err
	}
	return sess, nil
}

package files

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"blockchain-transactions/internal/env"
	"blockchain-transactions/internal/files_s3"
	"blockchain-transactions/internal/models"

	"github.com/google/uuid"
)

// s estructura de conexión s3
type s3 struct {
	user *models.User
	TxID string
}

func newDocumentFileS3Repository(user *models.User, txID string) *s3 {
	return &s3{
		user: user,
		TxID: txID,
	}
}

func (s *s3) upload(documentID int64, file *File) (*File, error) {
	c := env.NewConfiguration()
	var fullPath strings.Builder
	if file.Encoding == "" || file.OriginalFile == "" {
		return file, fmt.Errorf("couldn't create encoded file does not exist")
	}
	fl, err := base64.StdEncoding.DecodeString(file.Encoding)
	if err != nil {
		return file, err
	}
	file.Encoding = ""
	if fl == nil {
		return file, fmt.Errorf("couldn't create encoded file is null")
	}
	r := bytes.NewReader(fl)
	file.Path, file.FileName = s.getFullPath(file.OriginalFile, documentID)
	fullPath.WriteString(file.Path)
	fullPath.WriteString(file.FileName)
	file.Hash = s.getHashFromFile(fl)
	file.FileSize = int(r.Size())
	//TODO getNumberPage
	file.NumberPage = 1
	file.Bucket = c.Files.S3.Bucket
	err = files_s3.UploadFile(r, fullPath.String(), file.Bucket)
	if err != nil {
		return file, err
	}
	return file, nil
}

func (s *s3) getHashFromFile(file []byte) string {
	h := sha256.Sum256(file)
	return fmt.Sprintf("%x", h)
}

func (s *s3) getFullPath(originalFile string, documentID int64) (string, string) {
	fPath := fmt.Sprintf("/%d/%d/%d/%d/%d/", documentID, time.Now().Year(), time.Now().YearDay(), time.Now().Hour(), time.Now().Minute())
	fileName := fmt.Sprintf("%s%s", uuid.New(), filepath.Ext(originalFile))
	return fPath, fileName
}

func (s *s3) getFile(bucket, path, fileName string) (string, error) {
	return files_s3.GetObjectS3(bucket, path, fileName)
}

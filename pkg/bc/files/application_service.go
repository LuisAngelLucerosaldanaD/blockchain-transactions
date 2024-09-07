package files

import (
	"blockchain-transactions/internal/env"
	"blockchain-transactions/internal/models"
)

type PortsServerFile interface {
	UploadFile(idDocument int64, originalFile string, encoding string) (*File, error)
	GetFileByPath(path, fileName string) (*ResponseFile, int, error)
}

type service struct {
	repositoryS3 ServicesFileDocumentsRepository
	user         *models.User
	txID         string
}

func NewFileService(repositoryS3 ServicesFileDocumentsRepository, user *models.User, TxID string) PortsServerFile {
	return &service{repositoryS3: repositoryS3, user: user, txID: TxID}
}

func (s *service) UploadFile(idDocument int64, originalFile string, encoding string) (*File, error) {
	file := NewUploadFile(idDocument, originalFile, encoding)
	return s.repositoryS3.upload(idDocument, file)
}

func (s *service) GetFileByPath(path, fileName string) (*ResponseFile, int, error) {
	e := env.NewConfiguration()

	rf := ResponseFile{}
	file, err := s.repositoryS3.getFile(e.Files.S3.Bucket, path, fileName)
	if err != nil {
		return nil, 0, err
	}
	rf.Encoding = file
	rf.NameDocument = fileName
	return &rf, 29, nil
}

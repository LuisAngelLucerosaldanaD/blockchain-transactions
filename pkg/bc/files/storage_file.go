package files

import (
	"strings"

	"blockchain-transactions/internal/env"
	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/models"
)

const (
	S3 = "s3"
)

type ServicesFileDocumentsRepository interface {
	upload(documentID int64, file *File) (*File, error)
	getFile(bucket, path, fileName string) (string, error)
}

func FactoryFileDocumentRepository(user *models.User, txID string) ServicesFileDocumentsRepository {
	var s ServicesFileDocumentsRepository
	c := env.NewConfiguration()
	repo := strings.ToLower(c.Files.Repo)
	switch repo {
	case S3:
		return newDocumentFileS3Repository(user, txID)
	default:
		logger.Error.Println("el repositorio de documentos no está implementado.", repo)
	}
	return s
}

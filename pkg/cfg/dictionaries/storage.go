package dictionaries

import (
	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
)

type ServicesDictionariesRepository interface {
	create(m *Dictionary) error
	update(m *Dictionary) error
	delete(id int) error
	getByID(id int) (*Dictionary, error)
	getAll() ([]*Dictionary, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesDictionariesRepository {
	var s ServicesDictionariesRepository

	engine := db.DriverName()

	switch engine {
	case Postgresql:
		return newDictionaryPsqlRepository(db, user, txID)

	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}

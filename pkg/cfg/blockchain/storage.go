package blockchain

import (
	"github.com/jmoiron/sqlx"

	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesBlockchainRepository interface {
	create(m *Blockchain) error
	update(m *Blockchain) error
	delete(id string) error
	getByID(id string) (*Blockchain, error)
	getAll() ([]*Blockchain, error)
	getLasted() (*Blockchain, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesBlockchainRepository {
	var s ServicesBlockchainRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newBlockchainPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}

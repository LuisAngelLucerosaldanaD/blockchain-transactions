package transaction

import (
	"github.com/jmoiron/sqlx"

	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesTransactionRepository interface {
	create(m *Transaction) error
	update(m *Transaction) error
	delete(id string) error
	getByID(id string) (*Transaction, error)
	getAll(toID string, blockID int64) ([]*Transaction, error)
	getByBlockId(block int64) ([]*Transaction, error)
	GetCountTransactionByID(block int64) int
	getByIds(dni string, ids []string) ([]*Transaction, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesTransactionRepository {
	var s ServicesTransactionRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newTransactionPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}

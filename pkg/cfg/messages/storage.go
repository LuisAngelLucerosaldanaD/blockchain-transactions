package messages

import (
	"github.com/jmoiron/sqlx"

	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesMessagesRepository interface {
	create(m *Messages) error
	update(m *Messages) error
	delete(id int) error
	getByID(id int) (*Messages, error)
	getAll() ([]*Messages, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesMessagesRepository {
	var s ServicesMessagesRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newMessagesPsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}

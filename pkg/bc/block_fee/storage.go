package block_fee

import (
	"github.com/jmoiron/sqlx"

	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/models"
)

const (
	Postgresql = "postgres"
)

type ServicesBlockFeeRepository interface {
	create(m *BlockFee) error
	update(m *BlockFee) error
	delete(id string) error
	getByID(id string) (*BlockFee, error)
	getAll() ([]*BlockFee, error)
	getByBlockID(blockId int64) (*BlockFee, error)
}

func FactoryStorage(db *sqlx.DB, user *models.User, txID string) ServicesBlockFeeRepository {
	var s ServicesBlockFeeRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newBlockFeePsqlRepository(db, user, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}

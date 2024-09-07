package blockchain

import (
	"database/sql"
	"fmt"
	"time"

	"blockchain-transactions/internal/models"
	"github.com/jmoiron/sqlx"
)

// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newBlockchainPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *Blockchain) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	const psqlInsert = `INSERT INTO cfg.blockchain (id ,fee_blion, fee_miner, fee_validator, fee_node, ttl_block, max_transactions, max_validators, max_miners, tickets_price, lottery_ttl, wallet_main, deleted_at, created_at, updated_at) VALUES (:id ,:fee_blion, :fee_miner, :fee_validator, :fee_node, :ttl_block, :max_transactions, :max_validators, :max_miners, :tickets_price, :lottery_ttl, :wallet_main, :deleted_at,:created_at, :updated_at) `
	rs, err := s.DB.NamedExec(psqlInsert, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *psql) update(m *Blockchain) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE cfg.blockchain SET fee_blion = :fee_blion, fee_miner = :fee_miner, fee_validator = :fee_validator, fee_node = :fee_node, ttl_block = :ttl_block, max_transactions = :max_transactions, max_validators = :max_validators, max_miners = :max_miners, tickets_price = :tickets_price, lottery_ttl = :lottery_ttl, wallet_main = :wallet_main, deleted_at = :deleted_at, updated_at = :updated_at WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Delete elimina un registro de la BD
func (s *psql) delete(id string) error {
	const psqlDelete = `DELETE FROM cfg.blockchain WHERE id = :id `
	m := Blockchain{ID: id}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *psql) getByID(id string) (*Blockchain, error) {
	const psqlGetByID = `SELECT id , fee_blion, fee_miner, fee_validator, fee_node, ttl_block, max_transactions, max_validators, max_miners, tickets_price, lottery_ttl, wallet_main, deleted_at, created_at, updated_at FROM cfg.blockchain WHERE id = $1 `
	mdl := Blockchain{}
	err := s.DB.Get(&mdl, psqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) getAll() ([]*Blockchain, error) {
	var ms []*Blockchain
	const psqlGetAll = ` SELECT id , fee_blion, fee_miner, fee_validator, fee_node, ttl_block, max_transactions, max_validators, max_miners, tickets_price, lottery_ttl, wallet_main, deleted_at, created_at, updated_at FROM cfg.blockchain `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

func (s *psql) getLasted() (*Blockchain, error) {
	const psqlGetByID = `SELECT id , fee_blion, fee_miner, fee_validator, fee_node, ttl_block, max_transactions, max_validators, max_miners, tickets_price, lottery_ttl, wallet_main, deleted_at, created_at, updated_at FROM cfg.blockchain order by created_at desc`
	mdl := Blockchain{}
	err := s.DB.Get(&mdl, psqlGetByID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

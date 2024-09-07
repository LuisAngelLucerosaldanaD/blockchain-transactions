package transaction

import (
	"blockchain-transactions/internal/utils"
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

func newTransactionPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *Transaction) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date
	const psqlInsert = `INSERT INTO bc.transaction (id ,from_id, to_id, amount,type_id, data, files, block, created_at, updated_at) VALUES (:id ,:from_id, :to_id, :amount, :type_id, :data, :files, :block,:created_at, :updated_at) `
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
func (s *psql) update(m *Transaction) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE bc.transaction SET from_id = :from_id, to_id = :to_id, amount = :amount, type_id = :type_id, data = :data, files = :files, block = :block, updated_at = :updated_at WHERE id = :id `
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
	const psqlDelete = `DELETE FROM bc.transaction WHERE id = :id `
	m := Transaction{ID: id}
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
func (s *psql) getByID(id string) (*Transaction, error) {
	const psqlGetByID = `SELECT id , from_id, to_id, amount, type_id, data, files, block, created_at, updated_at FROM bc.transaction WHERE id = $1 `
	mdl := Transaction{}
	err := s.DB.Get(&mdl, psqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	//mdl.Data = ciphers.Decrypt(mdl.Data)
	return &mdl, nil
}

// GetAll consulta todos los registros de la BD
func (s *psql) getAll(toID string, blockID int64) ([]*Transaction, error) {
	var ms []*Transaction
	const psqlGetAll = ` SELECT DISTINCT t.id , t.from_id, t.to_id, t.amount, t.type_id, t.data, t.files, t.block, t.created_at, t.updated_at FROM bc.transaction t
                            JOIN auth.wallet  w ON (t.to_id = w.id or t.from_id = w.id) WHERE w.dni = $1 and t.block > $2`

	err := s.DB.Select(&ms, psqlGetAll, toID, blockID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

// getByBlockId consulta un registro por block id
func (s *psql) getByBlockId(block int64) ([]*Transaction, error) {
	const psqlGetByID = `SELECT id , from_id, to_id, amount, type_id, data, files, block, created_at, updated_at FROM bc.transaction WHERE block = $1 `
	var mdl []*Transaction
	err := s.DB.Select(&mdl, psqlGetByID, block)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return mdl, err
	}
	return mdl, nil
}

func (s *psql) GetCountTransactionByID(block int64) int {
	const queryGetCountTransactionByID = `select count(*) from bc.transaction t where t.block = $1;`
	var totalTransaction int
	err := s.DB.Get(&totalTransaction, queryGetCountTransactionByID, block)
	if err != nil {
		if err == sql.ErrNoRows {
			return totalTransaction
		}
		return totalTransaction
	}
	return totalTransaction
}

func (s *psql) getByIds(dni string, ids []string) ([]*Transaction, error) {
	var ms []*Transaction
	const psqlGetAll = ` SELECT DISTINCT t.id , t.from_id, t.to_id, t.amount, t.type_id, t.data, t.files, t.block, t.created_at, t.updated_at FROM bc.transaction t
                            JOIN auth.wallet w ON (t.to_id = w.id or t.from_id = w.id) WHERE w.dni = $1 and t.id in ($2)`

	err := s.DB.Select(&ms, psqlGetAll, dni, utils.SliceToString(ids))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

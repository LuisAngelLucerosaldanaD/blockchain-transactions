package messages

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

func newMessagesPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *Messages) error {

	const psqlInsert = `INSERT INTO cfg.Messages (id, spa, eng, type_message) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.ID,
		m.Spa,
		m.Eng,
		m.TypeMessage,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil

}

// Update actualiza un registro en la BD
func (s *psql) update(m *Messages) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE cfg.Messages SET spa = :spa, eng = :eng, type_message = :type_message, updated_at = :updated_at WHERE id = :id `
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
func (s *psql) delete(id int) error {
	const psqlDelete = `DELETE FROM cfg.Messages WHERE id = :id `
	m := Messages{ID: id}
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
func (s *psql) getByID(id int) (*Messages, error) {
	const psqlGetByID = `SELECT id, spa, eng, type_message, created_at, updated_at FROM cfg.Messages WHERE id = $1 `
	mdl := Messages{}
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
func (s *psql) getAll() ([]*Messages, error) {
	var ms []*Messages
	const psqlGetAll = ` SELECT id, spa, eng, type_message, created_at, updated_at FROM cfg.Messages `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

package categories

import (
	"blockchain-transactions/internal/models"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

// psql estructura de conexión a la BD de postgresql
type psql struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newCategoriesPsqlRepository(db *sqlx.DB, user *models.User, txID string) *psql {
	return &psql{
		DB:   db,
		user: user,
		TxID: txID,
	}
}

// Create registra en la BD
func (s *psql) create(m *Categories) error {

	const psqlInsert = `INSERT INTO cfg.categories (id, name, color, icon) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	stmt, err := s.DB.Prepare(psqlInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		m.ID,
		m.Name,
		m.Color,
		m.Icon,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// update actualiza un registro en la BD
func (s *psql) update(m *Categories) error {
	date := time.Now()
	m.UpdatedAt = date
	const psqlUpdate = `UPDATE cfg.categories SET name = :name, color = :color, icon = :icon, updated_at = :updated_at WHERE id = :id `
	rs, err := s.DB.NamedExec(psqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// delete elimina un registro de la BD
func (s *psql) delete(id string) error {
	const psqlDelete = `DELETE FROM cfg.categories WHERE id = :id `
	m := Categories{ID: id}
	rs, err := s.DB.NamedExec(psqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// getByID consulta un registro por su ID
func (s *psql) getByID(id string) (*Categories, error) {
	const psqlGetByID = `SELECT id , name, icon, color, created_at, updated_at FROM cfg.categories WHERE id = $1 `
	mdl := Categories{}
	err := s.DB.Get(&mdl, psqlGetByID, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	return &mdl, nil
}

// getAll consulta todos los registros de la BD
func (s *psql) getAll() ([]*Categories, error) {
	var ms []*Categories
	const psqlGetAll = ` SELECT id , name, icon, color, created_at, updated_at FROM cfg.categories `

	err := s.DB.Select(&ms, psqlGetAll)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return ms, err
	}
	return ms, nil
}

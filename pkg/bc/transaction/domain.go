package transaction

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de Transaction
type Transaction struct {
	ID        string    `json:"id" db:"id" valid:"required,uuid"`
	From      string    `json:"from" db:"from_id" valid:"required"`
	To        string    `json:"to" db:"to_id" valid:"required"`
	Amount    float64   `json:"amount" db:"amount" valid:"required"`
	TypeID    int       `json:"type_id" db:"type_id" valid:"required"`
	Data      string    `json:"data" db:"data" valid:"required"`
	Files     string    `json:"files" db:"files" valid:"-"`
	Block     int64     `json:"block" db:"block" valid:"-"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewTransaction(id string, from string, to string, amount float64, typeId int, data, files string, block int64) *Transaction {
	return &Transaction{
		ID:     id,
		From:   from,
		To:     to,
		Amount: amount,
		TypeID: typeId,
		Data:   data,
		Files:  files,
		Block:  block,
	}
}

func (m *Transaction) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

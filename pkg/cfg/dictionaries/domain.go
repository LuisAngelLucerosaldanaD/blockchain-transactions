package dictionaries

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de Dictionaries
type Dictionary struct {
	ID          int       `json:"id" db:"id" valid:"-"`
	Name        string    `json:"name" db:"name" valid:"required"`
	Value       string    `json:"value" db:"value" valid:"required"`
	Description string    `json:"description" db:"description" valid:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewDictionaries(id int, name string, value string, description string) *Dictionary {
	return &Dictionary{
		ID:          id,
		Name:        name,
		Value:       value,
		Description: description,
	}
}

func NewCreateDictionaries(name string, value string, description string) *Dictionary {
	return &Dictionary{
		Name:        name,
		Value:       value,
		Description: description,
	}
}

func (m *Dictionary) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

package categories

import (
	"github.com/asaskevich/govalidator"
	"time"
)

// Categories estructura de Categories
type Categories struct {
	ID        string    `json:"id" db:"id" valid:"-"`
	Name      string    `json:"name" db:"name" valid:"required"`
	Icon      string    `json:"icon" db:"icon" valid:"required"`
	Color     string    `json:"color" db:"color" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewCategories(id, name, icon, color string) *Categories {
	return &Categories{
		ID:    id,
		Name:  name,
		Icon:  icon,
		Color: color,
	}
}

func NewCreateCategories(name, icon, color string) *Categories {
	return &Categories{
		Name:  name,
		Icon:  icon,
		Color: color,
	}
}

func (m *Categories) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

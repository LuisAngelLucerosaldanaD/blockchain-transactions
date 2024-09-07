package messages

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de Dictionaries
type Messages struct {
	ID          int       `json:"id" db:"id" valid:"-"`
	Spa         string    `json:"spa" db:"spa" valid:"required"`
	Eng         string    `json:"eng" db:"eng" valid:"required"`
	TypeMessage int       `json:"type_message" db:"type_message" valid:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewMessages(id int, spa string, eng string, typeMessage int) *Messages {
	return &Messages{
		ID:          id,
		Spa:         spa,
		Eng:         eng,
		TypeMessage: typeMessage,
	}
}

func (m *Messages) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

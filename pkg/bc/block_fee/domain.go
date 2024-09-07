package block_fee

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Model estructura de BlockFee
type BlockFee struct {
	ID        string    `json:"id" db:"id" valid:"required,uuid"`
	BlockId   int64     `json:"block_id" db:"block_id" valid:"required"`
	Fee       float64   `json:"fee" db:"fee" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewBlockFee(id string, blockId int64, fee float64) *BlockFee {
	return &BlockFee{
		ID:      id,
		BlockId: blockId,
		Fee:     fee,
	}
}

func (m *BlockFee) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

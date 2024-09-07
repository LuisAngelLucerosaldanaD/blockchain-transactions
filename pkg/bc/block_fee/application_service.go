package block_fee

import (
	"fmt"

	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/models"
	"github.com/asaskevich/govalidator"
)

type PortsServerBlockFee interface {
	CreateBlockFee(id string, blockId int64, fee float64) (*BlockFee, int, error)
	UpdateBlockFee(id string, blockId int64, fee float64) (*BlockFee, int, error)
	DeleteBlockFee(id string) (int, error)
	GetBlockFeeByID(id string) (*BlockFee, int, error)
	GetAllBlockFee() ([]*BlockFee, error)
	GetBlockFeeByBlockID(blockId int64) (*BlockFee, int, error)
}

type service struct {
	repository ServicesBlockFeeRepository
	user       *models.User
	txID       string
}

func NewBlockFeeService(repository ServicesBlockFeeRepository, user *models.User, TxID string) PortsServerBlockFee {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateBlockFee(id string, blockId int64, fee float64) (*BlockFee, int, error) {
	m := NewBlockFee(id, blockId, fee)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create BlockFee :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateBlockFee(id string, blockId int64, fee float64) (*BlockFee, int, error) {
	m := NewBlockFee(id, blockId, fee)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update BlockFee :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteBlockFee(id string) (int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
	}

	if err := s.repository.delete(id); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetBlockFeeByID(id string) (*BlockFee, int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetBlockFeeByBlockID(blockId int64) (*BlockFee, int, error) {
	if blockId <= 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("blockId is required"))
		return nil, 15, fmt.Errorf("blockId is required")
	}
	m, err := s.repository.getByBlockID(blockId)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByBlockID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllBlockFee() ([]*BlockFee, error) {
	return s.repository.getAll()
}

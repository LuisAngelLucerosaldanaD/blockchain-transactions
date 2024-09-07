package transaction

import (
	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/models"
	"fmt"

	"github.com/asaskevich/govalidator"
)

type PortsServerTransaction interface {
	CreateTransaction(id string, from string, to string, amount float64, typeId int, data, files string, block int64) (*Transaction, int, error)
	UpdateTransaction(id string, from string, to string, amount float64, typeId int, data, files string, block int64) (*Transaction, int, error)
	DeleteTransaction(id string) (int, error)
	GetTransactionByID(id string) (*Transaction, int, error)
	GetAllTransaction(toID string, blockID int64) ([]*Transaction, int, error)
	GetTransactionsByBlockID(blockId int64) ([]*Transaction, error)
	GetTransactionByIds(dni string, ids []string) ([]*Transaction, int, error)
}

type service struct {
	repository ServicesTransactionRepository
	user       *models.User
	txID       string
}

func NewTransactionService(repository ServicesTransactionRepository, user *models.User, TxID string) PortsServerTransaction {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateTransaction(id string, from string, to string, amount float64, typeId int, data, files string, block int64) (*Transaction, int, error) {
	m := NewTransaction(id, from, to, amount, typeId, data, files, block)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Transaction :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateTransaction(id string, from string, to string, amount float64, typeId int, data, files string, block int64) (*Transaction, int, error) {
	m := NewTransaction(id, from, to, amount, typeId, data, files, block)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Transaction :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteTransaction(id string) (int, error) {
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

func (s *service) GetTransactionByID(id string) (*Transaction, int, error) {
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

func (s *service) GetAllTransaction(toID string, blockID int64) ([]*Transaction, int, error) {
	m, err := s.repository.getAll(toID, blockID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetTransactionsByBlockID(blockId int64) ([]*Transaction, error) {
	return s.repository.getByBlockId(blockId)
}

func (s *service) GetTransactionByIds(dni string, ids []string) ([]*Transaction, int, error) {
	m, err := s.repository.getByIds(dni, ids)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByIds row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

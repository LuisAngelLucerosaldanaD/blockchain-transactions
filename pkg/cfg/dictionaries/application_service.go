package dictionaries

import (
	"fmt"

	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/models"
)

type PortsServerDictionaries interface {
	CreateDictionaries(name string, value string, description string) (*Dictionary, int, error)
	UpdateDictionaries(id int, name string, value string, description string) (*Dictionary, int, error)
	DeleteDictionaries(id int) (int, error)
	GetDictionariesByID(id int) (*Dictionary, int, error)
	GetAllDictionaries() ([]*Dictionary, error)
}

type service struct {
	repository ServicesDictionariesRepository
	user       *models.User
	txID       string
}

func NewDictionariesService(repository ServicesDictionariesRepository, user *models.User, TxID string) PortsServerDictionaries {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateDictionaries(name string, value string, description string) (*Dictionary, int, error) {
	m := NewCreateDictionaries(name, value, description)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Dictionaries :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateDictionaries(id int, name string, value string, description string) (*Dictionary, int, error) {
	m := NewDictionaries(id, name, value, description)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Dictionaries :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteDictionaries(id int) (int, error) {
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return 15, fmt.Errorf("id is required")
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

func (s *service) GetDictionariesByID(id int) (*Dictionary, int, error) {
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return nil, 15, fmt.Errorf("id is required")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllDictionaries() ([]*Dictionary, error) {
	return s.repository.getAll()
}

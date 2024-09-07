package categories

import (
	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/models"
	"fmt"
)

type PortsServerCategories interface {
	CreateCategories(id, name, icon, color string) (*Categories, int, error)
	UpdateCategories(id, name, icon, color string) (*Categories, int, error)
	DeleteCategories(id string) (int, error)
	GetCategoriesByID(id string) (*Categories, int, error)
	GetAllCategories() ([]*Categories, error)
}

type service struct {
	repository ServicesCategoriesRepository
	user       *models.User
	txID       string
}

func NewCategoriesService(repository ServicesCategoriesRepository, user *models.User, TxID string) PortsServerCategories {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateCategories(id, name, icon, color string) (*Categories, int, error) {
	m := NewCategories(id, name, icon, color)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Categories :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateCategories(id, name, icon, color string) (*Categories, int, error) {
	m := NewCategories(id, name, icon, color)
	if id == "" {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Categories :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteCategories(id string) (int, error) {
	if id == "" {
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

func (s *service) GetCategoriesByID(id string) (*Categories, int, error) {
	if id == "" {
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

func (s *service) GetAllCategories() ([]*Categories, error) {
	return s.repository.getAll()
}

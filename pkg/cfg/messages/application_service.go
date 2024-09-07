package messages

import (
	"fmt"

	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/models"
)

type PortsServerMessages interface {
	CreateMessages(id int, name string, value string, typeMessage int) (*Messages, int, error)
	UpdateMessages(id int, name string, value string, typeMessage int) (*Messages, int, error)
	DeleteMessages(id int) (int, error)
	GetMessagesByID(id int) (*Messages, int, error)
	GetAllMessages() ([]*Messages, error)
}

type service struct {
	repository ServicesMessagesRepository
	user       *models.User
	txID       string
}

func NewMessagesService(repository ServicesMessagesRepository, user *models.User, TxID string) PortsServerMessages {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateMessages(id int, name string, value string, typeMessage int) (*Messages, int, error) {
	m := NewMessages(id, name, value, typeMessage)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Messages :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateMessages(id int, name string, value string, typeMessage int) (*Messages, int, error) {
	m := NewMessages(id, name, value, typeMessage)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Messages :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteMessages(id int) (int, error) {
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

func (s *service) GetMessagesByID(id int) (*Messages, int, error) {
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

func (s *service) GetAllMessages() ([]*Messages, error) {
	return s.repository.getAll()
}

package blockchain

import (
	"fmt"
	"time"

	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/models"
	"github.com/asaskevich/govalidator"
)

type PortsServerBlockchain interface {
	CreateBlockchain(id string, feeBlion float32, feeMiner float32, feeValidator float32, feeNode float32, ttlBlock int, maxTransactions int, deletedAt time.Time, maxMiners int, maxValidators int, ticketPrice int, lotteryTtl int, walletMain string) (*Blockchain, int, error)
	UpdateBlockchain(id string, feeBlion float32, feeMiner float32, feeValidator float32, feeNode float32, ttlBlock int, maxTransactions int, deletedAt time.Time, maxMiners int, maxValidators int, ticketPrice int, lotteryTtl int, walletMain string) (*Blockchain, int, error)
	DeleteBlockchain(id string) (int, error)
	GetBlockchainByID(id string) (*Blockchain, int, error)
	GetAllBlockchain() ([]*Blockchain, error)
	MustCloseBlock(TtlBlock time.Time, transactions int) bool
	GetFeeBLion(amount float64) float64
	GetLasted() (*Blockchain, error)
}

type service struct {
	repository ServicesBlockchainRepository
	user       *models.User
	txID       string
}

func NewBlockchainService(repository ServicesBlockchainRepository, user *models.User, TxID string) PortsServerBlockchain {
	return &service{repository: repository, user: user, txID: TxID}
}

func (s *service) CreateBlockchain(id string, feeBlion float32, feeMiner float32, feeValidator float32, feeNode float32, ttlBlock int, maxTransactions int, deletedAt time.Time, maxMiners int, maxValidators int, ticketPrice int, lotteryTtl int, walletMain string) (*Blockchain, int, error) {
	m := NewBlockchain(id, feeBlion, feeMiner, feeValidator, feeNode, ttlBlock, maxTransactions, deletedAt, maxMiners, maxValidators, ticketPrice, lotteryTtl, walletMain)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Blockchain :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateBlockchain(id string, feeBlion float32, feeMiner float32, feeValidator float32, feeNode float32, ttlBlock int, maxTransactions int, deletedAt time.Time, maxMiners int, maxValidators int, ticketPrice int, lotteryTtl int, walletMain string) (*Blockchain, int, error) {
	m := NewBlockchain(id, feeBlion, feeMiner, feeValidator, feeNode, ttlBlock, maxTransactions, deletedAt, maxMiners, maxValidators, ticketPrice, lotteryTtl, walletMain)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Blockchain :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteBlockchain(id string) (int, error) {
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

func (s *service) GetBlockchainByID(id string) (*Blockchain, int, error) {
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

func (s *service) GetAllBlockchain() ([]*Blockchain, error) {
	return s.repository.getAll()
}

func (s *service) MustCloseBlock(TtlBlock time.Time, transactions int) bool {
	config, err := s.repository.getLasted()
	if err != nil {
		logger.Error.Println("No se pudo obtener la configuracion de la blockchain, err: ", err.Error())
		return false
	}

	if config == nil {
		logger.Error.Println("No se pudo obtener la configuracion de la blockchain")
		return false
	}

	if transactions >= config.MaxTransactions {
		return true
	}

	lifeBlock := time.Now().Sub(TtlBlock).Seconds()

	if int(lifeBlock) > (config.TtlBlock * 1000) {
		return true
	}

	return false
}

func (s *service) GetFeeBLion(amount float64) float64 {
	config, err := s.repository.getLasted()
	if err != nil {
		logger.Error.Println("No se pudo obtener la configuracion de la blockchain, err: ", err.Error())
		return amount
	}

	if config == nil {
		logger.Error.Println("No se pudo obtener la configuracion de la blockchain")
		return amount
	}

	return amount * float64(config.FeeBlion)
}

func (s *service) GetLasted() (*Blockchain, error) {
	config, err := s.repository.getLasted()
	if err != nil {
		logger.Error.Println("No se pudo obtener la configuracion de la blockchain, err: ", err.Error())
		return nil, err
	}

	if config == nil {
		logger.Error.Println("No se pudo obtener la configuracion de la blockchain")
		return nil, fmt.Errorf("no se pudo obtener la configuracion de la blockchain")
	}

	return config, nil
}

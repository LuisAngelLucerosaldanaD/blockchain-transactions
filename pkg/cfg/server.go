package cfg

import (
	"blockchain-transactions/internal/models"
	"blockchain-transactions/pkg/cfg/blockchain"
	"blockchain-transactions/pkg/cfg/categories"
	"blockchain-transactions/pkg/cfg/dictionaries"
	"blockchain-transactions/pkg/cfg/messages"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvDictionaries dictionaries.PortsServerDictionaries
	SrvMessage      messages.PortsServerMessages
	SrvCategories   categories.PortsServerCategories
	SrvBlockchain   blockchain.PortsServerBlockchain
}

func NewServerCfg(db *sqlx.DB, user *models.User, txID string) *Server {

	repoDictionaries := dictionaries.FactoryStorage(db, user, txID)
	srvDictionaries := dictionaries.NewDictionariesService(repoDictionaries, user, txID)

	repoCategories := categories.FactoryStorage(db, user, txID)
	srvCategories := categories.NewCategoriesService(repoCategories, user, txID)

	repoMessage := messages.FactoryStorage(db, user, txID)
	srvMessage := messages.NewMessagesService(repoMessage, user, txID)

	repoBlockchain := blockchain.FactoryStorage(db, user, txID)
	srvBlockchain := blockchain.NewBlockchainService(repoBlockchain, user, txID)

	return &Server{
		SrvDictionaries: srvDictionaries,
		SrvCategories:   srvCategories,
		SrvMessage:      srvMessage,
		SrvBlockchain:   srvBlockchain,
	}
}

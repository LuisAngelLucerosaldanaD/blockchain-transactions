package bc

import (
	"blockchain-transactions/internal/models"
	"blockchain-transactions/pkg/bc/block_fee"
	"blockchain-transactions/pkg/bc/files"
	"blockchain-transactions/pkg/bc/transaction"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	SrvTransactions transaction.PortsServerTransaction
	SrvFiles        files.PortsServerFile
	SrvBlockFee     block_fee.PortsServerBlockFee
}

func NewServerBc(db *sqlx.DB, user *models.User, txID string) *Server {
	repoTransactions := transaction.FactoryStorage(db, user, txID)
	srvTransactions := transaction.NewTransactionService(repoTransactions, user, txID)

	repoS3File := files.FactoryFileDocumentRepository(user, txID)
	srvFiles := files.NewFileService(repoS3File, user, txID)

	repoBlockFee := block_fee.FactoryStorage(db, user, txID)
	srvBlockFee := block_fee.NewBlockFeeService(repoBlockFee, user, txID)

	return &Server{
		SrvTransactions: srvTransactions,
		SrvFiles:        srvFiles,
		SrvBlockFee:     srvBlockFee,
	}
}

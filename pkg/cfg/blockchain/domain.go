package blockchain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Blockchain  Model struct Blockchain
type Blockchain struct {
	ID              string    `json:"id" db:"id" valid:"required,uuid"`
	FeeBlion        float32   `json:"fee_blion" db:"fee_blion" valid:"required"`
	FeeMiner        float32   `json:"fee_miner" db:"fee_miner" valid:"required"`
	FeeValidator    float32   `json:"fee_validator" db:"fee_validator" valid:"required"`
	FeeNode         float32   `json:"fee_node" db:"fee_node" valid:"required"`
	TtlBlock        int       `json:"ttl_block" db:"ttl_block" valid:"required"`
	MaxTransactions int       `json:"max_transactions" db:"max_transactions" valid:"required"`
	MaxValidators   int       `json:"max_validators" db:"max_validators" valid:"required"`
	MaxMiners       int       `json:"max_miners" db:"max_miners" valid:"required"`
	TicketsPrice    int       `json:"tickets_price" db:"tickets_price" valid:"required"`
	LotteryTtl      int       `json:"lottery_ttl" db:"lottery_ttl" valid:"required"`
	WalletMain      string    `json:"wallet_main" db:"wallet_main" valid:"required"`
	DeletedAt       time.Time `json:"deleted_at" db:"deleted_at" valid:"required"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

func NewBlockchain(id string, feeBlion float32, feeMiner float32, feeValidator float32, feeNode float32, ttlBlock int, maxTransactions int, deletedAt time.Time, maxMiners int, maxValidators int, ticketPrice int, lotteryTtl int, walletMain string) *Blockchain {
	return &Blockchain{
		ID:              id,
		FeeBlion:        feeBlion,
		FeeMiner:        feeMiner,
		FeeValidator:    feeValidator,
		FeeNode:         feeNode,
		TtlBlock:        ttlBlock,
		MaxTransactions: maxTransactions,
		DeletedAt:       deletedAt,
		MaxMiners:       maxMiners,
		MaxValidators:   maxValidators,
		TicketsPrice:    ticketPrice,
		LotteryTtl:      lotteryTtl,
		WalletMain:      walletMain,
	}
}

func (m *Blockchain) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}

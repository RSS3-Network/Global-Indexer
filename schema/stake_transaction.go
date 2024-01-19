package schema

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type StakeTransactionType string

const (
	StakeTransactionTypeDeposit  StakeTransactionType = "deposit"
	StakeTransactionTypeWithdraw StakeTransactionType = "withdraw"
	StakeTransactionTypeStake    StakeTransactionType = "stake"
	StakeTransactionTypeUnstake  StakeTransactionType = "unstake"
)

type StakeTransactionImporter interface {
	Import(stakeTransaction StakeTransaction) error
}

type StakeTransactionExporter interface {
	Export() (*StakeTransaction, error)
}

type StakeTransactionTransformer interface {
	StakeTransactionImporter
	StakeTransactionExporter
}

type StakeTransaction struct {
	ID    common.Hash          `json:"id"`
	Type  StakeTransactionType `json:"type"`
	User  common.Address       `json:"sender"`
	Node  common.Address       `json:"receiver"`
	Value *big.Int             `json:"value"`
	Chips []*big.Int           `json:"chips"`
}

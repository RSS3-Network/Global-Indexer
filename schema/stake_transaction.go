package schema

import (
	"math/big"
	"time"

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
	ID               common.Hash          `json:"id"`
	Type             StakeTransactionType `json:"type"`
	User             common.Address       `json:"sender"`
	Node             common.Address       `json:"receiver"`
	Value            *big.Int             `json:"value"`
	Chips            []*big.Int           `json:"chips"`
	BlockTimestamp   time.Time            `json:"blockTimestamp"`
	BlockNumber      uint64               `json:"blockNumber"`
	TransactionIndex uint                 `json:"transactionIndex"`
}

type StakeTransactionQuery struct {
	ID      *common.Hash
	User    *common.Address
	Node    *common.Address
	Address *common.Address
	Type    *StakeTransactionType
}

type StakeTransactionsQuery struct {
	Cursor  *common.Hash
	IDs     []common.Hash
	User    *common.Address
	Node    *common.Address
	Address *common.Address
	Type    *StakeTransactionType
	Pending *bool
}

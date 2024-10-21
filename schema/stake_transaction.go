package schema

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type StakeTransactionType string

const (
	StakeTransactionTypeDeposit    StakeTransactionType = "deposit"
	StakeTransactionTypeWithdraw   StakeTransactionType = "withdraw"
	StakeTransactionTypeStake      StakeTransactionType = "stake"
	StakeTransactionTypeUnstake    StakeTransactionType = "unstake"
	StakeTransactionTypeMergeChips StakeTransactionType = "merge_chips"
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
	ID               common.Hash
	Type             StakeTransactionType
	User             common.Address
	Node             common.Address
	Value            *big.Int
	ChipIDs          []*big.Int
	BlockTimestamp   time.Time
	BlockNumber      uint64
	TransactionIndex uint
	Finalized        bool
}

type StakeTransactionQuery struct {
	ID      *common.Hash
	User    *common.Address
	Node    *common.Address
	Address *common.Address
	Type    *StakeTransactionType
}

type StakeTransactionsQuery struct {
	Cursor              *common.Hash
	IDs                 []common.Hash
	User                *common.Address
	Node                *common.Address
	Address             *common.Address
	Type                *StakeTransactionType
	AfterBlockTimestamp *time.Time
	BlockNumber         *uint64
	Pending             *bool
	Limit               int
	Order               string
	Finalized           *bool
}

type StakeRecentCount struct {
	StakerCount uint64
	StakeValue  decimal.Decimal
}

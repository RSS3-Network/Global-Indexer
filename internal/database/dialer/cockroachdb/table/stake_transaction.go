package table

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                        = (*StakeTransaction)(nil)
	_ schema.StakeTransactionTransformer = (*StakeTransaction)(nil)
)

type StakeTransaction struct {
	ID               string          `gorm:"column:id;primaryKey"`
	Type             string          `gorm:"column:type;primaryKey"`
	User             string          `gorm:"column:user"`
	Node             string          `gorm:"column:node"`
	Value            decimal.Decimal `gorm:"column:value"`
	ChipIDs          pq.Int64Array   `gorm:"column:chips;type:bigint[]"`
	BlockTimestamp   time.Time       `gorm:"column:block_timestamp"`
	BlockNumber      uint64          `gorm:"column:block_number"`
	TransactionIndex uint            `gorm:"column:transaction_index"`
	Finalized        bool            `gorm:"column:finalized"`
}

func (s *StakeTransaction) TableName() string {
	return "stake.transactions"
}

func (s *StakeTransaction) Export() (*schema.StakeTransaction, error) {
	var stakeTransaction = schema.StakeTransaction{
		ID:    common.HexToHash(s.ID),
		Type:  schema.StakeTransactionType(s.Type),
		User:  common.HexToAddress(s.User),
		Node:  common.HexToAddress(s.Node),
		Value: s.Value.BigInt(),
		ChipIDs: lo.Map(s.ChipIDs, func(value int64, _ int) *big.Int {
			return new(big.Int).SetInt64(value)
		}),
		BlockTimestamp:   s.BlockTimestamp,
		BlockNumber:      s.BlockNumber,
		TransactionIndex: s.TransactionIndex,
		Finalized:        s.Finalized,
	}

	return &stakeTransaction, nil
}

func (s *StakeTransaction) Import(stakeTransaction schema.StakeTransaction) error {
	s.ID = stakeTransaction.ID.String()
	s.Type = string(stakeTransaction.Type)
	s.User = stakeTransaction.User.String()
	s.Node = stakeTransaction.Node.String()
	s.Value = decimal.NewFromBigInt(stakeTransaction.Value, 0)
	s.ChipIDs = lo.Map(stakeTransaction.ChipIDs, func(value *big.Int, _ int) int64 {
		return value.Int64()
	})
	s.BlockTimestamp = stakeTransaction.BlockTimestamp
	s.BlockNumber = stakeTransaction.BlockNumber
	s.TransactionIndex = stakeTransaction.TransactionIndex
	s.Finalized = stakeTransaction.Finalized

	return nil
}

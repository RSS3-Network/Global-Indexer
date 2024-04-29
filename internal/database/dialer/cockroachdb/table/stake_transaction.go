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
	ID               string          `gorm:"column:id;type:text;not null;primaryKey;"`
	Type             string          `gorm:"column:type;type:text;not null;primaryKey;"`
	User             string          `gorm:"column:user;type:text;not null;index:idx_transactions_user;index:idx_transactions_address,priority:1;"`
	Node             string          `gorm:"column:node;type:text;not null;index:idx_transactions_node;index:idx_transactions_address,priority:2;"`
	Value            decimal.Decimal `gorm:"column:value;type:decimal;not null;"`
	Chips            pq.Int64Array   `gorm:"column:chips;type:bigint[];not null;"`
	BlockTimestamp   time.Time       `gorm:"column:block_timestamp;type:timestamp with time zone;not null;index:idx_transactions_order,priority:1,sort:desc;"`
	BlockNumber      uint64          `gorm:"column:block_number;type:bigint;not null;index:idx_transactions_order,priority:2,sort:desc;"`
	TransactionIndex uint            `gorm:"column:transaction_index;type:bigint;not null;index:idx_transactions_order,priority:3,sort:desc;"`
}

func (s *StakeTransaction) TableName() string {
	return "stake_transactions"
}

func (s *StakeTransaction) Export() (*schema.StakeTransaction, error) {
	var stakeTransaction = schema.StakeTransaction{
		ID:    common.HexToHash(s.ID),
		Type:  schema.StakeTransactionType(s.Type),
		User:  common.HexToAddress(s.User),
		Node:  common.HexToAddress(s.Node),
		Value: s.Value.BigInt(),
		Chips: lo.Map(s.Chips, func(value int64, _ int) *big.Int {
			return new(big.Int).SetInt64(value)
		}),
		BlockTimestamp:   s.BlockTimestamp,
		BlockNumber:      s.BlockNumber,
		TransactionIndex: s.TransactionIndex,
	}

	return &stakeTransaction, nil
}

func (s *StakeTransaction) Import(stakeTransaction schema.StakeTransaction) error {
	s.ID = stakeTransaction.ID.String()
	s.Type = string(stakeTransaction.Type)
	s.User = stakeTransaction.User.String()
	s.Node = stakeTransaction.Node.String()
	s.Value = decimal.NewFromBigInt(stakeTransaction.Value, 0)
	s.Chips = lo.Map(stakeTransaction.Chips, func(value *big.Int, _ int) int64 {
		return value.Int64()
	})
	s.BlockTimestamp = stakeTransaction.BlockTimestamp
	s.BlockNumber = stakeTransaction.BlockNumber
	s.TransactionIndex = stakeTransaction.TransactionIndex

	return nil
}

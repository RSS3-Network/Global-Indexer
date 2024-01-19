package table

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                        = (*StakeTransaction)(nil)
	_ schema.StakeTransactionTransformer = (*StakeTransaction)(nil)
)

type StakeTransaction struct {
	ID    string          `gorm:"column:id"`
	Type  string          `gorm:"column:type"`
	User  string          `gorm:"column:user"`
	Node  string          `gorm:"column:node"`
	Value decimal.Decimal `gorm:"column:value"`
	Chips pq.Int64Array   `gorm:"column:chips;type:bigint[]"`
}

func (s *StakeTransaction) TableName() string {
	return "stake.transactions"
}

func (s *StakeTransaction) Export() (*schema.StakeTransaction, error) {
	stakeTransaction := schema.StakeTransaction{
		ID:    common.HexToHash(s.ID),
		Type:  schema.StakeTransactionType(s.Type),
		User:  common.HexToAddress(s.User),
		Node:  common.HexToAddress(s.Node),
		Value: s.Value.BigInt(),
		Chips: lo.Map(s.Chips, func(value int64, _ int) *big.Int {
			return new(big.Int).SetInt64(value)
		}),
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

	return nil
}

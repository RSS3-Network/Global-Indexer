package table

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/shopspring/decimal"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                 = (*StakeChip)(nil)
	_ schema.StakeChipTransformer = (*StakeChip)(nil)
)

type StakeChip struct {
	ID    decimal.Decimal `gorm:"column:id"`
	Owner string          `gorm:"column:owner"`
	Node  string          `gorm:"column:node"`
}

func (s *StakeChip) TableName() string {
	return "stake.chips"
}

func (s *StakeChip) Import(stakeChip schema.StakeChip) error {
	s.ID = decimal.NewFromBigInt(stakeChip.ID, 0)
	s.Owner = stakeChip.Owner.String()
	s.Node = stakeChip.Node.String()

	return nil
}

func (s *StakeChip) Export() (*schema.StakeChip, error) {
	stakeChip := schema.StakeChip{
		ID:    s.ID.BigInt(),
		Owner: common.HexToAddress(s.Owner),
		Node:  common.HexToAddress(s.Node),
	}

	return &stakeChip, nil
}

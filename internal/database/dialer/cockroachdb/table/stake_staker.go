package table

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/shopspring/decimal"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                   = (*StakeStaker)(nil)
	_ schema.StakeStakerTransformer = (*StakeStaker)(nil)
)

type StakeStaker struct {
	User  string          `gorm:"column:user"`
	Node  string          `gorm:"column:node"`
	Value decimal.Decimal `gorm:"column:value"`
}

func (s *StakeStaker) TableName() string {
	return "stake.stakers"
}

func (s *StakeStaker) Import(stakeStaker schema.StakeStaker) error {
	s.User = stakeStaker.User.String()
	s.Node = stakeStaker.Node.String()
	s.Value = stakeStaker.Value

	return nil
}

func (s *StakeStaker) Export() (*schema.StakeStaker, error) {
	stakeStaker := schema.StakeStaker{
		User:  common.HexToAddress(s.User),
		Node:  common.HexToAddress(s.Node),
		Value: s.Value,
	}

	return &stakeStaker, nil
}

package table

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                    = (*StakeStaking)(nil)
	_ schema.StakeStakingTransformer = (*StakeStaking)(nil)
)

type StakeStaking struct {
	Staker string          `gorm:"column:staker"`
	Node   string          `gorm:"column:node"`
	Count  uint32          `gorm:"column:count"`
	Value  decimal.Decimal `gorm:"column:value"`
}

func (s *StakeStaking) TableName() string {
	return "stake.stakings"
}

func (s *StakeStaking) Import(stakeStaking schema.StakeStaking) error {
	s.Staker = stakeStaking.Staker.String()
	s.Node = stakeStaking.Node.String()
	s.Value = stakeStaking.Value

	return nil
}

func (s *StakeStaking) Export() (*schema.StakeStaking, error) {
	stakeStaker := schema.StakeStaking{
		Staker: common.HexToAddress(s.Staker),
		Node:   common.HexToAddress(s.Node),
		Value:  s.Value,
		Chips: schema.StakeStakingChips{
			Total: uint64(s.Count),
		},
	}

	return &stakeStaker, nil
}

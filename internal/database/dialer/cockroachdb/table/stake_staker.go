package table

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                   = (*StakeStaker)(nil)
	_ schema.StakeStakerTransformer = (*StakeStaker)(nil)
)

type StakeStaker struct {
	Address    string          `gorm:"column:address"`
	Nodes      uint64          `gorm:"column:nodes"`
	ChipNumber uint64          `gorm:"column:chip_number"`
	ChipValue  decimal.Decimal `gorm:"column:chip_value"`
}

func (s *StakeStaker) TableName() string {
	return "stake.stakers"
}

func (s *StakeStaker) Export() (*schema.StakeStaker, error) {
	stakeStaker := schema.StakeStaker{
		Address:           common.HexToAddress(s.Address),
		TotalStakedNodes:  s.Nodes,
		TotalOwnedChips:   s.ChipNumber,
		TotalStakedTokens: s.ChipValue,
	}

	return &stakeStaker, nil
}

package table

import (
	"encoding/json"
	"time"

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
	ID             decimal.Decimal `gorm:"column:id"`
	Owner          string          `gorm:"column:owner"`
	Node           string          `gorm:"column:node"`
	Value          decimal.Decimal `gorm:"column:value"`
	Metadata       json.RawMessage `gorm:"column:metadata"`
	BlockNumber    decimal.Decimal `gorm:"column:block_number"`
	BlockTimestamp time.Time       `gorm:"column:block_timestamp"`
}

func (s *StakeChip) TableName() string {
	return "stake.chips"
}

func (s *StakeChip) Import(stakeChip schema.StakeChip) error {
	s.ID = decimal.NewFromBigInt(stakeChip.ID, 0)
	s.Owner = stakeChip.Owner.String()
	s.Node = stakeChip.Node.String()
	s.Value = stakeChip.Value
	s.Metadata = stakeChip.Metadata
	s.BlockNumber = decimal.NewFromBigInt(stakeChip.BlockNumber, 0)
	s.BlockTimestamp = time.Unix(int64(stakeChip.BlockTimestamp), 0)

	return nil
}

func (s *StakeChip) Export() (*schema.StakeChip, error) {
	stakeChip := schema.StakeChip{
		ID:             s.ID.BigInt(),
		Owner:          common.HexToAddress(s.Owner),
		Node:           common.HexToAddress(s.Node),
		Value:          s.Value,
		Metadata:       s.Metadata,
		BlockNumber:    s.BlockNumber.BigInt(),
		BlockTimestamp: uint64(s.BlockTimestamp.Unix()),
	}

	return &stakeChip, nil
}

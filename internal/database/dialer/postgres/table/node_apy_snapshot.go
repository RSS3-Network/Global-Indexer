package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type NodeAPYSnapshot struct {
	ID          uint64          `gorm:"column:id"`
	Date        time.Time       `gorm:"column:date"`
	EpochID     uint64          `gorm:"column:epoch_id"`
	NodeAddress common.Address  `gorm:"column:node_address"`
	APY         decimal.Decimal `gorm:"column:apy"`
	CreatedAt   time.Time       `gorm:"column:created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at"`
}

func (s *NodeAPYSnapshot) TableName() string {
	return "node.apy_snapshots"
}

func (s *NodeAPYSnapshot) Import(nodeAPYSnapshot *schema.NodeAPYSnapshot) error {
	s.Date = nodeAPYSnapshot.Date
	s.EpochID = nodeAPYSnapshot.EpochID
	s.NodeAddress = nodeAPYSnapshot.NodeAddress
	s.APY = nodeAPYSnapshot.APY

	return nil
}

func (s *NodeAPYSnapshot) Export() (*schema.NodeAPYSnapshot, error) {
	return &schema.NodeAPYSnapshot{
		ID:          s.ID,
		Date:        s.Date,
		EpochID:     s.EpochID,
		NodeAddress: s.NodeAddress,
		APY:         s.APY,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}, nil
}

type NodeAPYSnapshots []NodeAPYSnapshot

func (s *NodeAPYSnapshots) Import(snapshots []*schema.NodeAPYSnapshot) error {
	for _, snapshot := range snapshots {
		var imported NodeAPYSnapshot

		if err := imported.Import(snapshot); err != nil {
			return err
		}

		*s = append(*s, imported)
	}

	return nil
}

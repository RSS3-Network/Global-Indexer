package table

import (
	"time"

	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                    = (*NodeSnapshot)(nil)
	_ schema.NodeSnapshotTransformer = (*NodeSnapshot)(nil)
)

type NodeSnapshot struct {
	EpochID        uint64    `gorm:"column:epoch_id"`
	Count          uint64    `gorm:"column:count"`
	BlockHash      string    `gorm:"column:block_hash"`
	BlockNumber    uint64    `gorm:"column:block_number"`
	BlockTimestamp time.Time `gorm:"column:block_timestamp"`
}

func (s *NodeSnapshot) TableName() string {
	return "node.snapshots"
}

func (s *NodeSnapshot) Import(stakeSnapshot schema.NodeSnapshot) error {
	s.EpochID = stakeSnapshot.EpochID
	s.Count = uint64(stakeSnapshot.Count)
	s.BlockHash = stakeSnapshot.BlockHash
	s.BlockNumber = stakeSnapshot.BlockNumber
	s.BlockTimestamp = stakeSnapshot.BlockTimestamp

	return nil
}

func (s *NodeSnapshot) Export() (*schema.NodeSnapshot, error) {
	stakeSnapshot := schema.NodeSnapshot{
		EpochID:        s.EpochID,
		Count:          int64(s.Count),
		BlockHash:      s.BlockHash,
		BlockNumber:    s.BlockNumber,
		BlockTimestamp: s.BlockTimestamp,
	}

	return &stakeSnapshot, nil
}

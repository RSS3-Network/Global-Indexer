package table

import (
	"time"

	"github.com/rss3-network/global-indexer/schema"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                    = (*NodeSnapshot)(nil)
	_ schema.NodeSnapshotTransformer = (*NodeSnapshot)(nil)
)

type NodeSnapshot struct {
	Date  time.Time `gorm:"column:date"`
	Count uint64    `gorm:"column:count"`
}

func (s *NodeSnapshot) TableName() string {
	return "node_count_snapshots"
}

func (s *NodeSnapshot) Import(stakeSnapshot schema.NodeSnapshot) error {
	s.Date = stakeSnapshot.Date
	s.Count = uint64(stakeSnapshot.Count)

	return nil
}

func (s *NodeSnapshot) Export() (*schema.NodeSnapshot, error) {
	stakeSnapshot := schema.NodeSnapshot{
		Date:  s.Date,
		Count: int64(s.Count),
	}

	return &stakeSnapshot, nil
}

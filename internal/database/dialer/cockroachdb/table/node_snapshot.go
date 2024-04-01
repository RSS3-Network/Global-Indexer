package table

import (
	"time"

	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"gorm.io/gorm"
	gormSchema "gorm.io/gorm/schema"
)

var (
	_ gormSchema.Tabler              = (*NodeSnapshot)(nil)
	_ schema.NodeSnapshotTransformer = (*NodeSnapshot)(nil)
)

type NodeSnapshot struct {
	gorm.Model
	Date  time.Time `gorm:"column:date"`
	Count uint64    `gorm:"column:count"`
}

func (s *NodeSnapshot) TableName() string {
	return "count_snapshots"
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

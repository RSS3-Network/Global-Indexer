package table

import (
	"time"

	"github.com/rss3-network/global-indexer/schema"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                     = (*StakerCountSnapshot)(nil)
	_ schema.StakeSnapshotTransformer = (*StakerCountSnapshot)(nil)
)

type StakerCountSnapshot struct {
	Date  time.Time `gorm:"column:date"`
	Count uint64    `gorm:"column:count"`
}

func (s *StakerCountSnapshot) TableName() string {
	return "stake.count_snapshots"
}

func (s *StakerCountSnapshot) Import(stakeSnapshot schema.StakerCountSnapshot) error {
	s.Date = stakeSnapshot.Date
	s.Count = uint64(stakeSnapshot.Count)

	return nil
}

func (s *StakerCountSnapshot) Export() (*schema.StakerCountSnapshot, error) {
	stakeSnapshot := schema.StakerCountSnapshot{
		Date:  s.Date,
		Count: int64(s.Count),
	}

	return &stakeSnapshot, nil
}

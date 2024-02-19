package table

import (
	"time"

	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                     = (*StakeSnapshot)(nil)
	_ schema.StakeSnapshotTransformer = (*StakeSnapshot)(nil)
)

type StakeSnapshot struct {
	Date  time.Time `gorm:"column:date"`
	Count uint64    `gorm:"column:count"`
}

func (s *StakeSnapshot) TableName() string {
	return "stake.snapshots"
}

func (s *StakeSnapshot) Import(stakeSnapshot schema.StakeSnapshot) error {
	s.Date = stakeSnapshot.Date
	s.Count = uint64(stakeSnapshot.Count)

	return nil
}

func (s *StakeSnapshot) Export() (*schema.StakeSnapshot, error) {
	stakeSnapshot := schema.StakeSnapshot{
		Date:  s.Date,
		Count: int64(s.Count),
	}

	return &stakeSnapshot, nil
}

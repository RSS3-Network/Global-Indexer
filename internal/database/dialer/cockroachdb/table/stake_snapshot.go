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
	EpochID        uint64    `gorm:"column:epoch_id"`
	Count          uint64    `gorm:"column:count"`
	BlockHash      string    `gorm:"column:block_hash"`
	BlockNumber    uint64    `gorm:"column:block_number"`
	BlockTimestamp time.Time `gorm:"column:block_timestamp"`
}

func (s *StakeSnapshot) TableName() string {
	return "stake.snapshots"
}

func (s *StakeSnapshot) Import(stakeSnapshot schema.StakeSnapshot) error {
	s.EpochID = stakeSnapshot.EpochID
	s.Count = uint64(stakeSnapshot.Count)
	s.BlockHash = stakeSnapshot.BlockHash
	s.BlockNumber = stakeSnapshot.BlockNumber
	s.BlockTimestamp = stakeSnapshot.BlockTimestamp

	return nil
}

func (s *StakeSnapshot) Export() (*schema.StakeSnapshot, error) {
	stakeSnapshot := schema.StakeSnapshot{
		EpochID:        s.EpochID,
		Count:          int64(s.Count),
		BlockHash:      s.BlockHash,
		BlockNumber:    s.BlockNumber,
		BlockTimestamp: s.BlockTimestamp,
	}

	return &stakeSnapshot, nil
}

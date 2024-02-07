package schema

import "time"

type StakeSnapshotImporter interface {
	Import(stakeSnapshot StakeSnapshot) error
}

type StakeSnapshotExporter interface {
	Export() (*StakeSnapshot, error)
}

type StakeSnapshotTransformer interface {
	StakeSnapshotImporter
	StakeSnapshotExporter
}

type StakeSnapshot struct {
	EpochID        uint64    `gorm:"column:epoch_id"`
	Count          int64     `gorm:"column:count"`
	BlockHash      string    `gorm:"column:block_hash"`
	BlockNumber    uint64    `gorm:"column:block_number"`
	BlockTimestamp time.Time `gorm:"column:block_timestamp"`
}

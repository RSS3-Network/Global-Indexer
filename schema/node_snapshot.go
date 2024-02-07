package schema

import "time"

type NodeSnapshotImporter interface {
	Import(nodeSnapshot NodeSnapshot) error
}

type NodeSnapshotExporter interface {
	Export() (*NodeSnapshot, error)
}

type NodeSnapshotTransformer interface {
	NodeSnapshotImporter
	NodeSnapshotExporter
}

type NodeSnapshot struct {
	EpochID        uint64    `gorm:"column:epoch_id"`
	Count          int64     `gorm:"column:count"`
	BlockHash      string    `gorm:"column:block_hash"`
	BlockNumber    uint64    `gorm:"column:block_number"`
	BlockTimestamp time.Time `gorm:"column:block_timestamp"`
}

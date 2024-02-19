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
	Date  time.Time `json:"date"`
	Count int64     `json:"count"`
}

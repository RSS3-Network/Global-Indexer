package schema

import "time"

type StakerCountSnapshotImporter interface {
	Import(stakeSnapshot StakerCountSnapshot) error
}

type StakerCountSnapshotExporter interface {
	Export() (*StakerCountSnapshot, error)
}

type StakeSnapshotTransformer interface {
	StakerCountSnapshotImporter
	StakerCountSnapshotExporter
}

type StakerCountSnapshot struct {
	Date  time.Time `json:"date"`
	Count int64     `json:"count"`
}

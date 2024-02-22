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
	Date  time.Time `json:"date"`
	Count int64     `json:"count"`
}

package schema

import "github.com/ethereum/go-ethereum/common"

type CheckpointImporter interface {
	Import(checkpoint Checkpoint) error
}

type CheckpointExporter interface {
	Export() (*Checkpoint, error)
}

type CheckpointTransformer interface {
	CheckpointImporter
	CheckpointExporter
}

type Checkpoint struct {
	ChainID     uint64      `json:"network"`
	BlockNumber uint64      `json:"block_number"`
	BlockHash   common.Hash `json:"block_hash"`
}

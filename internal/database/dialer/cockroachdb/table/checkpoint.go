package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/global-indexer/schema"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                  = (*Checkpoint)(nil)
	_ schema.CheckpointTransformer = (*Checkpoint)(nil)
)

type Checkpoint struct {
	ChainID     uint64    `gorm:"column:chain_id"`
	BlockNumber uint64    `gorm:"column:block_number"`
	BlockHash   string    `gorm:"column:block_hash"`
	CreatedAt   time.Time `gorm:"column:created_at;"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (c *Checkpoint) TableName() string {
	return "checkpoints"
}

func (c *Checkpoint) Import(checkpoint schema.Checkpoint) error {
	c.ChainID = checkpoint.ChainID
	c.BlockNumber = checkpoint.BlockNumber
	c.BlockHash = checkpoint.BlockHash.String()

	return nil
}

func (c *Checkpoint) Export() (*schema.Checkpoint, error) {
	checkpoint := schema.Checkpoint{
		ChainID:     c.ChainID,
		BlockNumber: c.BlockNumber,
		BlockHash:   common.HexToHash(c.BlockHash),
	}

	return &checkpoint, nil
}

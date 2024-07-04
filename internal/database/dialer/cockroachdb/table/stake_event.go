package table

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                  = (*StakeEvent)(nil)
	_ schema.StakeEventTransformer = (*StakeEvent)(nil)
)

type StakeEvent struct {
	ID                string    `gorm:"column:id"`
	Type              string    `gorm:"column:type"`
	TransactionHash   string    `gorm:"column:transaction_hash;primaryKey"`
	TransactionIndex  uint      `gorm:"column:transaction_index"`
	TransactionStatus uint64    `gorm:"column:transaction_status"`
	BlockHash         string    `gorm:"column:block_hash;primaryKey"`
	BlockNumber       uint64    `gorm:"column:block_number"`
	BlockTimestamp    time.Time `gorm:"column:block_timestamp"`
	Finalized         bool      `gorm:"column:finalized"`
}

func (b *StakeEvent) TableName() string {
	return "stake.events"
}

func (b *StakeEvent) Import(stakeEvent schema.StakeEvent) error {
	b.ID = stakeEvent.ID.String()
	b.Type = string(stakeEvent.Type)
	b.TransactionHash = stakeEvent.TransactionHash.String()
	b.TransactionIndex = stakeEvent.TransactionIndex
	b.TransactionStatus = stakeEvent.TransactionStatus
	b.BlockHash = stakeEvent.BlockHash.String()
	b.BlockNumber = stakeEvent.BlockNumber.Uint64()
	b.BlockTimestamp = stakeEvent.BlockTimestamp
	b.Finalized = stakeEvent.Finalized

	return nil
}

func (b *StakeEvent) Export() (*schema.StakeEvent, error) {
	stakeEvent := schema.StakeEvent{
		ID:                common.HexToHash(b.ID),
		Type:              schema.StakeEventType(b.Type),
		TransactionHash:   common.HexToHash(b.TransactionHash),
		TransactionIndex:  b.TransactionIndex,
		TransactionStatus: b.TransactionStatus,
		BlockHash:         common.HexToHash(b.BlockHash),
		BlockNumber:       new(big.Int).SetUint64(b.BlockNumber),
		BlockTimestamp:    b.BlockTimestamp,
		Finalized:         b.Finalized,
	}

	return &stakeEvent, nil
}

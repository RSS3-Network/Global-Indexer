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
	ID                string    `gorm:"column:id;type:text;not null;index:idx_id;"`
	Type              string    `gorm:"column:type;type:text;not null;"`
	TransactionHash   string    `gorm:"column:transaction_hash;type:text;not null;primaryKey;"`
	TransactionIndex  uint      `gorm:"column:transaction_index;type:bigint;not null;"`
	TransactionStatus uint64    `gorm:"column:transaction_status;type:bigint;not null;"`
	BlockHash         string    `gorm:"column:block_hash;type:text;not null;primaryKey"`
	BlockNumber       uint64    `gorm:"column:block_number;type:bigint;not null;"`
	BlockTimestamp    time.Time `gorm:"column:block_timestamp;type:timestamp with time zone;not null;"`
}

func (b *StakeEvent) TableName() string {
	return "stake_events"
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
	}

	return &stakeEvent, nil
}

package table

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                   = (*BridgeEvent)(nil)
	_ schema.BridgeEventTransformer = (*BridgeEvent)(nil)
)

type BridgeEvent struct {
	ID                string    `gorm:"column:id;type:text;not null;index:idx_id;"`
	Type              string    `gorm:"column:type;type:text;not null;"`
	TransactionHash   string    `gorm:"column:transaction_hash;type:text;primaryKey;"`
	TransactionIndex  uint      `gorm:"column:transaction_index;type:bigint;not null;"`
	TransactionStatus uint64    `gorm:"column:transaction_status;type:bigint;not null;"`
	ChainID           uint64    `gorm:"column:chain_id;type:bigint;not null;"`
	BlockHash         string    `gorm:"column:block_hash;type:text;primaryKey;"`
	BlockNumber       uint64    `gorm:"column:block_number;type:bigint;not null"`
	BlockTimestamp    time.Time `gorm:"column:block_timestamp;type:timestamp with time zone;not null;"`
}

func (b *BridgeEvent) TableName() string {
	return "bridge_events"
}

func (b *BridgeEvent) Import(bridgeEvent schema.BridgeEvent) error {
	b.ID = bridgeEvent.ID.String()
	b.Type = string(bridgeEvent.Type)
	b.TransactionHash = bridgeEvent.TransactionHash.String()
	b.TransactionIndex = bridgeEvent.TransactionIndex
	b.TransactionStatus = bridgeEvent.TransactionStatus
	b.BlockHash = bridgeEvent.BlockHash.String()
	b.BlockNumber = bridgeEvent.BlockNumber.Uint64()
	b.BlockTimestamp = bridgeEvent.BlockTimestamp

	return nil
}

func (b *BridgeEvent) Export() (*schema.BridgeEvent, error) {
	bridgeEvent := schema.BridgeEvent{
		ID:                common.HexToHash(b.ID),
		Type:              schema.BridgeEventType(b.Type),
		TransactionHash:   common.HexToHash(b.TransactionHash),
		TransactionIndex:  b.TransactionIndex,
		TransactionStatus: b.TransactionStatus,
		ChainID:           b.ChainID,
		BlockHash:         common.HexToHash(b.BlockHash),
		BlockNumber:       new(big.Int).SetUint64(b.BlockNumber),
		BlockTimestamp:    b.BlockTimestamp,
	}

	return &bridgeEvent, nil
}

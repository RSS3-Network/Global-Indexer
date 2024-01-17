package table

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/global-indexer/schema"
	gorm "gorm.io/gorm/schema"
)

var (
	_ gorm.Tabler                   = (*BridgeEvent)(nil)
	_ schema.BridgeEventTransformer = (*BridgeEvent)(nil)
)

type BridgeEvent struct {
	ID                string    `gorm:"column:id"`
	Type              string    `gorm:"column:type"`
	TransactionHash   string    `gorm:"column:transaction_hash"`
	TransactionIndex  uint      `gorm:"column:transaction_index"`
	TransactionStatus uint64    `gorm:"column:transaction_status"`
	BlockHash         string    `gorm:"column:block_hash"`
	BlockNumber       uint64    `gorm:"column:block_number"`
	BlockTimestamp    time.Time `gorm:"column:block_timestamp"`
}

func (b *BridgeEvent) TableName() string {
	return "bridge.events"
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
		BlockHash:         common.HexToHash(b.BlockHash),
		BlockNumber:       new(big.Int).SetUint64(b.BlockNumber),
		BlockTimestamp:    b.BlockTimestamp,
	}

	return &bridgeEvent, nil
}

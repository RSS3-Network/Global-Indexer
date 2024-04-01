package table

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"gorm.io/gorm"
	gormSchema "gorm.io/gorm/schema"
)

var (
	_ gormSchema.Tabler            = (*StakeEvent)(nil)
	_ schema.StakeEventTransformer = (*StakeEvent)(nil)
)

type StakeEvent struct {
	gorm.Model
	ID                string    `gorm:"column:id"`
	Type              string    `gorm:"column:type"`
	TransactionHash   string    `gorm:"column:transaction_hash"`
	TransactionIndex  uint      `gorm:"column:transaction_index"`
	TransactionStatus uint64    `gorm:"column:transaction_status"`
	BlockHash         string    `gorm:"column:block_hash"`
	BlockNumber       uint64    `gorm:"column:block_number"`
	BlockTimestamp    time.Time `gorm:"column:block_timestamp"`
}

func (b *StakeEvent) TableName() string {
	return "events"
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

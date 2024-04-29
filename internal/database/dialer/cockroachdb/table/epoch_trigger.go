package table

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
)

type EpochTrigger struct {
	TransactionHash string          `gorm:"column:transaction_hash;type:text;not null;primaryKey;"`
	EpochID         uint64          `gorm:"column:epoch_id;type:bigint;not null;index:idx_epoch_id;"`
	Data            json.RawMessage `gorm:"column:data;type:jsonb;not null;"`
	CreatedAt       time.Time       `gorm:"column:created_at;type:timestamp with time zone;not null;default:now();index:idx_created_at;"`
	UpdatedAt       time.Time       `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now()"`
}

func (e *EpochTrigger) TableName() string {
	return "epoch_trigger"
}

func (e *EpochTrigger) Import(epochTrigger *schema.EpochTrigger) (err error) {
	e.TransactionHash = epochTrigger.TransactionHash.String()
	e.EpochID = epochTrigger.EpochID
	e.CreatedAt = epochTrigger.CreatedAt
	e.UpdatedAt = epochTrigger.UpdatedAt

	e.Data, err = json.Marshal(epochTrigger.Data)

	return err
}

func (e *EpochTrigger) Export() (*schema.EpochTrigger, error) {
	var data schema.SettlementData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, err
	}

	return &schema.EpochTrigger{
		TransactionHash: common.HexToHash(e.TransactionHash),
		EpochID:         e.EpochID,
		Data:            data,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}, nil
}

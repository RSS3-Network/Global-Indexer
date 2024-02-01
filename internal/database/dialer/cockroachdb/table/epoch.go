package table

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type Epoch struct {
	ID                    uint64          `gorm:"column:id;primaryKey"`
	StartTimestamp        time.Time       `gorm:"column:start_timestamp"`
	EndTimestamp          time.Time       `gorm:"column:end_timestamp"`
	TransactionHash       string          `gorm:"column:transaction_hash"`
	BlockNumber           uint64          `gorm:"column:block_number"`
	TotalOperationRewards decimal.Decimal `gorm:"column:total_operation_rewards"`
	TotalStakingRewards   decimal.Decimal `gorm:"column:total_staking_rewards"`
	TotalRewardItems      int             `gorm:"column:total_reward_items"`
	Success               bool            `json:"success"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (e *Epoch) TableName() string {
	return "epoch"
}

func (e *Epoch) Import(epoch *schema.Epoch) error {
	e.ID = epoch.ID.Uint64()
	e.StartTimestamp = time.Unix(epoch.StartTimestamp, 0)
	e.EndTimestamp = time.Unix(epoch.EndTimestamp, 0)
	e.TransactionHash = epoch.TransactionHash.String()
	e.BlockNumber = epoch.BlockNumber.Uint64()
	e.TotalOperationRewards = lo.Must(decimal.NewFromString(epoch.TotalOperationRewards))
	e.TotalStakingRewards = lo.Must(decimal.NewFromString(epoch.TotalStakingRewards))
	e.TotalRewardItems = epoch.TotalRewardItems
	e.Success = epoch.Success

	return nil
}

func (e *Epoch) Export() (*schema.Epoch, error) {
	epoch := schema.Epoch{
		ID:                    new(big.Int).SetUint64(e.ID),
		StartTimestamp:        e.StartTimestamp.Unix(),
		EndTimestamp:          e.EndTimestamp.Unix(),
		TransactionHash:       common.HexToHash(e.TransactionHash),
		BlockNumber:           new(big.Int).SetUint64(e.BlockNumber),
		TotalOperationRewards: e.TotalOperationRewards.String(),
		TotalStakingRewards:   e.TotalStakingRewards.String(),
		TotalRewardItems:      e.TotalRewardItems,
		Success:               e.Success,
	}

	return &epoch, nil
}

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
	TransactionIndex      uint            `gorm:"column:transaction_index"`
	BlockNumber           uint64          `gorm:"column:block_number"`
	Success               bool            `gorm:"column:success"`
	TotalOperationRewards decimal.Decimal `gorm:"column:total_operation_rewards"`
	TotalStakingRewards   decimal.Decimal `gorm:"column:total_staking_rewards"`
	TotalRewardItems      int             `gorm:"column:total_reward_items"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (e *Epoch) TableName() string {
	return "epoch"
}

func (e *Epoch) Import(epoch *schema.Epoch) error {
	e.ID = epoch.ID
	e.StartTimestamp = time.Unix(epoch.StartTimestamp, 0)
	e.EndTimestamp = time.Unix(epoch.EndTimestamp, 0)
	e.TransactionHash = epoch.TransactionHash.String()
	e.TransactionIndex = epoch.TransactionIndex
	e.BlockNumber = epoch.BlockNumber.Uint64()
	e.Success = epoch.Success
	e.TotalOperationRewards = lo.Must(decimal.NewFromString(epoch.TotalOperationRewards))
	e.TotalStakingRewards = lo.Must(decimal.NewFromString(epoch.TotalStakingRewards))
	e.TotalRewardItems = epoch.TotalRewardItems

	return nil
}

func (e *Epoch) Export(epochItems []*schema.EpochItem) (*schema.Epoch, error) {
	epoch := schema.Epoch{
		ID:                    e.ID,
		StartTimestamp:        e.StartTimestamp.Unix(),
		EndTimestamp:          e.EndTimestamp.Unix(),
		TransactionHash:       common.HexToHash(e.TransactionHash),
		TransactionIndex:      e.TransactionIndex,
		BlockNumber:           new(big.Int).SetUint64(e.BlockNumber),
		TotalOperationRewards: e.TotalOperationRewards.String(),
		TotalStakingRewards:   e.TotalStakingRewards.String(),
		TotalRewardItems:      e.TotalRewardItems,
		Success:               e.Success,
		RewardItems:           epochItems,
	}

	return &epoch, nil
}

type Epochs []*Epoch

func (e *Epochs) Export() ([]*schema.Epoch, error) {
	epochs := make([]*schema.Epoch, 0, len(*e))

	for _, epoch := range *e {
		epoch, err := epoch.Export(nil)
		if err != nil {
			return nil, err
		}

		epochs = append(epochs, epoch)
	}

	return epochs, nil
}

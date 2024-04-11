package table

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type Epoch struct {
	ID                    uint64          `gorm:"column:id;primaryKey"`
	StartTimestamp        time.Time       `gorm:"column:start_timestamp"`
	EndTimestamp          time.Time       `gorm:"column:end_timestamp"`
	TransactionHash       string          `gorm:"column:transaction_hash"`
	TransactionIndex      uint            `gorm:"column:transaction_index"`
	BlockHash             string          `gorm:"column:block_hash"`
	BlockNumber           uint64          `gorm:"column:block_number"`
	BlockTimestamp        time.Time       `gorm:"column:block_timestamp"`
	TotalOperationRewards decimal.Decimal `gorm:"column:total_operation_rewards"`
	TotalStakingRewards   decimal.Decimal `gorm:"column:total_staking_rewards"`
	TotalRewardItems      int             `gorm:"column:total_reward_items"`
	TotalRequestCounts    decimal.Decimal `gorm:"column:total_request_counts"`

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
	e.BlockHash = epoch.BlockHash.String()
	e.BlockNumber = epoch.BlockNumber.Uint64()
	e.BlockTimestamp = time.Unix(epoch.BlockTimestamp, 0)
	e.TotalOperationRewards = epoch.TotalOperationRewards
	e.TotalStakingRewards = epoch.TotalStakingRewards
	e.TotalRewardItems = epoch.TotalRewardItems
	e.TotalRequestCounts = epoch.TotalRequestCounts

	return nil
}

func (e *Epoch) Export(epochItems []*schema.EpochItem) (*schema.Epoch, error) {
	epoch := schema.Epoch{
		ID:                    e.ID,
		StartTimestamp:        e.StartTimestamp.Unix(),
		EndTimestamp:          e.EndTimestamp.Unix(),
		TransactionHash:       common.HexToHash(e.TransactionHash),
		TransactionIndex:      e.TransactionIndex,
		BlockTimestamp:        e.BlockTimestamp.Unix(),
		BlockHash:             common.HexToHash(e.BlockHash),
		BlockNumber:           new(big.Int).SetUint64(e.BlockNumber),
		TotalOperationRewards: e.TotalOperationRewards,
		TotalStakingRewards:   e.TotalStakingRewards,
		TotalRewardItems:      e.TotalRewardItems,
		TotalRequestCounts:    e.TotalRequestCounts,
		RewardItems:           epochItems,
	}

	return &epoch, nil
}

type Epochs []*Epoch

func (e *Epochs) Export(epochItems []*schema.EpochItem) ([]*schema.Epoch, error) {
	if len(*e) == 0 {
		return nil, nil
	}

	itemsMap := make(map[common.Hash][]*schema.EpochItem, len(epochItems))

	for _, item := range epochItems {
		if _, ok := itemsMap[item.TransactionHash]; !ok {
			itemsMap[item.TransactionHash] = make([]*schema.EpochItem, 0, 1)
		}

		itemsMap[item.TransactionHash] = append(itemsMap[item.TransactionHash], item)
	}

	epochs := make([]*schema.Epoch, 0, len(*e))

	for _, epoch := range *e {
		epoch, err := epoch.Export(itemsMap[common.HexToHash(epoch.TransactionHash)])
		if err != nil {
			return nil, err
		}

		epochs = append(epochs, epoch)
	}

	return epochs, nil
}

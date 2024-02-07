package model

import (
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

type Epoch struct {
	ID                    uint64                      `json:"id"`
	StartTimestamp        int64                       `json:"startTimestamp"`
	EndTimestamp          int64                       `json:"endTimestamp"`
	Transaction           TransactionEventTransaction `json:"transaction"`
	Block                 TransactionEventBlock       `json:"block"`
	TotalOperationRewards string                      `json:"totalOperationRewards"`
	TotalStakingRewards   string                      `json:"totalStakingRewards"`
	TotalRewardItems      int                         `json:"totalRewardItems"`
	RewardItems           []*schema.EpochItem         `json:"rewardItems,omitempty"`
	CreatedAt             int64                       `json:"-"`
	UpdatedAt             int64                       `json:"-"`
}

func NewEpoch(epoch *schema.Epoch) *Epoch {
	return &Epoch{
		ID:             epoch.ID,
		StartTimestamp: epoch.StartTimestamp,
		EndTimestamp:   epoch.EndTimestamp,
		Transaction: TransactionEventTransaction{
			Hash:  epoch.TransactionHash,
			Index: epoch.TransactionIndex,
		},
		Block: TransactionEventBlock{
			Hash:      epoch.BlockHash,
			Number:    epoch.BlockNumber,
			Timestamp: epoch.BlockTimestamp,
		},
		TotalOperationRewards: epoch.TotalOperationRewards,
		TotalStakingRewards:   epoch.TotalStakingRewards,
		TotalRewardItems:      epoch.TotalRewardItems,
		RewardItems:           epoch.RewardItems,
		CreatedAt:             epoch.CreatedAt,
		UpdatedAt:             epoch.UpdatedAt,
	}
}

func NewEpochs(epochs []*schema.Epoch) []*Epoch {
	epochModels := make([]*Epoch, 0, len(epochs))
	for _, epoch := range epochs {
		epochModels = append(epochModels, NewEpoch(epoch))
	}

	return epochModels
}

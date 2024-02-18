package model

import (
	"sort"

	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

type Epoch struct {
	ID            uint64              `json:"id"`
	Distributions []*EpochTransaction `json:"distributions"`
}

type EpochTransaction struct {
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

func NewEpochs(epochs []*schema.Epoch) []*Epoch {
	epochMap := make(map[uint64][]*EpochTransaction)

	for _, epoch := range epochs {
		if _, ok := epochMap[epoch.ID]; !ok {
			epochMap[epoch.ID] = make([]*EpochTransaction, 0)
		}

		epochMap[epoch.ID] = append(epochMap[epoch.ID], NewEpochTransaction(epoch))
	}

	results := make([]*Epoch, 0, len(epochMap))

	for id, transactions := range epochMap {
		results = append(results, &Epoch{
			ID:            id,
			Distributions: transactions,
		})
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].ID > results[j].ID
	})

	return results
}

func NewEpochTransaction(epoch *schema.Epoch) *EpochTransaction {
	return &EpochTransaction{
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

func NewEpochTransactions(epochs []*schema.Epoch) []*EpochTransaction {
	epochModels := make([]*EpochTransaction, 0, len(epochs))
	for _, epoch := range epochs {
		epochModels = append(epochModels, NewEpochTransaction(epoch))
	}

	return epochModels
}

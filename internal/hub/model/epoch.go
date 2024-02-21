package model

import (
	"sort"

	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type Epoch struct {
	ID                    uint64              `json:"id"`
	StartTimestamp        int64               `json:"startTimestamp"`
	EndTimestamp          int64               `json:"endTimestamp"`
	TotalOperationRewards string              `json:"totalOperationRewards"`
	TotalStakingRewards   string              `json:"totalStakingRewards"`
	TotalRewardItems      int                 `json:"totalRewardItems"`
	Distributions         []*EpochTransaction `json:"distributions"`
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
	epochMap := make(map[uint64]*Epoch)

	for _, epoch := range epochs {
		if _, ok := epochMap[epoch.ID]; !ok {
			epochMap[epoch.ID] = &Epoch{
				ID:                    epoch.ID,
				StartTimestamp:        epoch.StartTimestamp,
				EndTimestamp:          epoch.EndTimestamp,
				TotalOperationRewards: "0",
				TotalStakingRewards:   "0",
				Distributions:         make([]*EpochTransaction, 0),
			}
		}

		epochMap[epoch.ID].TotalOperationRewards = lo.Must(decimal.NewFromString(epochMap[epoch.ID].TotalOperationRewards)).
			Add(lo.Must(decimal.NewFromString(epoch.TotalOperationRewards))).String()
		epochMap[epoch.ID].TotalStakingRewards = lo.Must(decimal.NewFromString(epochMap[epoch.ID].TotalStakingRewards)).
			Add(lo.Must(decimal.NewFromString(epoch.TotalStakingRewards))).String()
		epochMap[epoch.ID].TotalRewardItems += epoch.TotalRewardItems
		epochMap[epoch.ID].Distributions = append(epochMap[epoch.ID].Distributions, NewEpochTransaction(epoch))
	}

	results := make([]*Epoch, 0, len(epochMap))

	for _, epoch := range epochMap {
		results = append(results, epoch)
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].ID > results[j].ID
	})

	return results
}

func NewEpoch(id uint64, epochs []*schema.Epoch) *Epoch {
	epoch := &Epoch{
		ID:                    id,
		StartTimestamp:        epochs[0].StartTimestamp,
		EndTimestamp:          epochs[0].EndTimestamp,
		TotalOperationRewards: "0",
		TotalStakingRewards:   "0",
		Distributions:         make([]*EpochTransaction, 0),
	}

	for _, distributions := range epochs {
		epoch.TotalOperationRewards = lo.Must(decimal.NewFromString(epoch.TotalOperationRewards)).
			Add(lo.Must(decimal.NewFromString(distributions.TotalOperationRewards))).String()
		epoch.TotalStakingRewards = lo.Must(decimal.NewFromString(epoch.TotalStakingRewards)).
			Add(lo.Must(decimal.NewFromString(distributions.TotalStakingRewards))).String()
		epoch.TotalRewardItems += distributions.TotalRewardItems
		epoch.Distributions = append(epoch.Distributions, NewEpochTransaction(distributions))
	}

	return epoch
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

package nta

import (
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type GetEpochsRequest struct {
	Cursor *string `query:"cursor"`
	Limit  int     `query:"limit" validate:"min=1,max=50" default:"10"`
}

type GetEpochRequest struct {
	EpochID   uint64  `param:"epoch_id" validate:"required"`
	ItemLimit int     `query:"item_limit" validate:"min=1,max=50" default:"10"`
	Cursor    *string `query:"cursor"`
}

type GetEpochDistributionRequest struct {
	TransactionHash common.Hash `param:"transaction_hash" validate:"required"`
	ItemLimit       int         `query:"item_limit" validate:"min=1,max=50" default:"10"`
	Cursor          *string     `query:"cursor"`
}

type GetEpochNodeRewardsRequest struct {
	NodeAddress common.Address `param:"node_address" validate:"required"`
	Limit       int            `query:"limit" validate:"min=1,max=50" default:"10"`
	Cursor      *string        `query:"cursor"`
}

type GetEpochsResponseData []*Epoch

type GetEpochResponseData *Epoch

type GetEpochDistributionResponseData *EpochTransaction

type GetEpochNodeRewardsResponseData *Epoch

type Epoch struct {
	ID                    uint64          `json:"id"`
	StartTimestamp        int64           `json:"start_timestamp"`
	EndTimestamp          int64           `json:"end_timestamp"`
	TotalOperationRewards decimal.Decimal `json:"total_operation_rewards"`
	TotalStakingRewards   decimal.Decimal `json:"total_staking_rewards"`
	TotalRequestCounts    decimal.Decimal `json:"total_request_counts"`
	TotalRewardedNodes    int             `json:"total_rewarded_nodes"`

	Distributions []*EpochTransaction `json:"distributions"`
}

type EpochTransaction struct {
	ID                    uint64                      `json:"id"`
	StartTimestamp        int64                       `json:"start_timestamp"`
	EndTimestamp          int64                       `json:"end_timestamp"`
	Transaction           TransactionEventTransaction `json:"transaction"`
	Block                 TransactionEventBlock       `json:"block"`
	TotalOperationRewards decimal.Decimal             `json:"total_operation_rewards"`
	TotalStakingRewards   decimal.Decimal             `json:"total_staking_rewards"`
	TotalRequestCounts    decimal.Decimal             `json:"total_request_counts"`
	TotalRewardedNodes    int                         `json:"total_rewarded_nodes"`
	RewardedNodes         []*schema.RewardedNode      `json:"rewarded_nodes,omitempty"`
	CreatedAt             int64                       `json:"-"`
	UpdatedAt             int64                       `json:"-"`
}

func NewEpochs(epochs []*schema.Epoch) GetEpochsResponseData {
	epochMap := make(map[uint64]*Epoch)

	for _, epoch := range epochs {
		if _, ok := epochMap[epoch.ID]; !ok {
			epochMap[epoch.ID] = &Epoch{
				ID:                    epoch.ID,
				StartTimestamp:        epoch.StartTimestamp,
				EndTimestamp:          epoch.EndTimestamp,
				TotalOperationRewards: decimal.NewFromInt(0),
				TotalStakingRewards:   decimal.NewFromInt(0),
				TotalRequestCounts:    decimal.NewFromInt(0),
				Distributions:         make([]*EpochTransaction, 0),
			}
		}

		epochMap[epoch.ID].TotalOperationRewards = epochMap[epoch.ID].TotalOperationRewards.Add(epoch.TotalOperationRewards)
		epochMap[epoch.ID].TotalStakingRewards = epochMap[epoch.ID].TotalStakingRewards.Add(epoch.TotalStakingRewards)
		epochMap[epoch.ID].TotalRequestCounts = epochMap[epoch.ID].TotalRequestCounts.Add(epoch.TotalRequestCounts)
		epochMap[epoch.ID].TotalRewardedNodes += epoch.TotalRewardedNodes
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
		TotalOperationRewards: decimal.NewFromInt(0),
		TotalStakingRewards:   decimal.NewFromInt(0),
		TotalRequestCounts:    decimal.NewFromInt(0),
		Distributions:         make([]*EpochTransaction, 0),
	}

	for _, distributions := range epochs {
		epoch.TotalOperationRewards = epoch.TotalOperationRewards.Add(distributions.TotalOperationRewards)
		epoch.TotalStakingRewards = epoch.TotalStakingRewards.Add(distributions.TotalStakingRewards)
		epoch.TotalRequestCounts = epoch.TotalRequestCounts.Add(distributions.TotalRequestCounts)
		epoch.TotalRewardedNodes += distributions.TotalRewardedNodes
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
		TotalRewardedNodes:    epoch.TotalRewardedNodes,
		TotalRequestCounts:    epoch.TotalRequestCounts,
		RewardedNodes:         epoch.RewardedNodes,
		CreatedAt:             epoch.CreatedAt,
		UpdatedAt:             epoch.UpdatedAt,
	}
}

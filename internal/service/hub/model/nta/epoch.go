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
	ID         uint64  `param:"id" validate:"required"`
	ItemsLimit int     `query:"itemsLimit" validate:"min=1,max=50" default:"10"`
	Cursor     *string `query:"cursor"`
}

type GetEpochDistributionRequest struct {
	TransactionHash common.Hash `param:"transaction" validate:"required"`
	ItemsLimit      int         `query:"itemsLimit" validate:"min=1,max=50" default:"10"`
	Cursor          *string     `query:"cursor"`
}

type GetEpochNodeRewardsRequest struct {
	NodeAddress common.Address `param:"node" validate:"required"`
	Limit       int            `query:"limit" validate:"min=1,max=50" default:"10"`
	Cursor      *string        `query:"cursor"`
}

type GetEpochsResponseData []*Epoch

type GetEpochResponseData *Epoch

type GetEpochDistributionResponseData *EpochTransaction

type GetEpochNodeRewardsResponseData *Epoch

type Epoch struct {
	ID                    uint64          `json:"id"`
	StartTimestamp        int64           `json:"startTimestamp"`
	EndTimestamp          int64           `json:"endTimestamp"`
	TotalOperationRewards decimal.Decimal `json:"totalOperationRewards"`
	TotalStakingRewards   decimal.Decimal `json:"totalStakingRewards"`
	TotalRequestCounts    decimal.Decimal `json:"totalRequestCounts"`
	TotalRewardedNodes    int             `json:"totalRewardedNodes"`

	Distributions []*EpochTransaction `json:"distributions"`
}

type EpochTransaction struct {
	ID                    uint64                      `json:"id"`
	StartTimestamp        int64                       `json:"startTimestamp"`
	EndTimestamp          int64                       `json:"endTimestamp"`
	Transaction           TransactionEventTransaction `json:"transaction"`
	Block                 TransactionEventBlock       `json:"block"`
	TotalOperationRewards decimal.Decimal             `json:"totalOperationRewards"`
	TotalStakingRewards   decimal.Decimal             `json:"totalStakingRewards"`
	TotalRequestCounts    decimal.Decimal             `json:"totalRequestCounts"`
	TotalRewardedNodes    int                         `json:"totalRewardedNodes"`
	RewardedNodes         []*schema.RewardedNode      `json:"rewardedNodes,omitempty"`
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

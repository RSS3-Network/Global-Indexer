package schema

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type Epoch struct {
	ID                    uint64          `json:"id"`
	StartTimestamp        int64           `json:"startTimestamp"`
	EndTimestamp          int64           `json:"endTimestamp"`
	TransactionHash       common.Hash     `json:"transactionHash"`
	TransactionIndex      uint            `json:"transactionIndex"`
	BlockHash             common.Hash     `json:"blockHash"`
	BlockNumber           *big.Int        `json:"blockNumber"`
	BlockTimestamp        int64           `json:"blockTimestamp"`
	TotalOperationRewards decimal.Decimal `json:"totalOperationRewards"`
	TotalStakingRewards   decimal.Decimal `json:"totalStakingRewards"`
	TotalRewardNodes      int             `json:"totalRewardNodes"`
	TotalRequestCounts    decimal.Decimal `json:"totalRequestCounts"`
	RewardedNodes         []*RewardedNode `json:"rewardedNodes,omitempty"`
	CreatedAt             int64           `json:"-"`
	UpdatedAt             int64           `json:"-"`
}

type RewardedNode struct {
	EpochID          uint64          `json:"epochID"`
	Index            int             `json:"index"`
	TransactionHash  common.Hash     `json:"transactionHash"`
	NodeAddress      common.Address  `json:"nodeAddress"`
	OperationRewards decimal.Decimal `json:"operationRewards"`
	StakingRewards   decimal.Decimal `json:"stakingRewards"`
	TaxCollected     decimal.Decimal `json:"taxCollected"`
	RequestCount     decimal.Decimal `json:"requestCount"`
}

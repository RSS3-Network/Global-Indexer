package schema

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

// Epoch records an Epoch and its proof of rewards distribution
type Epoch struct {
	ID uint64 `json:"id"`
	// StartTimestamp when an Epoch begins.
	StartTimestamp int64 `json:"startTimestamp"`
	// EndTimestamp when an Epoch ends.
	EndTimestamp     int64       `json:"endTimestamp"`
	TransactionHash  common.Hash `json:"transactionHash"`
	TransactionIndex uint        `json:"transactionIndex"`
	BlockHash        common.Hash `json:"blockHash"`
	BlockNumber      *big.Int    `json:"blockNumber"`
	BlockTimestamp   int64       `json:"blockTimestamp"`
	// total Operation Rewards distributed.
	TotalOperationRewards decimal.Decimal `json:"totalOperationRewards"`
	// total Staking Rewards distributed.
	TotalStakingRewards decimal.Decimal `json:"totalStakingRewards"`
	// the number of Nodes that received rewards.
	TotalRewardedNodes int `json:"totalRewardedNodes"`
	// the list of Nodes that received rewards and the amount they received.
	RewardedNodes []*RewardedNode `json:"rewardedNodes,omitempty"`
	// the total number of DSL requests made during the Epoch.
	TotalRequestCounts decimal.Decimal `json:"totalRequestCounts"`
	CreatedAt          int64           `json:"-"`
	UpdatedAt          int64           `json:"-"`
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

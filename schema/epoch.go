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
	StartTimestamp int64 `json:"start_timestamp"`
	// EndTimestamp when an Epoch ends.
	EndTimestamp     int64       `json:"end_timestamp"`
	TransactionHash  common.Hash `json:"transaction_hash"`
	TransactionIndex uint        `json:"transaction_index"`
	BlockHash        common.Hash `json:"block_hash"`
	BlockNumber      *big.Int    `json:"block_number"`
	BlockTimestamp   int64       `json:"block_timestamp"`
	// total Operation Rewards distributed.
	TotalOperationRewards decimal.Decimal `json:"total_operation_rewards"`
	// total Staking Rewards distributed.
	TotalStakingRewards decimal.Decimal `json:"total_staking_rewards"`
	// the number of Nodes that received rewards.
	TotalRewardedNodes int `json:"total_rewarded_nodes"`
	// the list of Nodes that received rewards and the amount they received.
	RewardedNodes []*RewardedNode `json:"rewarded_nodes,omitempty"`
	// the total number of DSL requests made during the Epoch.
	TotalRequestCounts decimal.Decimal `json:"total_request_counts"`
	Finalized          bool            `json:"-"`
	CreatedAt          int64           `json:"-"`
	UpdatedAt          int64           `json:"-"`
}

type RewardedNode struct {
	EpochID          uint64          `json:"epoch_id"`
	Index            int             `json:"index"`
	TransactionHash  common.Hash     `json:"transaction_hash"`
	NodeAddress      common.Address  `json:"node_address"`
	OperationRewards decimal.Decimal `json:"operation_rewards"`
	StakingRewards   decimal.Decimal `json:"staking_rewards"`
	TaxCollected     decimal.Decimal `json:"tax_collected"`
	RequestCount     decimal.Decimal `json:"request_count"`
}

type FindEpochsQuery struct {
	EpochID     *uint64
	Distinct    *bool
	Limit       *int
	Cursor      *string
	BlockNumber *uint64
}

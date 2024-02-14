package schema

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Epoch struct {
	ID                    uint64       `json:"id"`
	StartTimestamp        int64        `json:"startTimestamp"`
	EndTimestamp          int64        `json:"endTimestamp"`
	TransactionHash       common.Hash  `json:"transactionHash"`
	TransactionIndex      uint         `json:"transactionIndex"`
	BlockHash             common.Hash  `json:"blockHash"`
	BlockNumber           *big.Int     `json:"blockNumber"`
	BlockTimestamp        int64        `json:"blockTimestamp"`
	TotalOperationRewards string       `json:"totalOperationRewards"`
	TotalStakingRewards   string       `json:"totalStakingRewards"`
	TotalRewardItems      int          `json:"totalRewardItems"`
	RewardItems           []*EpochItem `json:"rewardItems,omitempty"`
	CreatedAt             int64        `json:"-"`
	UpdatedAt             int64        `json:"-"`
}

type EpochItem struct {
	EpochID          uint64         `json:"epochID"`
	Index            int            `json:"index"`
	TransactionHash  common.Hash    `json:"transactionHash"`
	NodeAddress      common.Address `json:"nodeAddress"`
	RequestFees      string         `json:"requestFees"`
	OperationRewards string         `json:"operationRewards"`
	StakingRewards   string         `json:"stakingRewards"`
	TaxAmounts       string         `json:"taxAmounts"`
}

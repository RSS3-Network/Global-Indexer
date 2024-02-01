package schema

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Epoch struct {
	ID                    *big.Int     `json:"id"`
	StartTimestamp        int64        `json:"startTimestamp"`
	EndTimestamp          int64        `json:"endTimestamp"`
	TransactionHash       common.Hash  `json:"transactionHash"`
	BlockNumber           *big.Int     `json:"blockNumber"`
	TotalOperationRewards string       `json:"totalOperationRewards"`
	TotalStakingRewards   string       `json:"totalStakingRewards"`
	TotalRewardItems      int          `json:"totalRewardItems"`
	RewardItems           []*EpochItem `json:"rewardItems"`
	Success               bool         `json:"success"`
	CreatedAt             int64        `json:"-"`
	UpdatedAt             int64        `json:"-"`
}

type EpochItem struct {
	EpochID          *big.Int       `json:"-"`
	Index            int            `json:"index"`
	NodeAddress      common.Address `json:"nodeAddress"`
	RequestFees      string         `json:"requestFees"`
	OperationRewards string         `json:"operationRewards"`
	StakingRewards   string         `json:"stakingRewards"`
	TaxAmounts       string         `json:"taxAmounts"`
}

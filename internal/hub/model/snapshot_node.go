package model

import (
	"math/big"

	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

type NodeSnapshot struct {
	EpochID        uint64   `json:"epochID"`
	Count          int64    `json:"count"`
	BlockHash      string   `json:"blockHash"`
	BlockNumber    *big.Int `json:"blockNumber"`
	BlockTimestamp uint64   `json:"blockTimestamp"`
}

func NewNodeSnapshots(nodeSnapshots []*schema.NodeSnapshot) []*NodeSnapshot {
	return lo.Map(nodeSnapshots, func(nodeSnapshot *schema.NodeSnapshot, _ int) *NodeSnapshot {
		return &NodeSnapshot{
			EpochID:        nodeSnapshot.EpochID,
			Count:          nodeSnapshot.Count,
			BlockHash:      nodeSnapshot.BlockHash,
			BlockNumber:    new(big.Int).SetUint64(nodeSnapshot.BlockNumber),
			BlockTimestamp: uint64(nodeSnapshot.BlockTimestamp.Unix()),
		}
	})
}

type StakeSnapshot NodeSnapshot

func NewStakeSnapshots(stakeSnapshots []*schema.StakeSnapshot) []*StakeSnapshot {
	return lo.Map(stakeSnapshots, func(stakeSnapshot *schema.StakeSnapshot, _ int) *StakeSnapshot {
		return &StakeSnapshot{
			EpochID:        stakeSnapshot.EpochID,
			Count:          stakeSnapshot.Count,
			BlockHash:      stakeSnapshot.BlockHash,
			BlockNumber:    new(big.Int).SetUint64(stakeSnapshot.BlockNumber),
			BlockTimestamp: uint64(stakeSnapshot.BlockTimestamp.Unix()),
		}
	})
}

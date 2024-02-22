package model

import (
	"time"

	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

type NodeSnapshot struct {
	Date  string `json:"date"`
	Count uint64 `json:"count"`
}

func NewNodeSnapshots(nodeSnapshots []*schema.NodeSnapshot) []*NodeSnapshot {
	return lo.Map(nodeSnapshots, func(nodeSnapshot *schema.NodeSnapshot, _ int) *NodeSnapshot {
		return &NodeSnapshot{
			Date:  nodeSnapshot.Date.Format(time.DateOnly),
			Count: uint64(nodeSnapshot.Count),
		}
	})
}

type StakeSnapshot NodeSnapshot

func NewStakeSnapshots(stakeSnapshots []*schema.StakeSnapshot) []*StakeSnapshot {
	return lo.Map(stakeSnapshots, func(stakeSnapshot *schema.StakeSnapshot, _ int) *StakeSnapshot {
		return &StakeSnapshot{
			Date:  stakeSnapshot.Date.Format(time.DateOnly),
			Count: uint64(stakeSnapshot.Count),
		}
	})
}

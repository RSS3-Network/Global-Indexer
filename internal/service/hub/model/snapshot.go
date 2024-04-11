package model

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
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

func NewStakeSnapshots(stakeSnapshots []*schema.StakerCountSnapshot) []*StakeSnapshot {
	return lo.Map(stakeSnapshots, func(stakeSnapshot *schema.StakerCountSnapshot, _ int) *StakeSnapshot {
		return &StakeSnapshot{
			Date:  stakeSnapshot.Date.Format(time.DateOnly),
			Count: uint64(stakeSnapshot.Count),
		}
	})
}

type NodeMinTokensToStakeSnapshots struct {
	NodeAddress common.Address                         `json:"nodeAddress"`
	Snapshots   []*schema.NodeMinTokensToStakeSnapshot `json:"snapshots"`
}

func NewNodeMinTokensToStakeSnapshots(nodeMinTokensToStakeSnapshots []*schema.NodeMinTokensToStakeSnapshot) []*NodeMinTokensToStakeSnapshots {
	nodeMap := make(map[common.Address][]*schema.NodeMinTokensToStakeSnapshot)

	for _, nodeMinTokensToStakeSnapshot := range nodeMinTokensToStakeSnapshots {
		if _, ok := nodeMap[nodeMinTokensToStakeSnapshot.NodeAddress]; !ok {
			nodeMap[nodeMinTokensToStakeSnapshot.NodeAddress] = make([]*schema.NodeMinTokensToStakeSnapshot, 0)
		}

		nodeMap[nodeMinTokensToStakeSnapshot.NodeAddress] = append(nodeMap[nodeMinTokensToStakeSnapshot.NodeAddress], nodeMinTokensToStakeSnapshot)
	}

	data := make([]*NodeMinTokensToStakeSnapshots, 0)

	for nodeAddress, snapshots := range nodeMap {
		data = append(data, &NodeMinTokensToStakeSnapshots{
			NodeAddress: nodeAddress,
			Snapshots:   snapshots,
		})
	}

	return data
}

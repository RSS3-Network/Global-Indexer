package nta

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
)

type BatchNodeMinTokensToStakeRequest struct {
	NodeAddresses   []*common.Address `json:"node_addresses" validate:"required"`
	OnlyStartAndEnd bool              `json:"only_start_and_end"`
}

type GetStakerProfitSnapshotsRequest struct {
	StakerAddress common.Address `query:"staker" validate:"required"`
	Limit         *int           `query:"limit"`
	Cursor        *string        `query:"cursor"`
	BeforeDate    *time.Time     `query:"before_date"`
	AfterDate     *time.Time     `query:"after_date"`
}

type GetNodeOperationProfitSnapshotsRequest struct {
	NodeAddress common.Address `query:"node" validate:"required"`
	Limit       *int           `query:"limit"`
	Cursor      *string        `query:"cursor"`
	BeforeDate  *time.Time     `query:"before_date"`
	AfterDate   *time.Time     `query:"after_date"`
}

type GetNodeCountSnapshotsResponseData []*CountSnapshot

type BatchGetNodeMinTokensToStakeSnapshotsResponseData []*NodeMinTokensToStakeSnapshots

type GetStakerProfitSnapshotsResponseData []*CountSnapshot

type GetOperatorProfitsSnapshotsResponseData []*schema.OperatorProfitSnapshot

type CountSnapshot struct {
	Date  string `json:"date"`
	Count uint64 `json:"count"`
}

type NodeMinTokensToStakeSnapshots struct {
	NodeAddress common.Address                         `json:"node_address"`
	Snapshots   []*schema.NodeMinTokensToStakeSnapshot `json:"snapshots"`
}

func NewNodeCountSnapshots(nodeSnapshots []*schema.NodeSnapshot) GetNodeCountSnapshotsResponseData {
	return lo.Map(nodeSnapshots, func(nodeSnapshot *schema.NodeSnapshot, _ int) *CountSnapshot {
		return &CountSnapshot{
			Date:  nodeSnapshot.Date.Format(time.DateOnly),
			Count: uint64(nodeSnapshot.Count),
		}
	})
}

func NewStakeSnapshots(stakeSnapshots []*schema.StakerCountSnapshot) GetStakerProfitSnapshotsResponseData {
	return lo.Map(stakeSnapshots, func(stakeSnapshot *schema.StakerCountSnapshot, _ int) *CountSnapshot {
		return &CountSnapshot{
			Date:  stakeSnapshot.Date.Format(time.DateOnly),
			Count: uint64(stakeSnapshot.Count),
		}
	})
}

func NewNodeMinTokensToStakeSnapshots(nodeMinTokensToStakeSnapshots []*schema.NodeMinTokensToStakeSnapshot) BatchGetNodeMinTokensToStakeSnapshotsResponseData {
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

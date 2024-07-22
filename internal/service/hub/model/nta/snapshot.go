package nta

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
)

type GetStakerProfitSnapshotsRequest struct {
	StakerAddress common.Address `query:"staker_address" validate:"required"`
	Limit         *int           `query:"limit"`
	Cursor        *string        `query:"cursor"`
	BeforeDate    *time.Time     `query:"before_date"`
	AfterDate     *time.Time     `query:"after_date"`
}

type GetNodeOperationProfitSnapshotsRequest struct {
	NodeAddress common.Address `query:"node_address" validate:"required"`
	Limit       *int           `query:"limit"`
	Cursor      *string        `query:"cursor"`
	BeforeDate  *time.Time     `query:"before_date"`
	AfterDate   *time.Time     `query:"after_date"`
}

type GetNodeCountSnapshotsResponseData []*CountSnapshot

type GetStakerCountSnapshotsResponseData []*CountSnapshot

type GetOperatorProfitsSnapshotsResponseData []*schema.OperatorProfitSnapshot

type CountSnapshot struct {
	Date  string `json:"date"`
	Count uint64 `json:"count"`
}

func NewNodeCountSnapshots(nodeSnapshots []*schema.NodeSnapshot) GetNodeCountSnapshotsResponseData {
	return lo.Map(nodeSnapshots, func(nodeSnapshot *schema.NodeSnapshot, _ int) *CountSnapshot {
		return &CountSnapshot{
			Date:  nodeSnapshot.Date.Format(time.DateOnly),
			Count: uint64(nodeSnapshot.Count),
		}
	})
}

func NewStakerCountSnapshots(stakeSnapshots []*schema.StakerCountSnapshot) GetStakerCountSnapshotsResponseData {
	return lo.Map(stakeSnapshots, func(stakeSnapshot *schema.StakerCountSnapshot, _ int) *CountSnapshot {
		return &CountSnapshot{
			Date:  stakeSnapshot.Date.Format(time.DateOnly),
			Count: uint64(stakeSnapshot.Count),
		}
	})
}

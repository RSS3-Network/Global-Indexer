package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Stat struct {
	Address              common.Address `json:"address"`
	Endpoint             string         `json:"-"`
	Score                float64        `json:"score"`
	IsPublicGood         bool           `json:"isPublicGood"`
	IsFullNode           bool           `json:"isFullNode"`
	IsRssNode            bool           `json:"isRssNode"`
	Staking              float64        `json:"staking"`
	Epoch                int64          `json:"epoch"`
	TotalRequest         int64          `json:"totalRequest"`
	EpochRequest         int64          `json:"epochRequest"`
	EpochInvalidRequest  int64          `json:"epochInvalidRequest"`
	DecentralizedNetwork int            `json:"decentralizedNetwork"`
	FederatedNetwork     int            `json:"federatedNetwork"`
	Indexer              int            `json:"indexer"`
	ResetAt              time.Time      `json:"resetAt"`
}

type StatQuery struct {
	Address      *common.Address  `query:"address" form:"address,omitempty"`
	AddressList  []common.Address `query:"addressList" form:"addressList,omitempty"`
	IsFullNode   *bool            `query:"isFullNode" form:"isFullNode,omitempty"`
	IsRssNode    *bool            `query:"isRssNode" form:"isRssNode,omitempty"`
	PointsOrder  *string          `query:"pointsOrder" form:"pointsOrder,omitempty"`
	ValidRequest *int             `query:"validRequest" form:"validRequest,omitempty"`
	Limit        *int             `query:"limit" form:"limit,omitempty"`
	Cursor       *string          `query:"cursor" form:"cursor,omitempty"`
}

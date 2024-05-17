package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Stat struct {
	Address              common.Address `json:"address"`
	Endpoint             string         `json:"-"`
	Score                float64        `json:"score"`
	IsPublicGood         bool           `json:"is_public_good"`
	IsFullNode           bool           `json:"is_full_node"`
	IsRssNode            bool           `json:"is_rss_node"`
	Staking              float64        `json:"staking"`
	Epoch                int64          `json:"epoch"`
	TotalRequest         int64          `json:"total_request"`
	EpochRequest         int64          `json:"epoch_request"`
	EpochInvalidRequest  int64          `json:"epoch_invalid_request"`
	DecentralizedNetwork int            `json:"decentralized_network"`
	FederatedNetwork     int            `json:"federated_network"`
	Indexer              int            `json:"indexer"`
	ResetAt              time.Time      `json:"reset_at"`
}

type StatQuery struct {
	Address      *common.Address  `query:"address" form:"address,omitempty"`
	Addresses    []common.Address `query:"Addresses" form:"addresses,omitempty"`
	IsFullNode   *bool            `query:"isFullNode" form:"isFullNode,omitempty"`
	IsRssNode    *bool            `query:"isRssNode" form:"isRssNode,omitempty"`
	PointsOrder  *string          `query:"pointsOrder" form:"pointsOrder,omitempty"`
	ValidRequest *int             `query:"validRequest" form:"validRequest,omitempty"`
	Limit        *int             `query:"limit" form:"limit,omitempty"`
	Cursor       *string          `query:"cursor" form:"cursor,omitempty"`
}

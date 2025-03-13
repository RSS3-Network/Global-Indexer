package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Stat struct {
	Address              common.Address `json:"address"`
	Endpoint             string         `json:"-"`
	AccessToken          string         `json:"-"`
	Score                float64        `json:"score"`
	IsPublicGood         bool           `json:"is_public_good"`
	IsFullNode           bool           `json:"is_full_node"`
	IsRssNode            bool           `json:"is_rss_node"`
	IsAINode             bool           `json:"is_ai_node"`
	Staking              float64        `json:"staking"`
	Epoch                int64          `json:"epoch"`
	TotalRequest         int64          `json:"total_request"`
	EpochRequest         int64          `json:"epoch_request"`
	EpochInvalidRequest  int64          `json:"epoch_invalid_request"`
	DecentralizedNetwork int            `json:"decentralized_network"`
	FederatedNetwork     int            `json:"federated_network"`
	Indexer              int            `json:"indexer"`
	ResetAt              time.Time      `json:"reset_at"`

	Status   NodeStatus `json:"-"`
	HearBeat NodeStatus `json:"-"`
	Version  string     `json:"-"`
}

type StatQuery struct {
	Address      *common.Address  `query:"address" form:"address,omitempty"`
	Addresses    []common.Address `query:"addresses" form:"addresses,omitempty"`
	IsFullNode   *bool            `query:"is_full_node" form:"is_full_node,omitempty"`
	IsRssNode    *bool            `query:"is_rss_node" form:"is_rss_node,omitempty"`
	IsAINode     *bool            `query:"is_ai_node" form:"is_ai_node,omitempty"`
	PointsOrder  *string          `query:"points_order" form:"points_order,omitempty"`
	ValidRequest *int             `query:"valid_request" form:"valid_request,omitempty"`
	Limit        *int             `query:"limit" form:"limit,omitempty"`
	Cursor       *string          `query:"cursor" form:"cursor,omitempty"`
}

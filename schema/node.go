package schema

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/shopspring/decimal"
)

type Node struct {
	Address                common.Address         `json:"address"`
	Name                   string                 `json:"name"`
	Description            string                 `json:"description"`
	TaxRateBasisPoints     uint64                 `json:"taxRateBasisPoints"`
	IsPublicGood           bool                   `json:"isPublicGood"`
	OperationPoolTokens    string                 `json:"operationPoolTokens"`
	StakingPoolTokens      string                 `json:"stakingPoolTokens"`
	TotalShares            string                 `json:"totalShares"`
	SlashedTokens          string                 `json:"slashedTokens"`
	Endpoint               string                 `json:"-"`
	Stream                 json.RawMessage        `json:"-"`
	Config                 json.RawMessage        `json:"-"`
	Status                 Status                 `json:"status"`
	LastHeartbeatTimestamp int64                  `json:"lastHeartbeat"`
	Local                  []*NodeLocal           `json:"local"`
	Avatar                 *l2.ChipsTokenMetadata `json:"avatar"`
	MinTokensToStake       decimal.Decimal        `json:"minTokensToStake"`
	CreatedAt              int64                  `json:"createdAt"`
}

type NodeLocal struct {
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

//go:generate go run --mod=mod github.com/dmarkham/enumer@v1.5.9 --values --type=Status --linecomment --output node_status_string.go --json --yaml --sql
type Status int64

const (
	StatusOnline  Status = iota // online
	StatusOffline               // offline
)

type Stat struct {
	Address              common.Address `json:"address"`
	Endpoint             string         `json:"-"`
	Points               float64        `json:"points"`
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

type Indexer struct {
	Address common.Address `json:"address"`
	Network string         `json:"network"`
	Worker  string         `json:"worker"`
}

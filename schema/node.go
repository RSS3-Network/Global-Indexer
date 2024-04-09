package schema

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/shopspring/decimal"
)

type Node struct {
	ID                     *big.Int               `json:"id"`
	Address                common.Address         `json:"address"`
	Name                   string                 `json:"name"`
	Description            string                 `json:"description"`
	TaxRateBasisPoints     *uint64                `json:"taxRateBasisPoints"`
	HideTaxRate            bool                   `json:"-"`
	IsPublicGood           bool                   `json:"isPublicGood"`
	OperationPoolTokens    string                 `json:"operationPoolTokens"`
	StakingPoolTokens      string                 `json:"stakingPoolTokens"`
	TotalShares            string                 `json:"totalShares"`
	SlashedTokens          string                 `json:"slashedTokens"`
	Alpha                  bool                   `json:"alpha"`
	Endpoint               string                 `json:"-"`
	Stream                 json.RawMessage        `json:"-"`
	Config                 json.RawMessage        `json:"-"`
	Status                 NodeStatus             `json:"status"`
	LastHeartbeatTimestamp int64                  `json:"lastHeartbeat"`
	Local                  []*NodeLocal           `json:"local"`
	Avatar                 *l2.ChipsTokenMetadata `json:"avatar"`
	MinTokensToStake       decimal.Decimal        `json:"minTokensToStake"`
	APY                    decimal.Decimal        `json:"apy"`
	Score                  decimal.Decimal        `json:"score"`
	CreatedAt              int64                  `json:"createdAt"`
}

type NodeLocal struct {
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

//go:generate go run --mod=mod github.com/dmarkham/enumer@v1.5.9 --values --type=NodeStatus --linecomment --output node_status_string.go --json --yaml --sql
type NodeStatus int64

const (
	NodeStatusRegistered NodeStatus = iota // registered
	NodeStatusOnline                       // online
	NodeStatusOffline                      // offline
	NodeStatusExited                       // exiting
)

type BatchUpdateNode struct {
	Address          common.Address
	Apy              decimal.Decimal
	MinTokensToStake decimal.Decimal
}

type FindNodesQuery struct {
	NodeAddresses []common.Address
	Status        *NodeStatus
	Cursor        *string
	Limit         *int
}

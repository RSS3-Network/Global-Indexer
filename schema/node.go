package schema

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/shopspring/decimal"
)

type Node struct {
	ID                     *big.Int               `json:"id"`
	Address                common.Address         `json:"address"`
	Name                   string                 `json:"name"`
	Description            string                 `json:"description"`
	TaxRateBasisPoints     *uint64                `json:"tax_rate_basis_points"`
	HideTaxRate            bool                   `json:"-"`
	IsPublicGood           bool                   `json:"is_public_good"`
	OperationPoolTokens    string                 `json:"operation_pool_tokens"`
	StakingPoolTokens      string                 `json:"staking_pool_tokens"`
	TotalShares            string                 `json:"total_shares"`
	SlashedTokens          string                 `json:"slashed_tokens"`
	Alpha                  bool                   `json:"alpha"`
	Endpoint               string                 `json:"-"`
	Stream                 json.RawMessage        `json:"-"`
	Config                 json.RawMessage        `json:"-"`
	Status                 NodeStatus             `json:"status"`
	LastHeartbeatTimestamp int64                  `json:"last_heartbeat"`
	Location               []*NodeLocation        `json:"location"`
	Avatar                 *l2.ChipsTokenMetadata `json:"avatar"`
	APY                    decimal.Decimal        `json:"apy"`
	ActiveScore            decimal.Decimal        `json:"active_score"`
	ReliabilityScore       decimal.Decimal        `json:"reliability_score"`
	Version                string                 `json:"version"`
	AccessToken            string                 `json:"-"`
	CreatedAt              int64                  `json:"created_at"`
}

type NodeLocation struct {
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

//go:generate go run --mod=mod github.com/dmarkham/enumer@v1.5.9 --values --type=NodeStatus --linecomment --output node_status_string.go --json --yaml --sql
type NodeStatus int64

const (
	// NodeStatusRegistered
	// Node has been registered but does not meet the minimum requirements to be enter NodeStatusOnline.
	// Possible reasons:
	// - Node is not reachable by the Network.
	// - Node Operator has not deposited the minimum amount of tokens required.
	NodeStatusRegistered NodeStatus = iota // registered

	// NodeStatusOnline
	// Node is online and fully operational.
	NodeStatusOnline // online

	// NodeStatusOffline
	// Node was previously in NodeStatusOnline, but is currently offline.
	// Possible reasons:
	// - [Alpha only] Node missed a heartbeat.
	// - Node was slashed in the previous epoch, and was kicked out of the Network, the Operator did not acknowledge the slash and rejoin the Network.
	// - Node did not perform the mandatory upgrade before the deadline required by the Network.
	NodeStatusOffline // offline

	// NodeStatusExited
	// Node was previously in NodeStatusOnline, but is not anymore.
	// Possible reasons:
	// - Node announced its intention to leave the Network and gracefully exited after the mandatory waiting period.
	// - Node has been offline for a long time and is considered as exited.
	NodeStatusExited // exited

	// NodeStatusSlashed
	// Node was slashed in the current epoch, and was kicked out of the Network.
	NodeStatusSlashed // slashed

	// NodeStatusExiting
	// Node announced its intention to leave the Network, and is now in the mandatory waiting period.
	NodeStatusExiting // exiting
)

type BatchUpdateNode struct {
	Address common.Address
	Apy     decimal.Decimal
}

type FindNodesQuery struct {
	NodeAddresses []common.Address
	Status        *NodeStatus
	Cursor        *string
	Limit         *int
	OrderByScore  bool
}

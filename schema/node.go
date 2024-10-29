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
	ReliabilityScore       decimal.Decimal        `json:"reliability_score"`
	Version                string                 `json:"version"`
	Type                   string                 `json:"type"`
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
	// NodeStatusNone
	// The Node has been created on the VSL, but has not been registered yet.
	NodeStatusNone NodeStatus = iota // none

	// NodeStatusRegistered
	// The Node is registered on the VSL with a sufficient deposit.
	NodeStatusRegistered // registered

	// NodeStatusInitializing
	// The Node is operating on the DSL.
	// Automated tasks will be executed at this stage to ensure the Node is in a healthy condition.
	// This state applies to the initial startup or the first startup following any change in the Node’s coverage.
	NodeStatusInitializing // initializing

	// NodeStatusOutdated
	// The Node is outdated and needs to be updated to the minimum required version.
	NodeStatusOutdated // outdated

	// NodeStatusOnline
	// The Node is online and fully operational.
	NodeStatusOnline // online

	// NodeStatusOffline
	// The Node is not operational and not participating in network activities on the DSL.
	NodeStatusOffline // offline

	// NodeStatusSlashing
	// The Node has been reached the demotion threshold and is currently in the appeal period.
	NodeStatusSlashing // slashing

	// NodeStatusSlashed
	// The Node has been slashed due to a violation of network rules or malicious behavior on the VSL.
	NodeStatusSlashed // slashed

	// NodeStatusExiting
	// The Node is in the process of exiting the Network on the VSL.
	NodeStatusExiting // exiting

	// NodeStatusExited
	// The Node has successfully exited the Network on the VSL.
	NodeStatusExited // exited

)

//go:generate go run --mod=mod github.com/dmarkham/enumer@v1.5.9 --values --type=NodeType --linecomment --output node_type_string.go --json --yaml --sql
type NodeType int

const (
	// NodeTypeAlpha
	// Nodes in the alpha phase of the network will receive staking rewards, but they do not actually contribute to the information network.
	NodeTypeAlpha NodeType = iota // alpha

	// NodeTypeBeta
	// Nodes in the beta phase of the network do not require staking and will not receive rewards, but they do contribute to the information network and are referred to as public good nodes.
	NodeTypeBeta // beta

	// NodeTypeProduction
	// Nodes in the production phase of the network are required to contribute to the information network. All nodes, except for public good nodes that do not require staking, will receive staking and operation rewards.
	// Upon entering the production phase, nodes from both the alpha and beta phases are required to upgrade to production node version.
	NodeTypeProduction // production
)

type BatchUpdateNode struct {
	Address common.Address
	Apy     decimal.Decimal
}

type FindNodesQuery struct {
	NodeAddresses []common.Address
	Status        *NodeStatus
	Type          *NodeType
	Cursor        *string
	Limit         *int
}

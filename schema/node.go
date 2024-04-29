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
	Location               []*NodeLocation        `json:"location"`
	Avatar                 *l2.ChipsTokenMetadata `json:"avatar"`
	MinTokensToStake       decimal.Decimal        `json:"minTokensToStake"`
	APY                    decimal.Decimal        `json:"apy"`
	ActiveScore            decimal.Decimal        `json:"activeScore"`
	ReliabilityScore       decimal.Decimal        `json:"reliabilityScore"`
	Type                   string                 `json:"type"`
	CreatedAt              int64                  `json:"createdAt"`
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
	Address          common.Address
	Apy              decimal.Decimal
	MinTokensToStake decimal.Decimal
}

type FindNodesQuery struct {
	NodeAddresses []common.Address
	Status        *NodeStatus
	Cursor        *string
	Limit         *int
	OrderByScore  bool
}

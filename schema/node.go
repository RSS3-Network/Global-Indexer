package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/serving-node/config"
)

type Node struct {
	Address                common.Address `json:"address"`
	Name                   string         `json:"name"`
	Description            string         `json:"description"`
	TaxRateBasisPoints     uint64         `json:"taxRateBasisPoints"`
	IsPublicGood           bool           `json:"isPublicGood"`
	OperationPoolTokens    string         `json:"operationPoolTokens"`
	StakingPoolTokens      string         `json:"stakingPoolTokens"`
	TotalShares            string         `json:"totalShares"`
	SlashedTokens          string         `json:"slashedTokens"`
	Endpoint               string         `json:"-"`
	Stream                 *config.Stream `json:"-"`
	Config                 *config.Node   `json:"-"`
	Status                 Status         `json:"status"`
	LastHeartbeatTimestamp int64          `json:"lastHeartbeat"`
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
	TotalRequest         int64          `json:"totalRequest"`
	EpochRequest         int64          `json:"epochRequest"`
	EpochInvalidRequest  int64          `json:"epochInvalidRequest"`
	DecentralizedNetwork int            `json:"decentralizedNetwork"`
	FederatedNetwork     int            `json:"federatedNetwork"`
	Indexer              int            `json:"indexer"`
	ResetAt              time.Time      `json:"resetAt"`
}

type Indexer struct {
	Address common.Address `json:"address"`
	Network string         `json:"network"`
	Worker  string         `json:"worker"`
}

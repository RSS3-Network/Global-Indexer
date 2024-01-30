package schema

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

type Node struct {
	Address                common.Address  `json:"address"`
	Name                   string          `json:"name"`
	Description            string          `json:"description"`
	TaxRateBasisPoints     uint64          `json:"taxRateBasisPoints"`
	IsPublicGood           bool            `json:"isPublicGood"`
	OperationPoolTokens    string          `json:"operationPoolTokens"`
	StakingPoolTokens      string          `json:"stakingPoolTokens"`
	TotalShares            string          `json:"totalShares"`
	SlashedTokens          string          `json:"slashedTokens"`
	Endpoint               string          `json:"-"`
	Stream                 json.RawMessage `json:"-"`
	Config                 json.RawMessage `json:"-"`
	Status                 Status          `json:"status"`
	LastHeartbeatTimestamp int64           `json:"lastHeartbeat"`
}

//go:generate go run --mod=mod github.com/dmarkham/enumer@v1.5.9 --values --type=Status --linecomment --output node_status_string.go --json --yaml --sql
type Status int64

const (
	StatusOnline  Status = iota // online
	StatusOffline               // offline
)

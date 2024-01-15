package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-node/config"
)

type Node struct {
	Address      common.Address `json:"address"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Endpoint     string         `json:"-"`
	TaxFraction  uint64         `json:"taxFraction"`
	IsPublicGood bool           `json:"isPublicGood"`
	Stream       *config.Stream `json:"stream"`
	Config       *config.Node   `json:"-"`
}

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

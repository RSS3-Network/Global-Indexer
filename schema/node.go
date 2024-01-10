package schema

import (
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
	Stream       *config.Stream `json:"-"`
	Config       *config.Node   `json:"-"`
}

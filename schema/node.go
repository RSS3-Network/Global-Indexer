package schema

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-node/config"
)

type Node struct {
	Address             common.Address `json:"address"`
	Name                string         `json:"name"`
	Description         string         `json:"description"`
	TaxFraction         uint64         `json:"taxFraction"`
	IsPublicGood        bool           `json:"isPublicGood"`
	OperatingPoolTokens string         `json:"operatingPoolTokens"`
	StakingPoolTokens   string         `json:"stakingPoolTokens"`
	TotalShares         string         `json:"totalShares"`
	SlashedTokens       string         `json:"slashedTokens"`
	Endpoint            string         `json:"-"`
	Stream              *config.Stream `json:"-"`
	Config              *config.Node   `json:"-"`
}

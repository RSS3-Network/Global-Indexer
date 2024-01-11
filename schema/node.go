package schema

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-node/config"
)

type Node struct {
	Address             common.Address `json:"address"`
	Name                string         `json:"name"`
	Description         string         `json:"description"`
	TaxFraction         uint64         `json:"taxFraction"`
	IsPublicGood        bool           `json:"isPublicGood"`
	OperatingPoolTokens *big.Int       `json:"operatingPoolTokens"`
	StakingPoolTokens   *big.Int       `json:"stakingPoolTokens"`
	TotalShares         *big.Int       `json:"totalShares"`
	SlashedTokens       *big.Int       `json:"slashedTokens"`
	Endpoint            string         `json:"-"`
	Stream              *config.Stream `json:"-"`
	Config              *config.Node   `json:"-"`
}

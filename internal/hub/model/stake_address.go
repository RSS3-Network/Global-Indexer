package model

import (
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

type StakeAddress struct {
	Node   *common.Address   `json:"node,omitempty"`
	Staker *common.Address   `json:"staker,omitempty"`
	Chips  *StakeAddressChip `json:"chips"`
}

type StakeAddressChip struct {
	Total    uint64       `json:"total"`
	Showcase []*StakeChip `json:"showcase"`
}

func NewStakeAddress(stakeAddress *schema.StakeAddress, baseURL url.URL) *StakeAddress {
	return &StakeAddress{
		Node:   stakeAddress.Node,
		Staker: stakeAddress.Staker,
		Chips: &StakeAddressChip{
			Total:    uint64(stakeAddress.Chips.Total),
			Showcase: NewStakeChips(stakeAddress.Chips.Showcase, baseURL),
		},
	}
}

func NewStakeAddresses(stakeAddresses []*schema.StakeAddress, baseURL url.URL) []*StakeAddress {
	return lo.Map(stakeAddresses, func(stakeAddress *schema.StakeAddress, _ int) *StakeAddress {
		return NewStakeAddress(stakeAddress, baseURL)
	})
}

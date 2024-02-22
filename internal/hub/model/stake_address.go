package model

import (
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

type StakeAddress struct {
	Address common.Address    `json:"address"`
	Chips   *StakeAddressChip `json:"chips"`
}

type StakeAddressChip struct {
	Total    uint64       `json:"total"`
	Showcase []*StakeChip `json:"showcase"`
}

func NewStakeAddress(stakeAddress *schema.StakeAddress, baseURL url.URL) *StakeAddress {
	return &StakeAddress{
		Address: stakeAddress.Address,
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

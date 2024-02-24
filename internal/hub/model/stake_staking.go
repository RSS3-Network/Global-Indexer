package model

import (
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

type StakeStaking struct {
	Staker common.Address    `json:"staker,omitempty"`
	Node   common.Address    `json:"node,omitempty"`
	Chips  StakeStakingChips `json:"chips"`
}

type StakeStakingChips struct {
	Total    uint64       `json:"total"`
	Showcase []*StakeChip `json:"showcase"`
}

func NewStakeAddress(stakeAddress *schema.StakeStaking, baseURL url.URL) *StakeStaking {
	return &StakeStaking{
		Staker: stakeAddress.Staker,
		Node:   stakeAddress.Node,
		Chips: StakeStakingChips{
			Total:    stakeAddress.Chips.Total,
			Showcase: NewStakeChips(stakeAddress.Chips.Showcase, baseURL),
		},
	}
}

func NewStakeStaking(stakeStakings []*schema.StakeStaking, baseURL url.URL) []*StakeStaking {
	return lo.Map(stakeStakings, func(stakeStaking *schema.StakeStaking, _ int) *StakeStaking {
		return NewStakeAddress(stakeStaking, baseURL)
	})
}

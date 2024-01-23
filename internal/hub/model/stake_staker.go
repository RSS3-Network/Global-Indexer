package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

type StakeStaker struct {
	User  common.Address `json:"user"`
	Chips []*big.Int     `json:"chips"`
}

func NewStakeStakers(stakeChips []*schema.StakeChip) []*StakeStaker {
	stakeStakerMap := lo.GroupBy(stakeChips, func(stakeChip *schema.StakeChip) common.Address {
		return stakeChip.Owner
	})

	stakeStakerModels := make([]*StakeStaker, 0, len(stakeStakerMap))

	for user, chips := range stakeStakerMap {
		stakeStakerModel := StakeStaker{
			User: user,
			Chips: lo.Map(chips, func(chip *schema.StakeChip, _ int) *big.Int {
				return chip.ID
			}),
		}

		stakeStakerModels = append(stakeStakerModels, &stakeStakerModel)
	}

	return stakeStakerModels
}

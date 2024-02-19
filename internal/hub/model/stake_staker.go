package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

type StakeStaker struct {
	User  common.Address `json:"user"`
	Chips []*StakeChip   `json:"chips"`
}

type StakeChip struct {
	Node common.Address `json:"node"`
	IDs  []*big.Int     `json:"ids"`
}

func NewStakeStakers(stakeChips []*schema.StakeChip) []*StakeStaker {
	stakeStakerMap := lo.GroupBy(stakeChips, func(stakeChip *schema.StakeChip) common.Address {
		return stakeChip.Owner
	})

	stakeStakerModels := make([]*StakeStaker, 0, len(stakeStakerMap))

	for user, chips := range stakeStakerMap {
		result := make(map[common.Address][]*big.Int)

		for _, chip := range chips {
			if _, exists := result[chip.Node]; !exists {
				result[chip.Node] = make([]*big.Int, 0)
			}

			result[chip.Node] = append(result[chip.Node], chip.ID)
		}

		stakeStakerModel := StakeStaker{
			User: user,
			Chips: lo.MapToSlice(result, func(node common.Address, ids []*big.Int) *StakeChip {
				return &StakeChip{
					Node: node,
					IDs:  ids,
				}
			}),
		}

		stakeStakerModels = append(stakeStakerModels, &stakeStakerModel)
	}

	return stakeStakerModels
}

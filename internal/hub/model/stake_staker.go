package model

import (
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

type StakeStaker struct {
	User  common.Address                  `json:"user"`
	Chips map[common.Address][]*StakeChip `json:"chips"`
}

func NewStakeStakers(stakeChips []*schema.StakeChip, baseURL url.URL) []*StakeStaker {
	stakeStakerMap := lo.GroupBy(stakeChips, func(stakeChip *schema.StakeChip) common.Address {
		return stakeChip.Owner
	})

	stakeStakerModels := make([]*StakeStaker, 0, len(stakeStakerMap))

	for user, chips := range stakeStakerMap {
		result := make(map[common.Address][]*StakeChip)

		for _, chip := range chips {
			if _, exists := result[chip.Node]; !exists {
				result[chip.Node] = make([]*StakeChip, 0)
			}

			metadata, _ := BuildStakeChipMetadata(chip.ID, chip.Metadata, baseURL)

			result[chip.Node] = append(result[chip.Node], &StakeChip{
				ID:       chip.ID,
				Metadata: metadata,
			})
		}

		stakeStakerModel := StakeStaker{
			User:  user,
			Chips: result,
		}

		stakeStakerModels = append(stakeStakerModels, &stakeStakerModel)
	}

	return stakeStakerModels
}

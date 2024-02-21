package model

import (
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

type StakeNode struct {
	Node  common.Address `json:"node"`
	Chips []*StakeChip   `json:"chips"`
}

func NewStakeNodes(stakeChips []*schema.StakeChip, baseURL url.URL) []*StakeNode {
	stakeNodeMap := lo.GroupBy(stakeChips, func(stakeChip *schema.StakeChip) common.Address {
		return stakeChip.Node
	})

	stakeNodeModels := make([]*StakeNode, 0, len(stakeNodeMap))

	for node, chips := range stakeNodeMap {
		stakeNodeModel := StakeNode{
			Node: node,
			Chips: lo.Map(chips, func(chip *schema.StakeChip, _ int) *StakeChip {
				metadata, _ := BuildStakeChipMetadata(chip.ID, chip.Metadata, baseURL)

				return &StakeChip{
					ID:       chip.ID,
					Metadata: metadata,
				}
			}),
		}

		stakeNodeModels = append(stakeNodeModels, &stakeNodeModel)
	}

	return stakeNodeModels
}

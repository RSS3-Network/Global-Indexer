package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

type StakeNode struct {
	Node  common.Address `json:"node"`
	Chips []*big.Int     `json:"chips"`
}

func NewStakeNodes(stakeChips []*schema.StakeChip) []*StakeNode {
	stakeNodeMap := lo.GroupBy(stakeChips, func(stakeChip *schema.StakeChip) common.Address {
		return stakeChip.Node
	})

	stakeNodeModels := make([]*StakeNode, 0, len(stakeNodeMap))

	for node, chips := range stakeNodeMap {
		stakeNodeModel := StakeNode{
			Node: node,
			Chips: lo.Map(chips, func(chip *schema.StakeChip, _ int) *big.Int {
				return chip.ID
			}),
		}

		stakeNodeModels = append(stakeNodeModels, &stakeNodeModel)
	}

	return stakeNodeModels
}

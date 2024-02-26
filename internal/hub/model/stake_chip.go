package model

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

type StakeChip struct {
	ID       *big.Int        `json:"id"`
	Node     common.Address  `json:"node"`
	Owner    common.Address  `json:"owner"`
	Metadata json.RawMessage `json:"metadata"`
}

func NewStakeChip(stakeChip *schema.StakeChip, baseURL url.URL) *StakeChip {
	result := StakeChip{
		ID:       stakeChip.ID,
		Node:     stakeChip.Node,
		Owner:    stakeChip.Owner,
		Metadata: stakeChip.Metadata,
	}

	var tokenMetadata l2.ChipsTokenMetadata
	_ = json.Unmarshal(stakeChip.Metadata, &tokenMetadata)

	tokenMetadata.Image = baseURL.JoinPath(fmt.Sprintf("/chips/%d/image.svg", result.ID)).String()

	result.Metadata, _ = json.Marshal(tokenMetadata)

	return &result
}

func NewStakeChips(stakeChips []*schema.StakeChip, baseURL url.URL) []*StakeChip {
	return lo.Map(stakeChips, func(stakeChip *schema.StakeChip, _ int) *StakeChip {
		return NewStakeChip(stakeChip, baseURL)
	})
}

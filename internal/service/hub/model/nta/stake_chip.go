package nta

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type GetStakeChipsRequest struct {
	Cursor *big.Int        `query:"cursor"`
	IDs    []*big.Int      `query:"id"`
	Node   *common.Address `query:"node"`
	Owner  *common.Address `query:"owner"`
	Limit  int             `query:"limit" default:"10" min:"1" max:"10"`
}

type GetStakeChipRequest struct {
	ID *big.Int `param:"id"`
}

type GetStakeChipsImageRequest struct {
	ID *big.Int `param:"id"`
}

type GetStakeChipsResponseData []*StakeChip

type GetStakeChipResponseData *StakeChip

type StakeChip struct {
	ID          *big.Int        `json:"id"`
	Node        common.Address  `json:"node"`
	Owner       common.Address  `json:"owner"`
	Metadata    json.RawMessage `json:"metadata"`
	Value       decimal.Decimal `json:"value"`
	LatestValue decimal.Decimal `json:"latestValue"`
}

func NewStakeChip(stakeChip *schema.StakeChip, baseURL url.URL) GetStakeChipResponseData {
	result := StakeChip{
		ID:          stakeChip.ID,
		Node:        stakeChip.Node,
		Owner:       stakeChip.Owner,
		Metadata:    stakeChip.Metadata,
		Value:       stakeChip.Value,
		LatestValue: stakeChip.LatestValue,
	}

	var tokenMetadata l2.ChipsTokenMetadata
	_ = json.Unmarshal(stakeChip.Metadata, &tokenMetadata)

	tokenMetadata.Image = baseURL.JoinPath(fmt.Sprintf("/chips/%d/image.svg", result.ID)).String()

	result.Metadata, _ = json.Marshal(tokenMetadata)

	return &result
}

func NewStakeChips(stakeChips []*schema.StakeChip, baseURL url.URL) GetStakeChipsResponseData {
	return lo.Map(stakeChips, func(stakeChip *schema.StakeChip, _ int) *StakeChip {
		return NewStakeChip(stakeChip, baseURL)
	})
}

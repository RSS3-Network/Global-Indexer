package nta

import (
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type GetStakeStakingsRequest struct {
	Cursor *string         `query:"cursor"`
	Staker *common.Address `query:"staker"`
	Node   *common.Address `query:"node"`
	Limit  int             `query:"limit" default:"10" min:"1" max:"20"`
}

type GetStakeOwnerProfitRequest struct {
	Owner common.Address `param:"owner" validate:"required"`
}

type GetStakeStakingsResponseData []*StakeStaking

type GetStakeOwnerProfitResponseData struct {
	Owner            common.Address                               `json:"owner"`
	TotalChipAmounts decimal.Decimal                              `json:"totalChipAmounts"`
	TotalChipValues  decimal.Decimal                              `json:"totalChipValues"`
	OneDay           *GetStakeOwnerProfitChangesSinceResponseData `json:"oneDay"`
	OneWeek          *GetStakeOwnerProfitChangesSinceResponseData `json:"oneWeek"`
	OneMonth         *GetStakeOwnerProfitChangesSinceResponseData `json:"oneMonth"`
}

type GetStakeOwnerProfitChangesSinceResponseData struct {
	Date             time.Time       `json:"date"`
	TotalChipAmounts decimal.Decimal `json:"totalChipAmounts"`
	TotalChipValues  decimal.Decimal `json:"totalChipValues"`
	PNL              decimal.Decimal `json:"pnl"`
}

type StakeStaking struct {
	Staker common.Address    `json:"staker,omitempty"`
	Node   common.Address    `json:"node,omitempty"`
	Value  decimal.Decimal   `json:"value"`
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
		Value:  stakeAddress.Value,
		Chips: StakeStakingChips{
			Total:    stakeAddress.Chips.Total,
			Showcase: NewStakeChips(stakeAddress.Chips.Showcase, baseURL),
		},
	}
}

func NewStakeStaking(stakeStakings []*schema.StakeStaking, baseURL url.URL) GetStakeStakingsResponseData {
	return lo.Map(stakeStakings, func(stakeStaking *schema.StakeStaking, _ int) *StakeStaking {
		return NewStakeAddress(stakeStaking, baseURL)
	})
}

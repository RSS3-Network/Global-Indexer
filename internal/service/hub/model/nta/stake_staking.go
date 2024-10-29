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
	Cursor        *string         `query:"cursor"`
	StakerAddress *common.Address `query:"staker_address"`
	NodeAddress   *common.Address `query:"node_address"`
	Limit         int             `query:"limit" default:"50" min:"1" max:"100"`
}

type GetStakerProfitRequest struct {
	StakerAddress common.Address `param:"staker_address" validate:"required"`
}

type GetStakeStakingsResponseData []*StakeStaking

type GetStakerProfitResponseData struct {
	Owner           common.Address                           `json:"owner"`
	TotalChipAmount decimal.Decimal                          `json:"total_chip_amount"`
	TotalChipValue  decimal.Decimal                          `json:"total_chip_value"`
	OneDay          *GetStakerProfitChangesSinceResponseData `json:"one_day"`
	OneWeek         *GetStakerProfitChangesSinceResponseData `json:"one_week"`
	OneMonth        *GetStakerProfitChangesSinceResponseData `json:"one_month"`
}

type GetStakerProfitChangesSinceResponseData struct {
	Date            time.Time       `json:"date"`
	TotalChipAmount decimal.Decimal `json:"total_chip_amount"`
	TotalChipValue  decimal.Decimal `json:"total_chip_value"`
	ProfitAndLoss   decimal.Decimal `json:"profit_and_loss"`
}

type GetStakingStatRequest struct {
	Address common.Address `param:"staker_address" validate:"required"`
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

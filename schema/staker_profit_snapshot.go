package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type StakerProfitSnapshot struct {
	Date             time.Time       `json:"date"`
	EpochID          uint64          `json:"epochID"`
	OwnerAddress     common.Address  `json:"ownerAddress"`
	TotalChipAmounts decimal.Decimal `json:"totalChipAmounts"`
	TotalChipValues  decimal.Decimal `json:"totalChipValues"`
	ID               uint64          `json:"-"`
	CreatedAt        time.Time       `json:"-"`
	UpdatedAt        time.Time       `json:"-"`
}

type StakerProfitSnapshotsQuery struct {
	Cursor       *string         `json:"cursor"`
	Limit        int             `json:"limit"`
	OwnerAddress *common.Address `json:"ownerAddress"`
	EpochID      *uint64         `json:"epochID"`
}

package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type StakerProfitSnapshot struct {
	Date            time.Time       `json:"date"`
	EpochID         uint64          `json:"epoch_id"`
	OwnerAddress    common.Address  `json:"owner_address"`
	TotalChipAmount decimal.Decimal `json:"total_chip_amount"`
	TotalChipValue  decimal.Decimal `json:"total_chip_value"`
	ID              uint64          `json:"-"`
	CreatedAt       time.Time       `json:"-"`
	UpdatedAt       time.Time       `json:"-"`
}

type StakerProfitSnapshotsQuery struct {
	Cursor       *string         `json:"cursor"`
	Limit        *int            `json:"limit"`
	OwnerAddress *common.Address `json:"owner_address"`
	EpochID      *uint64         `json:"epoch_id"`
	Dates        []time.Time     `json:"dates"`
	BeforeDate   *time.Time      `json:"before_date"`
	AfterDate    *time.Time      `json:"after_date"`
}

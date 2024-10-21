package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type StakerCumulativeEarningSnapshot struct {
	ID                uint64          `json:"id"`
	Date              time.Time       `json:"date"`
	EpochID           uint64          `json:"epoch_id"`
	OwnerAddress      common.Address  `json:"owner_address"`
	CumulativeEarning decimal.Decimal `json:"cumulative_earning"`
	CreatedAt         time.Time       `json:"-"`
	UpdatedAt         time.Time       `json:"-"`
}

type StakerCumulativeEarningSnapshotsQuery struct {
	Cursor       *string         `json:"cursor"`
	Limit        *int            `json:"limit"`
	OwnerAddress *common.Address `json:"owner_address"`
	EpochID      *uint64         `json:"epoch_id"`
	EpochIDs     []uint64        `json:"epoch_ids"`
	Dates        []time.Time     `json:"dates"`
	BeforeDate   *time.Time      `json:"before_date"`
	AfterDate    *time.Time      `json:"after_date"`
}

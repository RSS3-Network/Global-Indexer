package schema

import (
	"time"

	"github.com/shopspring/decimal"
)

type EpochAPYSnapshot struct {
	Date      time.Time       `json:"date"`
	EpochID   uint64          `json:"epoch_id"`
	APY       decimal.Decimal `json:"apy"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
}

type EpochAPYSnapshotQuery struct {
	EpochID *uint64
	Limit   *int
}

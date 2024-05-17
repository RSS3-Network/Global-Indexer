package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type OperatorProfitSnapshot struct {
	Date          time.Time       `json:"date"`
	EpochID       uint64          `json:"epoch_id"`
	Operator      common.Address  `json:"operator"`
	OperationPool decimal.Decimal `json:"operation_pool"`
	ID            uint64          `json:"-"`
	CreatedAt     time.Time       `json:"-"`
	UpdatedAt     time.Time       `json:"-"`
}

type OperatorProfitSnapshotsQuery struct {
	Operator   *common.Address `json:"operator"`
	Limit      *int            `json:"limit"`
	Cursor     *string         `json:"cursor"`
	BeforeDate *time.Time      `json:"before_date"`
	AfterDate  *time.Time      `json:"after_date"`
	Dates      []time.Time     `json:"dates"`
}

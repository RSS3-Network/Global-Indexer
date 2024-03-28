package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type OperatorProfitSnapshot struct {
	Date          time.Time       `json:"date"`
	EpochID       uint64          `json:"epochID"`
	Operator      common.Address  `json:"operator"`
	OperationPool decimal.Decimal `json:"operationPool"`
	ID            uint64          `json:"-"`
	CreatedAt     time.Time       `json:"-"`
	UpdatedAt     time.Time       `json:"-"`
}

type OperatorProfitSnapshotsQuery struct {
	Operator   *common.Address `json:"operator"`
	Limit      *int            `json:"limit"`
	Cursor     *string         `json:"cursor"`
	BeforeDate *time.Time      `json:"BeforeDate"`
	AfterDate  *time.Time      `json:"AfterDate"`
}

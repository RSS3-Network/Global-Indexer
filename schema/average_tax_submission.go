package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type AverageTaxSubmission struct {
	ID              uint64          `json:"id"`
	EpochID         uint64          `json:"epoch_id"`
	TransactionHash common.Hash     `json:"transaction_hash"`
	AverageTax      decimal.Decimal `json:"average_tax"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type AverageTaxSubmissionQuery struct {
	EpochID *uint64 `json:"epoch_id"`
	Limit   *int    `json:"limit"`
}

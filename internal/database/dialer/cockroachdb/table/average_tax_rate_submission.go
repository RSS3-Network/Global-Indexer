package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
)

type AverageTaxRateSubmission struct {
	ID              uint64          `gorm:"id"`
	EpochID         uint64          `gorm:"epoch_id"`
	AverageTaxRate  decimal.Decimal `gorm:"average_tax_rate"`
	TransactionHash string          `gorm:"transaction_hash"`
	CreatedAt       time.Time       `gorm:"created_at"`
	UpdatedAt       time.Time       `gorm:"updated_at"`
}

func (a *AverageTaxRateSubmission) TableName() string {
	return "average_tax_rate_submissions"
}

func (a *AverageTaxRateSubmission) Import(submission *schema.AverageTaxRateSubmission) error {
	a.EpochID = submission.EpochID
	a.AverageTaxRate = submission.AverageTaxRate
	a.CreatedAt = submission.CreatedAt
	a.UpdatedAt = submission.UpdatedAt
	a.TransactionHash = submission.TransactionHash.String()

	return nil
}

func (a *AverageTaxRateSubmission) Export() (*schema.AverageTaxRateSubmission, error) {
	return &schema.AverageTaxRateSubmission{
		ID:              a.ID,
		EpochID:         a.EpochID,
		AverageTaxRate:  a.AverageTaxRate,
		TransactionHash: common.HexToHash(a.TransactionHash),
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}, nil
}

type AverageTaxSubmissions []AverageTaxRateSubmission

func (a *AverageTaxSubmissions) Import(submissions []*schema.AverageTaxRateSubmission) error {
	for _, submission := range submissions {
		var imported AverageTaxRateSubmission

		if err := imported.Import(submission); err != nil {
			return err
		}

		*a = append(*a, imported)
	}

	return nil
}

func (a *AverageTaxSubmissions) Export() ([]*schema.AverageTaxRateSubmission, error) {
	exported := make([]*schema.AverageTaxRateSubmission, 0)

	for _, submission := range *a {
		exportedSubmission, err := submission.Export()
		if err != nil {
			return nil, err
		}

		exported = append(exported, exportedSubmission)
	}

	return exported, nil
}

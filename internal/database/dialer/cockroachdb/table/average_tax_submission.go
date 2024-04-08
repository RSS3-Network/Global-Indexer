package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/shopspring/decimal"
)

type AverageTaxSubmission struct {
	ID              uint64          `gorm:"id"`
	EpochID         uint64          `gorm:"epoch_id"`
	AverageTax      decimal.Decimal `gorm:"average_tax"`
	TransactionHash string          `gorm:"transaction_hash"`
	CreatedAt       time.Time       `gorm:"created_at"`
	UpdatedAt       time.Time       `gorm:"updated_at"`
}

func (a *AverageTaxSubmission) TableName() string {
	return "average_tax_submissions"
}

func (a *AverageTaxSubmission) Import(submission *schema.AverageTaxSubmission) error {
	a.EpochID = submission.EpochID
	a.AverageTax = submission.AverageTax
	a.CreatedAt = submission.CreatedAt
	a.UpdatedAt = submission.UpdatedAt
	a.TransactionHash = submission.TransactionHash.String()

	return nil
}

func (a *AverageTaxSubmission) Export() (*schema.AverageTaxSubmission, error) {
	return &schema.AverageTaxSubmission{
		ID:              a.ID,
		EpochID:         a.EpochID,
		AverageTax:      a.AverageTax,
		TransactionHash: common.HexToHash(a.TransactionHash),
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}, nil
}

type AverageTaxSubmissions []AverageTaxSubmission

func (a *AverageTaxSubmissions) Import(submissions []*schema.AverageTaxSubmission) error {
	for _, submission := range submissions {
		var imported AverageTaxSubmission

		if err := imported.Import(submission); err != nil {
			return err
		}

		*a = append(*a, imported)
	}

	return nil
}

func (a *AverageTaxSubmissions) Export() ([]*schema.AverageTaxSubmission, error) {
	exported := make([]*schema.AverageTaxSubmission, 0)

	for _, submission := range *a {
		exportedSubmission, err := submission.Export()
		if err != nil {
			return nil, err
		}

		exported = append(exported, exportedSubmission)
	}

	return exported, nil
}

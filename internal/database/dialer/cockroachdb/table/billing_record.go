package table

import (
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	gormSchema "gorm.io/gorm/schema"
)

var (
	_ gormSchema.Tabler = (*BillingRecordDeposited)(nil)
	_ gormSchema.Tabler = (*BillingRecordWithdrawal)(nil)
	_ gormSchema.Tabler = (*BillingRecordCollected)(nil)
)

type BillingRecordDeposited struct {
	schema.BillingRecordBase
}

type BillingRecordWithdrawal struct {
	schema.BillingRecordBase

	Fee float64
}

type BillingRecordCollected struct {
	schema.BillingRecordBase
}

func (r *BillingRecordDeposited) TableName() string {
	return "billing.record.deposited"
}

func (r *BillingRecordWithdrawal) TableName() string {
	return "billing.record.withdrawn"
}

func (r *BillingRecordCollected) TableName() string {
	return "billing.record.collected"
}

func (r *BillingRecordDeposited) Import(billingRecord schema.BillingRecordDeposited) error {
	r.TxHash = billingRecord.TxHash
	r.Index = billingRecord.Index
	r.BlockTimestamp = billingRecord.BlockTimestamp
	r.User = billingRecord.User
	r.Amount = billingRecord.Amount

	return nil
}

func (r *BillingRecordDeposited) Export() (*schema.BillingRecordDeposited, error) {
	billingRecord := schema.BillingRecordDeposited{
		BillingRecordBase: schema.BillingRecordBase{
			TxHash:         r.TxHash,
			Index:          r.Index,
			BlockTimestamp: r.BlockTimestamp,
			User:           r.User,
			Amount:         r.Amount,
		},
	}

	return &billingRecord, nil
}

func (r *BillingRecordWithdrawal) Import(billingRecord schema.BillingRecordWithdrawal) error {
	r.TxHash = billingRecord.TxHash
	r.Index = billingRecord.Index
	r.BlockTimestamp = billingRecord.BlockTimestamp
	r.User = billingRecord.User
	r.Amount = billingRecord.Amount
	r.Fee = billingRecord.Fee

	return nil
}

func (r *BillingRecordWithdrawal) Export() (*schema.BillingRecordWithdrawal, error) {
	billingRecord := schema.BillingRecordWithdrawal{
		BillingRecordBase: schema.BillingRecordBase{
			TxHash:         r.TxHash,
			Index:          r.Index,
			BlockTimestamp: r.BlockTimestamp,
			User:           r.User,
			Amount:         r.Amount,
		},
		Fee: r.Fee,
	}

	return &billingRecord, nil
}

func (r *BillingRecordCollected) Import(billingRecord schema.BillingRecordCollected) error {
	r.TxHash = billingRecord.TxHash
	r.Index = billingRecord.Index
	r.BlockTimestamp = billingRecord.BlockTimestamp
	r.User = billingRecord.User
	r.Amount = billingRecord.Amount

	return nil
}

func (r *BillingRecordCollected) Export() (*schema.BillingRecordCollected, error) {
	billingRecord := schema.BillingRecordCollected{
		BillingRecordBase: schema.BillingRecordBase{
			TxHash:         r.TxHash,
			Index:          r.Index,
			BlockTimestamp: r.BlockTimestamp,
			User:           r.User,
			Amount:         r.Amount,
		},
	}

	return &billingRecord, nil
}

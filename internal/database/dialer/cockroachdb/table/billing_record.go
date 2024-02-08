package table

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/shopspring/decimal"
	gormSchema "gorm.io/gorm/schema"
	"time"
)

var (
	_ gormSchema.Tabler = (*BillingRecordDeposited)(nil)
	_ gormSchema.Tabler = (*BillingRecordWithdrawal)(nil)
	_ gormSchema.Tabler = (*BillingRecordCollected)(nil)
)

type BillingRecordBase struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	TxHash         common.Hash `gorm:"primaryKey;type:bytea;column:tx_hash"`
	Index          uint        `gorm:"column:index"`
	BlockTimestamp time.Time   `gorm:"index;column:block_timestamp"`

	User   common.Address  `gorm:"type:bytea;column:user"`
	Amount decimal.Decimal `gorm:"column:amount"`
}

type BillingRecordDeposited struct {
	BillingRecordBase
}

type BillingRecordWithdrawal struct {
	BillingRecordBase

	Fee decimal.Decimal
}

type BillingRecordCollected struct {
	BillingRecordBase
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
	r.Amount = decimal.NewFromBigInt(billingRecord.Amount, 0)

	return nil
}

func (r *BillingRecordDeposited) Export() (*schema.BillingRecordDeposited, error) {
	billingRecord := schema.BillingRecordDeposited{
		BillingRecordBase: schema.BillingRecordBase{
			TxHash:         r.TxHash,
			Index:          r.Index,
			BlockTimestamp: r.BlockTimestamp,
			User:           r.User,
			Amount:         r.Amount.BigInt(),
		},
	}

	return &billingRecord, nil
}

func (r *BillingRecordWithdrawal) Import(billingRecord schema.BillingRecordWithdrawal) error {
	r.TxHash = billingRecord.TxHash
	r.Index = billingRecord.Index
	r.BlockTimestamp = billingRecord.BlockTimestamp
	r.User = billingRecord.User
	r.Amount = decimal.NewFromBigInt(billingRecord.Amount, 0)
	r.Fee = decimal.NewFromBigInt(billingRecord.Fee, 0)

	return nil
}

func (r *BillingRecordWithdrawal) Export() (*schema.BillingRecordWithdrawal, error) {
	billingRecord := schema.BillingRecordWithdrawal{
		BillingRecordBase: schema.BillingRecordBase{
			TxHash:         r.TxHash,
			Index:          r.Index,
			BlockTimestamp: r.BlockTimestamp,
			User:           r.User,
			Amount:         r.Amount.BigInt(),
		},
		Fee: r.Fee.BigInt(),
	}

	return &billingRecord, nil
}

func (r *BillingRecordCollected) Import(billingRecord schema.BillingRecordCollected) error {
	r.TxHash = billingRecord.TxHash
	r.Index = billingRecord.Index
	r.BlockTimestamp = billingRecord.BlockTimestamp
	r.User = billingRecord.User
	r.Amount = decimal.NewFromBigInt(billingRecord.Amount, 0)

	return nil
}

func (r *BillingRecordCollected) Export() (*schema.BillingRecordCollected, error) {
	billingRecord := schema.BillingRecordCollected{
		BillingRecordBase: schema.BillingRecordBase{
			TxHash:         r.TxHash,
			Index:          r.Index,
			BlockTimestamp: r.BlockTimestamp,
			User:           r.User,
			Amount:         r.Amount.BigInt(),
		},
	}

	return &billingRecord, nil
}

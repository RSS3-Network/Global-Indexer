package table

import (
	gormSchema "gorm.io/gorm/schema"
	"math/big"
	"time"
)

var (
	_ gormSchema.Tabler = (*BillingRecordDeposited)(nil)
	_ gormSchema.Tabler = (*BillingRecordWithdrawn)(nil)
	_ gormSchema.Tabler = (*BillingRecordCollected)(nil)
)

type BillingRecordBase struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	TxHash         string `gorm:"primaryKey"`
	Index          uint
	BlockTimestamp time.Time `gorm:"index"`

	User   string
	Amount float64
}

type BillingRecordDeposited struct {
	BillingRecordBase
}

type BillingRecordWithdrawn struct {
	BillingRecordBase

	Fee *big.Int
}

type BillingRecordCollected struct {
	BillingRecordBase
}

func (r *BillingRecordDeposited) TableName() string {
	return "billing.record.deposited"
}

func (r *BillingRecordWithdrawn) TableName() string {
	return "billing.record.withdrawn"
}

func (r *BillingRecordCollected) TableName() string {
	return "billing.record.collected"
}

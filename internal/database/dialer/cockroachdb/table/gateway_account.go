package table

import (
	"gorm.io/gorm"
	gormSchema "gorm.io/gorm/schema"
	"time"
)

var (
	_ gormSchema.Tabler = (*GatewayAccount)(nil)
)

type GatewayAccount struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Address     string  `gorm:"primaryKey;column:address"`
	RuLimit     int64   `gorm:"column:ru_limit"`
	IsPaused    bool    `gorm:"column:is_paused"`
	BillingRate float64 `gorm:"column:billing_rate"`

	Keys []GatewayKey `gorm:"foreignKey:AccountAddress"` // Has many
}

func (r *GatewayAccount) TableName() string {
	return "gateway.account"
}

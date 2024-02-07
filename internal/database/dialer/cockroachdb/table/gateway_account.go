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

	Address     string `gorm:"primaryKey"`
	RuLimit     int64
	IsPaused    bool
	BillingRate float64

	Keys []GatewayKey `gorm:"foreignKey:AccountAddress"` // Has many
}

func (r *GatewayAccount) TableName() string {
	return "gateway.account"
}

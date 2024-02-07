package table

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	gormSchema "gorm.io/gorm/schema"
	"time"
)

var (
	_ gormSchema.Tabler = (*GatewayKey)(nil)
)

type GatewayKey struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Key uuid.UUID `gorm:"primaryKey"`

	RuUsedTotal     int64
	RuUsedCurrent   int64
	ApiCallsTotal   int64
	ApiCallsCurrent int64

	Name string

	AccountAddress  string                  `gorm:"index"`          // Foreign key of GatewayAccount
	ConsumptionLogs []GatewayConsumptionLog `gorm:"foreignKey:Key"` // Has many
}

func (r *GatewayKey) TableName() string {
	return "gateway.key"
}

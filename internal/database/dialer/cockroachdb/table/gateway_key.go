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

	Key uuid.UUID `gorm:"primaryKey;column:key"`

	RuUsedTotal     int64 `gorm:"column:ru_used_total"`
	RuUsedCurrent   int64 `gorm:"column:ru_used_current"`
	ApiCallsTotal   int64 `gorm:"column:api_calls_total"`
	ApiCallsCurrent int64 `gorm:"column:api_calls_current"`

	Name string `gorm:"column:name"`

	AccountAddress  string                  `gorm:"index"`          // Foreign key of GatewayAccount
	ConsumptionLogs []GatewayConsumptionLog `gorm:"foreignKey:Key"` // Has many
}

func (r *GatewayKey) TableName() string {
	return "gateway.key"
}

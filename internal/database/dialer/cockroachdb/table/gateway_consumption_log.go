package table

import (
	"github.com/google/uuid"
	gormSchema "gorm.io/gorm/schema"
	"time"
)

var (
	_ gormSchema.Tabler = (*GatewayConsumptionLog)(nil)
)

type GatewayConsumptionLog struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ConsumptionDate time.Time `gorm:"index;column:consumption_date"`
	RuUsed          int64     `gorm:"column:ru_used"`
	ApiCalls        int64     `gorm:"column:api_calls"`

	Key uuid.UUID `gorm:"index"` // Foreign key of GatewayKey
}

func (r *GatewayConsumptionLog) TableName() string {
	return "gateway.consumption_log"
}

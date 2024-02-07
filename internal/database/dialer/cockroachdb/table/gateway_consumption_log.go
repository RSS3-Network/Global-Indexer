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

	ConsumptionDate time.Time `gorm:"index"`
	RuUsed          int64
	ApiCalls        int64

	Key uuid.UUID `gorm:"index"` // Foreign key of GatewayKey
}

func (r *GatewayConsumptionLog) TableName() string {
	return "gateway.consumption_log"
}

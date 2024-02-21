package table

import (
	"time"

	gormSchema "gorm.io/gorm/schema"
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
	APICalls        int64     `gorm:"column:api_calls"`

	KeyID uint64     `gorm:"index;column:key_id"` // Foreign key of GatewayKey
	Key   GatewayKey `gorm:"foreignKey:KeyID"`
}

func (r *GatewayConsumptionLog) TableName() string {
	return "gateway.consumption_log"
}

package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
	gormSchema "gorm.io/gorm/schema"
)

var (
	_ gormSchema.Tabler = (*GatewayKey)(nil)
)

type GatewayKey struct {
	ID        uint64 `gorm:"primaryKey;column:id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Key uuid.UUID `gorm:"uniqueIndex;column:key"`

	RuUsedTotal     int64 `gorm:"column:ru_used_total"`
	RuUsedCurrent   int64 `gorm:"index;column:ru_used_current"`
	APICallsTotal   int64 `gorm:"column:api_calls_total"`
	APICallsCurrent int64 `gorm:"column:api_calls_current"`

	Name string `gorm:"column:name"`

	AccountAddress common.Address `gorm:"index;type:bytea;column:account_address"` // Foreign key of GatewayAccount
	Account        GatewayAccount `gorm:"foreignKey:AccountAddress"`
}

func (r *GatewayKey) TableName() string {
	return "gateway.key"
}

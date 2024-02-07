package table

import (
	gormSchema "gorm.io/gorm/schema"
	"time"
)

var (
	_ gormSchema.Tabler = (*GatewayPendingWithdrawRequest)(nil)
)

type GatewayPendingWithdrawRequest struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	Amount float64

	AccountAddress string         // Foreign key of GatewayAccount
	Account        GatewayAccount `gorm:"foreignKey:AccountAddress"` // Belongs to GatewayAccount
}

func (r *GatewayPendingWithdrawRequest) TableName() string {
	return "gateway.pending_withdraw_request"
}

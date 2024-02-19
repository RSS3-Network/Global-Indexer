package table

import (
	"github.com/ethereum/go-ethereum/common"
	gormSchema "gorm.io/gorm/schema"
	"time"
)

var (
	_ gormSchema.Tabler = (*GatewayPendingWithdrawRequest)(nil)
)

type GatewayPendingWithdrawRequest struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	Amount float64 `gorm:"column:amount"`

	AccountAddress common.Address `gorm:"primarykey;type:bytea;column:account_address"` // Foreign key of GatewayAccount
	Account        GatewayAccount `gorm:"foreignKey:AccountAddress"`                    // Belongs to GatewayAccount
}

func (r *GatewayPendingWithdrawRequest) TableName() string {
	return "gateway.pending_withdraw_request"
}

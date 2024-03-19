package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type NodeMinTokensToStakeSnapshot struct {
	Date             time.Time       `json:"date"`
	EpochID          uint64          `json:"epochID"`
	NodeAddress      common.Address  `json:"nodeAddress"`
	MinTokensToStake decimal.Decimal `json:"minTokensToStake"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
	ID               uint64          `json:"-"`
}

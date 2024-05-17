package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type NodeMinTokensToStakeSnapshot struct {
	Date             time.Time       `json:"date"`
	EpochID          uint64          `json:"epoch_id"`
	NodeAddress      common.Address  `json:"node_address"`
	MinTokensToStake decimal.Decimal `json:"min_tokens_to_stake"`
	ID               uint64          `json:"-"`
	CreatedAt        time.Time       `json:"-"`
	UpdatedAt        time.Time       `json:"-"`
}

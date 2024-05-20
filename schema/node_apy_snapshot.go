package schema

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type NodeAPYSnapshot struct {
	ID          uint64          `json:"id"`
	Date        time.Time       `json:"date"`
	EpochID     uint64          `json:"epoch_id"`
	NodeAddress common.Address  `json:"node_address"`
	APY         decimal.Decimal `json:"apy"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type NodeAPYSnapshotQuery struct {
	NodeAddress *common.Address
}

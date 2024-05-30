package schema

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type StakeChipImporter interface {
	Import(stakeChip StakeChip) error
}

type StakeChipExporter interface {
	Export() (*StakeChip, error)
}

type StakeChipTransformer interface {
	StakeChipImporter
	StakeChipExporter
}

type StakeChip struct {
	ID             *big.Int        `json:"id"`
	Owner          common.Address  `json:"owner"`
	Node           common.Address  `json:"node"`
	Value          decimal.Decimal `json:"value"`
	LatestValue    decimal.Decimal `json:"latest_value,omitempty"`
	Metadata       json.RawMessage `json:"metadata"`
	BlockNumber    *big.Int        `json:"block_number"`
	BlockTimestamp uint64          `json:"block_timestamp"`
}

type StakeChipQuery struct {
	ID *big.Int `query:"id"`
}

type StakeChipsQuery struct {
	Cursor        *big.Int
	IDs           []*big.Int
	Node          *common.Address
	Owner         *common.Address
	Limit         *int
	DistinctOwner bool
	BlockNumber   *big.Int
}

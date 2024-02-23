package schema

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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
	Metadata       json.RawMessage `json:"metadata"`
	BlockNumber    *big.Int        `json:"blockNumber"`
	BlockTimestamp uint64          `json:"blockTimestamp"`
}

type StakeChipQuery struct {
	ID *big.Int `query:"id"`
}

type StakeChipsQuery struct {
	Cursor *big.Int
	IDs    []*big.Int
	Node   *common.Address
	Owner  *common.Address
	Limit  int
}

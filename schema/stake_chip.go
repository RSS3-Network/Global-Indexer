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

type StakeChipsQuery struct {
	Cursor *common.Address `query:"cursor"`
	Direct bool            `query:"direct"`
	ID     *big.Int        `query:"id"`
	Owner  *common.Address `query:"owner"`
	Node   *common.Address `query:"node"`
}

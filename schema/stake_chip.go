package schema

import (
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
	ID    *big.Int       `json:"id"`
	Owner common.Address `json:"owner"`
	Node  common.Address `json:"node"`
}

type StakeChipsQuery struct {
	ID    *big.Int        `query:"id"`
	Owner *common.Address `query:"owner"`
	Node  *common.Address `query:"node"`
}

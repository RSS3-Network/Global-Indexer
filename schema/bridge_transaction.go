package schema

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type BridgeTransactionType string

const (
	BridgeTransactionTypeDeposit  BridgeTransactionType = "deposit"
	BridgeTransactionTypeWithdraw BridgeTransactionType = "withdraw"
)

type BridgeTransactionImporter interface {
	Import(bridgeTransaction BridgeTransaction) error
}

type BridgeTransactionExporter interface {
	Export() (*BridgeTransaction, error)
}

type BridgeTransactionTransformer interface {
	BridgeTransactionImporter
	BridgeTransactionExporter
}

type BridgeTransaction struct {
	ID             common.Hash           `json:"id"`
	Type           BridgeTransactionType `json:"type"`
	Sender         common.Address        `json:"sender"`
	Receiver       common.Address        `json:"receiver"`
	TokenAddressL1 *common.Address       `json:"tokenAddressL1"`
	TokenAddressL2 *common.Address       `json:"tokenAddressL2"`
	TokenValue     *big.Int              `json:"tokenValue"`
	TokenDecimal   uint                  `json:"tokenDecimal"`
	Data           string                `json:"data"`
}

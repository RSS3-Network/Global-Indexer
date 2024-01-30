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
	Data           string                `json:"data"`
}

type BridgeTransactionQuery struct {
	ID       *common.Hash    `query:"id"`
	Sender   *common.Address `query:"address"`
	Receiver *common.Address `query:"receiver"`
	Address  *common.Address `query:"address"`
	Type     *string         `query:"type"`
}

type BridgeTransactionsQuery struct {
	ID       *common.Hash    `query:"id"`
	Sender   *common.Address `query:"address"`
	Receiver *common.Address `query:"receiver"`
	Address  *common.Address `query:"address"`
	Type     *string         `query:"type"`
}

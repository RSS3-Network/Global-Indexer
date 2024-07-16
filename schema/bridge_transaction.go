package schema

import (
	"math/big"
	"time"

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
	ID               common.Hash           `json:"id"`
	Type             BridgeTransactionType `json:"type"`
	Sender           common.Address        `json:"sender"`
	Receiver         common.Address        `json:"receiver"`
	TokenAddressL1   *common.Address       `json:"token_address_l1"`
	TokenAddressL2   *common.Address       `json:"token_address_l2"`
	TokenValue       *big.Int              `json:"token_value"`
	Data             string                `json:"data"`
	ChainID          uint64                `json:"chain_id"`
	BlockTimestamp   time.Time             `json:"block_timestamp"`
	BlockNumber      uint64                `json:"block_number"`
	TransactionIndex uint                  `json:"transaction_index"`
	Finalized        bool                  `json:"finalized"`
}

type BridgeTransactionQuery struct {
	ID       *common.Hash           `query:"id"`
	Sender   *common.Address        `query:"sender"`
	Receiver *common.Address        `query:"receiver"`
	Address  *common.Address        `query:"address"`
	Type     *BridgeTransactionType `query:"type"`
}

type BridgeTransactionsQuery struct {
	Cursor   *common.Hash           `query:"cursor"`
	ID       *common.Hash           `query:"id"`
	Sender   *common.Address        `query:"sender"`
	Receiver *common.Address        `query:"receiver"`
	Address  *common.Address        `query:"address"`
	Type     *BridgeTransactionType `query:"type"`
}

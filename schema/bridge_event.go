package schema

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BridgeEventType string

const (
	BridgeEventTypeDepositInitialized BridgeEventType = "initialized"
	BridgeEventTypeDepositFinalized   BridgeEventType = "finalized"

	BridgeEventTypeWithdrawalInitialized BridgeEventType = "initialized"
	BridgeEventTypeWithdrawalProved      BridgeEventType = "proved"
	BridgeEventTypeWithdrawalFinalized   BridgeEventType = "finalized"
)

type BridgeEventImporter interface {
	Import(bridgeEvent BridgeEvent) error
}

type BridgeEventExporter interface {
	Export() (*BridgeEvent, error)
}

type BridgeEventTransformer interface {
	BridgeEventImporter
	BridgeEventExporter
}

type BridgeEvent struct {
	ID                common.Hash     `json:"id"`
	Type              BridgeEventType `json:"type"`
	TransactionHash   common.Hash     `json:"transactionHash"`
	TransactionIndex  uint            `json:"transactionIndex"`
	TransactionStatus uint64          `json:"transactionStatus"`
	BlockHash         common.Hash     `json:"blockHash"`
	BlockNumber       *big.Int        `json:"blockNumber"`
	BlockTimestamp    time.Time       `json:"blockTimestamp"`
}

func NewBridgeEvent(id common.Hash, eventType BridgeEventType, header *types.Header, transaction *types.Transaction, receipt *types.Receipt) *BridgeEvent {
	bridgeEvent := BridgeEvent{
		ID:                id,
		Type:              eventType,
		TransactionHash:   transaction.Hash(),
		TransactionIndex:  receipt.TransactionIndex,
		TransactionStatus: receipt.Status,
		BlockHash:         header.Hash(),
		BlockNumber:       header.Number,
		BlockTimestamp:    time.Unix(int64(header.Time), 0),
	}

	return &bridgeEvent
}

package table

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	gormSchema "gorm.io/gorm/schema"
)

var (
	_ gormSchema.Tabler                   = (*BridgeTransaction)(nil)
	_ schema.BridgeTransactionTransformer = (*BridgeTransaction)(nil)
)

type BridgeTransaction struct {
	gorm.Model
	ID               string          `gorm:"column:id;primaryKey"`
	Type             string          `gorm:"column:type;primaryKey"`
	Sender           string          `gorm:"column:sender"`
	Receiver         string          `gorm:"column:receiver"`
	TokenAddressL1   *string         `gorm:"column:token_address_l1"`
	TokenAddressL2   *string         `gorm:"column:token_address_l2"`
	TokenValue       decimal.Decimal `gorm:"column:token_value"`
	Data             string          `gorm:"column:data"`
	ChainID          uint64          `gorm:"column:chain_id"`
	BlockTimestamp   time.Time       `gorm:"column:block_timestamp"`
	BlockNumber      uint64          `gorm:"column:block_number"`
	TransactionIndex uint            `gorm:"column:transaction_index"`
}

func (b *BridgeTransaction) TableName() string {
	return "transactions"
}

func (b *BridgeTransaction) Import(bridgeTransaction schema.BridgeTransaction) error {
	b.ID = bridgeTransaction.ID.String()
	b.Type = string(bridgeTransaction.Type)
	b.Sender = bridgeTransaction.Sender.String()
	b.Receiver = bridgeTransaction.Receiver.String()
	b.TokenAddressL1 = lo.ToPtr(bridgeTransaction.TokenAddressL1.String())
	b.TokenAddressL2 = lo.ToPtr(bridgeTransaction.TokenAddressL2.String())
	b.TokenValue = decimal.NewFromBigInt(bridgeTransaction.TokenValue, 0)
	b.Data = bridgeTransaction.Data
	b.ChainID = bridgeTransaction.ChainID
	b.BlockTimestamp = bridgeTransaction.BlockTimestamp
	b.BlockNumber = bridgeTransaction.BlockNumber
	b.TransactionIndex = bridgeTransaction.TransactionIndex

	return nil
}

func (b *BridgeTransaction) Export() (*schema.BridgeTransaction, error) {
	bridgeTransaction := schema.BridgeTransaction{
		ID:       common.HexToHash(b.ID),
		Type:     schema.BridgeTransactionType(b.Type),
		Sender:   common.HexToAddress(b.Sender),
		Receiver: common.HexToAddress(b.Receiver),
		TokenAddressL1: func(tokenAddress *string) *common.Address {
			if tokenAddress == nil {
				return nil
			}

			return lo.ToPtr(common.HexToAddress(*tokenAddress))
		}(b.TokenAddressL1),
		TokenAddressL2: func(tokenAddress *string) *common.Address {
			if tokenAddress == nil {
				return nil
			}

			return lo.ToPtr(common.HexToAddress(*tokenAddress))
		}(b.TokenAddressL2),
		TokenValue:       b.TokenValue.BigInt(),
		Data:             b.Data,
		ChainID:          b.ChainID,
		BlockTimestamp:   b.BlockTimestamp,
		BlockNumber:      b.BlockNumber,
		TransactionIndex: b.TransactionIndex,
	}

	return &bridgeTransaction, nil
}

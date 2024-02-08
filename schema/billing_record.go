package schema

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/utils"
	"math/big"
	"time"
)

type BillingRecordBase struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	TxHash         common.Hash `gorm:"primaryKey;type:bytea"`
	Index          uint
	BlockTimestamp time.Time `gorm:"index"`

	User   common.Address `gorm:"type:bytea"`
	Amount float64
}

type BillingRecordDeposited struct {
	BillingRecordBase
}

type BillingRecordWithdrawal struct {
	BillingRecordBase

	Fee float64
}

type BillingRecordCollected struct {
	BillingRecordBase
}

func parseBase(header *types.Header, transaction *types.Transaction, receipt *types.Receipt, user common.Address, amount *big.Int) BillingRecordBase {
	amountParsed, _ := utils.ParseAmount(amount).Float64()
	return BillingRecordBase{
		TxHash:         transaction.Hash(),
		Index:          receipt.TransactionIndex,
		BlockTimestamp: time.Unix(int64(header.Time), 0),

		User:   user,
		Amount: amountParsed,
	}
}

func NewBillingRecordDeposited(header *types.Header, transaction *types.Transaction, receipt *types.Receipt, user common.Address, amount *big.Int) *BillingRecordDeposited {
	return &BillingRecordDeposited{
		BillingRecordBase: parseBase(header, transaction, receipt, user, amount),
	}
}

func NewBillingRecordWithdrawal(header *types.Header, transaction *types.Transaction, receipt *types.Receipt, user common.Address, amount *big.Int, fee *big.Int) *BillingRecordWithdrawal {
	feeParsed, _ := utils.ParseAmount(fee).Float64()
	return &BillingRecordWithdrawal{
		BillingRecordBase: parseBase(header, transaction, receipt, user, amount),
		Fee:               feeParsed,
	}
}

func NewBillingRecordCollected(header *types.Header, transaction *types.Transaction, receipt *types.Receipt, user common.Address, amount *big.Int) *BillingRecordCollected {
	return &BillingRecordCollected{
		BillingRecordBase: parseBase(header, transaction, receipt, user, amount),
	}
}

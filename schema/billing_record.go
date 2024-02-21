package schema

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BillingRecordBase struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	TxHash         common.Hash
	Index          uint
	BlockTimestamp time.Time

	User   common.Address
	Amount *big.Int
}

type BillingRecordDeposited struct {
	BillingRecordBase
}

type BillingRecordWithdrawal struct {
	BillingRecordBase

	Fee *big.Int
}

type BillingRecordCollected struct {
	BillingRecordBase
}

func parseBase(header *types.Header, transaction *types.Transaction, receipt *types.Receipt, user common.Address, amount *big.Int) BillingRecordBase {
	return BillingRecordBase{
		TxHash:         transaction.Hash(),
		Index:          receipt.TransactionIndex,
		BlockTimestamp: time.Unix(int64(header.Time), 0),

		User:   user,
		Amount: amount,
	}
}

func NewBillingRecordDeposited(header *types.Header, transaction *types.Transaction, receipt *types.Receipt, user common.Address, amount *big.Int) *BillingRecordDeposited {
	return &BillingRecordDeposited{
		BillingRecordBase: parseBase(header, transaction, receipt, user, amount),
	}
}

func NewBillingRecordWithdrawal(header *types.Header, transaction *types.Transaction, receipt *types.Receipt, user common.Address, amount *big.Int, fee *big.Int) *BillingRecordWithdrawal {
	return &BillingRecordWithdrawal{
		BillingRecordBase: parseBase(header, transaction, receipt, user, amount),
		Fee:               fee,
	}
}

func NewBillingRecordCollected(header *types.Header, transaction *types.Transaction, receipt *types.Receipt, user common.Address, amount *big.Int) *BillingRecordCollected {
	return &BillingRecordCollected{
		BillingRecordBase: parseBase(header, transaction, receipt, user, amount),
	}
}

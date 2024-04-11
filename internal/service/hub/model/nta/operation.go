package nta

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type GetOperatorProfitRequest struct {
	Operator common.Address `param:"operator" validate:"required"`
}

type GetOperatorProfitRepsonseData struct {
	Operator      common.Address                             `json:"operator"`
	OperationPool decimal.Decimal                            `json:"operationPool"`
	OneDay        *GetOperatorProfitChangesSinceResponseData `json:"oneDay"`
	OneWeek       *GetOperatorProfitChangesSinceResponseData `json:"oneWeek"`
	OneMonth      *GetOperatorProfitChangesSinceResponseData `json:"oneMonth"`
}

type GetOperatorProfitChangesSinceResponseData struct {
	Date          time.Time       `json:"date"`
	OperationPool decimal.Decimal `json:"operationPool"`
	PNL           decimal.Decimal `json:"pnl"`
}

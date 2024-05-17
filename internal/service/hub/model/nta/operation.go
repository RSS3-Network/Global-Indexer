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
	OperationPool decimal.Decimal                            `json:"operation_pool"`
	OneDay        *GetOperatorProfitChangesSinceResponseData `json:"one_day"`
	OneWeek       *GetOperatorProfitChangesSinceResponseData `json:"one_week"`
	OneMonth      *GetOperatorProfitChangesSinceResponseData `json:"one_month"`
}

type GetOperatorProfitChangesSinceResponseData struct {
	Date          time.Time       `json:"date"`
	OperationPool decimal.Decimal `json:"operation_pool"`
	ProfitAndLoss decimal.Decimal `json:"profit_and_loss"`
}

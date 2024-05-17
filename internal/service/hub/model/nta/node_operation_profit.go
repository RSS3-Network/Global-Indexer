package nta

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type GetNodeOperationProfitRequest struct {
	NodeAddress common.Address `param:"node_address" validate:"required"`
}

type GetNodeOperationProfitResponse struct {
	NodeAddress   common.Address          `json:"node_address"`
	OperationPool decimal.Decimal         `json:"operation_pool"`
	OneDay        *NodeProfitChangeDetail `json:"one_day"`
	OneWeek       *NodeProfitChangeDetail `json:"one_week"`
	OneMonth      *NodeProfitChangeDetail `json:"one_month"`
}

type NodeProfitChangeDetail struct {
	Date          time.Time       `json:"date"`
	OperationPool decimal.Decimal `json:"operation_pool"`
	ProfitAndLoss decimal.Decimal `json:"profit_and_loss"`
}

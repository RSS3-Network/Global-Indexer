package nta

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/shopspring/decimal"
)

func (n *NTA) GetNodeOperationProfit(c echo.Context) error {
	var request nta.GetNodeOperationProfitRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	node, err := n.stakingContract.GetNode(&bind.CallOpts{}, request.NodeAddress)
	if err != nil {
		return errorx.InternalError(c, fmt.Errorf("get Node from rpc: %w", err))
	}

	data := nta.GetNodeOperationProfitResponse{
		NodeAddress:   request.NodeAddress,
		OperationPool: decimal.NewFromBigInt(node.OperationPoolTokens, 0),
	}

	changes, err := n.findNodeOperationProfitSnapshots(c.Request().Context(), request.NodeAddress, &data)
	if err != nil {
		return errorx.InternalError(c, fmt.Errorf("find operator history profit snapshots: %w", err))
	}

	data.OneDay, data.OneWeek, data.OneMonth = changes[0], changes[1], changes[2]

	return c.JSON(http.StatusOK, nta.Response{
		Data: data,
	})
}

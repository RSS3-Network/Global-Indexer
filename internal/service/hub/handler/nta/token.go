package nta

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type GetTokensSupplyResponse struct {
	TotalSupply decimal.Decimal `json:"total_supply"`
}

func (n *NTA) GetTokenSupply(c echo.Context) error {
	totalSupply, err := n.getTokensTotalSupply()
	if err != nil {
		zap.L().Error("get token total supply", zap.Error(err))
		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: GetTokensSupplyResponse{
			TotalSupply: totalSupply,
		},
	})
}

func (n *NTA) getTokensTotalSupply() (decimal.Decimal, error) {
	totalSupply, err := n.contractGovernanceToken.TotalSupply(nil)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get total supply: %w", err)
	}

	return decimal.NewFromBigInt(totalSupply, -18), nil
}

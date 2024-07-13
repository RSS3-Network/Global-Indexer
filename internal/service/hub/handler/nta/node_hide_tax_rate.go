package nta

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"go.uber.org/zap"
)

func (n *NTA) PostNodeHideTaxRate(c echo.Context) error {
	var request nta.NodeHideTaxRateRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	message := fmt.Sprintf(hideTaxRateMessage, strings.ToLower(request.NodeAddress.String()))

	if err := n.checkSignature(c.Request().Context(), request.NodeAddress, message, request.Signature); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("check signature: %w", err))
	}

	// Cache the hide tax rate status
	if err := n.cacheClient.Set(c.Request().Context(), n.buildNodeHideTaxRateKey(request.NodeAddress), true, 0); err != nil {
		zap.L().Error("cache hide tax value", zap.Error(err))

		return errorx.InternalError(c)
	}

	// If the Node exists and is not a public good Node, update the hide tax rate status
	if node, err := n.getNode(c.Request().Context(), request.NodeAddress); err == nil && !node.IsPublicGood {
		if err := n.databaseClient.UpdateNodesHideTaxRate(c.Request().Context(), request.NodeAddress, true); err != nil {
			zap.L().Error("update node hide tax rate", zap.Error(err))

			return errorx.InternalError(c)
		}
	}

	return c.NoContent(http.StatusOK)
}

func (n *NTA) buildNodeHideTaxRateKey(address common.Address) string {
	return fmt.Sprintf("node::%s::hideTaxRate", strings.ToLower(address.String()))
}

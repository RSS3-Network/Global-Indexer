package nta

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
)

func (n *NTA) PostNodeHideTaxRate(c echo.Context) error {
	var request nta.PostNodeHideTaxRateRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	message := fmt.Sprintf(hideTaxRateMessage, strings.ToLower(request.NodeAddress.String()))

	if err := n.checkSignature(c.Request().Context(), request.NodeAddress, message, hexutil.MustDecode(request.Signature)); err != nil {
		return errorx.BadRequestError(c, fmt.Errorf("check signature: %w", err))
	}

	// Cache the hide tax rate status
	if err := n.cacheClient.Set(c.Request().Context(), n.buildNodeHideTaxRateKey(request.NodeAddress), true); err != nil {
		return errorx.InternalError(c, fmt.Errorf("cache hide tax value: %w", err))
	}

	// If the Node exists, update the hide tax rate status
	if _, err := n.getNode(c.Request().Context(), request.NodeAddress); err == nil {
		if err := n.databaseClient.UpdateNodesHideTaxRate(c.Request().Context(), request.NodeAddress, true); err != nil {
			return errorx.InternalError(c, fmt.Errorf("confirmation to hide tax rate: %w", err))
		}
	}

	return c.NoContent(http.StatusOK)
}

func (n *NTA) buildNodeHideTaxRateKey(address common.Address) string {
	return fmt.Sprintf("node::%s::hideTaxRate", strings.ToLower(address.String()))
}

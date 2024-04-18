package nta

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/distributor"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/protocol-go/schema/filter"
)

func (n *NTA) GetNetworkWorkers(c echo.Context) error {
	var request nta.NetworkRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	nid, err := filter.NetworkString(request.Network)

	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("network: %w", err))
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: distributor.NetworkToWorkersMap[nid],
	})
}

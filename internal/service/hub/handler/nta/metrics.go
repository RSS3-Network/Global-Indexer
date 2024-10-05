package nta

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
)

type DslTotalRequestsResponse struct {
	TotalRequests int64 `json:"total_requests"`
}

func (n *NTA) GetDslTotalRequests(c echo.Context) error {
	totalRequests, err := n.getNodeTotalRequests(c.Request().Context())

	if err != nil {
		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: DslTotalRequestsResponse{
			TotalRequests: totalRequests,
		},
	})
}

func (n *NTA) getNodeTotalRequests(ctx context.Context) (int64, error) {
	stats, err := n.databaseClient.FindNodeStats(ctx, &schema.StatQuery{})

	if err != nil {
		return 0, fmt.Errorf("failed to find node stats: %w", err)
	}

	var totalRequests int64

	for _, stat := range stats {
		totalRequests += stat.TotalRequest
	}

	return totalRequests, nil
}

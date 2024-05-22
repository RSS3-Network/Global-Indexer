package nta

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/creasty/defaults"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (n *NTA) GetNodeEvents(c echo.Context) error {
	var request nta.NodeEventsRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.BadRequestError(c, fmt.Errorf("set default failed: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	events, err := n.databaseClient.FindNodeEvents(c.Request().Context(), request.NodeAddress, request.Cursor, request.Limit)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("get Node events failed", zap.Error(err))

		return errorx.InternalError(c)
	}

	var cursor string

	if len(events) > 0 && len(events) == request.Limit {
		last, _ := lo.Last(events)
		cursor = fmt.Sprintf("%s:%d:%d", last.TransactionHash, last.TransactionIndex, last.LogIndex)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data:   nta.NewNodeEvents(events),
		Cursor: cursor,
	})
}

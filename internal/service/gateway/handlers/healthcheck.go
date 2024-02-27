package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
)

// HealthCheck implements oapi.ServerInterface
func (app *App) HealthCheck(ctx echo.Context, params oapi.HealthCheckParams) error {
	if params.Type == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	switch *params.Type {
	case "liveness":
		return ctx.NoContent(http.StatusOK)
	case "readiness":
		return ctx.NoContent(http.StatusOK)
	default:
		return ctx.NoContent(http.StatusBadRequest)
	}
}

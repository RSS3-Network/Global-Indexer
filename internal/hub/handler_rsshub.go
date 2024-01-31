package hub

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model/response"
)

func (h *Hub) GetRSSHubHandler(c echo.Context) error {
	path := c.Param("*")
	query := c.Request().URL.RawQuery

	data, err := h.routerRSSHubData(c.Request().Context(), path, query)

	if err != nil {
		return response.InternalError(c, err)
	}

	return c.JSONBlob(http.StatusOK, data)
}

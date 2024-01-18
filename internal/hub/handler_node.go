package hub

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-node/config"
)

func (h *Hub) GetNodesHandler(c echo.Context) error {
	var request BatchNodeRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	if err := defaults.Set(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("set default failed: %v", err))
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("validate failed: %v", err))
	}

	node, err := h.getNodes(c.Request().Context(), &request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("get failed: %v", err))
	}

	var cursor string
	if len(node) > 0 && len(node) == request.Limit {
		cursor = node[len(node)-1].Address.String()
	}

	return c.JSON(http.StatusOK, Response{
		Data:   node,
		Cursor: cursor,
	})
}

func (h *Hub) GetNodeHandler(c echo.Context) error {
	var request NodeRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("validate failed: %v", err))
	}

	node, err := h.getNode(c.Request().Context(), request.Address)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("get failed: %v", err))
	}

	return c.JSON(http.StatusOK, Response{
		Data: node,
	})
}

func (h *Hub) GetNodeChallengeHandler(c echo.Context) error {
	var request NodeRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("validate failed: %v", err))
	}

	return c.JSON(http.StatusOK, Response{
		Data: fmt.Sprintf(message, strings.ToLower(request.Address.String())),
	})
}

func (h *Hub) RegisterNodeHandler(c echo.Context) error {
	var request RegisterNodeRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("validate failed: %v", err))
	}

	if err := h.registerNode(c.Request().Context(), &request); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("register node failed: %v", err))
	}

	return c.JSON(http.StatusOK, Response{
		Data: fmt.Sprintf("node registered: %v", request.Address),
	})
}

func (h *Hub) NodeHeartbeatHandler(c echo.Context) error {
	var request NodeHeartbeatRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("validate failed: %v", err))
	}

	// TODO: resolve node heartbeat logic

	return c.JSON(http.StatusOK, Response{
		Data: fmt.Sprintf("node heartbeat: %v", request.Address),
	})
}

type RegisterNodeRequest struct {
	Address   common.Address `json:"address" validate:"required"`
	Signature string         `json:"signature" validate:"required"`
	Endpoint  string         `json:"endpoint" validate:"required"`
	Stream    *config.Stream `json:"stream,omitempty"`
	Config    *config.Node   `json:"config,omitempty"`
}

type NodeHeartbeatRequest struct {
	Address   common.Address `json:"address" validate:"required"`
	Signature string         `json:"signature" validate:"required"`
	Endpoint  string         `json:"endpoint" validate:"required"`
	Timestamp int64          `json:"timestamp" validate:"required"`
}

type NodeRequest struct {
	Address common.Address `param:"id" validate:"required"`
}

type BatchNodeRequest struct {
	Cursor      *string          `query:"cursor"`
	Limit       int              `query:"limit" validate:"min=1,max=50" default:"10"`
	NodeAddress []common.Address `query:"nodeAddress"`
}

type Response struct {
	Data   any    `json:"data"`
	Cursor string `json:"cursor,omitempty"`
}

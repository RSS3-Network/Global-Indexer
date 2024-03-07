package hub

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/samber/lo"
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

	nodes, err := h.getNodes(c.Request().Context(), &request)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("get failed: %v", err))
	}

	var cursor string
	if len(nodes) > 0 && len(nodes) == request.Limit {
		cursor = nodes[len(nodes)-1].Address.String()
	}

	return c.JSON(http.StatusOK, Response{
		Data:   model.NewNodes(nodes, baseURL(c)),
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
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("get failed: %v", err))
	}

	return c.JSON(http.StatusOK, Response{
		Data: model.NewNode(node, baseURL(c)),
	})
}

func (h *Hub) GetNodeEventsHandler(c echo.Context) error {
	var request NodeEventsRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	if err := defaults.Set(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("set default failed: %v", err))
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("validate failed: %v", err))
	}

	events, err := h.databaseClient.FindNodeEvents(c.Request().Context(), request.Address, request.Cursor, request.Limit)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("get failed: %v", err))
	}

	var cursor string

	if len(events) > 0 && len(events) == request.Limit {
		last, _ := lo.Last(events)
		cursor = fmt.Sprintf("%s:%d:%d", last.TransactionHash, last.TransactionIndex, last.LogIndex)
	}

	return c.JSON(http.StatusOK, Response{
		Data:   model.NewNodeEvents(events),
		Cursor: cursor,
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

func (h *Hub) GetNodeAvatarHandler(c echo.Context) error {
	var request NodeRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("validate failed: %v", err))
	}

	avatar, err := h.getNodeAvatar(c.Request().Context(), request.Address)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("get failed: %v", err))
	}

	return c.Blob(http.StatusOK, "image/svg+xml", avatar)
}

func (h *Hub) RegisterNodeHandler(c echo.Context) error {
	var request RegisterNodeRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("validate failed: %v", err))
	}

	ip, err := h.parseRequestIP(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("parse request ip failed: %v", err))
	}

	if err := h.register(c.Request().Context(), &request, ip.String()); err != nil {
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

	ip, err := h.parseRequestIP(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("parse request ip failed: %v", err))
	}

	if err := h.heartbeat(c.Request().Context(), &request, ip.String()); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("heartbeat failed: %v", err))
	}

	return c.JSON(http.StatusOK, Response{
		Data: fmt.Sprintf("node heartbeat: %v", request.Address),
	})
}

func (h *Hub) parseRequestIP(c echo.Context) (net.IP, error) {
	if ip := net.ParseIP(c.RealIP()); ip != nil {
		return ip, nil
	}

	ip, _, err := net.SplitHostPort(c.Request().RemoteAddr)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(ip), nil
}

type NodeRequest struct {
	Address common.Address `param:"id" validate:"required"`
}

type NodeEventsRequest struct {
	Address common.Address `param:"id" validate:"required"`
	Cursor  *string        `query:"cursor"`
	Limit   int            `query:"limit" validate:"min=1,max=100" default:"20"`
}

type RegisterNodeRequest struct {
	Address   common.Address  `json:"address" validate:"required"`
	Signature string          `json:"signature" validate:"required"`
	Endpoint  string          `json:"endpoint" validate:"required"`
	Stream    json.RawMessage `json:"stream,omitempty"`
	Config    json.RawMessage `json:"config,omitempty"`
}

type NodeHeartbeatRequest struct {
	Address   common.Address `json:"address" validate:"required"`
	Signature string         `json:"signature" validate:"required"`
	Endpoint  string         `json:"endpoint" validate:"required"`
	Timestamp int64          `json:"timestamp" validate:"required"`
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

package hub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/model"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/response"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

func (h *Hub) GetNodes(c echo.Context) error {
	var request BatchNodeRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return response.BadRequestError(c, fmt.Errorf("set default failed: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	nodes, err := h.getNodes(c.Request().Context(), &request)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return response.InternalError(c, fmt.Errorf("get nodes: %w", err))
	}

	var cursor string
	if len(nodes) > 0 && len(nodes) == request.Limit {
		cursor = nodes[len(nodes)-1].Address.String()
	}

	// If the score is the same, sort by staking pool size.
	// TODO: Since node's StakingPoolTokens needs to be obtained from vsl.
	//  Now only the nodes of the current page can be sorted.
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Score.Cmp(nodes[j].Score) == 0 {
			iTokens, _ := new(big.Int).SetString(nodes[i].StakingPoolTokens, 10)
			jTokens, _ := new(big.Int).SetString(nodes[j].StakingPoolTokens, 10)

			return iTokens.Cmp(jTokens) > 0
		}

		return nodes[i].Score.Cmp(nodes[j].Score) > 0
	})

	return c.JSON(http.StatusOK, Response{
		Data:   model.NewNodes(nodes, baseURL(c)),
		Cursor: cursor,
	})
}

func (h *Hub) GetNode(c echo.Context) error {
	var request NodeRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	node, err := h.getNode(c.Request().Context(), request.Address)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return response.InternalError(c, fmt.Errorf("get node: %w", err))
	}

	return c.JSON(http.StatusOK, Response{
		Data: model.NewNode(node, baseURL(c)),
	})
}

func (h *Hub) GetNodeEvents(c echo.Context) error {
	var request NodeEventsRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return response.BadRequestError(c, fmt.Errorf("set default failed: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	events, err := h.databaseClient.FindNodeEvents(c.Request().Context(), request.Address, request.Cursor, request.Limit)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return response.InternalError(c, fmt.Errorf("get node events: %w", err))
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

func (h *Hub) GetNodeChallenge(c echo.Context) error {
	var request NodeChallengeRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	switch request.Type {
	case "":
		return c.JSON(http.StatusOK, Response{
			Data: fmt.Sprintf(registerMessage, strings.ToLower(request.Address.String())),
		})
	case "hideTaxRate":
		return c.JSON(http.StatusOK, Response{
			Data: fmt.Sprintf(hideTaxRateMessage, strings.ToLower(request.Address.String())),
		})
	default:
		return response.BadRequestError(c, fmt.Errorf("invalid challenge type: %s", request.Type))
	}
}

func (h *Hub) PostNodeHideTaxRate(c echo.Context) error {
	var request NodeHideTaxRateRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	message := fmt.Sprintf(hideTaxRateMessage, strings.ToLower(request.Address.String()))

	if err := h.checkSignature(c.Request().Context(), request.Address, message, hexutil.MustDecode(request.Signature)); err != nil {
		return response.BadRequestError(c, fmt.Errorf("check signature: %w", err))
	}

	// Cache the hide tax rate status
	if err := h.cacheClient.Set(c.Request().Context(), h.buildNodeHideTaxRateKey(request.Address), true); err != nil {
		return response.InternalError(c, fmt.Errorf("cache hide tax value: %w", err))
	}

	// If the node exists, update the hide tax rate status
	if _, err := h.getNode(c.Request().Context(), request.Address); err == nil {
		if err := h.databaseClient.UpdateNodesHideTaxRate(c.Request().Context(), request.Address, true); err != nil {
			return response.InternalError(c, fmt.Errorf("confirmation to hide tax rate: %w", err))
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h *Hub) GetNodeAvatar(c echo.Context) error {
	var request NodeRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	avatar, err := h.getNodeAvatar(c.Request().Context(), request.Address)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return response.InternalError(c, fmt.Errorf("get node avatar: %w", err))
	}

	return c.Blob(http.StatusOK, "image/svg+xml", avatar)
}

func (h *Hub) RegisterNode(c echo.Context) error {
	var request RegisterNodeRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	ip, err := h.parseRequestIP(c)
	if err != nil {
		return response.InternalError(c, fmt.Errorf("parse request ip: %w", err))
	}

	if err := h.register(c.Request().Context(), &request, ip.String()); err != nil {
		return response.InternalError(c, fmt.Errorf("register failed: %w", err))
	}

	return c.JSON(http.StatusOK, Response{
		Data: fmt.Sprintf("node registered: %v", request.Address),
	})
}

func (h *Hub) NodeHeartbeat(c echo.Context) error {
	var request NodeHeartbeatRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	ip, err := h.parseRequestIP(c)
	if err != nil {
		return response.InternalError(c, fmt.Errorf("parse request ip: %w", err))
	}

	if err := h.heartbeat(c.Request().Context(), &request, ip.String()); err != nil {
		return response.InternalError(c, fmt.Errorf("heartbeat failed: %w", err))
	}

	return c.JSON(http.StatusOK, Response{
		Data: fmt.Sprintf("node heartbeat: %v", request.Address),
	})
}

func (h *Hub) GetOperatorProfit(c echo.Context) error {
	var request GetOperatorProfitRequest

	if err := c.Bind(&request); err != nil {
		return response.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	node, err := h.stakingContract.GetNode(&bind.CallOpts{}, request.Operator)
	if err != nil {
		return response.InternalError(c, fmt.Errorf("get node from rpc: %w", err))
	}

	data := GetOperatorProfitRepsonse{
		Operator:      request.Operator,
		OperationPool: decimal.NewFromBigInt(node.OperationPoolTokens, 0),
	}

	changes, err := h.findOperatorHistoryProfitSnapshots(c.Request().Context(), request.Operator, &data)
	if err != nil {
		return response.InternalError(c, fmt.Errorf("find operator history profit snapshots: %w", err))
	}

	data.OneDay, data.OneWeek, data.OneMonth = changes[0], changes[1], changes[2]

	return c.JSON(http.StatusOK, Response{
		Data: data,
	})
}

func (h *Hub) findOperatorHistoryProfitSnapshots(ctx context.Context, operator common.Address, profit *GetOperatorProfitRepsonse) ([]*GetOperatorProfitChangesSinceResponse, error) {
	if profit == nil {
		return nil, nil
	}

	now := time.Now()
	query := schema.OperatorProfitSnapshotsQuery{
		Operator: lo.ToPtr(operator),
		Dates: []time.Time{
			now.Add(-24 * time.Hour),      // 1 day
			now.Add(-7 * 24 * time.Hour),  // 1 week
			now.Add(-30 * 24 * time.Hour), // 1 month
		},
	}

	snapshots, err := h.databaseClient.FindOperatorProfitSnapshots(ctx, query)
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		return nil, fmt.Errorf("find operator profit snapshots: %w", err)
	}

	data := make([]*GetOperatorProfitChangesSinceResponse, len(query.Dates))

	for _, snapshot := range snapshots {
		if snapshot.OperationPool.IsZero() {
			continue
		}

		var index int

		if snapshot.Date.After(query.Dates[2]) && snapshot.Date.Before(query.Dates[1]) {
			index = 2
		} else if snapshot.Date.After(query.Dates[1]) && snapshot.Date.Before(query.Dates[0]) {
			index = 1
		}

		data[index] = &GetOperatorProfitChangesSinceResponse{
			Date:          snapshot.Date,
			OperationPool: snapshot.OperationPool,
			PNL:           profit.OperationPool.Sub(snapshot.OperationPool).Div(snapshot.OperationPool),
		}
	}

	return data, nil
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

func (h *Hub) buildNodeHideTaxRateKey(address common.Address) string {
	return fmt.Sprintf("node::%s::hideTaxRate", strings.ToLower(address.String()))
}

type NodeRequest struct {
	Address common.Address `param:"id" validate:"required"`
}

type NodeEventsRequest struct {
	Address common.Address `param:"id" validate:"required"`
	Cursor  *string        `query:"cursor"`
	Limit   int            `query:"limit" validate:"min=1,max=100" default:"20"`
}

type NodeChallengeRequest struct {
	Address common.Address `param:"id" validate:"required"`
	Type    string         `query:"type"`
}

type NodeHideTaxRateRequest struct {
	Address   common.Address `param:"id" validate:"required"`
	Signature string         `json:"signature" validate:"required"`
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

type GetOperatorProfitRequest struct {
	Operator common.Address `param:"operator" validate:"required"`
}

type GetOperatorProfitRepsonse struct {
	Operator      common.Address                         `json:"operator"`
	OperationPool decimal.Decimal                        `json:"operationPool"`
	OneDay        *GetOperatorProfitChangesSinceResponse `json:"oneDay"`
	OneWeek       *GetOperatorProfitChangesSinceResponse `json:"oneWeek"`
	OneMonth      *GetOperatorProfitChangesSinceResponse `json:"oneMonth"`
}

type GetOperatorProfitChangesSinceResponse struct {
	Date          time.Time       `json:"date"`
	OperationPool decimal.Decimal `json:"operationPool"`
	PNL           decimal.Decimal `json:"pnl"`
}

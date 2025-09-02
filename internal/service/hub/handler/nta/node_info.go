package nta

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/contract/l2"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (n *NTA) GetNodes(c echo.Context) error {
	var request nta.BatchNodeRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.BadRequestError(c, fmt.Errorf("set default failed: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	nodes, err := n.getNodes(c.Request().Context(), &request)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("get Nodes failed", zap.Error(err))

		return errorx.InternalError(c)
	}

	var cursor string
	if len(nodes) > 0 && len(nodes) == request.Limit {
		cursor = nodes[len(nodes)-1].Address.String()
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data:   nta.NewNodes(nodes, n.baseURL(c)),
		Cursor: cursor,
	})
}

func (n *NTA) GetNode(c echo.Context) error {
	var request nta.NodeRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	node, err := n.getNode(c.Request().Context(), request.Address)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("get Node failed", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: nta.NewNode(node, n.baseURL(c)),
	})
}

func (n *NTA) GetNodeAvatar(c echo.Context) error {
	var request nta.NodeRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	avatar, err := n.getNodeAvatar(c.Request().Context(), request.Address)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		zap.L().Error("get Node avatar failed", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.Blob(http.StatusOK, "image/svg+xml", avatar)
}

func (n *NTA) getNode(ctx context.Context, address common.Address) (*schema.Node, error) {
	node, err := n.databaseClient.FindNode(ctx, address)
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		return nil, fmt.Errorf("get Node %s: %w", address, err)
	}

	if node == nil {
		node = &schema.Node{
			Status: schema.NodeStatusRegistered,
			Avatar: &l2.ChipsTokenMetadata{
				Name: "Node Avatar",
			},
			StakingPoolTokens: decimal.Zero.String(),
		}
	}

	nodeInfo, err := n.stakingContract.GetNode(&bind.CallOpts{}, address)
	if err != nil {
		if node.Type == schema.NodeTypeRSSHub.String() && node.ID.Cmp(big.NewInt(10000)) >= 0 {
			node.ReliabilityScore = decimal.NewFromFloat(0)
			node.Status = schema.NodeStatusRegistered

			return node, nil
		}

		return nil, fmt.Errorf("get Node from chain: %w", err)
	}

	var reliabilityScore decimal.Decimal

	nodeStat, err := n.databaseClient.FindNodeStat(ctx, address)
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		return nil, fmt.Errorf("get Node Stat %s: %w", address, err)
	}

	if nodeStat != nil {
		reliabilityScore = decimal.NewFromFloat(nodeStat.Score)
	}

	if reliabilityScore.IsZero() {
		// set baseline score
		node.ReliabilityScore = setReliabilityBaselineScore(nodeInfo.StakingPoolTokens)
	}

	if nodeInfo.PublicGood {
		publicPool, err := n.stakingContract.GetPublicPool(&bind.CallOpts{})
		if err != nil {
			return nil, fmt.Errorf("get Public Pool from chain: %w", err)
		}

		nodeInfo.TaxRateBasisPoints = publicPool.TaxRateBasisPoints
		nodeInfo.OperationPoolTokens = publicPool.OperationPoolTokens
		nodeInfo.StakingPoolTokens = publicPool.StakingPoolTokens
	}

	node.ID = nodeInfo.NodeId
	node.Address = nodeInfo.Account
	node.Name = nodeInfo.Name
	node.Description = nodeInfo.Description
	node.TaxRateBasisPoints = &nodeInfo.TaxRateBasisPoints
	node.OperationPoolTokens = nodeInfo.OperationPoolTokens.String()
	node.StakingPoolTokens = nodeInfo.StakingPoolTokens.String()
	node.TotalShares = nodeInfo.TotalShares.String()
	node.SlashedTokens = big.NewInt(0).Add(nodeInfo.SlashedStakingPoolTokens, nodeInfo.SlashedOperationPoolTokens).String()
	node.Alpha = nodeInfo.Alpha
	node.ReliabilityScore = reliabilityScore
	node.Status = schema.NodeStatus(nodeInfo.Status)

	return node, nil
}

func (n *NTA) getNodes(ctx context.Context, request *nta.BatchNodeRequest) ([]*schema.Node, error) {
	nodes, err := n.databaseClient.FindNodes(ctx, schema.FindNodesQuery{
		NodeAddresses: request.NodeAddresses,
		Cursor:        request.Cursor,
		Limit:         lo.ToPtr(request.Limit),
	})
	if err != nil {
		return nil, fmt.Errorf("get Nodes: %w", err)
	}

	onChainNodes := lo.Filter(nodes, func(node *schema.Node, _ int) bool {
		return node.ID != nil && node.ID.Cmp(big.NewInt(10000)) < 0
	})

	addresses := lo.Map(onChainNodes, func(node *schema.Node, _ int) common.Address {
		return node.Address
	})

	// Get uncertain node from event.
	uncertainNodeEvents, err := n.databaseClient.FindNodeEvents(ctx, &schema.NodeEventsQuery{
		Finalized: lo.ToPtr(false),
		Type:      lo.ToPtr(schema.NodeEventNodeCreated),
	})
	if err != nil {
		return nil, fmt.Errorf("get Node Events: %w", err)
	}

	for _, event := range uncertainNodeEvents {
		if event.Metadata.NodeCreatedMetadata == nil {
			zap.L().Error("invalid NodeCreatedMetadata", zap.Any("event", event))

			continue
		}

		if _, exists := lo.Find(addresses, func(item common.Address) bool {
			return item == event.Metadata.NodeCreatedMetadata.Address
		}); exists {
			continue
		}

		nodes = append(nodes, &schema.Node{
			Address: event.Metadata.NodeCreatedMetadata.Address,
			ID:      event.Metadata.NodeCreatedMetadata.NodeID,
			Status:  schema.NodeStatusRegistered,
			Avatar: &l2.ChipsTokenMetadata{
				Name: "Node Avatar",
			},
		})

		addresses = append(addresses, event.Metadata.NodeCreatedMetadata.Address)
	}

	// Get node info from VSL.
	nodeInfoList, err := n.stakingContract.GetNodes(&bind.CallOpts{}, addresses)
	if err != nil {
		return nil, fmt.Errorf("get Nodes from chain: %w", err)
	}

	nodeInfoMap := lo.SliceToMap(nodeInfoList, func(node stakingv2.Node) (common.Address, stakingv2.Node) {
		return node.Account, node
	})

	// Get node stats from DB.
	nodeStats, err := n.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		Addresses: addresses,
	})
	if err != nil {
		return nil, fmt.Errorf("get Node Stats: %w", err)
	}

	nodeStatsMap := lo.SliceToMap(nodeStats, func(stat *schema.Stat) (common.Address, float64) {
		return stat.Address, stat.Score
	})

	var publicGoodPool *stakingv2.Node

	for _, node := range nodes {
		if score, exists := nodeStatsMap[node.Address]; exists {
			node.ReliabilityScore = decimal.NewFromFloat(score)
		}

		if node.ReliabilityScore.IsZero() {
			node.ReliabilityScore = decimal.NewFromFloat(0)
			// set baseline score
			if nodeInfo, exists := nodeInfoMap[node.Address]; exists {
				node.ReliabilityScore = setReliabilityBaselineScore(nodeInfo.StakingPoolTokens)
			}
		}

		if nodeInfo, exists := nodeInfoMap[node.Address]; exists {
			node.ID = nodeInfo.NodeId
			node.Name = nodeInfo.Name
			node.IsPublicGood = nodeInfo.PublicGood
			node.Description = nodeInfo.Description
			node.TaxRateBasisPoints = &nodeInfo.TaxRateBasisPoints
			node.OperationPoolTokens = nodeInfo.OperationPoolTokens.String()
			node.StakingPoolTokens = nodeInfo.StakingPoolTokens.String()
			node.TotalShares = nodeInfo.TotalShares.String()
			node.SlashedTokens = big.NewInt(0).Add(nodeInfo.SlashedStakingPoolTokens, nodeInfo.SlashedOperationPoolTokens).String()
			node.Alpha = nodeInfo.Alpha
			// get node status from chain
			node.Status = schema.NodeStatus(nodeInfo.Status)
		}

		if node.IsPublicGood {
			if publicGoodPool == nil {
				publicPool, err := n.stakingContract.GetPublicPool(&bind.CallOpts{})
				if err != nil {
					return nil, fmt.Errorf("get Public Pool from chain: %w", err)
				}

				publicGoodPool = &publicPool
			}

			node.TaxRateBasisPoints = &publicGoodPool.TaxRateBasisPoints
			node.OperationPoolTokens = publicGoodPool.OperationPoolTokens.String()
			node.StakingPoolTokens = publicGoodPool.StakingPoolTokens.String()
		}
	}

	return nodes, nil
}

func setReliabilityBaselineScore(stakingTokens *big.Int) decimal.Decimal {
	staking, _ := big.NewInt(0).Div(stakingTokens, big.NewInt(1e18)).Float64()
	return decimal.NewFromFloat(math.Min(math.Log(staking/100000+1)/math.Log(2), 0.2))
}

func (n *NTA) getNodeAvatar(ctx context.Context, address common.Address) ([]byte, error) {
	avatar, err := n.databaseClient.FindNodeAvatar(ctx, address)
	if err != nil {
		zap.L().Error("get Node avatar failed", zap.Error(err))

		avatar, err = n.buildNodeAvatar(ctx, address)
		if err != nil {
			return nil, fmt.Errorf("get Node avatar %s: %w", address, err)
		}
	}

	data, ok := strings.CutPrefix(avatar.Image, "data:image/svg+xml;base64,")
	if !ok {
		return nil, fmt.Errorf("invalid avatar")
	}

	return base64.StdEncoding.DecodeString(data)
}

// parseEndpoint parses the given endpoint string.
// If it does not start with "https://" or "http://", it returns an error.
// Then, it parses the endpoint URL, ignoring any query parameters.
// It returns the parsed URL as a string.
func (n *NTA) parseEndpoint(_ context.Context, endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parse endpoint: %w", err)
	}

	if (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
		return "", errors.New("invalid endpoint")
	}

	u.ForceQuery = false
	u.Path, u.RawQuery = "", ""

	return u.String(), nil
}

func (n *NTA) parseRequestIP(c echo.Context) (net.IP, error) {
	if ip := net.ParseIP(c.RealIP()); ip != nil {
		return ip, nil
	}

	ip, _, err := net.SplitHostPort(c.Request().RemoteAddr)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(ip), nil
}

func (n *NTA) buildNodeAvatar(_ context.Context, address common.Address) (*l2.ChipsTokenMetadata, error) {
	avatar, err := n.stakingContract.GetNodeAvatar(&bind.CallOpts{}, address)
	if err != nil {
		return nil, fmt.Errorf("get Node avatar from chain: %w", err)
	}

	encodedMetadata, ok := strings.CutPrefix(avatar, "data:application/json;base64,")
	if !ok {
		return nil, fmt.Errorf("invalid avatar: %s", avatar)
	}

	metadata, err := base64.StdEncoding.DecodeString(encodedMetadata)
	if err != nil {
		return nil, fmt.Errorf("decode avatar metadata: %w", err)
	}

	var avatarMetadata l2.ChipsTokenMetadata

	if err = json.Unmarshal(metadata, &avatarMetadata); err != nil {
		return nil, fmt.Errorf("unmarshal avatar metadata: %w", err)
	}

	return &avatarMetadata, nil
}

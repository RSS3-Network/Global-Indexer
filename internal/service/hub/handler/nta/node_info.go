package nta

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
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
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	nodes, err := n.getNodes(c.Request().Context(), &request)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return errorx.InternalError(c, fmt.Errorf("get Nodes: %w", err))
	}

	var cursor string
	if len(nodes) > 0 && len(nodes) == request.Limit {
		cursor = nodes[len(nodes)-1].Address.String()
	}

	// If the ActiveScore is the same, sort by staking pool size.
	// TODO: Since Node's StakingPoolTokens needs to be obtained from vsl.
	//  Now only the Nodes of the current page can be sorted.
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].ActiveScore.Cmp(nodes[j].ActiveScore) == 0 {
			iTokens, _ := new(big.Int).SetString(nodes[i].StakingPoolTokens, 10)
			jTokens, _ := new(big.Int).SetString(nodes[j].StakingPoolTokens, 10)

			return iTokens.Cmp(jTokens) > 0
		}

		return nodes[i].ActiveScore.Cmp(nodes[j].ActiveScore) > 0
	})

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
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	node, err := n.getNode(c.Request().Context(), request.Address)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return errorx.InternalError(c, fmt.Errorf("get Node: %w", err))
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
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	avatar, err := n.getNodeAvatar(c.Request().Context(), request.Address)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNotFound)
		}

		return errorx.InternalError(c, fmt.Errorf("get Node avatar: %w", err))
	}

	return c.Blob(http.StatusOK, "image/svg+xml", avatar)
}

func (n *NTA) getNode(ctx context.Context, address common.Address) (*schema.Node, error) {
	node, err := n.databaseClient.FindNode(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("get Node %s: %w", address, err)
	}

	nodeInfo, err := n.stakingContract.GetNode(&bind.CallOpts{}, address)
	if err != nil {
		return nil, fmt.Errorf("get Node from chain: %w", err)
	}

	var reliabilityScore decimal.Decimal

	nodeStat, err := n.databaseClient.FindNodeStat(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("get Node Stat %s: %w", address, err)
	}

	if nodeStat != nil {
		reliabilityScore = decimal.NewFromFloat(nodeStat.Score)
	}

	node.Name = nodeInfo.Name
	node.Description = nodeInfo.Description
	node.TaxRateBasisPoints = &nodeInfo.TaxRateBasisPoints
	node.OperationPoolTokens = nodeInfo.OperationPoolTokens.String()
	node.StakingPoolTokens = nodeInfo.StakingPoolTokens.String()
	node.TotalShares = nodeInfo.TotalShares.String()
	node.SlashedTokens = nodeInfo.SlashedTokens.String()
	node.Alpha = nodeInfo.Alpha
	node.ReliabilityScore = reliabilityScore

	return node, nil
}

func (n *NTA) getNodes(ctx context.Context, request *nta.BatchNodeRequest) ([]*schema.Node, error) {
	nodes, err := n.databaseClient.FindNodes(ctx, schema.FindNodesQuery{
		NodeAddresses: request.NodeAddress,
		Cursor:        request.Cursor,
		Limit:         lo.ToPtr(request.Limit),
		OrderByScore:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("get Nodes: %w", err)
	}

	addresses := lo.Map(nodes, func(node *schema.Node, _ int) common.Address {
		return node.Address
	})

	// Get node info from VSL.
	nodeInfo, err := n.stakingContract.GetNodes(&bind.CallOpts{}, addresses)
	if err != nil {
		return nil, fmt.Errorf("get Nodes from chain: %w", err)
	}

	nodeInfoMap := lo.SliceToMap(nodeInfo, func(node l2.DataTypesNode) (common.Address, l2.DataTypesNode) {
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

	for _, node := range nodes {
		if score, exists := nodeStatsMap[node.Address]; exists {
			node.ReliabilityScore = decimal.NewFromFloat(score)
		}

		if nodeInfo, exists := nodeInfoMap[node.Address]; exists {
			node.Name = nodeInfo.Name
			node.Description = nodeInfo.Description
			node.TaxRateBasisPoints = &nodeInfo.TaxRateBasisPoints
			node.OperationPoolTokens = nodeInfo.OperationPoolTokens.String()
			node.StakingPoolTokens = nodeInfo.StakingPoolTokens.String()
			node.TotalShares = nodeInfo.TotalShares.String()
			node.SlashedTokens = nodeInfo.SlashedTokens.String()
			node.Alpha = nodeInfo.Alpha
		}
	}

	return nodes, nil
}

func (n *NTA) getNodeAvatar(ctx context.Context, address common.Address) ([]byte, error) {
	avatar, err := n.databaseClient.FindNodeAvatar(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("get Node avatar %s: %w", address, err)
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

package nta

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/go-version"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/ethereum"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (n *NTA) RegisterNode(c echo.Context) error {
	var request nta.RegisterNodeRequest

	ctx := c.Request().Context()

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		zap.L().Error("set default values for request", zap.Error(err))

		return errorx.InternalError(c)
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	// Parse request IP.
	ip, err := n.parseRequestIP(c)
	if err != nil {
		zap.L().Error("parse request ip", zap.Error(err))

		return errorx.InternalError(c)
	}

	// Validate signature.
	if err = n.validateSignature(ctx, request.Address, request.Signature); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validate signature: %w", err))
	}

	// Validate Node info.
	nodeInfo, err := n.stakingContract.GetNode(&bind.CallOpts{}, request.Address)
	if err != nil {
		zap.L().Error("get the Node from VSL", zap.Error(err))

		return errorx.InternalError(c)
	}

	if nodeInfo.Account == ethereum.AddressGenesis {
		return errorx.ValidationFailedError(c, fmt.Errorf("node: %s has not been registered on the VSL", request.Address.String()))
	}

	if !nodeInfo.PublicGood && strings.Compare(nodeInfo.OperationPoolTokens.String(), MinDeposit.String()) < 0 {
		return errorx.ValidationFailedError(c, fmt.Errorf("insufficient operation pool tokens, expected min deposit %s, actual %s", MinDeposit.String(), nodeInfo.OperationPoolTokens.String()))
	}

	// Validate endpoint.
	if err = n.validateEndpoint(ctx, request.Address, request.Type, request.Version, request.Endpoint); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validate endpoint: %w", err))
	}

	// Register Node.
	if err = n.register(ctx, &request, ip.String(), nodeInfo); err != nil {
		zap.L().Error("register failed",
			zap.String("address", request.Address.String()),
			zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: fmt.Sprintf("successfully registered node: %v", request.Address),
	})
}

func (n *NTA) NodeHeartbeat(c echo.Context) error {
	var request nta.NodeHeartbeatRequest

	ctx := c.Request().Context()

	// Validate request.
	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	// Parse request IP.
	ip, err := n.parseRequestIP(c)
	if err != nil {
		zap.L().Error("parse request ip", zap.Error(err))

		return errorx.InternalError(c)
	}

	// Validate signature.
	if err = n.validateSignature(ctx, request.Address, request.Signature); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("check signature: %w", err))
	}

	// Validate Node.
	node, err := n.databaseClient.FindNode(c.Request().Context(), request.Address)
	if err != nil {
		zap.L().Error("find the node",
			zap.String("address", request.Address.String()),
			zap.Error(err))

		return errorx.InternalError(c)
	}

	if node == nil {
		return errorx.BadParamsError(c, fmt.Errorf("node %s not found", request.Address.String()))
	}

	// Validate endpoint.
	if err = n.validateEndpoint(ctx, request.Address, node.Type, node.Version, request.Endpoint); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validate endpoint: %w", err))
	}

	// Save Node heartbeat.
	if err = n.saveHeartbeat(ctx, node, ip.String()); err != nil {
		zap.L().Error("save heartbeat", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: fmt.Sprintf("successfully sent node heartbeat: %v", request.Address),
	})
}

func (n *NTA) RSSHubNodeHeartbeat(c echo.Context) error {
	var request nta.RSSHubNodeHeartbeatRequest

	ctx := c.Request().Context()

	// Validate and parse request
	if err := c.Bind(request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	// Parse request IP
	ip, err := n.parseRequestIP(c)
	if err != nil {
		zap.L().Error("parse request ip", zap.Error(err))
		return errorx.InternalError(c)
	}

	// Set default name
	if request.Name == "" {
		request.Name = ip.String()
	}

	// Check if this is an incomplete node
	if request.Signature == "" || request.Endpoint == "" || request.Address == (common.Address{}) {
		if err := n.handleIncompleteNode(ctx, request, ip.String()); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, nta.Response{
			Data: fmt.Sprintf("successfully sent RSSHub node heartbeat: %v", request.Name),
		})
	}

	// Handle complete node
	if err := n.handleCompleteNode(ctx, request, ip.String()); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: fmt.Sprintf("successfully sent RSSHub node heartbeat: %v", request.Name),
	})
}

// validateSignature validates the signature.
func (n *NTA) validateSignature(ctx context.Context, address common.Address, signature string) error {
	message := fmt.Sprintf(registrationMessage, strings.ToLower(address.String()))

	return n.checkSignature(ctx, address, message, signature)
}

// validateEndpoint validates the endpoint whether it's valid and available.
func (n *NTA) validateEndpoint(ctx context.Context, address common.Address, nodeType, nodeVersion, endpoint string) error {
	if nodeType == schema.NodeTypeAlpha.String() {
		return nil
	}

	var err error
	endpoint, err = n.parseEndpoint(ctx, endpoint)

	if err != nil {
		return fmt.Errorf("failed to parse endpoint: %w", err)
	}

	if err = n.checkAvailable(ctx, nodeVersion, endpoint, address); err != nil {
		return fmt.Errorf("failed to check endpoint available: %w", err)
	}

	return nil
}

// validateRSSHubEndpoint validates the endpoint whether it's valid and available.
func (n *NTA) validateRSSHubEndpoint(ctx context.Context, endpoint, accessToken string) error {
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("invalid RSS endpoint: %w", err)
	}

	baseURL.Path = path.Join(baseURL.Path, "healthz")
	if accessToken != "" {
		query := baseURL.Query()
		query.Set("key", accessToken)
		baseURL.RawQuery = query.Encode()
	}

	body, _, err := n.httpClient.FetchWithMethod(ctx, http.MethodGet, baseURL.String(), "", nil)

	if err != nil {
		return fmt.Errorf("failed to fetch RSS healthz: %w", err)
	}

	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if strings.Contains(strings.ToLower(strings.TrimSpace(string(data))), "ok") {
		return nil
	}

	return fmt.Errorf("invalid RSS healthz response, expected 'ok' but got: %s", string(data))
}

// register registers the Node to the database.
func (n *NTA) register(ctx context.Context, request *nta.RegisterNodeRequest, requestIP string, nodeInfo stakingv2.Node) error {
	// Find node from the database.
	node, err := n.databaseClient.FindNode(ctx, request.Address)
	if err != nil {
		node = &schema.Node{
			Address: request.Address,
		}

		// Get Node's avatar from the VSL
		if node.Avatar, err = n.buildNodeAvatar(ctx, request.Address); err != nil {
			return fmt.Errorf("build node avatar: %w", err)
		}

		// Get from redis if the tax rate of the Node needs to be hidden.
		if err = n.cacheClient.Get(ctx, n.buildNodeHideTaxRateKey(request.Address), &node.HideTaxRate); err != nil && !errors.Is(err, redis.Nil) {
			return fmt.Errorf("get hide tax rate: %w", err)
		}
	}

	node.Endpoint = request.Endpoint
	node.Stream = request.Stream
	node.Config = request.Config
	node.ID = nodeInfo.NodeId
	node.IsPublicGood = nodeInfo.PublicGood
	node.LastHeartbeatTimestamp = time.Now().Unix()
	node.Type = request.Type
	node.Version = request.Version
	// Compatible with the v1.0.0 version of the node.
	if node.Type != schema.NodeTypeAlpha.String() && node.Version == "v0.1.0" {
		node.Version = "v1.0.0"
	}
	// Implement RSS3 node authentication using Bearer tokens.
	node.AccessToken = fmt.Sprintf("Bearer %s", request.AccessToken)
	node.Status = schema.NodeStatusOnline
	node.Location, err = n.geoLite2.LookupNodeLocation(ctx, requestIP)

	if err != nil {
		zap.L().Error("get Node local error", zap.Error(err))
	}

	// Save Node to database.
	if err = n.databaseClient.SaveNode(ctx, node); err != nil {
		return fmt.Errorf("save Node: %s, %w", node.Address.String(), err)
	}

	if node.Type != schema.NodeTypeAlpha.String() {
		if err = n.updateNodeStats(ctx, node, nodeInfo); err != nil {
			return err
		}
	}

	return nil
}

// updateNodeStats updates node stats on nodes registered during the non-alpha phase.
func (n *NTA) updateNodeStats(ctx context.Context, node *schema.Node, nodeInfo stakingv2.Node) error {
	stat, err := n.updateNodeStat(ctx, node, nodeInfo)
	if err != nil {
		return fmt.Errorf("update Node stat: %w", err)
	}

	return n.databaseClient.SaveNodeStat(ctx, stat)
}

// updateNodeStat updates the Node stat.
func (n *NTA) updateNodeStat(ctx context.Context, node *schema.Node, nodeInfo stakingv2.Node) (*schema.Stat, error) {
	stat, err := n.databaseClient.FindNodeStat(ctx, node.Address)
	if err != nil {
		return nil, fmt.Errorf("find Node stat: %w", err)
	}

	// Convert the staking to float64.
	staking, _ := nodeInfo.StakingPoolTokens.Div(nodeInfo.StakingPoolTokens, big.NewInt(1e18)).Float64()
	isRsshubNode := false

	if node.Type == schema.NodeTypeRSSHub.String() {
		isRsshubNode = true
	}

	if stat == nil {
		stat = &schema.Stat{
			Address:      node.Address,
			Endpoint:     node.Endpoint,
			AccessToken:  node.AccessToken,
			IsPublicGood: node.IsPublicGood,
			IsRsshubNode: isRsshubNode,
			Staking:      staking,
			ResetAt:      time.Now(),
		}
	} else {
		stat.Endpoint = node.Endpoint
		stat.AccessToken = node.AccessToken
		stat.IsPublicGood = node.IsPublicGood
		stat.IsRsshubNode = isRsshubNode
		stat.Staking = staking
		stat.ResetAt = time.Now()
	}

	return stat, nil
}

// saveHeartbeat saves the heartbeat to the database.
func (n *NTA) saveHeartbeat(ctx context.Context, node *schema.Node, requestIP string) error {
	var err error
	// Get node local info.
	if len(node.Location) == 0 {
		node.Location, err = n.geoLite2.LookupNodeLocation(ctx, requestIP)
		if err != nil {
			zap.L().Error("failed to get Node local", zap.Error(err))
		}
	}

	if node.Address != (common.Address{}) {
		// Get Node's avatar from the VSL.
		if node.Avatar == nil || node.Avatar.Name == "" {
			node.Avatar, err = n.buildNodeAvatar(ctx, node.Address)
			if err != nil {
				return fmt.Errorf("failed to build Node avatar: %w", err)
			}
		}
	}

	node.LastHeartbeatTimestamp = time.Now().Unix()
	node.Status = schema.NodeStatusOnline

	if err != nil {
		return fmt.Errorf("failed to update Node status: %w", err)
	}

	// Save Node to database.
	return n.databaseClient.SaveNode(ctx, node)
}

// checkSignature checks the signature.
func (n *NTA) checkSignature(_ context.Context, address common.Address, message string, param string) error {
	signature, err := hexutil.Decode(param)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	data := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(data)).Bytes()

	if signature[crypto.RecoveryIDOffset] == 27 || signature[crypto.RecoveryIDOffset] == 28 {
		signature[crypto.RecoveryIDOffset] -= 27
	}

	pubKey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return fmt.Errorf("failed to parse signature: %w", err)
	}

	result := crypto.PubkeyToAddress(*pubKey)
	if address != result {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

// checkAvailable checks if the endpoint is available and contains the node's address.
func (n *NTA) checkAvailable(ctx context.Context, nodeVersion, endpoint string, address common.Address) error {
	curVersion, _ := version.NewVersion(nodeVersion)

	prefix := ""
	if minVersion, _ := version.NewVersion("1.1.2"); curVersion.GreaterThanOrEqual(minVersion) {
		prefix = "operators"
	}

	if prefix != "" {
		endpoint = strings.TrimSuffix(endpoint, "/") + "/" + prefix
	}

	response, _, err := n.httpClient.FetchWithMethod(ctx, http.MethodGet, endpoint, "", nil)
	if err != nil {
		return fmt.Errorf("failed to fetch node endpoint %s: %w", endpoint, err)
	}

	defer lo.Try(response.Close)

	// Use a limited reader to avoid reading too much data.
	content, err := io.ReadAll(io.LimitReader(response, 4096))
	if err != nil {
		return fmt.Errorf("failed to parse node response: %w", err)
	}

	// Check if the node's address is in the response.
	// This is a simple check to ensure the node is responding correctly.
	// The content sample is: "This is an RSS3 Node operated by 0x0000000000000000000000000000000000000000.".
	if !strings.Contains(string(content), address.String()) {
		return fmt.Errorf("invalid node response, expected response contains: %s, actual response: %s", address.String(), string(content))
	}

	return nil
}

// maskIPAddress masks the IP address to show only part of it (e.g., "111.222.xxx.xxx")
// func (n *NTA) maskIPAddress(ip net.IP) string {
// 	if ip == nil {
// 		return "unknown"
// 	}

// 	ipStr := ip.String()
// 	parts := strings.Split(ipStr, ".")

// 	if len(parts) == 4 {
// 		// Show first two octets, mask the last two
// 		return fmt.Sprintf("%s.%s.xxx.xxx", parts[0], parts[1])
// 	}

// 	// For non-IPv4 addresses, return as is
// 	return ipStr
// }

// handleIncompleteNode handles incomplete RSSHub nodes
func (n *NTA) handleIncompleteNode(ctx context.Context, request nta.RSSHubNodeHeartbeatRequest, ip string) error {
	endpoint := request.Endpoint
	if endpoint == "" {
		endpoint = ip
	}

	// Find existing nodes
	nodes, err := n.databaseClient.FindNodes(ctx, schema.FindNodesQuery{
		Endpoint: lo.ToPtr(endpoint),
	})
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			if len(nodes) == 0 {
				// Create new node
				return n.createNewRSSHubNode(ctx, request, endpoint)
			}
		}

		return fmt.Errorf("failed to find nodes by endpoint: %w", err)
	}

	// Update existing node
	return n.updateExistingNode(ctx, nodes[0])
}

// createNewRSSHubNode creates a new RSSHub node
func (n *NTA) createNewRSSHubNode(ctx context.Context, request nta.RSSHubNodeHeartbeatRequest, endpoint string) error {
	// Get next node ID
	nodeID, err := n.getNextNodeID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next node ID: %w", err)
	}

	// Generate node address
	address := n.generateAddressFromID(big.NewInt(nodeID))

	// Create node object
	node := &schema.Node{
		Address:                address,
		ID:                     big.NewInt(nodeID),
		Name:                   request.Name,
		Endpoint:               endpoint,
		IsPublicGood:           true,
		LastHeartbeatTimestamp: time.Now().Unix(),
		Type:                   schema.NodeTypeRSSHub.String(),
		AccessToken:            request.AccessToken,
		Status:                 schema.NodeStatusOnline,
	}

	// Get location information
	if err := n.setNodeLocation(ctx, node, endpoint); err != nil {
		zap.L().Warn("failed to get node location", zap.Error(err))
	}

	// Save node to database
	if err = n.databaseClient.SaveNode(ctx, node); err != nil {
		return fmt.Errorf("failed to save new RSSHub node: %w", err)
	}

	zap.L().Info("created new RSSHub node",
		zap.String("address", address.String()),
		zap.Int64("id", nodeID),
		zap.String("endpoint", endpoint))

	return nil
}

// updateExistingNode updates existing node
func (n *NTA) updateExistingNode(ctx context.Context, node *schema.Node) error {
	node.LastHeartbeatTimestamp = time.Now().Unix()
	node.Status = schema.NodeStatusOnline

	if err := n.databaseClient.SaveNode(ctx, node); err != nil {
		return fmt.Errorf("failed to update existing node: %w", err)
	}

	zap.L().Info("updated existing RSSHub node",
		zap.String("address", node.Address.String()),
		zap.String("endpoint", node.Endpoint))

	return nil
}

// handleCompleteNode handles complete RSSHub nodes
func (n *NTA) handleCompleteNode(ctx context.Context, request nta.RSSHubNodeHeartbeatRequest, ip string) error {
	address := request.Address
	signature := request.Signature
	endpoint := request.Endpoint

	// Validate signature
	if err := n.validateSignature(ctx, address, signature); err != nil {
		return errorx.ValidationFailedError(nil, fmt.Errorf("check signature: %w", err))
	}

	// Parse and validate endpoint
	parsedEndpoint, err := n.parseEndpoint(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("failed to parse endpoint: %w", err)
	}

	if err := n.validateRSSHubEndpoint(ctx, parsedEndpoint, request.AccessToken); err != nil {
		return errorx.ValidationFailedError(nil, fmt.Errorf("validate endpoint: %w", err))
	}

	// Find or create node
	node, err := n.findOrCreateCompleteNode(ctx, address, parsedEndpoint, request)
	if err != nil {
		return err
	}

	// Save heartbeat
	if err := n.saveHeartbeat(ctx, node, ip); err != nil {
		zap.L().Error("save heartbeat", zap.Error(err))
		return errorx.InternalError(nil)
	}

	return nil
}

// findOrCreateCompleteNode finds or creates a complete node
func (n *NTA) findOrCreateCompleteNode(ctx context.Context, address common.Address, endpoint string, request nta.RSSHubNodeHeartbeatRequest) (*schema.Node, error) {
	node, err := n.databaseClient.FindNode(ctx, address)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			// Node not found, create new node
			return n.createCompleteNodeFromVSL(ctx, address, endpoint, request)
		}

		return nil, err
	}

	node.Endpoint = endpoint
	node.AccessToken = request.AccessToken
	node.LastHeartbeatTimestamp = time.Now().Unix()
	node.Status = schema.NodeStatusOnline

	return node, nil
}

// createCompleteNodeFromVSL creates a complete node from VSL
func (n *NTA) createCompleteNodeFromVSL(ctx context.Context, address common.Address, endpoint string, request nta.RSSHubNodeHeartbeatRequest) (*schema.Node, error) {
	// Validate node info
	nodeInfo, err := n.stakingContract.GetNode(&bind.CallOpts{}, address)
	if err != nil {
		zap.L().Error("get the Node from VSL", zap.Error(err))
		return nil, errorx.InternalError(nil)
	}

	if nodeInfo.Account == ethereum.AddressGenesis {
		return nil, errorx.ValidationFailedError(nil, fmt.Errorf("node: %s has not been registered on the VSL", address.String()))
	}

	if !nodeInfo.PublicGood && strings.Compare(nodeInfo.OperationPoolTokens.String(), MinDeposit.String()) < 0 {
		return nil, errorx.ValidationFailedError(nil, fmt.Errorf("insufficient operation pool tokens, expected min deposit %s, actual %s", MinDeposit.String(), nodeInfo.OperationPoolTokens.String()))
	}

	// Create node object
	node := &schema.Node{
		Address: address,
	}

	// Get node avatar
	if node.Avatar, err = n.buildNodeAvatar(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to build node avatar: %w", err)
	}

	// Get hide tax rate setting
	if err = n.cacheClient.Get(ctx, n.buildNodeHideTaxRateKey(address), &node.HideTaxRate); err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to get hide tax rate: %w", err)
	}

	// Set node properties
	node.Endpoint = endpoint
	node.ID = nodeInfo.NodeId
	node.IsPublicGood = nodeInfo.PublicGood
	node.LastHeartbeatTimestamp = time.Now().Unix()
	node.Type = schema.NodeTypeRSSHub.String()
	node.AccessToken = request.AccessToken
	node.Status = schema.NodeStatusOnline

	// Get location
	if err := n.setNodeLocation(ctx, node, endpoint); err != nil {
		zap.L().Warn("failed to get node location", zap.Error(err))
	}

	// Save node
	if err = n.databaseClient.SaveNode(ctx, node); err != nil {
		return nil, fmt.Errorf("failed to save node: %w", err)
	}

	// Update node stats
	if err = n.updateNodeStats(ctx, node, nodeInfo); err != nil {
		return nil, err
	}

	zap.L().Info("created complete RSSHub node from VSL",
		zap.String("address", address.String()),
		zap.String("endpoint", endpoint))

	return node, nil
}

// setNodeLocation sets node location
func (n *NTA) setNodeLocation(ctx context.Context, node *schema.Node, endpoint string) error {
	location, err := n.geoLite2.LookupNodeLocation(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("failed to lookup node location: %w", err)
	}

	node.Location = location

	return nil
}

// getNextNodeID gets the next available node ID
func (n *NTA) getNextNodeID(ctx context.Context) (int64, error) {
	key := "rsshub_node_id_counter"

	// Get and increment counter
	result, err := n.cacheClient.Incr(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to increment node ID counter: %w", err)
	}

	// Start from 10000
	return 10000 + result - 1, nil
}

// generateAddressFromID generates a unique address based on node ID
func (n *NTA) generateAddressFromID(nodeID *big.Int) common.Address {
	// Create deterministic address based on ID
	data := fmt.Sprintf("rsshub_internal_%s", nodeID.String())
	hash := crypto.Keccak256Hash([]byte(data))

	// Convert to address (take first 20 bytes)
	var address common.Address

	copy(address[:], hash[:20])

	// Set special prefix to distinguish internal RSSHub nodes
	// Use 0xFF as first byte to indicate this is an internal address
	address[0] = 0xFF

	return address
}

package nta

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/ethereum"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (n *NTA) RegisterNode(c echo.Context) error {
	var request nta.RegisterNodeRequest

	// Validate request.
	if err := n.validateRegisterNodeRequest(c, &request); err != nil {
		return err
	}

	// Parse request IP.
	ip, err := n.parseRequestIP(c)
	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("failed to parse request ip: %w", err))
	}

	// Validate signature.
	if err = n.validateSignature(c, request.Address, request.Signature); err != nil {
		return err
	}

	// Validate Node info.
	nodeInfo, err := n.validateNodeInfo(c, &request)
	if err != nil {
		return err
	}

	// Validate endpoint.
	if err = n.validateEndpoint(c, request.Address, request.Type, request.Endpoint); err != nil {
		return err
	}

	// Register Node.
	if err = n.register(c.Request().Context(), &request, ip.String(), nodeInfo); err != nil {
		zap.L().Error("register failed", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: fmt.Sprintf("successfully registered node: %v", request.Address),
	})
}

func (n *NTA) NodeHeartbeat(c echo.Context) error {
	var request nta.NodeHeartbeatRequest

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
		return errorx.BadParamsError(c, fmt.Errorf("failed to parse request ip: %w", err))
	}

	// Validate signature.
	if err = n.validateSignature(c, request.Address, request.Signature); err != nil {
		return err
	}

	// Validate Node.
	node, err := n.validateHeartbeatNode(c, request.Address)
	if err != nil {
		return err
	}

	// Validate endpoint.
	if err = n.validateEndpoint(c, request.Address, node.Type, request.Endpoint); err != nil {
		return err
	}

	// Save Node heartbeat.
	if err = n.saveHeartbeat(c.Request().Context(), node, ip.String()); err != nil {
		zap.L().Error("failed to save heartbeat", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: fmt.Sprintf("successfully sent node heartbeat: %v", request.Address),
	})
}

// validateRegisterNodeRequest validates the request.
func (n *NTA) validateRegisterNodeRequest(c echo.Context, request *nta.RegisterNodeRequest) error {
	if err := c.Bind(request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := defaults.Set(request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("failed to set default values: %w", err))
	}

	if err := c.Validate(request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	return nil
}

// validateSignature validates the signature.
func (n *NTA) validateSignature(c echo.Context, address common.Address, signature string) error {
	message := fmt.Sprintf(registrationMessage, strings.ToLower(address.String()))
	if err := n.checkSignature(c.Request().Context(), address, message, signature); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("failed to check signature: %w", err))
	}

	return nil
}

// validateNodeInfo validates the Node info from the VSL.
func (n *NTA) validateNodeInfo(c echo.Context, request *nta.RegisterNodeRequest) (*stakingv2.Node, error) {
	nodeInfo, err := n.stakingContract.GetNode(&bind.CallOpts{}, request.Address)
	if err != nil {
		return nil, errorx.ValidationFailedError(c, fmt.Errorf("failed to get Node from VSL: %w", err))
	}

	if nodeInfo.Account == ethereum.AddressGenesis {
		return nil, errorx.ValidationFailedError(c, fmt.Errorf("node: %s has not been registered on the VSL",
			request.Address.String()))
	}

	if !nodeInfo.PublicGood && strings.Compare(nodeInfo.OperationPoolTokens.String(), MinDeposit.String()) < 0 {
		return nil, errorx.ValidationFailedError(c, fmt.Errorf("insufficient operation pool tokens"))
	}

	return &nodeInfo, nil
}

// validateEndpoint validates the endpoint whether it's valid and available.
func (n *NTA) validateEndpoint(c echo.Context, address common.Address, nodeType, endpoint string) error {
	if nodeType == schema.NodeTypeAlpha.String() {
		return nil
	}

	var err error
	endpoint, err = n.parseEndpoint(c.Request().Context(), endpoint)

	if err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("failed to parse endpoint: %w", err))
	}

	if err = n.checkAvailable(c.Request().Context(), endpoint, address); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("failed to check endpoint available: %w", err))
	}

	return nil
}

// validateHeartbeatNode validates the Node from the database.
func (n *NTA) validateHeartbeatNode(c echo.Context, address common.Address) (*schema.Node, error) {
	ctx := c.Request().Context()

	node, err := n.databaseClient.FindNode(ctx, address)
	if err != nil {
		zap.L().Error("failed to find the node",
			zap.String("address", address.String()),
			zap.Error(err))

		return nil, errorx.InternalError(c)
	}

	if node == nil {
		return nil, errorx.BadParamsError(c, fmt.Errorf("node %s not found", address.String()))
	}

	return node, nil
}

// register registers the Node to the database.
func (n *NTA) register(ctx context.Context, request *nta.RegisterNodeRequest, requestIP string, nodeInfo *stakingv2.Node) error {
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
func (n *NTA) updateNodeStats(ctx context.Context, node *schema.Node, nodeInfo *stakingv2.Node) error {
	stat, err := n.updateNodeStat(ctx, node, nodeInfo)
	if err != nil {
		return fmt.Errorf("update Node stat: %w", err)
	}

	return n.databaseClient.SaveNodeStat(ctx, stat)
}

// updateNodeStat updates the Node stat.
func (n *NTA) updateNodeStat(ctx context.Context, node *schema.Node, nodeInfo *stakingv2.Node) (*schema.Stat, error) {
	stat, err := n.databaseClient.FindNodeStat(ctx, node.Address)
	if err != nil {
		return nil, fmt.Errorf("find Node stat: %w", err)
	}

	// Convert the staking to float64.
	staking, _ := nodeInfo.StakingPoolTokens.Div(nodeInfo.StakingPoolTokens, big.NewInt(1e18)).Float64()

	if stat == nil {
		stat = &schema.Stat{
			Address:      node.Address,
			Endpoint:     node.Endpoint,
			AccessToken:  node.AccessToken,
			IsPublicGood: node.IsPublicGood,
			Staking:      staking,
			ResetAt:      time.Now(),
		}
	} else {
		stat.Endpoint = node.Endpoint
		stat.AccessToken = node.AccessToken
		stat.IsPublicGood = node.IsPublicGood
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

	// Get Node's avatar from the VSL.
	if node.Avatar == nil || node.Avatar.Name == "" {
		node.Avatar, err = n.buildNodeAvatar(ctx, node.Address)
		if err != nil {
			return fmt.Errorf("failed to build Node avatar: %w", err)
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
func (n *NTA) checkAvailable(ctx context.Context, endpoint string, address common.Address) error {
	response, err := n.httpClient.FetchWithMethod(ctx, http.MethodGet, endpoint, "", nil)
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

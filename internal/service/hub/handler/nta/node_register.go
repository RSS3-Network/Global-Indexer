package nta

import (
	"context"
	"encoding/json"
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
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/enforcer"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (n *NTA) RegisterNode(c echo.Context) error {
	var request nta.RegisterNodeRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.InternalError(c, fmt.Errorf("set default values for reqeust: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	ip, err := n.parseRequestIP(c)
	if err != nil {
		return errorx.InternalError(c, fmt.Errorf("parse request ip: %w", err))
	}

	if err := n.register(c.Request().Context(), &request, ip.String()); err != nil {
		return errorx.InternalError(c, fmt.Errorf("register failed: %w", err))
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: fmt.Sprintf("node registered: %v", request.Address),
	})
}

func (n *NTA) NodeHeartbeat(c echo.Context) error {
	var request nta.NodeHeartbeatRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	ip, err := n.parseRequestIP(c)
	if err != nil {
		return errorx.InternalError(c, fmt.Errorf("parse request ip: %w", err))
	}

	if err := n.heartbeat(c.Request().Context(), &request, ip.String()); err != nil {
		return errorx.InternalError(c, fmt.Errorf("heartbeat failed: %w", err))
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: fmt.Sprintf("node heartbeat: %v", request.Address),
	})
}

func (n *NTA) register(ctx context.Context, request *nta.RegisterNodeRequest, requestIP string) error {
	// Check signature.
	message := fmt.Sprintf(registerMessage, strings.ToLower(request.Address.String()))

	if err := n.checkSignature(ctx, request.Address, message, hexutil.MustDecode(request.Signature)); err != nil {
		return err
	}

	// Check Node from the VSL.
	nodeInfo, err := n.stakingContract.GetNode(&bind.CallOpts{}, request.Address)
	if err != nil {
		return fmt.Errorf("get Node from chain: %w", err)
	}

	if nodeInfo.Account == ethereum.AddressGenesis {
		return fmt.Errorf("node: %s has not been registered on the VSL", strings.ToLower(request.Address.String()))
	}

	if strings.Compare(nodeInfo.OperationPoolTokens.String(), decimal.NewFromInt(10000).Mul(decimal.NewFromInt(1e18)).String()) < 0 {
		return fmt.Errorf("insufficient operation pool tokens")
	}

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

	node.Endpoint, err = n.parseEndpoint(ctx, request.Endpoint)
	if err != nil {
		zap.L().Error("parse endpoint", zap.Error(err), zap.String("endpoint", request.Endpoint))

		return fmt.Errorf("parse endpoint: %w", err)
	}

	node.Stream = request.Stream
	node.Config = request.Config
	node.ID = nodeInfo.NodeId
	node.IsPublicGood = nodeInfo.PublicGood
	node.LastHeartbeatTimestamp = time.Now().Unix()
	node.Type = request.Type

	// Checks begin from the beta stage.
	if node.Type == "beta" {
		// Check if the endpoint is available and contains the node's address before update the node's status to online.
		if err = n.checkAvailable(ctx, node.Endpoint, node.Address); err != nil {
			return fmt.Errorf("check endpoint available: %w", err)
		}
	}

	err = nta.UpdateNodeStatus(node, schema.NodeStatusOnline)
	if err != nil {
		return fmt.Errorf("update node status: %w", err)
	}

	minTokensToStake, err := n.stakingContract.MinTokensToStake(&bind.CallOpts{}, request.Address)
	if err != nil {
		return fmt.Errorf("get min token to stake from chain: %w", err)
	}

	node.MinTokensToStake = decimal.NewFromBigInt(minTokensToStake, 0)

	node.Location, err = n.geoLite2.LookupNodeLocation(ctx, requestIP)
	if err != nil {
		zap.L().Error("get Node local error", zap.Error(err))
	}

	// Save Node to database.
	if err = n.databaseClient.SaveNode(ctx, node); err != nil {
		return fmt.Errorf("save Node: %s, %w", node.Address.String(), err)
	}

	// TODO: The transitional implementation during the beta phase, which will soon be deprecated.
	if !nodeInfo.Alpha {
		if err = n.updateBetaNodeStats(ctx, request.Config, node, nodeInfo); err != nil {
			return err
		}
	}

	return nil
}

func (n *NTA) updateBetaNodeStats(ctx context.Context, config json.RawMessage, node *schema.Node, nodeInfo l2.DataTypesNode) error {
	var nodeConfig NodeConfig

	if err := json.Unmarshal(config, &nodeConfig); err != nil {
		return fmt.Errorf("unmarshal node config: %w", err)
	}

	// Check if the Node is a full node.
	fullNode, err := isFullNode(nodeConfig.Decentralized)
	if err != nil {
		return fmt.Errorf("check full node error: %w", err)
	}

	stat, err := n.updateNodeStat(ctx, node, nodeConfig, fullNode, nodeInfo)
	if err != nil {
		return fmt.Errorf("update node stat: %w", err)
	}

	return n.databaseClient.WithTransaction(ctx, func(ctx context.Context, client database.Client) error {
		// Save Node stat to database
		if err = client.SaveNodeStat(ctx, stat); err != nil {
			return fmt.Errorf("save Node stat: %s, %w", node.Address.String(), err)
		}

		zap.L().Info("save Node stat", zap.Any("node", node.Address.String()))

		// If the Node is a full node,
		// then delete the record from the table.
		// Otherwise, add the worker to the table.
		if err = client.DeleteNodeWorkers(ctx, node.Address); err != nil {
			return fmt.Errorf("delete node workers: %s, %w", node.Address.String(), err)
		}

		// Save light node workers to database.
		if !fullNode {
			workers := updateNodeWorkers(node.Address, nodeConfig)
			if err = client.SaveNodeWorkers(ctx, workers); err != nil {
				return fmt.Errorf("save Node workers: %s, %w", node.Address.String(), err)
			}

			zap.L().Info("save Node worker", zap.Any("node", node.Address.String()))
		}

		return nil
	})
}

// isFullNode returns true if the Node is a full Node: has every worker on all possible networks.
func isFullNode(workers []*NodeConfigModule) (bool, error) {
	if len(workers) < len(model.WorkerToNetworksMap) {
		return false, nil
	}

	workerToNetworksMap := make(map[filter.Name]map[string]struct{})

	for _, worker := range workers {
		wid, err := filter.NameString(worker.Worker.String())

		if err != nil {
			return false, err
		}

		if _, exists := workerToNetworksMap[wid]; !exists {
			workerToNetworksMap[wid] = make(map[string]struct{})
		}

		workerToNetworksMap[wid][worker.Network.String()] = struct{}{}
	}

	// Ensure all networks for each worker are present
	for wid, requiredNetworks := range model.WorkerToNetworksMap {
		networks, exists := workerToNetworksMap[wid]
		if !exists || len(networks) != len(requiredNetworks) {
			return false, nil
		}

		for _, network := range requiredNetworks {
			if _, exists = networks[network]; !exists {
				return false, nil
			}
		}
	}

	return true, nil
}

func (n *NTA) updateNodeStat(ctx context.Context, node *schema.Node, nodeConfig NodeConfig, fullNode bool, nodeInfo l2.DataTypesNode) (*schema.Stat, error) {
	var (
		stat *schema.Stat
		err  error
	)

	stat, err = n.databaseClient.FindNodeStat(ctx, node.Address)
	if err != nil {
		return nil, fmt.Errorf("find Node stat: %w", err)
	}

	// Convert the staking to float64.
	staking, _ := nodeInfo.StakingPoolTokens.Div(nodeInfo.StakingPoolTokens, big.NewInt(1e18)).Float64()

	if stat == nil {
		stat = &schema.Stat{
			Address:      node.Address,
			Endpoint:     node.Endpoint,
			IsPublicGood: node.IsPublicGood,
			Staking:      staking,
			ResetAt:      time.Now(),
			IsFullNode:   fullNode,
			IsRssNode:    len(nodeConfig.RSS) > 0,
			DecentralizedNetwork: len(lo.UniqBy(nodeConfig.Decentralized, func(module *NodeConfigModule) filter.Network {
				return module.Network
			})),
			FederatedNetwork: len(nodeConfig.Federated),
			Indexer:          len(nodeConfig.Decentralized),
		}
	} else {
		stat.Endpoint = node.Endpoint
		stat.IsPublicGood = node.IsPublicGood
		stat.Staking = staking
		stat.IsFullNode = fullNode
		stat.IsRssNode = len(nodeConfig.RSS) > 0
		stat.DecentralizedNetwork = len(lo.UniqBy(nodeConfig.Decentralized, func(module *NodeConfigModule) filter.Network {
			return module.Network
		}))
		stat.FederatedNetwork = len(nodeConfig.Federated)
		stat.Indexer = len(nodeConfig.Decentralized)
	}

	// Calculate the reliability score.
	_ = enforcer.CalculateReliabilityScore(stat)

	return stat, nil
}

func updateNodeWorkers(address common.Address, nodeConfig NodeConfig) []*schema.Worker {
	workers := make([]*schema.Worker, 0, len(nodeConfig.Decentralized))

	for _, worker := range nodeConfig.Decentralized {
		workers = append(workers, &schema.Worker{
			Address: address,
			Network: worker.Network.String(),
			Name:    worker.Worker.String(),
		})
	}

	return workers
}

func (n *NTA) heartbeat(ctx context.Context, request *nta.NodeHeartbeatRequest, requestIP string) error {
	// Check signature.
	message := fmt.Sprintf(registerMessage, strings.ToLower(request.Address.String()))

	if err := n.checkSignature(ctx, request.Address, message, hexutil.MustDecode(request.Signature)); err != nil {
		return fmt.Errorf("check signature: %w", err)
	}

	// Check Node from database.
	node, err := n.databaseClient.FindNode(ctx, request.Address)
	if err != nil {
		return fmt.Errorf("get Node %s from database: %w", request.Address, err)
	}

	if node == nil {
		return fmt.Errorf("node %s not found", request.Address)
	}

	if node.Type == "beta" {
		// Check if the endpoint is available and contains the node's address.
		if err := n.checkAvailable(ctx, node.Endpoint, node.Address); err != nil {
			return fmt.Errorf("check endpoint available: %w", err)
		}
	}

	// Get node local info.
	if len(node.Location) == 0 {
		node.Location, err = n.geoLite2.LookupNodeLocation(ctx, requestIP)
		if err != nil {
			zap.L().Error("get Node local error", zap.Error(err))
		}
	}

	// Get Node's avatar from the VSL.
	if node.Avatar == nil || node.Avatar.Name == "" {
		node.Avatar, err = n.buildNodeAvatar(ctx, request.Address)
		if err != nil {
			return fmt.Errorf("build node avatar: %w", err)
		}
	}

	node.LastHeartbeatTimestamp = time.Now().Unix()
	err = nta.UpdateNodeStatus(node, schema.NodeStatusOnline)

	if err != nil {
		return fmt.Errorf("update node status: %w", err)
	}

	// Save Node to database.
	return n.databaseClient.SaveNode(ctx, node)
}

func (n *NTA) checkSignature(_ context.Context, address common.Address, message string, signature []byte) error {
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
	response, err := n.httpClient.Fetch(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("fetch node endpoint %s: %w", endpoint, err)
	}

	defer lo.Try(response.Close)

	// Use a limited reader to avoid reading too much data.
	content, err := io.ReadAll(io.LimitReader(response, 4096))
	if err != nil {
		return fmt.Errorf("parse node response: %w", err)
	}

	// Check if the node's address is in the response.
	// This is a simple check to ensure the node is responding correctly.
	// The content sample is: "This is an RSS3 Node operated by 0x0000000000000000000000000000000000000000.".
	if !strings.Contains(string(content), address.String()) {
		return fmt.Errorf("invalid node response")
	}

	return nil
}

type NodeConfig struct {
	RSS           []*NodeConfigModule `json:"rss"`
	Federated     []*NodeConfigModule `json:"federated"`
	Decentralized []*NodeConfigModule `json:"decentralized"`
}

type NodeConfigModule struct {
	Network  filter.Network `json:"network"`
	Endpoint string         `json:"endpoint"`
	Worker   filter.Name    `json:"worker"`
}

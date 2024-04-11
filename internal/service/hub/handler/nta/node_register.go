package nta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/ethereum"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/distributor"
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

	// Check node from the chain.
	nodeInfo, err := n.stakingContract.GetNode(&bind.CallOpts{}, request.Address)
	if err != nil {
		return fmt.Errorf("get node from chain: %w", err)
	}

	if nodeInfo.Account == ethereum.AddressGenesis {
		return fmt.Errorf("node: %s has not been registered on the chain", strings.ToLower(request.Address.String()))
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

		// Get node's avatar from the chain
		if node.Avatar, err = n.buildNodeAvatar(ctx, request.Address); err != nil {
			return fmt.Errorf("build node avatar: %w", err)
		}

		// Get from redis if the tax rate of the node needs to be hidden.
		if err = n.cacheClient.Get(ctx, n.buildNodeHideTaxRateKey(request.Address), &node.HideTaxRate); err != nil && !errors.Is(err, redis.Nil) {
			return fmt.Errorf("get hide tax rate: %w", err)
		}
	}

	node.Endpoint = n.parseEndpoint(ctx, request.Endpoint)
	node.Stream = request.Stream
	node.Config = request.Config
	node.ID = nodeInfo.NodeId
	node.IsPublicGood = nodeInfo.PublicGood
	node.LastHeartbeatTimestamp = time.Now().Unix()
	node.Status = schema.NodeStatusOnline

	minTokensToStake, err := n.stakingContract.MinTokensToStake(&bind.CallOpts{}, request.Address)
	if err != nil {
		return fmt.Errorf("get min token to stake from chain: %w", err)
	}

	node.MinTokensToStake = decimal.NewFromBigInt(minTokensToStake, 0)

	node.Local, err = n.geoLite2.LookupLocal(ctx, requestIP)
	if err != nil {
		zap.L().Error("get node local error", zap.Error(err))
	}

	var (
		nodeConfig NodeConfig
		indexers   []*schema.Indexer
	)

	if err = json.Unmarshal(request.Config, &nodeConfig); err != nil {
		return fmt.Errorf("unmarshal node config: %w", err)
	}

	fullNode, err := n.verifyFullNode(nodeConfig.Decentralized)
	if err != nil {
		return fmt.Errorf("check full node error: %w", err)
	}

	stat, err := n.updateNodeStat(ctx, request, nodeConfig, fullNode, nodeInfo.PublicGood)
	if err != nil {
		return fmt.Errorf("update node stat: %w", err)
	}

	if !fullNode {
		indexers = n.updateNodeIndexers(ctx, request.Address, nodeConfig)
	}

	// Save node info to the database.
	return n.databaseClient.WithTransaction(ctx, func(ctx context.Context, client database.Client) error {
		// Save node to database.
		if err = client.SaveNode(ctx, node); err != nil {
			return fmt.Errorf("save node: %s, %w", node.Address.String(), err)
		}

		zap.L().Info("save node", zap.Any("node", node.Address.String()))

		// Save node stat to database
		if err = client.SaveNodeStat(ctx, stat); err != nil {
			return fmt.Errorf("save node stat: %s, %w", node.Address.String(), err)
		}

		zap.L().Info("save node stat", zap.Any("node", node.Address.String()))

		// If the node is a full node,
		// then delete the record from the table.
		// Otherwise, add the indexers to the table.
		if err = client.DeleteNodeIndexers(ctx, node.Address); err != nil {
			return fmt.Errorf("delete node indexers: %s, %w", node.Address.String(), err)
		}

		if !fullNode {
			if err = client.SaveNodeIndexers(ctx, indexers); err != nil {
				return fmt.Errorf("save node indexers: %s, %w", node.Address.String(), err)
			}

			zap.L().Info("save node indexer", zap.Any("node", node.Address.String()))
		}

		return nil
	})
}

func (n *NTA) updateNodeStat(ctx context.Context, request *nta.RegisterNodeRequest, nodeConfig NodeConfig, fullNode, publicNode bool) (*schema.Stat, error) {
	var (
		stat *schema.Stat
		err  error
	)

	stat, err = n.databaseClient.FindNodeStat(ctx, request.Address)
	if err != nil {
		return nil, fmt.Errorf("find node stat: %w", err)
	}

	if stat == nil {
		stat = &schema.Stat{
			Address:      request.Address,
			Endpoint:     request.Endpoint,
			IsPublicGood: publicNode,
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
		stat.Endpoint = request.Endpoint
		stat.IsPublicGood = publicNode
		stat.IsFullNode = fullNode
		stat.IsRssNode = len(nodeConfig.RSS) > 0
		stat.DecentralizedNetwork = len(lo.UniqBy(nodeConfig.Decentralized, func(module *NodeConfigModule) filter.Network {
			return module.Network
		}))
		stat.FederatedNetwork = len(nodeConfig.Federated)
		stat.Indexer = len(nodeConfig.Decentralized)
	}

	return stat, nil
}

func (n *NTA) updateNodeIndexers(_ context.Context, address common.Address, nodeConfig NodeConfig) []*schema.Indexer {
	indexers := make([]*schema.Indexer, 0, len(nodeConfig.Decentralized))

	for _, indexer := range nodeConfig.Decentralized {
		indexers = append(indexers, &schema.Indexer{
			Address: address,
			Network: indexer.Network.String(),
			Worker:  indexer.Worker.String(),
		})
	}

	return indexers
}

func (n *NTA) heartbeat(ctx context.Context, request *nta.NodeHeartbeatRequest, requestIP string) error {
	// Check signature.
	message := fmt.Sprintf(registerMessage, strings.ToLower(request.Address.String()))

	if err := n.checkSignature(ctx, request.Address, message, hexutil.MustDecode(request.Signature)); err != nil {
		return fmt.Errorf("check signature: %w", err)
	}

	// Check node from database.
	node, err := n.databaseClient.FindNode(ctx, request.Address)
	if err != nil {
		return fmt.Errorf("get node %s from database: %w", request.Address, err)
	}

	if node == nil {
		return fmt.Errorf("node %s not found", request.Address)
	}

	// Get node local info.
	if len(node.Local) == 0 {
		node.Local, err = n.geoLite2.LookupLocal(ctx, requestIP)
		if err != nil {
			zap.L().Error("get node local error", zap.Error(err))
		}
	}

	// Get node's avatar from the chain.
	if node.Avatar == nil || node.Avatar.Name == "" {
		node.Avatar, err = n.buildNodeAvatar(ctx, request.Address)
		if err != nil {
			return fmt.Errorf("build node avatar: %w", err)
		}
	}

	node.LastHeartbeatTimestamp = time.Now().Unix()
	node.Status = schema.NodeStatusOnline

	// Save node to database.
	return n.databaseClient.SaveNode(ctx, node)
}

func (n *NTA) verifyFullNode(indexers []*NodeConfigModule) (bool, error) {
	if len(indexers) < len(distributor.WorkerToNetworksMap) {
		return false, nil
	}

	workerToNetworksMap := make(map[filter.Name]map[string]struct{})

	for _, indexer := range indexers {
		wid, err := filter.NameString(indexer.Worker.String())

		if err != nil {
			return false, err
		}

		if _, exists := workerToNetworksMap[wid]; !exists {
			workerToNetworksMap[wid] = make(map[string]struct{})
		}

		workerToNetworksMap[wid][indexer.Network.String()] = struct{}{}
	}

	for wid, requiredNetworks := range distributor.WorkerToNetworksMap {
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

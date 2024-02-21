package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/ethereum"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var message = "I, %s, am signing this message for registering my intention to operate an RSS3 Node."

func (h *Hub) getNode(ctx context.Context, address common.Address) (*schema.Node, error) {
	node, err := h.databaseClient.FindNode(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("get node %s: %w", address, err)
	}

	if node == nil {
		return nil, nil
	}

	nodeInfo, err := h.stakingContract.GetNode(&bind.CallOpts{}, address)
	if err != nil {
		return nil, fmt.Errorf("get node from chain: %w", err)
	}

	node.Name = nodeInfo.Name
	node.Description = nodeInfo.Description
	node.TaxRateBasisPoints = nodeInfo.TaxRateBasisPoints
	node.OperationPoolTokens = nodeInfo.OperationPoolTokens.String()
	node.StakingPoolTokens = nodeInfo.StakingPoolTokens.String()
	node.TotalShares = nodeInfo.TotalShares.String()
	node.SlashedTokens = nodeInfo.SlashedTokens.String()

	return node, nil
}

func (h *Hub) getNodes(ctx context.Context, request *BatchNodeRequest) ([]*schema.Node, error) {
	nodes, err := h.databaseClient.FindNodes(ctx, request.NodeAddress, nil, request.Cursor, request.Limit)
	if err != nil {
		return nil, fmt.Errorf("get nodes: %w", err)
	}

	addresses := lo.Map(nodes, func(node *schema.Node, _ int) common.Address {
		return node.Address
	})

	nodeInfo, err := h.stakingContract.GetNodes(&bind.CallOpts{}, addresses)
	if err != nil {
		return nil, fmt.Errorf("get nodes from chain: %w", err)
	}

	nodeInfoMap := lo.SliceToMap(nodeInfo, func(node l2.DataTypesNode) (common.Address, l2.DataTypesNode) {
		return node.Account, node
	})

	for _, node := range nodes {
		if nodeInfo, exists := nodeInfoMap[node.Address]; exists {
			node.Name = nodeInfo.Name
			node.Description = nodeInfo.Description
			node.TaxRateBasisPoints = nodeInfo.TaxRateBasisPoints
			node.OperationPoolTokens = nodeInfo.OperationPoolTokens.String()
			node.StakingPoolTokens = nodeInfo.StakingPoolTokens.String()
			node.TotalShares = nodeInfo.TotalShares.String()
			node.SlashedTokens = nodeInfo.SlashedTokens.String()
		}
	}

	return nodes, nil
}

func (h *Hub) register(ctx context.Context, request *RegisterNodeRequest, requestIP string) error {
	// Check signature.
	if err := h.checkSignature(ctx, request.Address, hexutil.MustDecode(request.Signature)); err != nil {
		return err
	}

	node := &schema.Node{
		Address:  request.Address,
		Endpoint: h.parseEndpoint(ctx, request.Endpoint),
		Stream:   request.Stream,
		Config:   request.Config,
	}

	// Check node from chain.
	nodeInfo, err := h.stakingContract.GetNode(&bind.CallOpts{}, request.Address)
	if err != nil {
		return fmt.Errorf("get node from chain: %w", err)
	}

	if nodeInfo.Account == ethereum.AddressGenesis {
		return fmt.Errorf("node: %s has not been registered on the chain", strings.ToLower(request.Address.String()))
	}

	if strings.Compare(nodeInfo.OperationPoolTokens.String(), decimal.NewFromInt(10000).Mul(decimal.NewFromInt(1e18)).String()) < 0 {
		return fmt.Errorf("insufficient operation pool tokens")
	}

	node.IsPublicGood = nodeInfo.PublicGood
	node.LastHeartbeatTimestamp = time.Now().Unix()
	node.Status = schema.StatusOnline

	node.Local, err = h.geoLite2.LookupLocal(ctx, requestIP)
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

	fullNode, err := h.verifyFullNode(nodeConfig.Decentralized)
	if err != nil {
		return fmt.Errorf("check full node error: %w", err)
	}

	stat, err := h.updateNodeStat(ctx, request, nodeConfig, fullNode, nodeInfo.PublicGood)
	if err != nil {
		return fmt.Errorf("update node stat: %w", err)
	}

	if !fullNode {
		indexers = h.updateNodeIndexers(ctx, request.Address, nodeConfig)
	}

	// Save node info to the database.
	return h.databaseClient.WithTransaction(ctx, func(ctx context.Context, client database.Client) error {
		// Save node to database.
		if err = h.databaseClient.SaveNode(ctx, node); err != nil {
			return fmt.Errorf("save node: %s, %w", node.Address.String(), err)
		}

		zap.L().Info("save node", zap.Any("node", node.Address.String()))

		// Save node stat to database
		if err = h.databaseClient.SaveNodeStat(ctx, stat); err != nil {
			return fmt.Errorf("save node stat: %s, %w", node.Address.String(), err)
		}

		zap.L().Info("save node stat", zap.Any("node", node.Address.String()))

		// If the node is a full node,
		// then delete the record from the table.
		// Otherwise, add the indexers to the table.
		if err = h.databaseClient.DeleteNodeIndexers(ctx, node.Address); err != nil {
			return fmt.Errorf("delete node indexers: %s, %w", node.Address.String(), err)
		}

		if !fullNode {
			if err = h.databaseClient.SaveNodeIndexers(ctx, indexers); err != nil {
				return fmt.Errorf("save node indexers: %s, %w", node.Address.String(), err)
			}

			zap.L().Info("save node indexer", zap.Any("node", node.Address.String()))
		}

		return nil
	})
}

func (h *Hub) updateNodeStat(ctx context.Context, request *RegisterNodeRequest, nodeConfig NodeConfig, fullNode, publicNode bool) (*schema.Stat, error) {
	var (
		stat *schema.Stat
		err  error
	)

	stat, err = h.databaseClient.FindNodeStat(ctx, request.Address)
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
			DecentralizedNetwork: len(lo.UniqBy(nodeConfig.Decentralized, func(module *Module) filter.Network {
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
		stat.DecentralizedNetwork = len(lo.UniqBy(nodeConfig.Decentralized, func(module *Module) filter.Network {
			return module.Network
		}))
		stat.FederatedNetwork = len(nodeConfig.Federated)
		stat.Indexer = len(nodeConfig.Decentralized)
	}

	return stat, nil
}

func (h *Hub) updateNodeIndexers(_ context.Context, address common.Address, nodeConfig NodeConfig) []*schema.Indexer {
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

func (h *Hub) heartbeat(ctx context.Context, request *NodeHeartbeatRequest) error {
	// Check signature.
	if err := h.checkSignature(ctx, request.Address, hexutil.MustDecode(request.Signature)); err != nil {
		return fmt.Errorf("check signature: %w", err)
	}

	// Check node from database.
	node, err := h.databaseClient.FindNode(ctx, request.Address)
	if err != nil {
		return fmt.Errorf("get node %s from database: %w", request.Address, err)
	}

	if node == nil {
		return fmt.Errorf("node %s not found", request.Address)
	}

	node.LastHeartbeatTimestamp = request.Timestamp
	node.Status = schema.StatusOnline

	// Save node to database.
	return h.databaseClient.SaveNode(ctx, node)
}

func (h *Hub) verifyFullNode(indexers []*Module) (bool, error) {
	if len(indexers) < len(model.WorkerToNetworksMap) {
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

func (h *Hub) checkSignature(_ context.Context, address common.Address, signature []byte) error {
	message := fmt.Sprintf(message, strings.ToLower(address.Hex()))
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

func (h *Hub) parseEndpoint(_ context.Context, endpoint string) string {
	if ip := net.ParseIP(endpoint); ip != nil {
		return endpoint
	}

	uri, err := url.Parse(endpoint)
	if err != nil {
		return endpoint
	}

	return uri.Hostname()
}

func (h *Hub) parseRequestIP(_ context.Context, request *http.Request) (net.IP, error) {
	if ip := request.Header.Get("X-Real-IP"); net.ParseIP(ip) != nil {
		return net.ParseIP(ip), nil
	}

	ip, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(ip), nil
}

type NodeConfig struct {
	RSS           []*Module `json:"rss"`
	Federated     []*Module `json:"federated"`
	Decentralized []*Module `json:"decentralized"`
}

type Module struct {
	Network  filter.Network `json:"network"`
	Endpoint string         `json:"endpoint"`
	Worker   filter.Name    `json:"worker"`
}

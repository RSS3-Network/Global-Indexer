package hub

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/ethereum"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cache"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/provider/node"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/node/config"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	message = "I, %s, am signing this message for registering my intention to operate an RSS3 Serving Node."

	messageNodeDataFailed = "request node data failed"

	defaultNodeCount   = 3
	defaultSlashCount  = 4
	defaultVerifyCount = 3
)

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
	nodes, err := h.databaseClient.FindNodes(ctx, request.NodeAddress, request.Cursor, request.Limit)
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

func (h *Hub) register(ctx context.Context, request *RegisterNodeRequest) error {
	// Check signature.
	if err := h.checkSignature(ctx, request.Address, hexutil.MustDecode(request.Signature)); err != nil {
		return err
	}

	node := &schema.Node{
		Address:  request.Address,
		Endpoint: request.Endpoint,
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

	var nodeConfig config.Node

	if err = json.Unmarshal(request.Config, &nodeConfig); err != nil {
		return fmt.Errorf("unmarshal node config: %w", err)
	}

	fullNode, err := h.isFullNode(nodeConfig.Decentralized)

	if err != nil {
		return fmt.Errorf("check full node error: %w", err)
	}

	stat := &schema.Stat{
		Address:      request.Address,
		Endpoint:     request.Endpoint,
		IsPublicGood: nodeInfo.PublicGood,
		ResetAt:      time.Now(),
		IsFullNode:   fullNode,
		IsRssNode:    len(nodeConfig.RSS) > 0,
		DecentralizedNetwork: len(lo.UniqBy(nodeConfig.Decentralized, func(module *config.Module) filter.Network {
			return module.Network
		})),
		FederatedNetwork: len(nodeConfig.Federated),
		Indexer:          len(nodeConfig.Decentralized),
	}

	indexers := make([]*schema.Indexer, 0, len(nodeConfig.Decentralized))

	if !fullNode {
		for _, indexer := range nodeConfig.Decentralized {
			indexers = append(indexers, &schema.Indexer{
				Address: request.Address,
				Network: indexer.Network.String(),
				Worker:  indexer.Worker.String(),
			})
		}
	}

	// Save node info to the database.
	if err = h.databaseClient.WithTransaction(ctx, func(ctx context.Context, client database.Client) error {
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
	}); err != nil {
		return err
	}

	return nil
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

func (h *Hub) isFullNode(indexers []*config.Module) (bool, error) {
	if len(indexers) < len(node.WorkerToNetworksMap) {
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

	for wid, requiredNetworks := range node.WorkerToNetworksMap {
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

func (h *Hub) routerRSSHubData(ctx context.Context, path, query string) ([]byte, error) {
	nodes, err := h.retrieveNodes(ctx, node.RssNodeCacheKey)

	if err != nil {
		return nil, err
	}

	nodeMap, err := h.pathBuilder.GetRSSHubPath(path, query, nodes)

	if err != nil {
		return nil, err
	}

	nodeRes, err := h.batchRequest(ctx, nodeMap, h.processRSSHubResults)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeRes.Address.String()))

	return nodeRes.Data, nil
}

func (h *Hub) routerActivityData(ctx context.Context, request node.ActivityRequest) ([]byte, error) {
	nodes, err := h.retrieveNodes(ctx, node.FullNodeCacheKey)

	if err != nil {
		return nil, err
	}

	nodeMap, err := h.pathBuilder.GetActivityByIDPath(request, nodes)

	if err != nil {
		return nil, err
	}

	nodeRes, err := h.batchRequest(ctx, nodeMap, h.processActivityResults)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeRes.Address.String()))

	return nodeRes.Data, nil
}

func (h *Hub) routerActivitiesData(ctx context.Context, request node.AccountActivitiesRequest) ([]byte, error) {
	nodes := make([]node.Cache, 0, defaultNodeCount)

	nodeAddresses, err := h.matchLightNodes(ctx, request)

	if err != nil {
		return nil, err
	}

	if len(nodeAddresses) > 0 {
		nodeStats, err := h.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
			AddressList: nodeAddresses,
			Limit:       lo.ToPtr(defaultNodeCount),
			PointsOrder: lo.ToPtr("DESC"),
		})

		if err != nil {
			return nil, err
		}

		num := lo.Ternary(len(nodeStats) > defaultNodeCount, defaultNodeCount, len(nodeStats))

		for i := 0; i < num; i++ {
			nodes = append(nodes, node.Cache{
				Address:  nodeStats[i].Address.String(),
				Endpoint: nodeStats[i].Endpoint,
			})
		}
	}

	if len(nodes) < defaultNodeCount {
		fullNodes, err := h.retrieveNodes(ctx, node.FullNodeCacheKey)
		if err != nil {
			return nil, err
		}

		nodesNeeded := defaultNodeCount - len(nodes)
		nodesToAdd := lo.Ternary(nodesNeeded > len(fullNodes), len(fullNodes), nodesNeeded)

		for i := 0; i < nodesToAdd; i++ {
			nodes = append(nodes, fullNodes[i])
		}
	}

	nodeMap, err := h.pathBuilder.GetAccountActivitiesPath(request, nodes)

	if err != nil {
		return nil, err
	}

	nodeRes, err := h.batchRequest(ctx, nodeMap, h.processActivitiesResults)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeRes.Address.String()))

	return nodeRes.Data, nil
}

func (h *Hub) matchLightNodes(ctx context.Context, request node.AccountActivitiesRequest) ([]common.Address, error) {
	tagWorkers := make(map[string]struct{})
	platformWorkers := make(map[string]struct{})

	// check network
	for i, network := range request.Network {
		nid, err := filter.NetworkString(network)
		if err != nil {
			return nil, err
		}

		request.Network[i] = nid.String()
	}

	// check tag
	for _, tag := range request.Tag {
		tid, err := filter.TagString(tag)

		if err != nil {
			return nil, err
		}

		tagWorker, exists := node.TagToWorkersMap[tid]

		if !exists {
			return nil, err
		}

		for _, worker := range tagWorker {
			tagWorkers[worker] = struct{}{}
		}
	}

	// check platform
	for _, platform := range request.Platform {
		pid, err := filter.PlatformString(platform)

		if err != nil {
			return nil, err
		}

		platformWorker, exists := node.PlatformToWorkerMap[pid]

		if !exists {
			return nil, err
		}

		platformWorkers[platformWorker] = struct{}{}
	}

	var (
		nodes        []common.Address
		err          error
		workers      []string
		needsWorker  = len(tagWorkers) > 0 || len(platformWorkers) > 0
		needsNetwork = len(request.Network) > 0
	)

	if len(tagWorkers) > 0 && len(platformWorkers) > 0 {
		workers = findCommonElements(lo.Keys(tagWorkers), lo.Keys(platformWorkers))
	} else if len(tagWorkers) > 0 {
		workers = lo.Keys(tagWorkers)
	} else if len(platformWorkers) > 0 {
		workers = lo.Keys(platformWorkers)
	}

	switch {
	case !needsWorker && needsNetwork:
		nodes, err = h.matchNetwork(ctx, request.Network)
	case needsWorker && !needsNetwork:
		nodes, err = h.matchWorker(ctx, workers)
	case needsWorker && needsNetwork:
		nodes, err = h.matchWorkerAndNetwork(ctx, workers, request.Network)
	default:
	}

	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (h *Hub) matchNetwork(ctx context.Context, requestNetworks []string) ([]common.Address, error) {
	nodes := make([]common.Address, 0)

	indexers, err := h.databaseClient.FindNodeIndexers(ctx, nil, requestNetworks, nil)

	if err != nil {
		return nil, err
	}

	nodeNetworkToWorkerMap := make(map[common.Address]map[string][]string)

	for _, indexer := range indexers {
		if _, exists := nodeNetworkToWorkerMap[indexer.Address]; !exists {
			nodeNetworkToWorkerMap[indexer.Address] = make(map[string][]string)
		}

		nodeNetworkToWorkerMap[indexer.Address][indexer.Network] = append(nodeNetworkToWorkerMap[indexer.Address][indexer.Network], indexer.Worker)
	}

	for address, networkToWorkersMap := range nodeNetworkToWorkerMap {
		if len(networkToWorkersMap) != len(requestNetworks) {
			// number of networks not match
			continue
		}

		flag := true

		// check every network workers
		for network, workers := range networkToWorkersMap {
			nid, _ := filter.NetworkString(network)

			requiredWorkers := node.NetworkToWorkersMap[nid]

			// number of workers not match
			if len(requiredWorkers) != len(workers) {
				flag = false

				break
			}

			workerMap := make(map[string]struct{}, len(workers))

			for _, worker := range workers {
				workerMap[worker] = struct{}{}
			}

			// check every worker
			for _, worker := range requiredWorkers {
				if _, exists := workerMap[worker]; !exists {
					flag = false

					break
				}
			}
		}

		if flag {
			nodes = append(nodes, address)
		}
	}

	return nodes, nil
}

func (h *Hub) matchWorker(ctx context.Context, workers []string) ([]common.Address, error) {
	nodes := make([]common.Address, 0)

	indexers, err := h.databaseClient.FindNodeIndexers(ctx, nil, nil, workers)

	if err != nil {
		return nil, err
	}

	nodeWorkerToNetworksMap := make(map[common.Address]map[string][]string)

	for _, indexer := range indexers {
		if _, exists := nodeWorkerToNetworksMap[indexer.Address]; !exists {
			nodeWorkerToNetworksMap[indexer.Address] = make(map[string][]string)
		}

		nodeWorkerToNetworksMap[indexer.Address][indexer.Worker] = append(nodeWorkerToNetworksMap[indexer.Address][indexer.Worker], indexer.Network)
	}

	for address, workerToNetworksMap := range nodeWorkerToNetworksMap {
		if len(nodeWorkerToNetworksMap) != len(workers) {
			// number of workers not match
			continue
		}

		flag := true

		// check every worker networks
		for worker, networks := range workerToNetworksMap {
			wid, _ := filter.NameString(worker)

			requiredNetworks := node.WorkerToNetworksMap[wid]

			// number of networks not match
			if len(requiredNetworks) != len(networks) {
				flag = false

				break
			}

			networkMap := make(map[string]struct{}, len(networks))

			for _, network := range networks {
				networkMap[network] = struct{}{}
			}

			// check every network
			for _, network := range requiredNetworks {
				if _, exists := networkMap[network]; !exists {
					flag = false

					break
				}
			}
		}

		if flag {
			nodes = append(nodes, address)
		}
	}

	return nodes, nil
}

func (h *Hub) matchWorkerAndNetwork(ctx context.Context, workers, requestNetworks []string) ([]common.Address, error) {
	nodes := make([]common.Address, 0)

	indexers, err := h.databaseClient.FindNodeIndexers(ctx, nil, requestNetworks, workers)

	if err != nil {
		return nil, err
	}

	nodeWorkerToNetworksMap := make(map[common.Address]map[string][]string)

	for _, indexer := range indexers {
		if _, exists := nodeWorkerToNetworksMap[indexer.Address]; !exists {
			nodeWorkerToNetworksMap[indexer.Address] = make(map[string][]string)
		}

		nodeWorkerToNetworksMap[indexer.Address][indexer.Worker] = append(nodeWorkerToNetworksMap[indexer.Address][indexer.Worker], indexer.Network)
	}

	for address, workerToNetworksMap := range nodeWorkerToNetworksMap {
		flag := true

		// check every worker networks
		for worker, networks := range workerToNetworksMap {
			wid, _ := filter.NameString(worker)

			workerRequiredNetworks := node.WorkerToNetworksMap[wid]

			requiredNetworks := findCommonElements(workerRequiredNetworks, requestNetworks)

			networkMap := make(map[string]struct{}, len(networks))

			for _, network := range networks {
				networkMap[network] = struct{}{}
			}

			// check every network
			for _, network := range requiredNetworks {
				if _, exists := networkMap[network]; !exists {
					flag = false

					break
				}
			}
		}

		if flag {
			nodes = append(nodes, address)
		}
	}

	return nodes, nil
}

func findCommonElements(slice1, slice2 []string) []string {
	elementMap := make(map[string]struct{})

	var commonElements []string

	for _, v := range slice1 {
		elementMap[v] = struct{}{}
	}

	for _, v := range slice2 {
		if _, found := elementMap[v]; found {
			commonElements = append(commonElements, v)
			delete(elementMap, v)
		}
	}

	return commonElements
}

func (h *Hub) retrieveNodes(ctx context.Context, key string) ([]node.Cache, error) {
	var (
		nodesCache []node.Cache
		nodes      []*schema.Stat
	)

	err := cache.Get(ctx, key, &nodesCache)

	if err == nil {
		return nodesCache, nil
	}

	zap.L().Info("not found nodes from cache", zap.String("key", key))

	if errors.Is(err, redis.Nil) {
		switch key {
		case node.RssNodeCacheKey:
			nodes, err = h.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
				IsRssNode:   lo.ToPtr(true),
				Limit:       lo.ToPtr(defaultNodeCount),
				PointsOrder: lo.ToPtr("DESC"),
			})

			if err != nil {
				return nil, err
			}
		case node.FullNodeCacheKey:
			nodes, err = h.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
				IsFullNode:  lo.ToPtr(true),
				Limit:       lo.ToPtr(defaultNodeCount),
				PointsOrder: lo.ToPtr("DESC"),
			})

			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown cache key: %s", key)
		}

		if err = h.setNodeCache(ctx, key, nodes); err != nil {
			return nil, err
		}

		zap.L().Info("set nodes to cache", zap.String("key", key))

		nodesCache = lo.Map(nodes, func(n *schema.Stat, _ int) node.Cache {
			return node.Cache{
				Address:  n.Address.String(),
				Endpoint: n.Endpoint,
			}
		})

		return nodesCache, nil
	}

	return nil, fmt.Errorf("get nodes from cache: %s, %w", key, err)
}

func (h *Hub) batchRequest(_ context.Context, nodeMap map[common.Address]string, processResults func([]node.DataResponse)) (node.DataResponse, error) {
	var (
		waitGroup   sync.WaitGroup
		firstResult = make(chan node.DataResponse, 1)
		results     []node.DataResponse
		mu          sync.Mutex
	)

	for address, endpoint := range nodeMap {
		waitGroup.Add(1)

		go func(address common.Address, endpoint string) {
			defer waitGroup.Done()

			data, err := h.fetch(context.Background(), endpoint)
			if err != nil {
				zap.L().Error("fetch request error", zap.Any("node", address.String()), zap.Error(err))

				mu.Lock()
				results = append(results, node.DataResponse{Address: address, Err: err})

				if len(results) == len(nodeMap) {
					firstResult <- node.DataResponse{Address: address, Data: []byte(messageNodeDataFailed)}
				}

				mu.Unlock()

				return
			}

			flagActivities, _ := h.validateActivities(data)
			flagActivity, _ := h.validateActivity(data)

			if !flagActivities && !flagActivity {
				zap.L().Error("response parse error", zap.Any("node", address.String()))

				mu.Lock()
				results = append(results, node.DataResponse{Address: address, Err: fmt.Errorf("invalid data")})

				if len(results) == len(nodeMap) {
					firstResult <- node.DataResponse{Address: address, Data: data}
				}
				mu.Unlock()

				return
			}

			mu.Lock()
			results = append(results, node.DataResponse{Address: address, Data: data, First: true})
			mu.Unlock()

			select {
			case firstResult <- node.DataResponse{Address: address, Data: data}:
			default:
			}
		}(address, endpoint)
	}

	go func() {
		waitGroup.Wait()
		close(firstResult)
		processResults(results)
	}()

	select {
	case result := <-firstResult:
		return result, nil
	case <-time.After(time.Second * 3):
		return node.DataResponse{Data: []byte(messageNodeDataFailed)}, fmt.Errorf("timeout waiting for results")
	}
}

func (h *Hub) validateActivities(data []byte) (bool, *node.ActivitiesResponse) {
	var (
		res    node.ActivitiesResponse
		errRes node.ErrResponse
	)

	if err := json.Unmarshal(data, &errRes); err != nil {
		return false, nil
	}

	if errRes.ErrorCode != "" {
		return false, nil
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return false, nil
	}

	return true, &res
}

func (h *Hub) validateActivity(data []byte) (bool, *node.ActivityResponse) {
	var (
		res    node.ActivityResponse
		errRes node.ErrResponse
	)

	if err := json.Unmarshal(data, &errRes); err != nil {
		return false, nil
	}

	if errRes.ErrorCode != "" {
		return false, nil
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return false, nil
	}

	return true, &res
}

func (h *Hub) processRSSHubResults(results []node.DataResponse) {
	if err := h.verifyData(context.Background(), results); err != nil {
		zap.L().Error("fail to verify rss hub request", zap.Any("results", len(results)))
	} else {
		zap.L().Info("complete rss hub request verify", zap.Any("results", len(results)))
	}
}

func (h *Hub) processActivityResults(results []node.DataResponse) {
	if err := h.verifyData(context.Background(), results); err != nil {
		zap.L().Error("fail to verify  feed id request ", zap.Any("results", len(results)))
	} else {
		zap.L().Info("complete feed id request verify", zap.Any("results", len(results)))
	}
}

func (h *Hub) processActivitiesResults(results []node.DataResponse) {
	ctx := context.Background()

	if err := h.verifyData(ctx, results); err != nil {
		zap.L().Error("fail feed request verify", zap.Any("results", len(results)))

		return
	}

	zap.L().Info("complete feed request verify", zap.Any("results", len(results)))

	if !results[0].First {
		return
	}

	var activities node.ActivitiesResponse

	data := results[0].Data

	if err := json.Unmarshal(data, &activities); err != nil {
		zap.L().Error("fail to unmarshall activities")

		return
	}

	// data is empty, no need to 2nd verify
	if activities.Data == nil {
		return
	}

	workingNodes := lo.Map(results, func(result node.DataResponse, _ int) common.Address {
		return result.Address
	})

	h.process2ndVerify(activities.Data, workingNodes)
}

func (h *Hub) process2ndVerify(feeds []*node.Feed, workingNodes []common.Address) {
	ctx := context.Background()
	platformMap := make(map[string]struct{})
	statMap := make(map[string]struct{})

	for _, feed := range feeds {
		if len(feed.Platform) == 0 {
			continue
		}

		h.verifyPlatform(ctx, feed, platformMap, statMap, workingNodes)

		if _, exists := platformMap[feed.Platform]; !exists {
			if len(platformMap) == defaultVerifyCount {
				break
			}
		}
	}
}

func (h *Hub) verifyPlatform(ctx context.Context, feed *node.Feed, platformMap, statMap map[string]struct{}, workingNodes []common.Address) {
	pid, err := filter.PlatformString(feed.Platform)
	if err != nil {
		return
	}

	worker := node.PlatformToWorkerMap[pid]

	indexers, err := h.databaseClient.FindNodeIndexers(ctx, nil, []string{feed.Network}, []string{worker})

	if err != nil {
		return
	}

	nodeAddresses := lo.Map(indexers, func(indexer *schema.Indexer, _ int) common.Address {
		return indexer.Address
	})

	nodeAddresses = lo.Filter(nodeAddresses, func(item common.Address, _ int) bool {
		return !lo.Contains(workingNodes, item)
	})

	if len(nodeAddresses) == 0 {
		return
	}

	stats, err := h.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList: nodeAddresses,
		PointsOrder: lo.ToPtr("DESC"),
	})

	if err != nil || len(stats) == 0 {
		return
	}

	h.verifyStat(ctx, feed, stats, statMap)

	platformMap[feed.Platform] = struct{}{}
}

func (h *Hub) verifyStat(ctx context.Context, feed *node.Feed, stats []*schema.Stat, statMap map[string]struct{}) {
	for _, stat := range stats {
		if stat.EpochInvalidRequest >= int64(defaultSlashCount) {
			continue
		}

		if _, exists := statMap[stat.Address.String()]; !exists {
			statMap[stat.Address.String()] = struct{}{}

			request := node.ActivityRequest{
				ID: feed.ID,
			}

			nodeMap, err := h.pathBuilder.GetActivityByIDPath(
				request,
				[]node.Cache{
					{Address: stat.Address.String(), Endpoint: stat.Endpoint},
				})

			if err != nil {
				continue
			}

			data, err := h.fetch(ctx, nodeMap[stat.Address])

			flag, res := h.validateActivity(data)

			if err != nil || !flag {
				stat.EpochInvalidRequest++
			} else {
				if !h.compareFeeds(feed, res.Data) {
					stat.EpochInvalidRequest++
				} else {
					stat.TotalRequest++
					stat.EpochRequest++
				}
			}

			_ = h.databaseClient.SaveNodeStat(ctx, stat)

			break
		}
	}
}

func (h *Hub) compareFeeds(src, des *node.Feed) bool {
	var flag bool

	if src.ID != des.ID ||
		src.Network != des.Network ||
		src.Index != des.Index ||
		src.From != des.From ||
		src.To != des.To ||
		src.Tag != des.Tag ||
		src.Type != des.Type ||
		src.Platform != des.Platform ||
		len(src.Actions) != len(des.Actions) {
		return false
	}

	if len(src.Actions) > 0 {
		srcAction := src.Actions[0]

		for _, action := range des.Actions {
			if srcAction.From == action.From &&
				srcAction.To == action.To &&
				srcAction.Tag == action.Tag &&
				srcAction.Type == action.Type {
				desMetadata, _ := json.Marshal(action.Metadata)
				srcMetadata, _ := json.Marshal(srcAction.Metadata)

				if compareData(srcMetadata, desMetadata) {
					flag = true
				}
			}
		}
	}

	return flag
}

func (h *Hub) verifyData(ctx context.Context, results []node.DataResponse) error {
	statsMap, err := h.getNodeStatsMap(ctx, results)
	if err != nil {
		return fmt.Errorf("find node stats: %w", err)
	}

	h.sortResults(results)

	if len(statsMap) < defaultNodeCount {
		for i := range results {
			if _, exists := statsMap[results[i].Address]; exists {
				if results[i].Err != nil {
					results[i].InvalidRequest = 1
				} else {
					results[i].Request = 1
				}
			}
		}
	} else {
		if !results[0].First {
			for i := range results {
				results[i].InvalidRequest = 1
			}
		} else {
			h.updateRequestsBasedOnDataCompare(results)
		}
	}

	h.updateStatsWithResults(statsMap, results)

	if err = h.databaseClient.SaveNodeStats(ctx, statsMapToSlice(statsMap)); err != nil {
		return fmt.Errorf("save node stats: %w", err)
	}

	return nil
}

func (h *Hub) getNodeStatsMap(ctx context.Context, results []node.DataResponse) (map[common.Address]*schema.Stat, error) {
	stats, err := h.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList: lo.Map(results, func(result node.DataResponse, _ int) common.Address {
			return result.Address
		}),
		PointsOrder: lo.ToPtr("DESC"),
	})

	if err != nil {
		return nil, err
	}

	statsMap := make(map[common.Address]*schema.Stat)

	for _, stat := range stats {
		statsMap[stat.Address] = stat
	}

	return statsMap, nil
}

func (h *Hub) sortResults(results []node.DataResponse) {
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].First && !results[j].First
	})
}

func (h *Hub) updateStatsWithResults(statsMap map[common.Address]*schema.Stat, results []node.DataResponse) {
	for _, result := range results {
		if stat, exists := statsMap[result.Address]; exists {
			stat.TotalRequest += int64(result.Request)
			stat.EpochRequest += int64(result.Request)
			stat.EpochInvalidRequest += int64(result.InvalidRequest)
		}
	}
}

func (h *Hub) updateRequestsBasedOnDataCompare(results []node.DataResponse) {
	diff01 := compareData(results[0].Data, results[1].Data)
	diff02 := compareData(results[0].Data, results[2].Data)
	diff12 := compareData(results[1].Data, results[2].Data)

	if diff01 && diff02 {
		results[0].Request = 2
		results[1].Request = 1
		results[2].Request = 1
	} else if !diff01 && diff12 {
		results[0].InvalidRequest = 1
		results[1].Request = 1
		results[2].Request = 1
	} else if !diff01 && diff02 {
		results[0].Request = 2
		results[1].InvalidRequest = 1
		results[2].Request = 1
	} else if diff01 && !diff02 {
		results[0].Request = 2
		results[1].Request = 1
		results[2].InvalidRequest = 1
	} else if !diff01 && !diff02 && !diff12 {
		for i := range results {
			if results[i].Data == nil && results[i].Err != nil {
				results[i].InvalidRequest = 1
			}

			if results[i].Data != nil && results[i].Err == nil {
				results[i].Request = 1
			}
		}
	}
}

func statsMapToSlice(statsMap map[common.Address]*schema.Stat) []*schema.Stat {
	statsSlice := make([]*schema.Stat, 0, len(statsMap))
	for _, stat := range statsMap {
		statsSlice = append(statsSlice, stat)
	}

	return statsSlice
}

func compareData(src, des []byte) bool {
	if src == nil || des == nil {
		return false
	}

	srcHash, destHash := sha256.Sum256(src), sha256.Sum256(des)

	return string(srcHash[:]) == string(destHash[:])
}

func (h *Hub) fetch(ctx context.Context, decodedURI string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, decodedURI, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	res, err := h.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	//if res.StatusCode != http.StatusOK {
	//	return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	//}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return data, nil
}

func (h *Hub) setNodeCache(ctx context.Context, key string, stats []*schema.Stat) error {
	nodesCache := lo.Map(stats, func(n *schema.Stat, _ int) node.Cache {
		return node.Cache{Address: n.Address.String(), Endpoint: n.Endpoint}
	})

	if err := cache.Set(ctx, key, nodesCache); err != nil {
		return fmt.Errorf("set nodes to cache: %s, %w", key, err)
	}

	return nil
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

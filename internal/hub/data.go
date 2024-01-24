package hub

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"math"
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
	"github.com/rss3-network/serving-node/config"
	"github.com/rss3-network/serving-node/schema/filter"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	rssNodeCacheKey  = "nodes:rss"
	fullNodeCacheKey = "nodes:full"
)

var message = "I, %s, am signing this message for registering my intention to operate an RSS3 Serving Node."

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

	fullNode := h.isFullNode(request.Config.Decentralized)

	stat := &schema.Stat{
		Address:      request.Address,
		Endpoint:     request.Endpoint,
		IsPublicGood: nodeInfo.PublicGood,
		ResetAt:      time.Now(),
		IsFullNode:   fullNode,
		IsRssNode:    len(request.Config.RSS) > 0,
		DecentralizedNetwork: len(lo.UniqBy(request.Config.Decentralized, func(module *config.Module) filter.Network {
			return module.Network
		})),
		FederatedNetwork: len(request.Config.Federated),
		Indexer:          len(request.Config.Decentralized),
	}

	indexers := make([]*schema.Indexer, 0, len(request.Config.Decentralized))

	if !fullNode {
		for _, indexer := range request.Config.Decentralized {
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

// Check if node is full node
func (h *Hub) isFullNode(indexers []*config.Module) bool {
	if len(indexers) < len(node.WorkerToNetworksMap) {
		return false
	}

	workerToNetworksMap := make(map[filter.Name]map[string]struct{})

	for _, indexer := range indexers {
		wid, err := filter.NameString(indexer.Worker.String())

		if err != nil {
			return false
		}

		if _, exists := workerToNetworksMap[wid]; !exists {
			workerToNetworksMap[wid] = make(map[string]struct{})
		}

		workerToNetworksMap[wid][indexer.Network.String()] = struct{}{}
	}

	for wid, requiredNetworks := range node.WorkerToNetworksMap {
		networks, exists := workerToNetworksMap[wid]
		if !exists || len(networks) != len(requiredNetworks) {
			return false
		}

		for _, network := range requiredNetworks {
			if _, exists := networks[network]; !exists {
				return false
			}
		}
	}

	return true
}

func (h *Hub) routerRSSHubData(ctx context.Context, path, query string) ([]byte, error) {
	nodes, err := h.filterNodes(ctx, rssNodeCacheKey)

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
	nodes, err := h.filterNodes(ctx, fullNodeCacheKey)

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
	var nodes []node.Cache

	nodes, err := h.filterNodes(ctx, fullNodeCacheKey)

	if err != nil {
		return nil, err
	}

	nodeAddresses, err := h.matchLightNodes(ctx, request)

	if err != nil {
		return nil, err
	}

	// Combine light nodes and full nodes
	if len(nodeAddresses) > 0 {
		nodeStats, err := h.databaseClient.FindNodeStats(ctx, nodeAddresses)

		if err != nil {
			return nil, err
		}

		for i, n := range nodeStats {
			nodes[i].Address = n.Address.String()
			nodes[i].Endpoint = n.Endpoint
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

func (h *Hub) filterNodes(ctx context.Context, key string) ([]node.Cache, error) {
	var (
		nodesCache []node.Cache
		nodes      []*schema.Stat
	)

	// Get nodes from cache.
	exists, err := cache.Get(ctx, key, &nodesCache)
	if err != nil {
		return nil, fmt.Errorf("get nodes from cache: %s, %w", key, err)
	}

	if exists {
		return nodesCache, nil
	}

	// Get nodes from database.
	switch key {
	case rssNodeCacheKey:
		nodes, err = h.databaseClient.FindNodeStatsByType(ctx, nil, lo.ToPtr(true), 3)

		if err != nil {
			return nil, err
		}
	case fullNodeCacheKey:
		nodes, err = h.databaseClient.FindNodeStatsByType(ctx, lo.ToPtr(true), nil, 3)

		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown cache key: %s", key)
	}

	if err = h.setNodeCache(ctx, key, nodes); err != nil {
		return nil, err
	}

	return nodesCache, nil
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
				zap.L().Error("fetch request error", zap.Any("node", address.String()))

				mu.Lock()
				results = append(results, node.DataResponse{Address: address, Err: err})
				mu.Unlock()

				return
			}

			if !h.validateData(data) {
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
		return node.DataResponse{}, fmt.Errorf("timeout waiting for results")
	}
}

func (h *Hub) validateData(data []byte) bool {
	var (
		activityRes node.ActivitiesResponse
		errRes      node.ErrResponse
	)

	if err := json.Unmarshal(data, &errRes); err != nil {
		return false
	}

	if errRes.ErrorCode != "" {
		return false
	}

	if err := json.Unmarshal(data, &activityRes); err != nil {
		return false
	}

	return true
}

func (h *Hub) processRSSHubResults(results []node.DataResponse) {
	ctx := context.Background()

	if err := h.verifyData(ctx, results); err != nil {
		zap.L().Error("fail rss request verify", zap.Any("results", len(results)))
	} else {
		zap.L().Info("complete rss request verify", zap.Any("results", len(results)))
	}
}

func (h *Hub) processActivityResults(results []node.DataResponse) {
	ctx := context.Background()

	if err := h.verifyData(ctx, results); err != nil {
		zap.L().Error("fail feed id request verify", zap.Any("results", len(results)))
	} else {
		zap.L().Info("complete feed id request verify", zap.Any("results", len(results)))
	}
}

func (h *Hub) processActivitiesResults(results []node.DataResponse) {
	ctx := context.Background()

	if err := h.verifyData(ctx, results); err != nil {
		zap.L().Error("fail feed request verify", zap.Any("results", len(results)))
	} else {
		zap.L().Info("complete feed request verify", zap.Any("results", len(results)))
	}

	if !results[0].First {
		return
	}

	// TODO 2nd request、verify
	var activities node.ActivitiesResponse

	data := results[0].Data

	if err := json.Unmarshal(data, &activities); err != nil {
		// TODO: Handle error
		fmt.Printf("fail to unmarshal: %w", err)
	}

	fmt.Println()
}

func (h *Hub) verifyData(ctx context.Context, results []node.DataResponse) error {
	if len(results) < 3 {
		return fmt.Errorf("insufficient data: expected 3 results, got %d", len(results))
	}

	statsMap, err := h.getNodeStatsMap(ctx, results)
	if err != nil {
		return fmt.Errorf("find node stats: %w", err)
	}

	h.sortResultsByFirst(results)

	if !results[0].First {
		for i := range results {
			results[i].InvalidRequest = 1
		}
	} else {
		h.updateRequestsBasedOnDataDiffs(results)
	}

	h.updateStatsWithResults(statsMap, results)

	if err = h.databaseClient.SaveNodeStats(ctx, statsMapToSlice(statsMap)); err != nil {
		return fmt.Errorf("save node stats: %w", err)
	}

	return nil
}

func (h *Hub) getNodeStatsMap(ctx context.Context, results []node.DataResponse) (map[common.Address]*schema.Stat, error) {
	stats, err := h.databaseClient.FindNodeStats(ctx, lo.Map(results, func(result node.DataResponse, _ int) common.Address {
		return result.Address
	}))

	if err != nil {
		return nil, err
	}

	statsMap := make(map[common.Address]*schema.Stat)

	for _, stat := range stats {
		statsMap[stat.Address] = stat
	}

	return statsMap, nil
}

func (h *Hub) sortResultsByFirst(results []node.DataResponse) {
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

func (h *Hub) updateRequestsBasedOnDataDiffs(results []node.DataResponse) {
	diff01 := diffData(results[0].Data, results[1].Data)
	diff02 := diffData(results[0].Data, results[2].Data)
	diff12 := diffData(results[1].Data, results[2].Data)

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
	}
}

func statsMapToSlice(statsMap map[common.Address]*schema.Stat) []*schema.Stat {
	statsSlice := make([]*schema.Stat, 0, len(statsMap))
	for _, stat := range statsMap {
		statsSlice = append(statsSlice, stat)
	}

	return statsSlice
}

func diffData(src, des []byte) bool {
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

// cron task
func (h *Hub) sortNodesTask() error {
	var (
		stats []*schema.Stat

		err error
	)

	ctx := context.Background()

	stats, err = h.databaseClient.FindNodeStats(ctx, []common.Address{})

	if err != nil {
		return err
	}

	// TODO: parallel
	for _, stat := range stats {
		if err = h.updateNodeEpochStats(stat); err != nil {
			return err
		}

		h.calcPoints(stat)

		if err = h.databaseClient.SaveNodeStat(ctx, stat); err != nil {
			return err
		}
	}

	// Update node cache.
	rssNodes, err := h.databaseClient.FindNodeStatsByType(ctx, nil, lo.ToPtr(true), 3)

	if err != nil {
		return err
	}

	if err = h.setNodeCache(ctx, rssNodeCacheKey, rssNodes); err != nil {
		return err
	}

	fullNodes, err := h.databaseClient.FindNodeStatsByType(ctx, lo.ToPtr(true), nil, 3)

	if err != nil {
		return err
	}

	if err = h.setNodeCache(ctx, fullNodeCacheKey, fullNodes); err != nil {
		return err
	}

	return nil
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

func (h *Hub) updateNodeEpochStats(stat *schema.Stat) error {
	nodeInfo, err := h.stakingContract.GetNode(&bind.CallOpts{}, stat.Address)

	if err != nil {
		return fmt.Errorf("get node info: %s,%w", stat.Address.String(), err)
	}

	stat.Staking = float64(nodeInfo.StakingPoolTokens.Uint64())
	stat.EpochRequest = 0
	stat.EpochInvalidRequest = 0

	return nil
}

// calculation rule https://docs.google.com/spreadsheets/d/1N7zEwUooiOjCIHzhoHuf8aM_lbF5bS0ZC-4luxc2qNU/edit?pli=1#gid=0
func (h *Hub) calcPoints(stat *schema.Stat) {
	// staking pool tokens
	stat.Points = math.Min(math.Log2(stat.Staking/100000)+1, 0.2)

	// public good
	stat.Points += float64(lo.Ternary(stat.IsPublicGood, 0, 1))

	// running time
	stat.Points += math.Min(math.Ceil(time.Since(stat.ResetAt).Hours()/18)/120, 0.3)

	// total requests
	stat.Points += math.Min(math.Log(float64(stat.TotalRequest)/100000+1)/math.Log(100), 0.3)

	// epoch requests
	stat.Points += math.Min(math.Log(float64(stat.EpochRequest)/1000000+1)/math.Log(5000), 1)

	// network count
	stat.Points += 0.1*float64(stat.DecentralizedNetwork+stat.FederatedNetwork) + 0.3*float64(lo.Ternary(stat.IsRssNode, 1, 0))

	// indexer count
	stat.Points += math.Min(float64(stat.Indexer)*0.05, 0.2)

	// epoch failure requests
	stat.Points -= 0.5 * float64(stat.EpochInvalidRequest)
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

package hub

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (h *Hub) routerRSSHubData(ctx context.Context, path, query string) ([]byte, error) {
	nodes, err := h.retrieveNodes(ctx, model.RssNodeCacheKey)

	if err != nil {
		return nil, err
	}

	nodeMap, err := h.buildRSSHubPath(path, query, nodes)

	if err != nil {
		return nil, err
	}

	nodeRes, err := h.simpleRouter.DistributeRequest(ctx, nodeMap, h.processRSSHubResults)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeRes.Address.String()))

	return nodeRes.Data, nil
}

func (h *Hub) routerActivityData(ctx context.Context, request ActivityRequest) ([]byte, error) {
	nodes, err := h.retrieveNodes(ctx, model.FullNodeCacheKey)

	if err != nil {
		return nil, err
	}

	nodeMap, err := h.buildActivityByIDPath(request, nodes)

	if err != nil {
		return nil, err
	}

	nodeRes, err := h.simpleRouter.DistributeRequest(ctx, nodeMap, h.processActivityResults)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeRes.Address.String()))

	return nodeRes.Data, nil
}

func (h *Hub) routerActivitiesData(ctx context.Context, request AccountActivitiesRequest) ([]byte, error) {
	nodes := make([]model.Cache, 0, model.DefaultNodeCount)

	nodeAddresses, err := h.matchLightNodes(ctx, request)

	if err != nil {
		return nil, err
	}

	if len(nodeAddresses) > 0 {
		nodeStats, err := h.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
			AddressList: nodeAddresses,
			Limit:       lo.ToPtr(model.DefaultNodeCount),
			PointsOrder: lo.ToPtr("DESC"),
		})

		if err != nil {
			return nil, err
		}

		num := lo.Ternary(len(nodeStats) > model.DefaultNodeCount, model.DefaultNodeCount, len(nodeStats))

		for i := 0; i < num; i++ {
			nodes = append(nodes, model.Cache{
				Address:  nodeStats[i].Address.String(),
				Endpoint: nodeStats[i].Endpoint,
			})
		}
	}

	if len(nodes) < model.DefaultNodeCount {
		fullNodes, err := h.retrieveNodes(ctx, model.FullNodeCacheKey)
		if err != nil {
			return nil, err
		}

		nodesNeeded := model.DefaultNodeCount - len(nodes)
		nodesToAdd := lo.Ternary(nodesNeeded > len(fullNodes), len(fullNodes), nodesNeeded)

		for i := 0; i < nodesToAdd; i++ {
			nodes = append(nodes, fullNodes[i])
		}
	}

	nodeMap, err := h.buildAccountActivitiesPath(request, nodes)

	if err != nil {
		return nil, err
	}

	nodeRes, err := h.simpleRouter.DistributeRequest(ctx, nodeMap, h.processActivitiesResults)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeRes.Address.String()))

	return nodeRes.Data, nil
}

func (h *Hub) buildRSSHubPath(param, query string, nodes []model.Cache) (map[common.Address]string, error) {
	endpointMap, err := model.BuildPath(fmt.Sprintf("/rss/%s?%s", param, query), nil, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

func (h *Hub) buildActivityByIDPath(query ActivityRequest, nodes []model.Cache) (map[common.Address]string, error) {
	endpointMap, err := model.BuildPath(fmt.Sprintf("/decentralized/tx/%s", query.ID), query, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

func (h *Hub) buildAccountActivitiesPath(query AccountActivitiesRequest, nodes []model.Cache) (map[common.Address]string, error) {
	endpointMap, err := model.BuildPath(fmt.Sprintf("/decentralized/%s", query.Account), query, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

func (h *Hub) matchLightNodes(ctx context.Context, request AccountActivitiesRequest) ([]common.Address, error) {
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

		tagWorker, exists := model.TagToWorkersMap[tid]

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

		platformWorker, exists := model.PlatformToWorkerMap[pid]

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
		workers = findCommonStr(lo.Keys(tagWorkers), lo.Keys(platformWorkers))
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

			requiredWorkers := model.NetworkToWorkersMap[nid]

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

			requiredNetworks := model.WorkerToNetworksMap[wid]

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

			workerRequiredNetworks := model.WorkerToNetworksMap[wid]

			requiredNetworks := findCommonStr(workerRequiredNetworks, requestNetworks)

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

func findCommonStr(slice1, slice2 []string) []string {
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

func (h *Hub) retrieveNodes(ctx context.Context, key string) ([]model.Cache, error) {
	var (
		nodesCache []model.Cache
		nodes      []*schema.Stat
	)

	err := h.cacheClient.Get(ctx, key, &nodesCache)

	if err == nil {
		return nodesCache, nil
	}

	zap.L().Info("not found nodes from cache", zap.String("key", key))

	if errors.Is(err, redis.Nil) {
		switch key {
		case model.RssNodeCacheKey:
			nodes, err = h.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
				IsRssNode:    lo.ToPtr(true),
				Limit:        lo.ToPtr(model.DefaultNodeCount),
				ValidRequest: lo.ToPtr(model.DefaultSlashCount),
				PointsOrder:  lo.ToPtr("DESC"),
			})

			if err != nil {
				return nil, err
			}
		case model.FullNodeCacheKey:
			nodes, err = h.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
				IsFullNode:   lo.ToPtr(true),
				Limit:        lo.ToPtr(model.DefaultNodeCount),
				ValidRequest: lo.ToPtr(model.DefaultSlashCount),
				PointsOrder:  lo.ToPtr("DESC"),
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

		nodesCache = lo.Map(nodes, func(n *schema.Stat, _ int) model.Cache {
			return model.Cache{
				Address:  n.Address.String(),
				Endpoint: n.Endpoint,
			}
		})

		return nodesCache, nil
	}

	return nil, fmt.Errorf("get nodes from cache: %s, %w", key, err)
}
func (h *Hub) processRSSHubResults(results []model.DataResponse) {
	if err := h.simpleEnforcer.Verify(context.Background(), results); err != nil {
		zap.L().Error("fail to verify rss hub request", zap.Any("results", len(results)))
	} else {
		zap.L().Info("complete rss hub request verify", zap.Any("results", len(results)))
	}
}
func (h *Hub) processActivityResults(results []model.DataResponse) {
	if err := h.simpleEnforcer.Verify(context.Background(), results); err != nil {
		zap.L().Error("fail to verify  feed id request ", zap.Any("results", len(results)))
	} else {
		zap.L().Info("complete feed id request verify", zap.Any("results", len(results)))
	}
}
func (h *Hub) processActivitiesResults(results []model.DataResponse) {
	ctx := context.Background()

	if err := h.simpleEnforcer.Verify(ctx, results); err != nil {
		zap.L().Error("fail feed request verify", zap.Any("results", len(results)))

		return
	}

	zap.L().Info("complete feed request verify", zap.Any("results", len(results)))

	_ = h.simpleEnforcer.PartialVerify(ctx, results)
}
func (h *Hub) setNodeCache(ctx context.Context, key string, stats []*schema.Stat) error {
	nodesCache := lo.Map(stats, func(n *schema.Stat, _ int) model.Cache {
		return model.Cache{Address: n.Address.String(), Endpoint: n.Endpoint}
	})

	if err := h.cacheClient.Set(ctx, key, nodesCache); err != nil {
		return fmt.Errorf("set nodes to cache: %s, %w", key, err)
	}

	return nil
}

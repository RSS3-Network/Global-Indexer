package enforcer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/rss3-network/node/schema/worker"
	"github.com/rss3-network/node/schema/worker/decentralized"
	"github.com/rss3-network/node/schema/worker/rss"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// maintainNodeWorkerWorker maintains the worker information for the network nodes at every new epoch.
func (e *SimpleEnforcer) maintainNodeWorker(ctx context.Context, epoch int64, stats []*schema.Stat) error {
	nodeToDataMap, fullNodeWorkerToNetworksMap, networkToWorkersMap, platformToWorkersMap, tagToWorkersMap := e.generateMaps(ctx, stats)

	mapTransform(fullNodeWorkerToNetworksMap, networkToWorkersMap, platformToWorkersMap, tagToWorkersMap)

	if err := e.setMapCache(ctx); err != nil {
		return err
	}

	return e.updateNodeWorkers(ctx, stats, nodeToDataMap, epoch)
}

// generateMaps generates the worker related maps.
func (e *SimpleEnforcer) generateMaps(ctx context.Context, stats []*schema.Stat) (map[common.Address]*ComponentInfo, map[decentralized.Worker]map[network.Network]struct{}, map[network.Network]map[decentralized.Worker]struct{}, map[decentralized.Platform]map[decentralized.Worker]struct{}, map[tag.Tag]map[decentralized.Worker]struct{}) {
	var (
		nodeToDataMap               = make(map[common.Address]*ComponentInfo, len(stats))
		fullNodeWorkerToNetworksMap = make(map[decentralized.Worker]map[network.Network]struct{}, len(decentralized.WorkerValues()))
		networkToWorkersMap         = make(map[network.Network]map[decentralized.Worker]struct{}, len(network.NetworkValues()))
		platformToWorkersMap        = make(map[decentralized.Platform]map[decentralized.Worker]struct{}, len(decentralized.PlatformValues()))
		tagToWorkersMap             = make(map[tag.Tag]map[decentralized.Worker]struct{}, len(tag.TagValues()))

		wg sync.WaitGroup
		mu sync.Mutex
	)

	for _, stat := range stats {
		wg.Add(1)

		go func(stat *schema.Stat) {
			defer wg.Done()

			workerStatus, err := e.getNodeWorkerStatus(ctx, stat.Endpoint)
			if err != nil {
				zap.L().Error("get node worker status", zap.Error(err), zap.String("node", stat.Address.String()))

				return
			}

			mu.Lock()
			nodeToDataMap[stat.Address] = workerStatus.Data
			mu.Unlock()

			for _, workerInfo := range workerStatus.Data.Decentralized {
				if workerInfo.Status != worker.StatusReady {
					continue
				}

				if _, ok := networkToWorkersMap[workerInfo.Network]; !ok {
					networkToWorkersMap[workerInfo.Network] = make(map[decentralized.Worker]struct{})
				}

				mu.Lock()
				networkToWorkersMap[workerInfo.Network][workerInfo.Worker] = struct{}{}
				mu.Unlock()

				if _, ok := platformToWorkersMap[workerInfo.Platform]; !ok && workerInfo.Platform != decentralized.PlatformUnknown {
					platformToWorkersMap[workerInfo.Platform] = make(map[decentralized.Worker]struct{})
				}

				mu.Lock()
				if workerInfo.Platform != decentralized.PlatformUnknown {
					platformToWorkersMap[workerInfo.Platform][workerInfo.Worker] = struct{}{}
				}
				mu.Unlock()

				for _, tagX := range workerInfo.Tags {
					if _, ok := tagToWorkersMap[tagX]; !ok {
						tagToWorkersMap[tagX] = make(map[decentralized.Worker]struct{})
					}

					mu.Lock()
					tagToWorkersMap[tagX][workerInfo.Worker] = struct{}{}
					mu.Unlock()
				}

				if _, ok := fullNodeWorkerToNetworksMap[workerInfo.Worker]; !ok {
					fullNodeWorkerToNetworksMap[workerInfo.Worker] = make(map[network.Network]struct{})
				}

				mu.Lock()
				fullNodeWorkerToNetworksMap[workerInfo.Worker][workerInfo.Network] = struct{}{}
				mu.Unlock()
			}
		}(stat)
	}

	wg.Wait()

	return nodeToDataMap, fullNodeWorkerToNetworksMap, networkToWorkersMap, platformToWorkersMap, tagToWorkersMap
}

// mapTransform transforms the map to slice.
func mapTransform(fullNodeWorkerToNetworksMap map[decentralized.Worker]map[network.Network]struct{}, networkToWorkersMap map[network.Network]map[decentralized.Worker]struct{}, platformToWorkersMap map[decentralized.Platform]map[decentralized.Worker]struct{}, tagToWorkersMap map[tag.Tag]map[decentralized.Worker]struct{}) {
	go func() {
		for workerX, networks := range fullNodeWorkerToNetworksMap {
			model.WorkerToNetworksMap[workerX.Name()] = lo.MapToSlice(networks, func(networkX network.Network, _ struct{}) string {
				return networkX.String()
			})
		}
	}()

	go func() {
		for networkX, workers := range networkToWorkersMap {
			model.NetworkToWorkersMap[networkX.String()] = lo.MapToSlice(workers, func(workerX decentralized.Worker, _ struct{}) string {
				return workerX.Name()
			})
		}
	}()

	go func() {
		for platformX, workers := range platformToWorkersMap {
			model.PlatformToWorkersMap[platformX.String()] = lo.MapToSlice(workers, func(workerX decentralized.Worker, _ struct{}) string {
				return workerX.Name()
			})
		}
	}()

	go func() {
		for tagX, workers := range tagToWorkersMap {
			model.TagToWorkersMap[tagX.String()] = lo.MapToSlice(workers, func(workerX decentralized.Worker, _ struct{}) string {
				return workerX.Name()
			})
		}
	}()
}

// setMapCache sets the cache for the worker related maps.
func (e *SimpleEnforcer) setMapCache(ctx context.Context) error {
	var wg sync.WaitGroup

	errChan := make(chan error, 4)

	setCache := func(key string, data interface{}) {
		defer wg.Done()

		if err := e.cacheClient.Set(ctx, key, data); err != nil {
			errChan <- err
		}
	}

	wg.Add(4)

	go setCache(model.WorkerToNetworksMapKey, model.WorkerToNetworksMap)
	go setCache(model.NetworkToWorkersMapKey, model.NetworkToWorkersMap)
	go setCache(model.PlatformToWorkersMapKey, model.PlatformToWorkersMap)
	go setCache(model.TagToWorkersMapKey, model.TagToWorkersMap)

	wg.Wait()

	close(errChan)

	for err := range errChan {
		if err != nil {
			zap.L().Error("Error setting cache", zap.Error(err))
			return err
		}
	}

	return nil
}

// updateNodeWorkers checks whether the node is a full node and updates the worker information to the database.
func (e *SimpleEnforcer) updateNodeWorkers(ctx context.Context, stats []*schema.Stat, nodeToDataMap map[common.Address]*ComponentInfo, epoch int64) error {
	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		workerList = make([]*schema.Worker, 0)
	)

	for i, stat := range stats {
		wg.Add(1)

		go func(i int, stat *schema.Stat) {
			defer wg.Done()

			workerInfo, exist := nodeToDataMap[stat.Address]
			if !exist {
				return
			}

			isFull := determineFullNode(workerInfo.Decentralized)
			stats[i].IsFullNode = isFull

			uniqueNetworks := make(map[network.Network]struct{})
			for _, info := range workerInfo.Decentralized {
				uniqueNetworks[info.Network] = struct{}{}
			}

			stats[i].DecentralizedNetwork = len(uniqueNetworks)
			stats[i].Indexer = len(workerInfo.Decentralized)
			stats[i].IsRssNode = determineRssNode(workerInfo.RSS)
			stats[i].FederatedNetwork = calculateFederatedNetwork(workerInfo.Federated)
			stats[i].Epoch = epoch
			stats[i].EpochRequest = 0
			stats[i].EpochInvalidRequest = 0

			if !isFull {
				mu.Lock()
				workerList = append(workerList, buildNodeWorkers(epoch, stat.Address, workerInfo.Decentralized)...)
				mu.Unlock()
			}
		}(i, stat)
	}

	wg.Wait()

	return e.databaseClient.WithTransaction(ctx, func(ctx context.Context, client database.Client) error {
		if err := client.UpdateNodeWorkerActive(ctx); err != nil {
			return fmt.Errorf("update node worker active: %w", err)
		}

		if err := client.SaveNodeWorkers(ctx, workerList); err != nil {
			return fmt.Errorf("save node workers: %w", err)
		}

		return nil
	})
}

// determineRssNode determines if the node is an RSS node.
func determineRssNode(workers []*RSSWorkerInfo) bool {
	for _, w := range workers {
		if w.Worker == rss.RSSHub && w.Status == worker.StatusReady {
			return true
		}
	}

	return false
}

// calculateFederatedNetwork calculates the federated network.
func calculateFederatedNetwork(_ []*FederatedInfo) int {
	return 0
}

// determineFullNode determines if the node is a full node.
func determineFullNode(workers []*DecentralizedWorkerInfo) bool {
	workers = lo.Filter(workers, func(workerInfo *DecentralizedWorkerInfo, _ int) bool {
		return workerInfo.Status == worker.StatusReady
	})

	if len(workers) < len(model.WorkerToNetworksMap) {
		return false
	}

	workerToNetworksMap := make(map[string]map[string]struct{})

	for _, w := range workers {
		if _, exists := workerToNetworksMap[w.Worker.Name()]; !exists {
			workerToNetworksMap[w.Worker.Name()] = make(map[string]struct{})
		}

		workerToNetworksMap[w.Worker.Name()][w.Network.String()] = struct{}{}
	}

	// Ensure all networks for each worker are present
	for w, requiredNetworks := range model.WorkerToNetworksMap {
		networks, exists := workerToNetworksMap[w]
		if !exists || len(networks) != len(requiredNetworks) {
			return false
		}

		for _, n := range requiredNetworks {
			if _, exists = networks[n]; !exists {
				return false
			}
		}
	}

	return true
}

// getNodeWorkerStatus gets the worker status of the node.
func (e *SimpleEnforcer) getNodeWorkerStatus(ctx context.Context, endpoint string) (*WorkerResponse, error) {
	fullURL := endpoint + "/workers_status"

	body, err := e.httpClient.Fetch(ctx, fullURL)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	response := &WorkerResponse{}

	if err = json.Unmarshal(data, response); err != nil {
		return nil, err
	}

	for i, w := range response.Data.Decentralized {
		if w.Network == network.Farcaster {
			response.Data.Decentralized[i].Platform = decentralized.PlatformFarcaster
			response.Data.Decentralized[i].Tags = []tag.Tag{tag.Social}
		}
	}

	return response, nil
}

// buildNodeWorkers builds the worker information for the node.
func buildNodeWorkers(epoch int64, address common.Address, workerInfo []*DecentralizedWorkerInfo) []*schema.Worker {
	workers := make([]*schema.Worker, 0, len(workerInfo))

	for _, w := range workerInfo {
		workers = append(workers, &schema.Worker{
			EpochID:  uint64(epoch),
			Address:  address,
			Network:  w.Network.String(),
			Name:     w.Worker.Name(),
			IsActive: true,
		})
	}

	return workers
}

func (e *SimpleEnforcer) initWorkerMap(ctx context.Context) error {
	var wg sync.WaitGroup

	errChan := make(chan error, 4)

	getCache := func(key string, data interface{}) {
		defer wg.Done()

		if err := e.cacheClient.Get(ctx, key, data); err != nil {
			errChan <- err
		}
	}

	wg.Add(4)

	go getCache(model.WorkerToNetworksMapKey, model.WorkerToNetworksMap)
	go getCache(model.NetworkToWorkersMapKey, model.NetworkToWorkersMap)
	go getCache(model.PlatformToWorkersMapKey, model.PlatformToWorkersMap)
	go getCache(model.TagToWorkersMapKey, model.TagToWorkersMap)

	wg.Wait()

	close(errChan)

	for err := range errChan {
		if err != nil {
			if errors.Is(err, redis.Nil) {
				epoch, err := e.getCurrentEpoch(ctx)
				if err != nil {
					return err
				}

				stats, err := e.getAllNodeStats(ctx)
				if err != nil {
					return err
				}

				return e.maintainNodeWorker(ctx, epoch, stats)
			}

			zap.L().Error("Error setting cache", zap.Error(err))

			return err
		}
	}

	return nil
}

type WorkerInfo struct {
	Network network.Network `json:"network"`
	Tags    []tag.Tag       `json:"tags"`
	Status  worker.Status   `json:"status"`
	//RemoteState  uint64          `json:"remote_state"`
	//IndexedState uint64          `json:"indexed_state"`
}

type DecentralizedWorkerInfo struct {
	WorkerInfo
	Worker   decentralized.Worker   `json:"worker"`
	Platform decentralized.Platform `json:"platform"`
}

type RSSWorkerInfo struct {
	WorkerInfo
	Worker rss.Worker `json:"worker"`
}

type FederatedInfo struct {
	WorkerInfo
}

type ComponentInfo struct {
	Decentralized []*DecentralizedWorkerInfo `json:"decentralized"`
	RSS           []*RSSWorkerInfo           `json:"rss"`
	Federated     []*FederatedInfo           `json:"federated"`
}

type WorkerResponse struct {
	Data *ComponentInfo `json:"data"`
}

// UpdateNodeCache updates the cache for all Nodes.
// 1. update the sorted set nodes.
// 2. update the cache for the Node subscription.
func (e *SimpleEnforcer) updateNodeCache(ctx context.Context, epoch int64) error {
	for _, key := range []string{model.RssNodeCacheKey, model.FullNodeCacheKey} {
		if err := e.updateSortedSetForNodeType(ctx, key); err != nil {
			return err
		}
	}

	return e.cacheClient.Set(ctx, model.SubscribeNodeCacheKey, epoch)
}

// updateSortedSetForNodeType updates the sorted set for different types of Nodes.
func (e *SimpleEnforcer) updateSortedSetForNodeType(ctx context.Context, key string) error {
	nodesEndpointCaches, err := retrieveNodeEndpointCaches(ctx, key, e.databaseClient)
	if err != nil {
		return err
	}

	nodesEndpointCachesMap := lo.SliceToMap(nodesEndpointCaches, func(node *model.NodeEndpointCache) (string, *model.NodeEndpointCache) {
		return node.Address, node
	})

	members, err := e.cacheClient.ZRevRangeWithScores(ctx, key, 0, -1)
	if err != nil {
		return err
	}

	membersToRemove := make([]string, 0)
	membersToAdd := make([]redis.Z, 0, len(nodesEndpointCachesMap))

	for _, member := range members {
		if _, ok := nodesEndpointCachesMap[member.Member.(string)]; !ok {
			membersToRemove = append(membersToRemove, member.Member.(string))
		}
	}

	for _, node := range nodesEndpointCaches {
		membersToAdd = append(membersToAdd, redis.Z{
			Member: node.Address,
			Score:  node.Score,
		})
	}

	if len(membersToAdd) > 0 {
		if err = e.cacheClient.ZAdd(ctx, key, membersToAdd...); err != nil {
			return err
		}
	}

	if len(membersToRemove) == 0 {
		return nil
	}

	return e.cacheClient.ZRem(ctx, key, membersToRemove)
}

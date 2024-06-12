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
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// maintainNodeWorkerWorker maintains the worker information for the network nodes at every new epoch.
func (e *SimpleEnforcer) maintainNodeWorker(ctx context.Context, epoch int64, stats []*schema.Stat) error {
	nodeToWorkersMap, fullNodeWorkerToNetworksMap, networkToWorkersMap, platformToWorkersMap, tagToWorkersMap := e.generateMaps(ctx, stats)

	mapTransform(fullNodeWorkerToNetworksMap, networkToWorkersMap, platformToWorkersMap, tagToWorkersMap)

	if err := e.setMapCache(ctx); err != nil {
		return err
	}

	return e.updateNodeWorkers(ctx, stats, nodeToWorkersMap, epoch)
}

// generateMaps generates the worker related maps.
func (e *SimpleEnforcer) generateMaps(ctx context.Context, stats []*schema.Stat) (map[common.Address][]*WorkerInfo, map[worker.Worker]map[network.Network]struct{}, map[network.Network]map[worker.Worker]struct{}, map[decentralized.Platform]map[worker.Worker]struct{}, map[tag.Tag]map[worker.Worker]struct{}) {
	var (
		nodeToWorkersMap            = make(map[common.Address][]*WorkerInfo, len(stats))
		fullNodeWorkerToNetworksMap = make(map[worker.Worker]map[network.Network]struct{}, len(decentralized.WorkerValues()))
		networkToWorkersMap         = make(map[network.Network]map[worker.Worker]struct{}, len(network.NetworkValues()))
		platformToWorkersMap        = make(map[decentralized.Platform]map[worker.Worker]struct{}, len(decentralized.PlatformValues()))
		tagToWorkersMap             = make(map[tag.Tag]map[worker.Worker]struct{}, len(tag.TagValues()))

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

			nodeToWorkersMap[stat.Address] = workerStatus.Data

			for _, workerInfo := range workerStatus.Data {
				if workerInfo.Status != worker.StatusReady {
					continue
				}

				if _, ok := networkToWorkersMap[workerInfo.Network]; !ok {
					networkToWorkersMap[workerInfo.Network] = make(map[worker.Worker]struct{})
				}

				mu.Lock()
				networkToWorkersMap[workerInfo.Network][workerInfo.Worker] = struct{}{}
				mu.Unlock()

				if _, ok := platformToWorkersMap[workerInfo.Platform]; !ok && workerInfo.Platform != decentralized.PlatformUnknown {
					platformToWorkersMap[workerInfo.Platform] = make(map[worker.Worker]struct{})
				}

				mu.Lock()
				if workerInfo.Platform != decentralized.PlatformUnknown {
					platformToWorkersMap[workerInfo.Platform][workerInfo.Worker] = struct{}{}
				}
				mu.Unlock()

				for _, tagX := range workerInfo.Tags {
					if _, ok := tagToWorkersMap[tagX]; !ok {
						tagToWorkersMap[tagX] = make(map[worker.Worker]struct{})
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

	return nodeToWorkersMap, fullNodeWorkerToNetworksMap, networkToWorkersMap, platformToWorkersMap, tagToWorkersMap
}

// mapTransform transforms the map to slice.
func mapTransform(fullNodeWorkerToNetworksMap map[worker.Worker]map[network.Network]struct{}, networkToWorkersMap map[network.Network]map[worker.Worker]struct{}, platformToWorkersMap map[decentralized.Platform]map[worker.Worker]struct{}, tagToWorkersMap map[tag.Tag]map[worker.Worker]struct{}) {
	go func() {
		for workerX, networks := range fullNodeWorkerToNetworksMap {
			model.WorkerToNetworksMap[workerX.Name()] = lo.MapToSlice(networks, func(networkX network.Network, _ struct{}) string {
				return networkX.String()
			})
		}
	}()

	go func() {
		for networkX, workers := range networkToWorkersMap {
			model.NetworkToWorkersMap[networkX.String()] = lo.MapToSlice(workers, func(workerX worker.Worker, _ struct{}) string {
				return workerX.Name()
			})
		}
	}()

	go func() {
		for platformX, workers := range platformToWorkersMap {
			model.PlatformToWorkersMap[platformX.String()] = lo.MapToSlice(workers, func(workerX worker.Worker, _ struct{}) string {
				return workerX.Name()
			})
		}
	}()

	go func() {
		for tagX, workers := range tagToWorkersMap {
			model.TagToWorkersMap[tagX.String()] = lo.MapToSlice(workers, func(workerX worker.Worker, _ struct{}) string {
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
func (e *SimpleEnforcer) updateNodeWorkers(ctx context.Context, stats []*schema.Stat, nodeToWorkersMap map[common.Address][]*WorkerInfo, epoch int64) error {
	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		workerList = make([]*schema.Worker, 0)
	)

	for i, stat := range stats {
		wg.Add(1)

		go func(i int, stat *schema.Stat) {
			defer wg.Done()

			workerInfo, exist := nodeToWorkersMap[stat.Address]
			if !exist {
				return
			}

			isFull := determineFullNode(workerInfo)
			stats[i].IsFullNode = isFull

			uniqueNetworks := make(map[network.Network]struct{})
			for _, info := range workerInfo {
				uniqueNetworks[info.Network] = struct{}{}
			}

			stats[i].DecentralizedNetwork = len(uniqueNetworks)
			stats[i].Indexer = len(workerInfo)
			stats[i].IsRssNode = determineRssNode(workerInfo)
			stats[i].FederatedNetwork = calculateFederatedNetwork(workerInfo)

			if !isFull {
				mu.Lock()
				workerList = append(workerList, buildNodeWorkers(epoch, stat.Address, workerInfo)...)
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
// rss node is enabled by default now.
func determineRssNode(_ []*WorkerInfo) bool {
	return true
}

// calculateFederatedNetwork calculates the federated network.
func calculateFederatedNetwork(_ []*WorkerInfo) int {
	return 0
}

// determineFullNode determines if the node is a full node.
func determineFullNode(workers []*WorkerInfo) bool {
	workers = lo.Filter(workers, func(workerInfo *WorkerInfo, _ int) bool {
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

	for i, w := range response.Data {
		if w.Network == network.Farcaster {
			response.Data[i].Platform = decentralized.PlatformFarcaster
			response.Data[i].Tags = []tag.Tag{tag.Social}
		}
	}

	return response, nil
}

// buildNodeWorkers builds the worker information for the node.
func buildNodeWorkers(epoch int64, address common.Address, workerInfo []*WorkerInfo) []*schema.Worker {
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

type WorkerResponse struct {
	Data []*WorkerInfo `json:"data"`
}

type WorkerInfo struct {
	Network  network.Network        `json:"network"`
	Worker   worker.Worker          `json:"worker"`
	Tags     []tag.Tag              `json:"tags"`
	Platform decentralized.Platform `json:"platform"`
	Status   worker.Status          `json:"status"`
}

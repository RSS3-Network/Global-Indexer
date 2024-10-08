package enforcer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-version"
	"github.com/redis/go-redis/v9"
	v2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"github.com/rss3-network/global-indexer/internal/cache"
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
	"golang.org/x/sync/errgroup"
)

// maintainNodeWorkerWorker maintains the worker information for network nodes at each new epoch.
func (e *SimpleEnforcer) maintainNodeWorker(ctx context.Context, epoch int64, stats []*schema.Stat) error {
	addresses := lo.Map(stats, func(stat *schema.Stat, _ int) common.Address {
		return stat.Address
	})

	var (
		nodeInfo      []*schema.Node
		nodeVSLInfo   []v2.Node
		minVersionStr string
		err           error
	)

	eg, cCtx := errgroup.WithContext(ctx)

	// get node info from database
	eg.Go(func() error {
		nodeInfo, err = e.databaseClient.FindNodes(cCtx, schema.FindNodesQuery{NodeAddresses: addresses})
		return err
	})

	// get node info from chain
	eg.Go(func() error {
		nodeVSLInfo, err = e.stakingContract.GetNodes(&bind.CallOpts{}, addresses)
		return err
	})

	// get min version
	eg.Go(func() error {
		minVersionStr, err = e.getNodeMinVersion()
		return err
	})

	if err = eg.Wait(); err != nil {
		return fmt.Errorf("get node info: %w", err)
	}

	nodeInfoMap := lo.SliceToMap(nodeInfo, func(node *schema.Node) (common.Address, *schema.Node) {
		return node.Address, node
	})

	originalStatusList := make([]schema.NodeStatus, len(stats))

	for i, stat := range stats {
		stat.Status = schema.NodeStatus(nodeVSLInfo[i].Status)
		stat.Version = nodeInfoMap[stat.Address].Version
		stat.HearBeat = nodeInfoMap[stat.Address].Status
		originalStatusList[i] = stat.Status
	}

	// Initialize maps related to worker data.
	nodeToDataMap, fullNodeWorkerToNetworksMap, networkToWorkersMap, platformToWorkersMap, tagToWorkersMap := e.generateMaps(ctx, stats, minVersionStr)
	// Transform the map and assigns the result to the global variable.
	mapTransformAssign(fullNodeWorkerToNetworksMap, networkToWorkersMap, platformToWorkersMap, tagToWorkersMap)

	// Set cache data to persist across program restarts or refresh at the start of each new epoch.
	if err := e.setMapCache(ctx); err != nil {
		return fmt.Errorf("set map cache: %w", err)
	}
	// Update node statistics and worker data.
	if err := e.updateNodeWorkers(ctx, stats, nodeToDataMap, epoch); err != nil {
		return fmt.Errorf("update node workers: %w", err)
	}
	// filter the exited node status and submit the demotion to the VSL.
	// Fixme: deprecated if there is no alpha type node
	nodeAddresses, nodeStatusList, err := e.filterNodeStatus(ctx, stats, originalStatusList)
	if err != nil {
		return fmt.Errorf("filter node status: %w", err)
	}

	return e.updateNodeStatusAndSubmitDemotionToVSL(ctx, nodeAddresses, nodeStatusList, nil, nil, nil)
}

func (e *SimpleEnforcer) filterNodeStatus(ctx context.Context, stats []*schema.Stat, originalStatusList []schema.NodeStatus) ([]common.Address, []uint8, error) {
	nodeAddresses := make([]common.Address, 0, len(stats))
	nodeStatusList := make([]uint8, 0, len(stats))

	for i := range stats {
		if stats[i].Status != originalStatusList[i] {
			nodeAddresses = append(nodeAddresses, stats[i].Address)
			nodeStatusList = append(nodeStatusList, uint8(stats[i].Status))
		}
	}

	// filter the exited node status and submit the demotion to the VSL.
	// Fixme: deprecated if there is no alpha type node
	alphaNodes, err := e.databaseClient.FindNodes(ctx, schema.FindNodesQuery{
		Type: lo.ToPtr(schema.NodeTypeAlpha),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("find alpha nodes: %w", err)
	}

	alphaNodesVSLInfo, err := e.stakingContract.GetNodes(&bind.CallOpts{}, lo.Map(alphaNodes, func(node *schema.Node, _ int) common.Address {
		return node.Address
	}))

	if err != nil {
		return nil, nil, fmt.Errorf("get alpha nodes from chain: %w", err)
	}

	for i := range alphaNodesVSLInfo {
		if alphaNodesVSLInfo[i].Status == uint8(schema.NodeStatusExiting) {
			nodeAddresses = append(nodeAddresses, alphaNodes[i].Address)
			nodeStatusList = append(nodeStatusList, uint8(schema.NodeStatusExited))
		}
	}

	return nodeAddresses, nodeStatusList, nil
}

// generateMaps generates maps related to worker data.
func (e *SimpleEnforcer) generateMaps(ctx context.Context, stats []*schema.Stat, minVersionStr string) (map[common.Address]*ComponentInfo, map[string]map[string]struct{}, map[string]map[string]struct{}, map[string]map[string]struct{}, map[string]map[string]struct{}) {
	var (
		// nodeToDataMap stores the API response from /workers_status for each node.
		nodeToDataMap = make(map[common.Address]*ComponentInfo, len(stats))
		// fullNodeWorkerToNetworksMap maps each worker in a full node to their respective networks that are fully supported by the entire network.
		// A node qualifies as a full node if it includes all workers, with each worker fully supporting its designated network.
		fullNodeWorkerToNetworksMap = make(map[string]map[string]struct{}, len(decentralized.WorkerValues()))
		// networkToWorkersMap maps networks to workers that are supported across the entire network.
		networkToWorkersMap = make(map[string]map[string]struct{}, len(network.NetworkValues()))
		// platformToWorkersMap maps platforms to workers that are supported across the entire network.
		platformToWorkersMap = make(map[string]map[string]struct{}, len(decentralized.PlatformValues()))
		// tagToWorkersMap maps tags to workers that are supported across the entire network.
		tagToWorkersMap = make(map[string]map[string]struct{}, len(tag.TagValues()))

		wg sync.WaitGroup
		mu sync.Mutex
	)

	minVersion, _ := version.NewVersion(minVersionStr)

	for _, stat := range stats {
		wg.Add(1)

		go func(stat *schema.Stat) {
			defer wg.Done()

			// Skip processing the node if its status is marked as exited, offline, slashing or slashed.
			if !isValidNodeStatus(stat.Status) {
				return
			}

			// If the node status is exiting, set it to exited.
			if stat.Status == schema.NodeStatusExiting {
				stat.Status = schema.NodeStatusExited
				return
			}

			// Set the node status to online.
			stat.Status = schema.NodeStatusOnline

			if nodeVersion, _ := version.NewVersion(stat.Version); nodeVersion.LessThan(minVersion) {
				// Set the node status to outdated.
				stat.Status = schema.NodeStatusOutdated

				// Disqualified the node from the current request distribution round
				// if retrieving the node info fails.
				return
			}

			// Retrieve the status of the node's worker,
			// including details like name, network, tags, and platform information.
			workerStatus, err := e.getNodeWorkerStatus(ctx, stat.Endpoint, stat.AccessToken)
			if err != nil || workerStatus == nil {
				zap.L().Error("get node worker status", zap.Error(err), zap.String("node", stat.Address.String()))

				if stat.HearBeat == schema.NodeStatusOffline {
					stat.Status = schema.NodeStatusOffline
				} else {
					stat.Status = schema.NodeStatusInitializing
				}

				return
			}

			mu.Lock()
			workerStatus.Data.Decentralized = filterDuplicateWorkers(workerStatus.Data.Decentralized)
			nodeToDataMap[stat.Address] = workerStatus.Data
			mu.Unlock()

			// if all workers are unhealthy, the node is registered
			isRegistered := true

			for _, workerInfo := range workerStatus.Data.Decentralized {
				// Skip processing the worker if its status is not marked as ready.
				if workerInfo.Status != worker.StatusReady {
					if workerInfo.Status == worker.StatusIndexing {
						isRegistered = false
						stat.Status = schema.NodeStatusInitializing
					}

					continue
				}

				isRegistered = false
				networkName := workerInfo.Network.String()
				platformName := workerInfo.Platform.String()
				workerName := workerInfo.Worker.String()

				if name, exist := model.RenameWorkerMap[workerInfo.Network]; exist {
					workerName = name
				}

				mu.Lock()
				if _, ok := networkToWorkersMap[networkName]; !ok {
					networkToWorkersMap[networkName] = make(map[string]struct{})
				}

				networkToWorkersMap[workerInfo.Network.String()][workerName] = struct{}{}
				mu.Unlock()

				mu.Lock()
				if _, ok := platformToWorkersMap[platformName]; !ok && platformName != decentralized.PlatformUnknown.String() {
					platformToWorkersMap[workerInfo.Platform.String()] = make(map[string]struct{})
				}

				if platformName != decentralized.PlatformUnknown.String() {
					platformToWorkersMap[platformName][workerName] = struct{}{}
				}
				mu.Unlock()

				for _, tagX := range workerInfo.Tags {
					mu.Lock()
					tagName := tagX.String()

					if _, ok := tagToWorkersMap[tagName]; !ok {
						tagToWorkersMap[tagName] = make(map[string]struct{})
					}

					tagToWorkersMap[tagName][workerName] = struct{}{}
					mu.Unlock()
				}

				mu.Lock()
				if _, ok := fullNodeWorkerToNetworksMap[workerName]; !ok {
					fullNodeWorkerToNetworksMap[workerName] = make(map[string]struct{})
				}

				fullNodeWorkerToNetworksMap[workerName][networkName] = struct{}{}
				mu.Unlock()
			}

			if isRegistered {
				stat.Status = schema.NodeStatusRegistered
			}
		}(stat)
	}

	wg.Wait()

	return nodeToDataMap, fullNodeWorkerToNetworksMap, networkToWorkersMap, platformToWorkersMap, tagToWorkersMap
}

// isValidNodeStatus checks if the node status is valid.
func isValidNodeStatus(status schema.NodeStatus) bool {
	switch status {
	case schema.NodeStatusOnline, schema.NodeStatusInitializing, schema.NodeStatusOutdated, schema.NodeStatusRegistered, schema.NodeStatusExiting:
		return true
	default:
		return false
	}
}

// filterDuplicateWorkers filters out duplicate workers.
func filterDuplicateWorkers(workers []*DecentralizedWorkerInfo) []*DecentralizedWorkerInfo {
	seen := make(map[string]*DecentralizedWorkerInfo, len(workers))

	for _, workerInfo := range workers {
		key := fmt.Sprintf("%s-%s", workerInfo.Network.String(), workerInfo.Worker.String())

		if existingWorker, ok := seen[key]; ok {
			if existingWorker.Status == worker.StatusReady && workerInfo.Status == worker.StatusReady {
				continue
			}

			if workerInfo.Status != worker.StatusReady {
				seen[key] = workerInfo
			}
		} else {
			seen[key] = workerInfo
		}
	}

	filteredWorkers := make([]*DecentralizedWorkerInfo, 0, len(seen))

	for _, workerInfo := range seen {
		filteredWorkers = append(filteredWorkers, workerInfo)
	}

	return filteredWorkers
}

// mapTransformAssign transforms the map and assigns the result to a global variable.
func mapTransformAssign(fullNodeWorkerToNetworksMap, networkToWorkersMap, platformToWorkersMap, tagToWorkersMap map[string]map[string]struct{}) {
	var (
		wg  sync.WaitGroup
		mux sync.Mutex
	)

	srcMaps := []map[string]map[string]struct{}{
		fullNodeWorkerToNetworksMap,
		networkToWorkersMap,
		platformToWorkersMap,
		tagToWorkersMap,
	}

	desMaps := []*map[string][]string{
		&model.WorkerToNetworksMap,
		&model.NetworkToWorkersMap,
		&model.PlatformToWorkersMap,
		&model.TagToWorkersMap,
	}

	transformAndAssign := func(srcMap map[string]map[string]struct{}, targetMap *map[string][]string) {
		localMap := make(map[string][]string)

		for key, value := range srcMap {
			localMap[key] = lo.MapToSlice(value, func(s string, _ struct{}) string {
				return s
			})
		}

		mux.Lock()
		*targetMap = localMap
		mux.Unlock()
	}

	for i := range srcMaps {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			transformAndAssign(srcMaps[i], desMaps[i])
		}(i)
	}

	wg.Wait()
}

// setMapCache caches worker-related maps for use in each epoch and retains them across program restarts.
func (e *SimpleEnforcer) setMapCache(ctx context.Context) error {
	var wg sync.WaitGroup

	keys := []string{
		model.WorkerToNetworksMapKey,
		model.NetworkToWorkersMapKey,
		model.PlatformToWorkersMapKey,
		model.TagToWorkersMapKey,
	}

	maps := []interface{}{
		&model.WorkerToNetworksMap,
		&model.NetworkToWorkersMap,
		&model.PlatformToWorkersMap,
		&model.TagToWorkersMap,
	}

	errChan := make(chan error, len(keys))

	for i := range keys {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			if err := e.cacheClient.Set(ctx, keys[i], maps[i], 0); err != nil {
				errChan <- err
			}
		}(i)
	}

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

// updateNodeWorkers checks if the node is a full node and updates the corresponding worker information in the database.
func (e *SimpleEnforcer) updateNodeWorkers(ctx context.Context, stats []*schema.Stat, nodeToDataMap map[common.Address]*ComponentInfo, epoch int64) error {
	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		workerList = make([]*schema.Worker, 0)
	)

	for i := range stats {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			workerInfo, exist := nodeToDataMap[stats[i].Address]
			// Disqualifies the node from the current request distribution round if it does not exist,
			// indicating a failure to retrieve the current epoch's status.
			if !exist {
				stats[i].EpochInvalidRequest = int64(model.DemotionCountBeforeSlashing)

				return
			}

			// Determine whether the node qualifies as a full node.
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

			// Reset the epoch, request count, and invalid request count if a new epoch is detected,
			// different from the previous one.
			if epoch != stats[i].Epoch {
				stats[i].Epoch = epoch
				stats[i].EpochInvalidRequest = 0
			}

			// Update worker information in the database if the node is not a full node.
			if !isFull {
				mu.Lock()
				workerList = append(workerList, buildNodeWorkers(epoch, stats[i].Address, workerInfo.Decentralized)...)
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	return e.databaseClient.WithTransaction(ctx, func(ctx context.Context, client database.Client) error {
		// Set the 'active' status to false for all workers from outdated epochs.
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
func determineRssNode(w *RSSWorkerInfo) bool {
	if w != nil && w.Worker == rss.RSSHub && w.Status == worker.StatusReady {
		return true
	}

	return false
}

// calculateFederatedNetwork calculates the federated network.
func calculateFederatedNetwork(_ []*FederatedInfo) int {
	return 0
}

// determineFullNode determines if the node is a full node based on
// whether the workers' information matches the WorkerToNetworksMap exactly.
func determineFullNode(workers []*DecentralizedWorkerInfo) bool {
	workers = lo.Filter(workers, func(workerInfo *DecentralizedWorkerInfo, _ int) bool {
		return workerInfo.Status == worker.StatusReady
	})

	if len(workers) < len(model.WorkerToNetworksMap) {
		return false
	}

	workerToNetworksMap := make(map[string]map[string]struct{})

	for _, w := range workers {
		workerName := w.Worker.String()
		networkName := w.Network.String()

		if name, exist := model.RenameWorkerMap[w.Network]; exist {
			workerName = name
		}

		if _, exists := workerToNetworksMap[workerName]; !exists {
			workerToNetworksMap[workerName] = make(map[string]struct{})
		}

		workerToNetworksMap[workerName][networkName] = struct{}{}
	}

	// Ensure each worker has all required networks present.
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

// buildNodeWorkers builds and populates worker information for the node.
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

// getWorkerMapFromCache retrieves the worker map from the cache.
func getWorkerMapFromCache(ctx context.Context, cacheClient cache.Client) chan error {
	var wg sync.WaitGroup

	keys := []string{
		model.WorkerToNetworksMapKey,
		model.NetworkToWorkersMapKey,
		model.PlatformToWorkersMapKey,
		model.TagToWorkersMapKey,
	}

	maps := []interface{}{
		&model.WorkerToNetworksMap,
		&model.NetworkToWorkersMap,
		&model.PlatformToWorkersMap,
		&model.TagToWorkersMap,
	}
	errChan := make(chan error, len(keys))

	for i := range keys {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			if err := cacheClient.Get(ctx, keys[i], maps[i]); err != nil {
				errChan <- err
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	return errChan
}

// initWorkerMap initializes the worker map on first startup or when cache data is lost.
// TODO Implement reverse recovery from the latest epoch's node data if cache data is lost.
func (e *SimpleEnforcer) initWorkerMap(ctx context.Context) error {
	errChan := getWorkerMapFromCache(ctx, e.cacheClient)

	for err := range errChan {
		if err != nil {
			if errors.Is(err, redis.Nil) {
				epoch, err := e.getCurrentEpoch(ctx)
				if err != nil {
					return err
				}

				stats, err := e.getAllNodeStats(ctx, &schema.StatQuery{
					Limit: lo.ToPtr(defaultLimit),
				})

				if err != nil {
					return err
				}

				if err = e.maintainNodeWorker(ctx, epoch, stats); err != nil {
					return err
				}

				return e.processNodeStats(ctx, stats, true)
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
	RSS           *RSSWorkerInfo             `json:"rss"`
	Federated     []*FederatedInfo           `json:"federated"`
}

type WorkersStatusResponse struct {
	Data *ComponentInfo `json:"data"`
}

type InfoResponse struct {
	Data struct {
		Version struct {
			Tag    string `json:"tag"`
			Commit string `json:"commit"`
		} `json:"version"`
	} `json:"data"`
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

	return e.cacheClient.Set(ctx, model.SubscribeNodeCacheKey, epoch, 0)
}

// updateSortedSetForNodeType updates the sorted set for different types of Nodes.
func (e *SimpleEnforcer) updateSortedSetForNodeType(ctx context.Context, key string) error {
	nodeStats, err := retrieveNodeStatsFromDB(ctx, key, e.databaseClient)
	if err != nil {
		return err
	}

	nodesEndpointCachesMap := lo.SliceToMap(nodeStats, func(stat *schema.Stat) (string, *EndpointCache) {
		return stat.Address.String(), &EndpointCache{
			Endpoint:    stat.Endpoint,
			AccessToken: stat.AccessToken,
		}
	})

	// Get current members from Redis sorted set
	members, err := e.cacheClient.ZRevRangeWithScores(ctx, key, 0, -1)
	if err != nil {
		return err
	}

	// Prepare batch operations for Redis Sorted Set
	membersToAdd, membersToRemove := prepareMembers(members, nodesEndpointCachesMap, nodeStats)

	if len(membersToRemove) > 0 {
		if err = e.cacheClient.ZRem(ctx, key, membersToRemove); err != nil {
			return fmt.Errorf("failed to remove old members: %w", err)
		}
	}

	if len(membersToAdd) > 0 {
		if err = e.cacheClient.ZAdd(ctx, key, membersToAdd...); err != nil {
			return fmt.Errorf("failed to add new members: %w", err)
		}
	}

	return nil
}

// prepareMembers prepares the members to add and remove in the sorted set.
func prepareMembers(members []redis.Z, nodesEndpointCachesMap map[string]*EndpointCache, nodeStats []*schema.Stat) ([]redis.Z, []string) {
	membersToRemove := filterMembersToRemove(members, nodesEndpointCachesMap)

	membersToAdd := make([]redis.Z, 0)

	for _, stat := range nodeStats {
		if stat.EpochInvalidRequest < int64(model.DemotionCountBeforeSlashing) {
			membersToAdd = append(membersToAdd, redis.Z{
				Member: stat.Address.String(),
				Score:  stat.Score,
			})
		}
	}

	return membersToAdd, membersToRemove
}

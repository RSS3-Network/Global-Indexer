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
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/rss3-network/protocol-go/schema/metadata"
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

	nodeRes, err := h.batchRequest(ctx, nodeMap, h.processRSSHubResults)

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

	nodeRes, err := h.batchRequest(ctx, nodeMap, h.processActivityResults)

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

	nodeRes, err := h.batchRequest(ctx, nodeMap, h.processActivitiesResults)

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

func (h *Hub) batchRequest(_ context.Context, nodeMap map[common.Address]string, processResults func([]DataResponse)) (DataResponse, error) {
	var (
		waitGroup   sync.WaitGroup
		firstResult = make(chan DataResponse, 1)
		results     []DataResponse
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
				results = append(results, DataResponse{Address: address, Err: err})

				if len(results) == len(nodeMap) {
					firstResult <- DataResponse{Address: address, Data: []byte(model.MessageNodeDataFailed)}
				}

				mu.Unlock()

				return
			}

			flagActivities, _ := h.validateActivities(data)
			flagActivity, _ := h.validateActivity(data)

			if !flagActivities && !flagActivity {
				zap.L().Error("response parse error", zap.Any("node", address.String()))

				mu.Lock()
				results = append(results, DataResponse{Address: address, Err: fmt.Errorf("invalid data")})

				if len(results) == len(nodeMap) {
					firstResult <- DataResponse{Address: address, Data: data}
				}
				mu.Unlock()

				return
			}

			mu.Lock()
			results = append(results, DataResponse{Address: address, Data: data, First: true})
			mu.Unlock()

			select {
			case firstResult <- DataResponse{Address: address, Data: data}:
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
		return DataResponse{Data: []byte(model.MessageNodeDataFailed)}, fmt.Errorf("timeout waiting for results")
	}
}

func (h *Hub) validateActivities(data []byte) (bool, *ActivitiesResponse) {
	var (
		res    ActivitiesResponse
		errRes ErrResponse
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

func (h *Hub) validateActivity(data []byte) (bool, *ActivityResponse) {
	var (
		res    ActivityResponse
		errRes ErrResponse
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

func (h *Hub) processRSSHubResults(results []DataResponse) {
	if err := h.verifyData(context.Background(), results); err != nil {
		zap.L().Error("fail to verify rss hub request", zap.Any("results", len(results)))
	} else {
		zap.L().Info("complete rss hub request verify", zap.Any("results", len(results)))
	}
}

func (h *Hub) processActivityResults(results []DataResponse) {
	if err := h.verifyData(context.Background(), results); err != nil {
		zap.L().Error("fail to verify  feed id request ", zap.Any("results", len(results)))
	} else {
		zap.L().Info("complete feed id request verify", zap.Any("results", len(results)))
	}
}

func (h *Hub) processActivitiesResults(results []DataResponse) {
	ctx := context.Background()

	if err := h.verifyData(ctx, results); err != nil {
		zap.L().Error("fail feed request verify", zap.Any("results", len(results)))

		return
	}

	zap.L().Info("complete feed request verify", zap.Any("results", len(results)))

	if !results[0].First {
		return
	}

	var activities ActivitiesResponse

	data := results[0].Data

	if err := json.Unmarshal(data, &activities); err != nil {
		zap.L().Error("fail to unmarshall activities")

		return
	}

	// data is empty, no need to 2nd verify
	if activities.Data == nil {
		return
	}

	workingNodes := lo.Map(results, func(result DataResponse, _ int) common.Address {
		return result.Address
	})

	h.process2ndVerify(activities.Data, workingNodes)
}

func (h *Hub) process2ndVerify(feeds []*Feed, workingNodes []common.Address) {
	ctx := context.Background()
	platformMap := make(map[string]struct{})
	statMap := make(map[string]struct{})

	for _, feed := range feeds {
		if len(feed.Platform) == 0 {
			continue
		}

		h.verifyPlatform(ctx, feed, platformMap, statMap, workingNodes)

		if _, exists := platformMap[feed.Platform]; !exists {
			if len(platformMap) == model.DefaultVerifyCount {
				break
			}
		}
	}
}

func (h *Hub) verifyPlatform(ctx context.Context, feed *Feed, platformMap, statMap map[string]struct{}, workingNodes []common.Address) {
	pid, err := filter.PlatformString(feed.Platform)
	if err != nil {
		return
	}

	worker := model.PlatformToWorkerMap[pid]

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

func (h *Hub) verifyStat(ctx context.Context, feed *Feed, stats []*schema.Stat, statMap map[string]struct{}) {
	for _, stat := range stats {
		if stat.EpochInvalidRequest >= int64(model.DefaultSlashCount) {
			continue
		}

		if _, exists := statMap[stat.Address.String()]; !exists {
			statMap[stat.Address.String()] = struct{}{}

			request := ActivityRequest{
				ID: feed.ID,
			}

			nodeMap, err := h.buildActivityByIDPath(
				request,
				[]model.Cache{
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

func (h *Hub) compareFeeds(src, des *Feed) bool {
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

func (h *Hub) verifyData(ctx context.Context, results []DataResponse) error {
	statsMap, err := h.getNodeStatsMap(ctx, results)
	if err != nil {
		return fmt.Errorf("find node stats: %w", err)
	}

	h.sortResults(results)

	if len(statsMap) < model.DefaultNodeCount {
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

func (h *Hub) getNodeStatsMap(ctx context.Context, results []DataResponse) (map[common.Address]*schema.Stat, error) {
	stats, err := h.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		AddressList: lo.Map(results, func(result DataResponse, _ int) common.Address {
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

func (h *Hub) sortResults(results []DataResponse) {
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].First && !results[j].First
	})
}

func (h *Hub) updateStatsWithResults(statsMap map[common.Address]*schema.Stat, results []DataResponse) {
	for _, result := range results {
		if stat, exists := statsMap[result.Address]; exists {
			stat.TotalRequest += int64(result.Request)
			stat.EpochRequest += int64(result.Request)
			stat.EpochInvalidRequest += int64(result.InvalidRequest)
		}
	}
}

func (h *Hub) updateRequestsBasedOnDataCompare(results []DataResponse) {
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

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return data, nil
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

type DataResponse struct {
	Address        common.Address
	Data           []byte
	First          bool
	Err            error
	Request        int
	InvalidRequest int
}

type ErrResponse struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type ActivityResponse struct {
	Data *Feed `json:"data"`
}

type ActivitiesResponse struct {
	Data []*Feed     `json:"data"`
	Meta *MetaCursor `json:"meta,omitempty"`
}

type MetaCursor struct {
	Cursor string `json:"cursor"`
}

type Feed struct {
	ID       string    `json:"id"`
	Owner    string    `json:"owner,omitempty"`
	Network  string    `json:"network"`
	Index    uint      `json:"index"`
	From     string    `json:"from"`
	To       string    `json:"to"`
	Tag      string    `json:"tag"`
	Type     string    `json:"type"`
	Platform string    `json:"platform,omitempty"`
	Actions  []*Action `json:"actions"`
}

type Action struct {
	Tag         string            `json:"tag"`
	Type        string            `json:"type"`
	Platform    string            `json:"platform,omitempty"`
	From        string            `json:"from"`
	To          string            `json:"to"`
	Metadata    metadata.Metadata `json:"metadata"`
	RelatedURLs []string          `json:"related_urls,omitempty"`
}

type Actions []*Action

var _ json.Unmarshaler = (*Action)(nil)

func (a *Action) UnmarshalJSON(bytes []byte) error {
	type ActionAlias Action

	type action struct {
		ActionAlias

		MetadataX json.RawMessage `json:"metadata"`
	}

	var temp action

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("unmarshal action: %w", err)
	}

	tag, err := filter.TagString(temp.Tag)
	if err != nil {
		return fmt.Errorf("invalid action tag: %w", err)
	}

	typeX, err := filter.TypeString(tag, temp.Type)
	if err != nil {
		return fmt.Errorf("invalid action type: %w", err)
	}

	temp.Metadata, err = metadata.Unmarshal(typeX, temp.MetadataX)
	if err != nil {
		return fmt.Errorf("invalid action metadata: %w", err)
	}

	*a = Action(temp.ActionAlias)

	return nil
}

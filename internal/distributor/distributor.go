package distributor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/form/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/dsl"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Distributor struct {
	cacheClient    cache.Client
	databaseClient database.Client
	httpClient     *http.Client
}

// RouterRSSHubData routes RSS Hub data retrieval requests.
// It takes a context, path, and query string as input parameters.
// It returns the retrieved data or an error if any occurred.
func (d *Distributor) RouterRSSHubData(ctx context.Context, path, query string) ([]byte, error) {
	nodes, err := d.retrieveNodes(ctx, RssNodeCacheKey)

	if err != nil {
		return nil, err
	}

	nodeMap, err := d.buildRSSHubPath(path, query, nodes)

	if err != nil {
		return nil, err
	}

	nodeRes, err := d.batchRequestNodes(ctx, nodeMap, d.processRSSHubResults)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeRes.Address.String()))

	return nodeRes.Data, nil
}

// RouterActivityData routes activity data retrieval requests.
// It takes a context and an activity request as input parameters.
// It returns the retrieved data or an error if any occurred.
func (d *Distributor) RouterActivityData(ctx context.Context, request dsl.ActivityRequest) ([]byte, error) {
	nodes, err := d.retrieveNodes(ctx, FullNodeCacheKey)

	if err != nil {
		return nil, err
	}

	nodeMap, err := d.buildActivityPathByID(request, nodes)

	if err != nil {
		return nil, err
	}

	nodeRes, err := d.batchRequestNodes(ctx, nodeMap, d.processActivityResults)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeRes.Address.String()))

	return nodeRes.Data, nil
}

// RouterActivitiesData routes account activities data retrieval requests.
// It takes a context and an account activities request as input parameters.
// It returns the retrieved data or an error if any occurred.
func (d *Distributor) RouterActivitiesData(ctx context.Context, request dsl.AccountActivitiesRequest) ([]byte, error) {
	nodes := make([]NodeEndpointCache, 0, DefaultNodeCount)

	nodeAddresses, err := d.matchLightNodes(ctx, request)

	if err != nil {
		return nil, err
	}

	if len(nodeAddresses) > 0 {
		nodeStats, err := d.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
			AddressList: nodeAddresses,
			Limit:       lo.ToPtr(DefaultNodeCount),
			PointsOrder: lo.ToPtr("DESC"),
		})

		if err != nil {
			return nil, err
		}

		num := lo.Ternary(len(nodeStats) > DefaultNodeCount, DefaultNodeCount, len(nodeStats))

		for i := 0; i < num; i++ {
			nodes = append(nodes, NodeEndpointCache{
				Address:  nodeStats[i].Address.String(),
				Endpoint: nodeStats[i].Endpoint,
			})
		}
	}

	if len(nodes) < DefaultNodeCount {
		fullNodes, err := d.retrieveNodes(ctx, FullNodeCacheKey)
		if err != nil {
			return nil, err
		}

		nodesNeeded := DefaultNodeCount - len(nodes)
		nodesToAdd := lo.Ternary(nodesNeeded > len(fullNodes), len(fullNodes), nodesNeeded)

		for i := 0; i < nodesToAdd; i++ {
			nodes = append(nodes, fullNodes[i])
		}
	}

	nodeMap, err := d.buildAccountActivitiesPath(request, nodes)

	if err != nil {
		return nil, err
	}

	nodeRes, err := d.batchRequestNodes(ctx, nodeMap, d.processActivitiesResults)

	if err != nil {
		return nil, err
	}

	zap.L().Info("first node return", zap.Any("address", nodeRes.Address.String()))

	return nodeRes.Data, nil
}

// retrieveNodes retrieves nodes from the cache or database.
// It takes a context and a cache key as input parameters.
// It returns the retrieved nodes or an error if any occurred.
func (d *Distributor) retrieveNodes(ctx context.Context, key string) ([]NodeEndpointCache, error) {
	var (
		nodesCache []NodeEndpointCache
		nodes      []*schema.Stat
	)

	err := d.cacheClient.Get(ctx, key, &nodesCache)
	if err == nil {
		return nodesCache, nil
	}

	zap.L().Info("not found nodes from cache", zap.String("key", key))

	if errors.Is(err, redis.Nil) {
		switch key {
		case RssNodeCacheKey:
			nodes, err = d.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
				IsRssNode:    lo.ToPtr(true),
				Limit:        lo.ToPtr(DefaultNodeCount),
				ValidRequest: lo.ToPtr(DefaultSlashCount),
				PointsOrder:  lo.ToPtr("DESC"),
			})

			if err != nil {
				return nil, err
			}
		case FullNodeCacheKey:
			nodes, err = d.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
				IsFullNode:   lo.ToPtr(true),
				Limit:        lo.ToPtr(DefaultNodeCount),
				ValidRequest: lo.ToPtr(DefaultSlashCount),
				PointsOrder:  lo.ToPtr("DESC"),
			})

			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown cache key: %s", key)
		}

		if err = d.setNodeCache(ctx, key, nodes); err != nil {
			return nil, err
		}

		zap.L().Info("set nodes to cache", zap.String("key", key))

		nodesCache = lo.Map(nodes, func(n *schema.Stat, _ int) NodeEndpointCache {
			return NodeEndpointCache{
				Address:  n.Address.String(),
				Endpoint: n.Endpoint,
			}
		})

		return nodesCache, nil
	}

	return nil, fmt.Errorf("get nodes from cache: %s, %w", key, err)
}

// batchRequestNodes requests nodes in parallel and returns the first response.
// It takes a context, a node map, and a function to process results as input parameters.
// It returns the first response or an error if any occurred.
func (d *Distributor) batchRequestNodes(_ context.Context, nodeMap map[common.Address]string, processResults func([]DataResponse)) (DataResponse, error) {
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

			data, err := d.fetch(context.Background(), endpoint)
			if err != nil {
				zap.L().Error("fetch request error", zap.Any("node", address.String()), zap.Error(err))

				mu.Lock()
				results = append(results, DataResponse{Address: address, Err: err})

				if len(results) == len(nodeMap) {
					firstResult <- DataResponse{Address: address, Data: []byte(MessageNodeDataFailed)}
				}

				mu.Unlock()

				return
			}

			flagActivities, _ := d.validateActivities(data)
			flagActivity, _ := d.validateActivity(data)

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
			results = append(results, DataResponse{Address: address, Data: data, Valid: true})
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
		return DataResponse{Data: []byte(MessageNodeDataFailed)}, fmt.Errorf("timeout waiting for results")
	}
}

// processRSSHubResults processes the RSS Hub results.
func (d *Distributor) processRSSHubResults(results []DataResponse) {
	if err := d.verifyData(context.Background(), results); err != nil {
		zap.L().Error("fail to verify rss hub request", zap.Any("results", len(results)))
	} else {
		zap.L().Info("complete rss hub request verify", zap.Any("results", len(results)))
	}
}

// processActivityResults processes activity data retrieval results.
func (d *Distributor) processActivityResults(results []DataResponse) {
	if err := d.verifyData(context.Background(), results); err != nil {
		zap.L().Error("fail to verify activity id request ", zap.Any("results", len(results)))
	} else {
		zap.L().Info("complete activity id request verify", zap.Any("results", len(results)))
	}
}

// processActivitiesResults processes account activities data retrieval results.
func (d *Distributor) processActivitiesResults(results []DataResponse) {
	ctx := context.Background()

	if err := d.verifyData(ctx, results); err != nil {
		zap.L().Error("fail to verify activity request", zap.Any("results", len(results)))

		return
	}

	zap.L().Info("complete activity request verify", zap.Any("results", len(results)))

	if !results[0].Valid {
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

	d.processSecondVerify(activities.Data, workingNodes)
}

// buildActivityPathByID builds the path for activity data retrieval by ID.
// It takes an activity request and a slice of cache nodes as input parameters.
// It returns a map of addresses to URLs or an error if any occurred.
func (d *Distributor) buildActivityPathByID(query dsl.ActivityRequest, nodes []NodeEndpointCache) (map[common.Address]string, error) {
	endpointMap, err := d.buildPath(fmt.Sprintf("/decentralized/tx/%s", query.ID), query, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

// buildAccountActivitiesPath builds the path for account activities data retrieval.
// It takes an account activities request and a slice of cache nodes as input parameters.
// It returns a map of addresses to URLs or an error if any occurred.
func (d *Distributor) buildAccountActivitiesPath(query dsl.AccountActivitiesRequest, nodes []NodeEndpointCache) (map[common.Address]string, error) {
	endpointMap, err := d.buildPath(fmt.Sprintf("/decentralized/%s", query.Account), query, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

// setNodeCache sets nodes to the cache.
// It takes a context, a cache key, and a slice of stats as input parameters.
// It returns an error if any occurred.
func (d *Distributor) setNodeCache(ctx context.Context, key string, stats []*schema.Stat) error {
	nodesCache := lo.Map(stats, func(n *schema.Stat, _ int) NodeEndpointCache {
		return NodeEndpointCache{Address: n.Address.String(), Endpoint: n.Endpoint}
	})

	if err := d.cacheClient.Set(ctx, key, nodesCache); err != nil {
		return fmt.Errorf("set nodes to cache: %s, %w", key, err)
	}

	return nil
}

// validateActivities validates the activities response.
// It takes a byte slice of data as an input parameter.
// It returns true if validation is successful, along with the activities response, or false otherwise.
func (d *Distributor) validateActivities(data []byte) (bool, *ActivitiesResponse) {
	var (
		res      ActivitiesResponse
		errRes   ErrResponse
		notFound NotFoundResponse
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

	if err := json.Unmarshal(data, &notFound); err != nil {
		return false, nil
	}

	if notFound.Message != "" {
		return false, nil
	}

	return true, &res
}

// validateActivity validates the activity response.
// It takes a byte slice of data as an input parameter.
// It returns true if validation is successful, along with the activity response, or false otherwise.
func (d *Distributor) validateActivity(data []byte) (bool, *ActivityResponse) {
	var (
		res      ActivityResponse
		errRes   ErrResponse
		notFound NotFoundResponse
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

	if err := json.Unmarshal(data, &notFound); err != nil {
		return false, nil
	}

	if notFound.Message != "" {
		return false, nil
	}

	return true, &res
}

// fetch fetches data from the endpoint.
// It takes a context and a decoded URI as input parameters.
// It returns the retrieved data or an error if any occurred.
func (d *Distributor) fetch(ctx context.Context, decodedURI string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, decodedURI, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	res, err := d.httpClient.Do(request)
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

// buildPath builds the path for nodes.
// It takes a path, a query, and a slice of cache nodes as input parameters.
// It returns a map of addresses to URLs or an error if any occurred.
func (d *Distributor) buildPath(path string, query any, nodes []NodeEndpointCache) (map[common.Address]string, error) {
	if query != nil {
		values, err := form.NewEncoder().Encode(query)

		if err != nil {
			return nil, fmt.Errorf("build params %w", err)
		}

		path = fmt.Sprintf("%s?%s", path, values.Encode())
	}

	urls := make(map[common.Address]string, len(nodes))

	for _, node := range nodes {
		fullURL, err := url.JoinPath(node.Endpoint, path)
		if err != nil {
			return nil, fmt.Errorf("failed to join path for node %s: %w", node.Address, err)
		}

		decodedURL, err := url.QueryUnescape(fullURL)
		if err != nil {
			return nil, fmt.Errorf("failed to unescape url for node %s: %w", node.Address, err)
		}

		urls[common.HexToAddress(node.Address)] = decodedURL
	}

	return urls, nil
}

// buildRSSHubPath builds the path for RSS Hub data retrieval.
// It takes a parameter, a query, and a slice of cache nodes as input parameters.
// It returns a map of addresses to URLs or an error if any occurred.
func (d *Distributor) buildRSSHubPath(param, query string, nodes []NodeEndpointCache) (map[common.Address]string, error) {
	endpointMap, err := d.buildPath(fmt.Sprintf("/rss/%s?%s", param, query), nil, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

// NewDistributor creates a new distributor.
func NewDistributor(_ context.Context, database database.Client, cache cache.Client) *Distributor {
	return &Distributor{
		cacheClient:    cache,
		databaseClient: database,
		httpClient:     http.DefaultClient,
	}
}

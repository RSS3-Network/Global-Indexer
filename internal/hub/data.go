package hub

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/global-indexer/internal/cache"
	"github.com/naturalselectionlabs/global-indexer/provider/node"
	"github.com/naturalselectionlabs/global-indexer/schema"
	"github.com/samber/lo"
)

var (
	rssNodeCacheKey  = "nodes:rss"
	fullNodeCacheKey = "nodes:full"
)

func (h *Hub) registerNode(ctx context.Context, request *RegisterNodeRequest) error {
	node := &schema.Node{
		Address:  request.Address,
		Endpoint: request.Endpoint,
		Stream:   request.Stream,
		Config:   request.Config,
	}

	// Query node from chain
	// TODO

	// Save node to database
	return h.databaseClient.SaveNode(ctx, node)
}

func (h *Hub) routerRSSHubData(ctx context.Context, path, query string) ([]byte, error) {
	nodes, err := h.filterNodes(ctx, rssNodeCacheKey)

	if err != nil {
		return nil, err
	}

	// TODO generate path
	nodeMap, err := h.pathBuilder.GetRSSHubPath(path, query, nodes)

	if err != nil {
		return nil, err
	}

	// TODO batch request
	node, err := h.batchRequest(ctx, nodeMap, processRSSHubResults)

	if err != nil {
		return nil, err
	}

	fmt.Printf("first node return %s\n", node.Address.String())

	return node.Data, nil
}

func (h *Hub) routerActivityData(ctx context.Context, request node.ActivityRequest) ([]byte, error) {
	nodes, err := h.filterNodes(ctx, fullNodeCacheKey)

	if err != nil {
		return nil, err
	}

	// TODO generate path
	nodeMap, err := h.pathBuilder.GetActivityByIDPath(request, nodes)

	if err != nil {
		return nil, err
	}

	// TODO batch request
	node, err := h.batchRequest(ctx, nodeMap, processActivityResults)

	if err != nil {
		return nil, err
	}

	fmt.Printf("first node return %s\n", node.Address.String())

	return node.Data, nil
}

func (h *Hub) routerActivitiesData(ctx context.Context, request node.AccountActivitiesRequest) ([]byte, error) {
	nodes, err := h.filterNodes(ctx, fullNodeCacheKey)

	if err != nil {
		return nil, err
	}

	// TODO generate path
	nodeMap, err := h.pathBuilder.GetAccountActivitiesPath(request, nodes)

	if err != nil {
		return nil, err
	}

	// TODO batch request
	node, err := h.batchRequest(ctx, nodeMap, processActivitiesResults)

	if err != nil {
		return nil, err
	}

	fmt.Printf("first node return %s\n", node.Address.String())

	return node.Data, nil
}

func (h *Hub) filterNodes(ctx context.Context, key string) ([]node.Cache, error) {
	var (
		nodesCache []node.Cache
		nodes      []*schema.Stat
	)

	// TODO: get nodes from cache
	exists, err := cache.Get(ctx, key, &nodesCache)
	if err != nil {
		return nil, fmt.Errorf("get nodes from cache: %s, %w", key, err)
	}

	if exists {
		return nodesCache, nil
	}

	// TODO: get nodes from database
	switch key {
	case rssNodeCacheKey:
		nodes, err = h.databaseClient.FindNodeStats(ctx)

		if err != nil {
			return nil, err
		}
	case fullNodeCacheKey:
		nodes, err = h.databaseClient.FindNodeStats(ctx)

		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown cache key: %s", key)
	}

	nodesCache = lo.Map(nodes, func(n *schema.Stat, _ int) node.Cache {
		return node.Cache{Address: n.Address, Endpoint: n.Endpoint}
	})

	if err := cache.Set(ctx, key, nodesCache); err != nil {
		return nil, fmt.Errorf("set nodes to cache: %s, %w", key, err)
	}

	return nodesCache, nil
}

func processRSSHubResults(results []node.DataResponse) {
	fmt.Printf("rss finish, %d\n", len(results))

	for _, node := range results {
		fmt.Println(node.Address.String())
	}

	// TODO
	time.Sleep(time.Second * 1)
	fmt.Println("step1: data verify")

	time.Sleep(time.Second * 1)
	// TODO
	fmt.Println("step2: data statistic")

	time.Sleep(time.Second * 1)
	// TODO
	fmt.Println("step3: 2nd request、verify")
}

func processActivityResults(results []node.DataResponse) {
	fmt.Printf("feed id finish, %d\n", len(results))

	for _, node := range results {
		fmt.Println(node.Address.String())
	}

	// TODO
	time.Sleep(time.Second * 1)
	fmt.Println("step1: data verify")

	time.Sleep(time.Second * 1)
	// TODO
	fmt.Println("step2: data statistic")

	time.Sleep(time.Second * 1)
	// TODO
	fmt.Println("step3: 2nd request、verify")
}

func processActivitiesResults(results []node.DataResponse) {
	fmt.Printf("feeds finish, %d\n", len(results))

	for _, node := range results {
		fmt.Println(node.Address.String())
	}

	// TODO
	time.Sleep(time.Second * 1)
	fmt.Println("step1: data verify")

	time.Sleep(time.Second * 1)
	// TODO
	fmt.Println("step2: data statistic")

	time.Sleep(time.Second * 1)
	// TODO
	fmt.Println("step3: 2nd request、verify")
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
				fmt.Printf("fetch error: %s, %v\n", address.String(), err)

				return
			}

			mu.Lock()
			results = append(results, node.DataResponse{Address: address, Data: data})
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

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return data, nil
}

func (h *Hub) verifyData(src, des []byte) bool {
	srcHash, destHash := sha256.Sum256(src), sha256.Sum256(des)

	return string(srcHash[:]) == string(destHash[:])
}

func (h *Hub) updateStats() error {
	return nil
}

// cron task
func (h *Hub) sortNodesTask() {
	// TODO nodes sort based on rules

	// TODO save to cache
}

package federatedhandles

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
	node_schemas "github.com/rss3-network/node/schema/worker/federated"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// maintainFederatedHandles maintains the federated handles in the cache.
func (s *server) maintainFederatedHandles(ctx context.Context) error {
	nodeStats, err := s.getFilteredNodeStats(ctx)
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	since, err := s.getLastUpdateTimestamp(ctx)

	if err != nil {
		return err
	}

	handleToNodes, err := s.collectFederatedHandles(ctx, nodeStats, since)
	if err != nil {
		return err
	}

	if err = s.updateCache(ctx, handleToNodes); err != nil {
		return err
	}

	if err = s.updateLastUpdateTimestamp(ctx, now); err != nil {
		return err
	}

	zap.L().Info("maintain federated handles completed", zap.Int("node_count", len(nodeStats)), zap.Int("handle_count", len(handleToNodes)))

	return nil
}

// getFilteredNodeStats retrieves the node statistics that are part of the federated network.
func (s *server) getFilteredNodeStats(ctx context.Context) ([]*schema.Stat, error) {
	nodeStats, err := s.getAllNodeStats(ctx, &schema.StatQuery{})
	if err != nil {
		return nil, err
	}

	return lo.Filter(nodeStats, func(stat *schema.Stat, _ int) bool {
		return stat.FederatedNetwork > 0 && stat.EpochInvalidRequest < int64(model.DemotionCountBeforeSlashing)
	}), nil
}

// getLastUpdateTimestamp retrieves the last update timestamp from the cache.
func (s *server) getLastUpdateTimestamp(ctx context.Context) (uint64, error) {
	var since uint64
	err := s.cacheClient.Get(ctx, fmt.Sprintf("%s%s", model.FederatedHandlesPrefixCacheKey, "since"), &since)

	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("get last update timestamp: %w", err)
	}

	return since, nil
}

// collectFederatedHandles collects the federated handles from the nodes.
func (s *server) collectFederatedHandles(ctx context.Context, nodeStats []*schema.Stat, since uint64) (map[string][]string, error) {
	handleToNodes := make(map[string][]string)

	var mu sync.Mutex

	g, ctx := errgroup.WithContext(ctx)

	for _, stat := range nodeStats {
		stat := stat

		g.Go(func() error {
			handles, err := s.getNodeFederatedHandles(ctx, stat.Endpoint, stat.AccessToken, since)
			if err != nil {
				zap.L().Error("get node federated handles", zap.Error(err), zap.String("endpoint", stat.Endpoint), zap.String("accessToken", stat.AccessToken), zap.Uint64("since", since))
				return nil
			}

			mu.Lock()
			for _, handle := range handles {
				handleToNodes[handle] = append(handleToNodes[handle], stat.Address.String())
			}
			mu.Unlock()

			return nil
		})
	}

	return handleToNodes, g.Wait()
}

// updateCache updates the cache with the new federated handles.
func (s *server) updateCache(ctx context.Context, handleToNodes map[string][]string) error {
	nodeHandleCount := make(map[string]int)

	for handle, nodeAddresses := range handleToNodes {
		key := fmt.Sprintf("%s%s", model.FederatedHandlesPrefixCacheKey, handle)
		existingAddresses, err := s.getExistingAddresses(ctx, key)

		if err != nil {
			return err
		}

		mergedAddresses := s.mergeAddresses(existingAddresses, nodeAddresses)

		if err = s.cacheClient.Set(ctx, key, mergedAddresses, 0); err != nil {
			return fmt.Errorf("set cache: %w", err)
		}

		newAddresses := lo.Filter(nodeAddresses, func(nodeAddress string, _ int) bool {
			return !lo.Contains(existingAddresses, nodeAddress)
		})

		for _, nodeAddress := range newAddresses {
			nodeHandleCount[nodeAddress]++
		}
	}

	countKey := fmt.Sprintf("%s%s", model.FederatedHandlesPrefixCacheKey, "count")
	pipeline := s.cacheClient.Pipeline(ctx)

	for nodeAddress, count := range nodeHandleCount {
		pipeline.ZIncrBy(ctx, countKey, float64(count), nodeAddress)
	}

	_, err := pipeline.Exec(ctx)

	if err != nil {
		return fmt.Errorf("update sorted set: %w", err)
	}

	return nil
}

// getExistingAddresses retrieves the existing addresses from the cache.
func (s *server) getExistingAddresses(ctx context.Context, key string) ([]string, error) {
	var existingAddresses []string

	err := s.cacheClient.Get(ctx, key, &existingAddresses)

	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("get cache: %w", err)
	}

	return existingAddresses, nil
}

// mergeAddresses merges the existing addresses with the new addresses.
func (s *server) mergeAddresses(existingAddresses, newAddresses []string) []string {
	if len(existingAddresses) == 0 {
		return newAddresses
	}

	return lo.Uniq(append(existingAddresses, newAddresses...))
}

// updateLastUpdateTimestamp updates the last update timestamp in the cache.
func (s *server) updateLastUpdateTimestamp(ctx context.Context, timestamp int64) error {
	return s.cacheClient.Set(ctx, fmt.Sprintf("%s%s", model.FederatedHandlesPrefixCacheKey, "since"), timestamp, 0)
}

// getAllNodeStats retrieves all node statistics matching the given query from the database.
func (s *server) getAllNodeStats(ctx context.Context, query *schema.StatQuery) ([]*schema.Stat, error) {
	stats := make([]*schema.Stat, 0)

	// Traverse the entire node.
	for {
		tempStats, err := s.databaseClient.FindNodeStats(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("find node stats: %w", err)
		}

		// If there are no stats, exit the loop.
		if len(tempStats) == 0 {
			break
		}

		stats = append(stats, tempStats...)
		query.Cursor = lo.ToPtr(tempStats[len(tempStats)-1].Address.String())
	}

	return stats, nil
}

type NodeFederatedHandlesResponse struct {
	Platform   string   `json:"platform"`
	Handles    []string `json:"handles"`
	Cursor     string   `json:"cursor"`
	TotalCount int      `json:"total_count"`
}

// getNodeFederatedHandles retrieves the federated handles from the node
func (s *server) getNodeFederatedHandles(ctx context.Context, endpoint, accessToken string, since uint64) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	handles := make([]string, 0, 100)

	var mu sync.Mutex

	g, ctx := errgroup.WithContext(ctx)
	platforms := node_schemas.PlatformStrings()[1:]

	for _, platform := range platforms {
		platform := platform

		g.Go(func() error {
			platformHandles, err := s.fetchPlatformHandles(ctx, endpoint, accessToken, since, platform)
			if err != nil {
				return fmt.Errorf("fetch %s handles: %w", platform, err)
			}

			mu.Lock()
			handles = append(handles, platformHandles...)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return handles, nil
}

func (s *server) fetchPlatformHandles(ctx context.Context, endpoint, accessToken string, since uint64, platform string) ([]string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse endpoint: %w", err)
	}

	u.Path = path.Join(u.Path, "federated/handles")
	q := u.Query()
	q.Set("platform", platform)

	if since > 0 {
		q.Set("since", strconv.FormatUint(since, 10))
	}

	var handles []string

	for {
		u.RawQuery = q.Encode()

		response, err := s.fetchAndDecodeResponse(ctx, u.String(), accessToken)
		if err != nil {
			return nil, err
		}

		handles = append(handles, response.Handles...)

		if response.Cursor == "" {
			break
		}

		q.Set("cursor", response.Cursor)
	}

	return handles, nil
}

// fetchAndDecodeResponse fetches and decodes the response from the given URL.
func (s *server) fetchAndDecodeResponse(ctx context.Context, url, accessToken string) (*NodeFederatedHandlesResponse, error) {
	body, err := s.httpClient.FetchWithMethod(ctx, http.MethodGet, url, accessToken, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch with method: %w", err)
	}

	defer body.Close()

	var response NodeFederatedHandlesResponse
	if err := json.NewDecoder(io.LimitReader(body, 1<<20)).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &response, nil
}

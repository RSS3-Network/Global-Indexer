package enforcer

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
)

// ScoreMaintainer is a structure that maintains a map for quick access.
type ScoreMaintainer struct {
	cacheClient        cache.Client
	nodeEndpointCaches map[string]*model.NodeEndpointCache
	lock               sync.Mutex
}

// addOrUpdateScore updates or adds a nodeEndpointCache in the data structure.
// If the invalidCount is greater than or equal to DemotionCountBeforeSlashing, the nodeEndpointCache is removed.
func (sm *ScoreMaintainer) addOrUpdateScore(ctx context.Context, setKey string, address string, score float64, invalidCount int64) error {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	nodeEndpointCache, ok := sm.nodeEndpointCaches[address]
	if invalidCount >= int64(model.DemotionCountBeforeSlashing) {
		if ok {
			// Remove from Redis.
			if err := sm.cacheClient.ZRem(ctx, setKey, address); err != nil {
				return err
			}
			// Remove from map.
			delete(sm.nodeEndpointCaches, address)
		}

		return nil
	}

	if ok {
		nodeEndpointCache.Score = score
		nodeEndpointCache.InvalidCount = invalidCount
	} else {
		sm.nodeEndpointCaches[address] = &model.NodeEndpointCache{
			Address:      address,
			Score:        score,
			InvalidCount: invalidCount,
		}
	}

	if err := sm.cacheClient.ZAdd(ctx, setKey, redis.Z{
		Member: address,
		Score:  score,
	}); err != nil {
		return err
	}

	return nil
}

// retrieveQualifiedNodes returns the top n NodeEndpointCaches from the sorted set.
func (sm *ScoreMaintainer) retrieveQualifiedNodes(ctx context.Context, setKey string, n int) ([]*model.NodeEndpointCache, error) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	result, err := sm.cacheClient.ZRevRangeWithScores(ctx, setKey, 0, int64(n-1))
	if err != nil {
		return nil, err
	}

	qualifiedNodes := make([]*model.NodeEndpointCache, 0, n)

	for _, item := range result {
		qualifiedNodes = append(qualifiedNodes, sm.nodeEndpointCaches[item.Member.(string)])
	}

	return qualifiedNodes, nil
}

// updateAllQualifiedNodes replaces the current nodeEndpointCaches and redis sorted set.
func (sm *ScoreMaintainer) updateAllQualifiedNodes(ctx context.Context, setKey string, nodeCaches []*model.NodeEndpointCache) error {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	validCaches := make(map[string]*model.NodeEndpointCache)
	members := make([]redis.Z, 0, len(nodeCaches))

	for _, nodeCache := range nodeCaches {
		if nodeCache.InvalidCount < int64(model.DemotionCountBeforeSlashing) {
			validCaches[nodeCache.Address] = nodeCache
			members = append(members, redis.Z{
				Member: nodeCache.Address,
				Score:  nodeCache.Score,
			})
		}
	}

	if err := sm.cacheClient.ZAdd(ctx, setKey, members...); err != nil {
		return err
	}

	tempCaches := sm.nodeEndpointCaches
	sm.nodeEndpointCaches = validCaches

	// Remove the invalid node caches.
	var needRemoveMembers []string

	for address := range tempCaches {
		if _, ok := validCaches[address]; !ok {
			needRemoveMembers = append(needRemoveMembers, address)
		}
	}

	if err := sm.cacheClient.ZRem(ctx, setKey, needRemoveMembers); err != nil {
		return err
	}

	return nil
}

// newScoreMaintainer creates a new ScoreMaintainer with the nodeEndpointCaches and redis sorted set.
func newScoreMaintainer(ctx context.Context, setKey string, nodeCaches []*model.NodeEndpointCache, cacheClient cache.Client) (*ScoreMaintainer, error) {
	validCaches := make(map[string]*model.NodeEndpointCache)
	members := make([]redis.Z, 0, len(nodeCaches))

	for _, nodeCache := range nodeCaches {
		if nodeCache.InvalidCount < int64(model.DemotionCountBeforeSlashing) {
			validCaches[nodeCache.Address] = nodeCache
			members = append(members, redis.Z{
				Member: nodeCache.Address,
				Score:  nodeCache.Score,
			})
		}
	}

	if err := cacheClient.ZAdd(ctx, setKey, members...); err != nil {
		return nil, err
	}

	return &ScoreMaintainer{
		cacheClient:        cacheClient,
		nodeEndpointCaches: validCaches,
	}, nil
}

package enforcer

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
)

// ScoreMaintainer is a structure used to maintain a sorted set and a quick lookup map.
// It uses Redis to keep a sorted set based on node scores,
// and a map in memory for fast access to each node endpoint's cached data.
// This structure helps in quickly and efficiently updating and retrieving scores and statuses of nodes in distributed systems.
type ScoreMaintainer struct {
	cacheClient        cache.Client
	nodeEndpointCaches map[string]*model.NodeEndpointCache
	lock               sync.Mutex
}

// addOrUpdateScore updates or adds a nodeEndpointCache in the data structure.
// If the invalidCount is greater than or equal to DemotionCountBeforeSlashing, the nodeEndpointCache is removed.
func (sm *ScoreMaintainer) addOrUpdateScore(ctx context.Context, setKey string, nodeCache *model.NodeEndpointCache) error {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	if nodeCache.InvalidCount >= int64(model.DemotionCountBeforeSlashing) {
		if _, ok := sm.nodeEndpointCaches[nodeCache.Address]; ok {
			// Remove from sorted set.
			if err := sm.cacheClient.ZRem(ctx, setKey, nodeCache.Address); err != nil {
				return err
			}
			// Remove from map.
			delete(sm.nodeEndpointCaches, nodeCache.Address)
		}

		return nil
	}

	sm.nodeEndpointCaches[nodeCache.Address] = nodeCache

	return sm.cacheClient.ZAdd(ctx, setKey, redis.Z{
		Member: nodeCache.Address,
		Score:  nodeCache.Score,
	})
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

// updateQualifiedNodesMap replaces the current nodeEndpointCaches.
func (sm *ScoreMaintainer) updateQualifiedNodesMap(nodeCaches []*model.NodeEndpointCache) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	validCaches := make(map[string]*model.NodeEndpointCache, len(nodeCaches))

	for _, nodeCache := range nodeCaches {
		if nodeCache.InvalidCount < int64(model.DemotionCountBeforeSlashing) {
			validCaches[nodeCache.Address] = nodeCache
		}
	}

	sm.nodeEndpointCaches = validCaches
}

// newScoreMaintainer creates a new ScoreMaintainer with the nodeEndpointCaches and redis sorted set.
func newScoreMaintainer(ctx context.Context, setKey string, nodeCaches []*model.NodeEndpointCache, cacheClient cache.Client) (*ScoreMaintainer, error) {
	validCaches, newMembers := prepareNodeCachesAndMembers(nodeCaches)
	if len(newMembers) > 0 {
		if err := cacheClient.ZAdd(ctx, setKey, newMembers...); err != nil {
			return nil, err
		}
	}

	members, err := cacheClient.ZRevRangeWithScores(ctx, setKey, 0, -1)
	if err != nil {
		return nil, err
	}

	membersToRemove := filterMembers(members, validCaches)
	if len(membersToRemove) > 0 {
		if err = cacheClient.ZRem(ctx, setKey, membersToRemove); err != nil {
			return nil, err
		}
	}

	return &ScoreMaintainer{
		cacheClient:        cacheClient,
		nodeEndpointCaches: validCaches,
	}, nil
}

// filterMembers filters out the members that are not in the validCaches.
func filterMembers(members []redis.Z, validCaches map[string]*model.NodeEndpointCache) []string {
	membersToRemove := make([]string, 0)

	for _, member := range members {
		if _, ok := validCaches[member.Member.(string)]; !ok {
			membersToRemove = append(membersToRemove, member.Member.(string))
		}
	}

	return membersToRemove
}

// prepareNodeCachesAndMembers filters out invalid node caches and prepares the members for the sorted set.
func prepareNodeCachesAndMembers(nodeCaches []*model.NodeEndpointCache) (map[string]*model.NodeEndpointCache, []redis.Z) {
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

	return validCaches, members
}

package enforcer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
)

// ScoreMaintainer is a structure used to maintain a sorted set and a quick lookup map.
// It uses Redis to keep a sorted set based on node scores,
// and a map in memory for fast access to each node endpoint's cached data.
// This structure helps in quickly and efficiently updating and retrieving scores and statuses of nodes in distributed systems.
type ScoreMaintainer struct {
	cacheClient        cache.Client
	nodeEndpointCaches map[string]string
	lock               sync.Mutex
}

// addOrUpdateScore updates or adds a nodeEndpointCache in the data structure.
// If the invalidCount is greater than or equal to DemotionCountBeforeSlashing, the nodeEndpointCache is removed.
func (sm *ScoreMaintainer) addOrUpdateScore(ctx context.Context, setKey string, nodeStat *schema.Stat) error {
	if nodeStat.EpochInvalidRequest >= int64(model.DemotionCountBeforeSlashing) {
		if _, ok := sm.nodeEndpointCaches[nodeStat.Address.String()]; ok {
			// Fixme: add redis lock
			// Remove from sorted set.
			if err := sm.cacheClient.ZRem(ctx, setKey, nodeStat.Address.String()); err != nil {
				return err
			}
			// Remove from map.
			delete(sm.nodeEndpointCaches, nodeStat.Address.String())
		}

		return nil
	}

	sm.lock.Lock()
	sm.nodeEndpointCaches[nodeStat.Address.String()] = nodeStat.Endpoint
	sm.lock.Unlock()

	return sm.cacheClient.ZAdd(ctx, setKey, redis.Z{
		Member: nodeStat.Address.String(),
		Score:  nodeStat.Score,
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
		if endpoint, ok := sm.nodeEndpointCaches[item.Member.(string)]; ok {
			qualifiedNodes = append(qualifiedNodes, &model.NodeEndpointCache{
				Address:  item.Member.(string),
				Endpoint: endpoint,
			})
		}
	}

	return qualifiedNodes, nil
}

// updateQualifiedNodesMap replaces the current nodeEndpointCaches.
func (sm *ScoreMaintainer) updateQualifiedNodesMap(ctx context.Context, stats []*schema.Stat) error {
	validCaches := make(map[string]string, len(stats))

	// Fixme: parallelize this
	for _, stat := range stats {
		if err := sm.cacheClient.Set(ctx, formatNodeStatRedisKey(model.InvalidRequestCount, stat.Address.String()), stat.EpochInvalidRequest); err != nil {
			return err
		}

		if err := sm.cacheClient.Set(ctx, formatNodeStatRedisKey(model.InvalidRequestCount, stat.Address.String()), stat.EpochInvalidRequest); err != nil {
			return err
		}

		if stat.EpochInvalidRequest < int64(model.DemotionCountBeforeSlashing) {
			validCaches[stat.Address.String()] = stat.Endpoint
		}
	}

	sm.lock.Lock()
	sm.nodeEndpointCaches = validCaches
	sm.lock.Unlock()

	return nil
}

// newScoreMaintainer creates a new ScoreMaintainer with the nodeEndpointCaches and redis sorted set.
func newScoreMaintainer(ctx context.Context, setKey string, nodeStats []*schema.Stat, cacheClient cache.Client) (*ScoreMaintainer, error) {
	validCaches, newMembers, err := prepareNodeCachesAndMembers(ctx, nodeStats, cacheClient)
	if err != nil {
		return nil, err
	}

	if err = adjustMembersToSet(ctx, setKey, newMembers, validCaches, cacheClient); err != nil {
		return nil, err
	}

	return &ScoreMaintainer{
		cacheClient:        cacheClient,
		nodeEndpointCaches: validCaches,
	}, nil
}

// prepareNodeCachesAndMembers filters out invalid node caches and prepares the members for the sorted set.
func prepareNodeCachesAndMembers(ctx context.Context, nodeStats []*schema.Stat, cacheClient cache.Client) (map[string]string, []redis.Z, error) {
	nodeEndpointMap := make(map[string]string, len(nodeStats))
	members := make([]redis.Z, 0, len(nodeStats))

	// Fixme: parallelize this
	for _, stat := range nodeStats {
		var (
			invalidCount int64
			validCount   int64
		)

		if err := getCacheCount(ctx, cacheClient, model.InvalidRequestCount, stat.Address, &invalidCount, stat.EpochInvalidRequest); err != nil {
			return nil, nil, err
		}

		if err := getCacheCount(ctx, cacheClient, model.ValidRequestCount, stat.Address, &validCount, stat.EpochRequest); err != nil {
			return nil, nil, err
		}

		stat.EpochInvalidRequest = invalidCount

		if stat.EpochRequest < validCount {
			stat.TotalRequest += validCount - stat.EpochRequest
		}

		stat.EpochRequest = validCount

		if invalidCount < int64(model.DemotionCountBeforeSlashing) {
			nodeEndpointMap[stat.Address.String()] = stat.Endpoint

			calculateReliabilityScore(stat)

			members = append(members, redis.Z{
				Member: stat.Address.String(),
				Score:  stat.Score,
			})
		}
	}

	return nodeEndpointMap, members, nil
}

func getCacheCount(ctx context.Context, cacheClient cache.Client, key string, address common.Address, resCount *int64, statCount int64) error {
	if err := cacheClient.Get(ctx, formatNodeStatRedisKey(key, address.String()), resCount); err != nil {
		if errors.Is(err, redis.Nil) {
			*resCount = statCount
			return cacheClient.Set(ctx, formatNodeStatRedisKey(key, address.String()), resCount)
		}

		return err
	}

	return nil
}

func adjustMembersToSet(ctx context.Context, setKey string, newMembers []redis.Z, validCaches map[string]string, cacheClient cache.Client) error {
	if len(newMembers) > 0 {
		if err := cacheClient.ZAdd(ctx, setKey, newMembers...); err != nil {
			return err
		}
	}

	members, err := cacheClient.ZRevRangeWithScores(ctx, setKey, 0, -1)
	if err != nil {
		return err
	}

	membersToRemove := filterMembers(members, validCaches)
	if len(membersToRemove) > 0 {
		if err = cacheClient.ZRem(ctx, setKey, membersToRemove); err != nil {
			return err
		}
	}

	return nil
}

// filterMembers filters out the members that are not in the validCaches.
func filterMembers(members []redis.Z, validCaches map[string]string) []string {
	membersToRemove := make([]string, 0)

	for _, member := range members {
		if _, ok := validCaches[member.Member.(string)]; !ok {
			membersToRemove = append(membersToRemove, member.Member.(string))
		}
	}

	return membersToRemove
}

// formatNodeStatRedisKey formats the redis key.
func formatNodeStatRedisKey(key string, address string) string {
	return fmt.Sprintf("%s:%s", key, address)
}

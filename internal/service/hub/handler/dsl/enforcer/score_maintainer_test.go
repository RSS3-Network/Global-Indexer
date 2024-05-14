package enforcer

import (
	"context"
	"testing"

	"github.com/orlangure/gnomock"
	redismock "github.com/orlangure/gnomock/preset/redis"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const setKey = "test"

func TestScoreMaintainer(t *testing.T) {
	t.Parallel()

	var (
		container *gnomock.Container
		err       error
	)

	for {
		container, err = createContainer(context.Background())
		if err == nil {
			break
		}
	}

	t.Cleanup(func() {
		require.NoError(t, gnomock.Stop(container))
	})

	cacheClient := cache.New(redis.NewClient(&redis.Options{
		Addr: container.DefaultAddress(),
	}))
	require.NotNil(t, cacheClient)

	nodeCaches := []*model.NodeEndpointCache{
		{
			Address:      "addr0",
			Score:        3.0,
			InvalidCount: 0,
		},
		{
			Address:      "addr1",
			Score:        2.0,
			InvalidCount: 0,
		},
		{
			Address:      "addr2",
			Score:        4.0,
			InvalidCount: 1,
		},
		{
			Address:      "addr3",
			Score:        4.5,
			InvalidCount: 0,
		},
		{
			Address:      "addr100",
			Score:        40.0,
			InvalidCount: 10,
		},
	}
	sm, err := newScoreMaintainer(context.Background(), setKey, nodeCaches, cacheClient)
	require.NoError(t, err)

	// Retrieve qualified nodes
	nodes, err := sm.retrieveQualifiedNodes(context.Background(), setKey, 3)
	require.NoError(t, err)
	assert.Equal(t, 3, len(nodes))
	assert.Equal(t, "addr3", nodes[0].Address)
	assert.Equal(t, "addr2", nodes[1].Address)
	assert.Equal(t, "addr0", nodes[2].Address)
	assert.Equal(t, 4, len(sm.nodeEndpointCaches))

	// Add a new node
	err = sm.addOrUpdateScore(context.Background(), setKey, &model.NodeEndpointCache{
		Address:      "addr4",
		Score:        5.0,
		InvalidCount: 0,
	})
	require.NoError(t, err)
	// Update the score of an existing node

	err = sm.addOrUpdateScore(context.Background(), setKey, &model.NodeEndpointCache{
		Address:      "addr0",
		Score:        6.0,
		InvalidCount: 1,
	})
	require.NoError(t, err)
	assert.Equal(t, 5, len(sm.nodeEndpointCaches))
	assert.Equal(t, 6.0, sm.nodeEndpointCaches["addr0"].Score)
	// Retrieve qualified nodes
	nodes, err = sm.retrieveQualifiedNodes(context.Background(), setKey, 10)
	require.NoError(t, err)
	assert.Equal(t, 5, len(nodes))
	assert.Equal(t, "addr0", nodes[0].Address)
	assert.Equal(t, "addr4", nodes[1].Address)
	assert.Equal(t, "addr3", nodes[2].Address)
	assert.Equal(t, "addr2", nodes[3].Address)
	assert.Equal(t, "addr1", nodes[4].Address)

	// Add a new node with invalid count greater than DemotionCountBeforeSlashing
	err = sm.addOrUpdateScore(context.Background(), setKey, &model.NodeEndpointCache{
		Address:      "addr4",
		Score:        7.0,
		InvalidCount: int64(model.DemotionCountBeforeSlashing),
	})
	require.NoError(t, err)
	assert.Equal(t, 4, len(sm.nodeEndpointCaches))

	// Retrieve qualified nodes
	nodes, err = sm.retrieveQualifiedNodes(context.Background(), setKey, 10)
	require.NoError(t, err)
	assert.Equal(t, 4, len(nodes))
	assert.Equal(t, "addr0", nodes[0].Address)
	assert.Equal(t, "addr3", nodes[1].Address)
	assert.Equal(t, "addr2", nodes[2].Address)
	assert.Equal(t, "addr1", nodes[3].Address)

	// Update all qualified nodes
	newNodeCaches := make([]*model.NodeEndpointCache, 0)
	node3 := &model.NodeEndpointCache{
		Address:      "addr5",
		Score:        7.0,
		InvalidCount: 1,
	}
	newNodeCaches = append(newNodeCaches, node3)
	sm.updateQualifiedNodesMap(newNodeCaches)
	require.NoError(t, err)
	assert.Equal(t, 1, len(sm.nodeEndpointCaches))
}

func createContainer(ctx context.Context) (container *gnomock.Container, err error) {
	preset := redismock.Preset()

	// Health check function to ensure the Redis server is ready.
	healthcheckFunc := func(ctx context.Context, container *gnomock.Container) error {
		client := redis.NewClient(&redis.Options{
			Addr: container.DefaultAddress(),
		})

		_, err = client.Ping(ctx).Result()
		if err != nil {
			return err
		}

		return nil
	}

	// Start the Gnomock container with the specified Redis preset and health check.
	container, err = gnomock.Start(preset, gnomock.WithContext(ctx), gnomock.WithHealthCheck(healthcheckFunc))
	if err != nil {
		return nil, err
	}

	return container, nil
}

package enforcer

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/orlangure/gnomock"
	redismock "github.com/orlangure/gnomock/preset/redis"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
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

	nodeStats := []*schema.Stat{
		{
			Address:             common.Address{0},
			Score:               3.0,
			Endpoint:            "addr0",
			EpochInvalidRequest: 0,
			EpochRequest:        0,
			TotalRequest:        200,
		},
		{
			Address:             common.Address{1},
			Score:               2.0,
			Endpoint:            "addr1",
			EpochInvalidRequest: 0,
			EpochRequest:        0,
			TotalRequest:        100,
		},
		{
			Address:             common.Address{2},
			Score:               4.0,
			Endpoint:            "addr2",
			EpochInvalidRequest: 0,
			EpochRequest:        0,
			TotalRequest:        300,
		},
		{
			Address:             common.Address{3},
			Score:               4.5,
			Endpoint:            "addr3",
			EpochInvalidRequest: 0,
			EpochRequest:        0,
			TotalRequest:        400,
		},
		{
			Address:             common.Address{4},
			Score:               40.0,
			Endpoint:            "addr4",
			EpochInvalidRequest: 10,
			EpochRequest:        10,
			TotalRequest:        400,
		},
	}
	sm, err := newScoreMaintainer(context.Background(), setKey, nodeStats, cacheClient)
	require.NoError(t, err)

	// Retrieve qualified nodes
	nodes, err := sm.retrieveQualifiedNodes(context.Background(), setKey, 3)
	require.NoError(t, err)
	assert.Equal(t, 3, len(nodes))
	assert.Equal(t, nodeStats[3].Address.String(), nodes[0].Address)
	assert.Equal(t, nodeStats[2].Address.String(), nodes[1].Address)
	assert.Equal(t, nodeStats[0].Address.String(), nodes[2].Address)
	assert.Equal(t, 4, len(sm.nodeEndpointCaches))

	// Add a new node
	err = sm.addOrUpdateScore(context.Background(), setKey, &schema.Stat{
		Address:  common.Address{5},
		Endpoint: "addr5",
		Score:    5.0,
	})
	require.NoError(t, err)
	// Update the score of an existing node
	err = sm.addOrUpdateScore(context.Background(), setKey, &schema.Stat{
		Address:  common.Address{0},
		Endpoint: "addr0",
		Score:    6.0,
	})
	require.NoError(t, err)
	assert.Equal(t, 5, len(sm.nodeEndpointCaches))
	nodes, err = sm.retrieveQualifiedNodes(context.Background(), setKey, 10)
	require.NoError(t, err)
	assert.Equal(t, common.Address{0}.String(), nodes[0].Address)
	assert.Equal(t, common.Address{5}.String(), nodes[1].Address)
	assert.Equal(t, common.Address{3}.String(), nodes[2].Address)
	assert.Equal(t, common.Address{2}.String(), nodes[3].Address)
	assert.Equal(t, common.Address{1}.String(), nodes[4].Address)

	//// Add a new node with invalid count greater than DemotionCountBeforeSlashing
	err = sm.addOrUpdateScore(context.Background(), setKey, &schema.Stat{
		Address:             common.Address{0},
		Score:               7.0,
		EpochInvalidRequest: int64(model.DemotionCountBeforeSlashing),
	})
	require.NoError(t, err)
	assert.Equal(t, 4, len(sm.nodeEndpointCaches))

	// Retrieve qualified nodes
	nodes, err = sm.retrieveQualifiedNodes(context.Background(), setKey, 10)
	require.NoError(t, err)
	assert.Equal(t, 4, len(nodes))
	assert.Equal(t, common.Address{5}.String(), nodes[0].Address)
	assert.Equal(t, common.Address{3}.String(), nodes[1].Address)
	assert.Equal(t, common.Address{2}.String(), nodes[2].Address)
	assert.Equal(t, common.Address{1}.String(), nodes[3].Address)

	// Update all qualified nodes
	newNodeCaches := make([]*schema.Stat, 0)
	node3 := &schema.Stat{
		Address:             common.Address{100},
		Endpoint:            "addr100",
		Score:               7.0,
		EpochInvalidRequest: 1,
	}
	newNodeCaches = append(newNodeCaches, node3)
	err = sm.updateQualifiedNodesMap(context.Background(), newNodeCaches)
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

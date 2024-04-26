package enforcer

import (
	"testing"

	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/stretchr/testify/assert"
)

func TestScoreMaintainer(t *testing.T) {
	t.Parallel()

	pq := []*model.NodeEndpointCache{
		{
			Address:      "addr0",
			Score:        3.0,
			InvalidCount: 0,
			Index:        0,
		},
		{
			Address:      "addr1",
			Score:        2.0,
			InvalidCount: 0,
			Index:        1,
		},
		{
			Address:      "addr2",
			Score:        4.0,
			InvalidCount: 1,
			Index:        2,
		},
		{
			Address:      "addr3",
			Score:        4.0,
			InvalidCount: 0,
			Index:        3,
		},
		{
			Address:      "addr100",
			Score:        40.0,
			InvalidCount: 10,
			Index:        3,
		},
	}
	sm := newScoreMaintainer(pq)

	// Retrieve qualified nodes
	nodes := sm.retrieveQualifiedNodes(3)
	assert.Equal(t, 3, len(nodes))
	assert.Equal(t, "addr3", nodes[0].Address)
	assert.Equal(t, "addr2", nodes[1].Address)
	assert.Equal(t, "addr0", nodes[2].Address)

	// Add a new node
	sm.addOrUpdateScore("addr4", 5.0, 0)
	// Update the score of an existing node
	sm.addOrUpdateScore("addr0", 6.0, 1)
	assert.Equal(t, 5, sm.queue.Len())
	assert.Equal(t, 6.0, (*sm.queue)[0].Score)
	assert.Equal(t, "addr0", (*sm.queue)[0].Address)

	// Add a new node with invalid count greater than DemotionCountBeforeSlashing
	sm.addOrUpdateScore("addr4", 7.0, int64(model.DemotionCountBeforeSlashing))
	assert.Equal(t, 4, sm.queue.Len())
	assert.Equal(t, "addr0", (*sm.queue)[0].Address)

	// Retrieve qualified nodes
	nodes = sm.retrieveQualifiedNodes(10)
	assert.Equal(t, 4, len(nodes))

	// Update all qualified nodes
	newPQ := make(priorityNodeQueue, 0)
	node3 := &model.NodeEndpointCache{
		Address:      "addr5",
		Score:        7.0,
		InvalidCount: 1,
	}
	newPQ = append(newPQ, node3)
	sm.updateAllQualifiedNodes(newPQ)
	assert.Equal(t, 1, sm.queue.Len())
	assert.Equal(t, "addr5", (*sm.queue)[0].Address)
}

package enforcer

import (
	"container/heap"
	"sync"

	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/samber/lo"
)

// priorityNodeQueue implements heap.Interface and holds NodeEndpointCaches.
type priorityNodeQueue []*model.NodeEndpointCache

func (pq priorityNodeQueue) Len() int {
	return len(pq)
}

func (pq priorityNodeQueue) Less(i, j int) bool {
	if pq[i].Score == pq[j].Score {
		// If Scores are the same, return true if pq[i] has a smaller InvalidCount than pq[j]
		return pq[i].InvalidCount < pq[j].InvalidCount
	}
	// Otherwise, return true if pq[i] has a greater score than pq[j]
	return pq[i].Score > pq[j].Score
}

func (pq priorityNodeQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *priorityNodeQueue) Push(x interface{}) {
	n := len(*pq)
	nodeEndpointCache := x.(*model.NodeEndpointCache)
	nodeEndpointCache.Index = n
	*pq = append(*pq, nodeEndpointCache)
}

func (pq *priorityNodeQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	nodeEndpointCache := old[n-1]
	old[n-1] = nil
	nodeEndpointCache.InvalidCount = -1
	*pq = old[:n-1]

	return nodeEndpointCache
}

type ScoreMaintainer struct {
	queue              *priorityNodeQueue
	nodeEndpointCaches map[string]*model.NodeEndpointCache
	lock               sync.Mutex
}

// AddOrUpdateScore updates or adds a nodeEndpointCache in the data structure
func (sm *ScoreMaintainer) addOrUpdateScore(address string, score float64, invalidCount int64) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	nodeEndpointCache, ok := sm.nodeEndpointCaches[address]
	if invalidCount >= int64(model.DemotionCountBeforeSlashing) {
		if ok {
			// Remove from heap.
			heap.Remove(sm.queue, nodeEndpointCache.Index)
			// Remove from map.
			delete(sm.nodeEndpointCaches, address)
		}

		return
	}

	if !ok {
		newNodeEndpointCache := &model.NodeEndpointCache{
			Address:      address,
			Score:        score,
			InvalidCount: invalidCount,
		}
		heap.Push(sm.queue, newNodeEndpointCache)
		sm.nodeEndpointCaches[address] = newNodeEndpointCache
	} else {
		nodeEndpointCache.Score = score
		nodeEndpointCache.InvalidCount = invalidCount
		heap.Fix(sm.queue, nodeEndpointCache.Index)
	}
}

func (sm *ScoreMaintainer) getQualifiedNodes(n int) []*model.NodeEndpointCache {
	var qualifiedNodes []*model.NodeEndpointCache
	// Temporary storage to hold elements popped from the heap.
	var tempHeap priorityNodeQueue

	// Continue until we have enough qualifiedNodes or the heap is empty
	for len(qualifiedNodes) < n && sm.queue.Len() > 0 {
		// Pop the highest score node from the heap
		qualifiedNode := heap.Pop(sm.queue).(*model.NodeEndpointCache)

		qualifiedNodes = append(qualifiedNodes, qualifiedNode)

		// Store the qualifiedNode to re-push later.
		tempHeap = append(tempHeap, qualifiedNode)

		// If we have enough qualifiedNodes, break.
		if len(qualifiedNodes) == n {
			break
		}
	}

	// Push all item back to restore the heap
	for _, item := range tempHeap {
		heap.Push(sm.queue, item)
	}

	return qualifiedNodes
}

func (sm *ScoreMaintainer) updateAllQualifiedNodes(pq priorityNodeQueue) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	heap.Init(&pq)

	sm.queue = &pq
	sm.nodeEndpointCaches = lo.SliceToMap(pq, func(n *model.NodeEndpointCache) (string, *model.NodeEndpointCache) {
		return n.Address, n
	})
}

func newScoreMaintainer(pq priorityNodeQueue) *ScoreMaintainer {
	heap.Init(&pq)

	return &ScoreMaintainer{
		queue: &pq,
		nodeEndpointCaches: lo.SliceToMap(pq, func(n *model.NodeEndpointCache) (string, *model.NodeEndpointCache) {
			return n.Address, n
		}),
	}
}

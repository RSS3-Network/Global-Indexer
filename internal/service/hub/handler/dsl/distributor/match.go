package distributor

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
)

// matchLightNodes matches light Nodes based on the given Activities request
// A light Node is a non-Full Node
func (d *Distributor) matchLightNodes(ctx context.Context, workers, networks []string) ([]common.Address, error) {
	var (
		nodes []common.Address
		err   error
	)

	switch {
	case len(workers) > 0 && len(networks) > 0:
		nodes, err = d.matchWorkerAndNetwork(ctx, workers, networks)
	case len(workers) > 0:
		nodes, err = d.matchWorker(ctx, workers)
	case len(networks) > 0:
		nodes, err = d.matchNetwork(ctx, networks)
	}

	return nodes, err
}

// matchNetwork matches nodes based on the given network requests,
// and returns the addresses of Nodes that match the requests.
func (d *Distributor) matchNetwork(ctx context.Context, networks []string) ([]common.Address, error) {
	// Find all indexers that match the networks.
	indexers, err := d.databaseClient.FindNodeWorkers(ctx, &schema.WorkerQuery{
		Networks: networks,
		IsActive: lo.ToPtr(true),
	})
	if err != nil {
		return nil, err
	}

	// Generate a map of Node addresses to network workers.
	nodeNetworkWorkersMap := generateNodeNetworkWorkersMap(indexers)
	// Filter nodes that match the network requests.
	return filterMatchingNetworkNodes(nodeNetworkWorkersMap, networks), nil
}

type NetworkWorkersMap struct {
	Workers map[string][]string
}

// generateNodeNetworkWorkersMap generates a map of Node addresses to network workers.
func generateNodeNetworkWorkersMap(workers []*schema.Worker) map[common.Address]NetworkWorkersMap {
	nodeNetworkWorkersMap := make(map[common.Address]NetworkWorkersMap)

	for _, worker := range workers {
		if _, exists := nodeNetworkWorkersMap[worker.Address]; !exists {
			nodeNetworkWorkersMap[worker.Address] = NetworkWorkersMap{Workers: make(map[string][]string)}
		}

		networkWorkersMap := nodeNetworkWorkersMap[worker.Address].Workers
		networkWorkersMap[worker.Network] = append(networkWorkersMap[worker.Network], worker.Name)
	}

	return nodeNetworkWorkersMap
}

// filterMatchingNetworkNodes filters nodes that match the network requests.
func filterMatchingNetworkNodes(nodeNetworkWorkersMap map[common.Address]NetworkWorkersMap, requestNetworks []string) []common.Address {
	var nodes []common.Address

	for address, networkWorkersMap := range nodeNetworkWorkersMap {
		if isValidNetworkNode(networkWorkersMap, requestNetworks) {
			nodes = append(nodes, address)
		}
	}

	return nodes
}

// isValidNetworkNode checks if the Node has the capability to serve the required networks.
func isValidNetworkNode(networkWorkersMap NetworkWorkersMap, requestNetworks []string) bool {
	// Check if the number of networks match the number of request networks.
	if len(networkWorkersMap.Workers) != len(requestNetworks) {
		return false
	}

	for n, workers := range networkWorkersMap.Workers {
		// Check if the workers match the required workers for the network.
		requiredWorkers := model.NetworkToWorkersMap[n]
		if !AreSliceElementsIdentical(workers, requiredWorkers) {
			return false
		}
	}

	return true
}

type WorkerNetworksMap struct {
	Networks map[string][]string
}

// matchWorker matches nodes based on the given worker names,
// and returns the addresses of Nodes that match the requests.
func (d *Distributor) matchWorker(ctx context.Context, workers []string) ([]common.Address, error) {
	indexers, err := d.databaseClient.FindNodeWorkers(ctx, &schema.WorkerQuery{
		Names:    workers,
		IsActive: lo.ToPtr(true),
	})
	if err != nil {
		return nil, err
	}

	nodeWorkerNetworksMap := generateNodeWorkerNetworksMap(indexers)

	return filterMatchingWorkerNodes(nodeWorkerNetworksMap, workers), nil
}

// generateNodeWorkerNetworksMap generates a map of Node addresses to worker networks.
func generateNodeWorkerNetworksMap(workers []*schema.Worker) map[common.Address]WorkerNetworksMap {
	nodeWorkerNetworksMap := make(map[common.Address]WorkerNetworksMap)

	for _, worker := range workers {
		if _, exists := nodeWorkerNetworksMap[worker.Address]; !exists {
			nodeWorkerNetworksMap[worker.Address] = WorkerNetworksMap{Networks: make(map[string][]string)}
		}

		workerNetworksMap := nodeWorkerNetworksMap[worker.Address].Networks
		workerNetworksMap[worker.Name] = append(workerNetworksMap[worker.Name], worker.Network)
	}

	return nodeWorkerNetworksMap
}

// filterMatchingWorkerNodes filters nodes that match the worker requests.
func filterMatchingWorkerNodes(nodeWorkerNetworksMap map[common.Address]WorkerNetworksMap, workers []string) []common.Address {
	var nodes []common.Address

	for address, workerNetworksMap := range nodeWorkerNetworksMap {
		if isValidWorkerNode(workerNetworksMap, workers) {
			nodes = append(nodes, address)
		}
	}

	return nodes
}

// isValidWorkerNode checks if the Node serves the required workers.
func isValidWorkerNode(workerNetworksMap WorkerNetworksMap, workers []string) bool {
	if len(workerNetworksMap.Networks) != len(workers) {
		return false
	}

	for w, networks := range workerNetworksMap.Networks {
		requiredNetworks := model.WorkerToNetworksMap[w]
		if !AreSliceElementsIdentical(networks, requiredNetworks) {
			return false
		}
	}

	return true
}

// matchWorkerAndNetwork matches nodes based on both worker and network.
func (d *Distributor) matchWorkerAndNetwork(ctx context.Context, workers, networks []string) ([]common.Address, error) {
	indexers, err := d.databaseClient.FindNodeWorkers(ctx, &schema.WorkerQuery{
		Names:    workers,
		Networks: networks,
		IsActive: lo.ToPtr(true),
	})

	if err != nil {
		return nil, err
	}

	nodeWorkerNetworksMap := generateNodeWorkerNetworksMap(indexers)

	return filterMatchingWorkerAndNetworkNodes(nodeWorkerNetworksMap, workers, networks), nil
}

// filterMatchingWorkerAndNetworkNodes filters nodes that match the worker and network requests.
func filterMatchingWorkerAndNetworkNodes(nodeWorkerNetworksMap map[common.Address]WorkerNetworksMap, workers, networks []string) []common.Address {
	var nodes []common.Address

	for address, workerNetworksMap := range nodeWorkerNetworksMap {
		if isValidWorkerAndNetworkNode(workerNetworksMap, workers, networks) {
			nodes = append(nodes, address)
		}
	}

	return nodes
}

// isValidWorkerAndNetworkNode checks if the Node matches the required workers and networks.
func isValidWorkerAndNetworkNode(workerNetworksMap WorkerNetworksMap, workers, requestNetworks []string) bool {
	if len(workerNetworksMap.Networks) != len(workers) {
		return false
	}

	for w, networks := range workerNetworksMap.Networks {
		workerRequiredNetworks := model.WorkerToNetworksMap[w]

		requiredNetworks := IntersectUnique(workerRequiredNetworks, requestNetworks)

		if !AreSliceElementsIdentical(networks, requiredNetworks) {
			return false
		}
	}

	return true
}

// AreSliceElementsIdentical checks if the elements of two string slices are identical.
func AreSliceElementsIdentical(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	slice1Set := make(map[string]struct{}, len(slice1))

	for _, str := range slice1 {
		slice1Set[str] = struct{}{}
	}

	for _, str := range slice2 {
		if _, exists := slice1Set[str]; !exists {
			return false
		}
	}

	return true
}

// IntersectUnique returns the unique common elements between two string slices.
func IntersectUnique(slice1, slice2 []string) []string {
	elementMap := make(map[string]struct{})

	var uniqueElements []string

	for _, v := range slice1 {
		elementMap[v] = struct{}{}
	}

	for _, v := range slice2 {
		if _, found := elementMap[v]; found {
			uniqueElements = append(uniqueElements, v)
			delete(elementMap, v)
		}
	}

	return uniqueElements
}

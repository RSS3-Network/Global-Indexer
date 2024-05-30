package distributor

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/dsl"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/rss3-network/node/schema/worker"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
	"github.com/samber/lo"
)

// matchLightNodes matches light Nodes based on the given Activities request
// A light Node is a non-Full Node
func (d *Distributor) matchLightNodes(ctx context.Context, request dsl.ActivitiesRequest) ([]common.Address, error) {
	// Find network nodes that match the network requests.
	networks, err := getNetworks(request.Network)
	if err != nil {
		return nil, err
	}

	// Find tag workers that match the tag requests.
	tagWorkers, err := getWorkersByTag(request.Tag)
	if err != nil {
		return nil, err
	}

	// Find platform workers that match the platform requests.
	platformWorkers, err := getWorkersByPlatform(request.Platform)
	if err != nil {
		return nil, err
	}

	// Find nodes that match the tag workers, platform workers, and networks.
	return d.findQualifiedNodes(ctx, tagWorkers, platformWorkers, networks)
}

type WorkerSet map[string]struct{}

// getNetworks returns a slice of networks based on the given network names.
func getNetworks(networks []string) ([]string, error) {
	for i, n := range networks {
		nid, err := network.NetworkString(n)
		if err != nil {
			return nil, err
		}

		networks[i] = nid.String()
	}

	return networks, nil
}

// getWorkersByTag returns a set of workers based on the given tags.
func getWorkersByTag(tags []string) (WorkerSet, error) {
	tagWorkers := make(WorkerSet)

	for _, tagX := range tags {
		tid, err := tag.TagString(tagX)
		if err != nil {
			return nil, err
		}

		tagWorker, exists := model.TagToWorkersMap[tid.String()]
		if !exists {
			return nil, fmt.Errorf("no workers found for tag: %s", tid)
		}

		for _, w := range tagWorker {
			tagWorkers[w] = struct{}{}
		}
	}

	return tagWorkers, nil
}

// getWorkersByPlatform returns a set of workers based on the given platforms.
func getWorkersByPlatform(platforms []string) (WorkerSet, error) {
	platformWorkers := make(WorkerSet)

	for _, platform := range platforms {
		pid, err := worker.PlatformString(platform)
		if err != nil {
			return nil, err
		}

		workers, exists := model.PlatformToWorkersMap[pid.String()]
		if !exists {
			return nil, fmt.Errorf("no worker found for platform: %s", pid)
		}

		for _, w := range workers {
			platformWorkers[w] = struct{}{}
		}
	}

	return platformWorkers, nil
}

// findQualifiedNodes finds nodes that match the given tag workers, platform workers, and networks.
func (d *Distributor) findQualifiedNodes(ctx context.Context, tagWorkers, platformWorkers WorkerSet, networks []string) ([]common.Address, error) {
	workers := combineTagAndPlatformWorkers(tagWorkers, platformWorkers)
	// If no common workers are found between tag workers and platform workers,
	// it indicates that tags and platforms are not compatible.
	if len(workers) == 0 && (len(tagWorkers) > 0 || len(platformWorkers) > 0) {
		return nil, fmt.Errorf("no workers found for tags and platforms")
	}

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
	default:
	}

	return nodes, err
}

// combineTagAndPlatformWorkers combines tag workers and platform workers.
func combineTagAndPlatformWorkers(tagWorkers, platformWorkers WorkerSet) []string {
	if len(tagWorkers) == 0 {
		return lo.Keys(platformWorkers)
	}

	if len(platformWorkers) == 0 {
		return lo.Keys(tagWorkers)
	}

	// Find common workers between tag workers and platform workers.
	return IntersectUnique(lo.Keys(tagWorkers), lo.Keys(platformWorkers))
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
		nid, _ := network.NetworkString(n)

		// Check if the workers match the required workers for the network.
		requiredWorkers := model.NetworkToWorkersMap[nid.String()]
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
		wid, _ := worker.WorkerString(w)

		requiredNetworks := model.WorkerToNetworksMap[wid.String()]
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
		wid, _ := worker.WorkerString(w)

		workerRequiredNetworks := model.WorkerToNetworksMap[wid.String()]

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

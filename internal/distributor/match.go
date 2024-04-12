package distributor

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/dsl"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/samber/lo"
)

// matchLightNodes matches light nodes based on the given account activities request,
// and returns the addresses of light nodes that match the request.
func (d *Distributor) matchLightNodes(ctx context.Context, request dsl.AccountActivitiesRequest) ([]common.Address, error) {
	tagWorkers := make(map[string]struct{})
	platformWorkers := make(map[string]struct{})

	// check network
	for i, network := range request.Network {
		nid, err := filter.NetworkString(network)
		if err != nil {
			return nil, err
		}

		request.Network[i] = nid.String()
	}

	// check tag
	for _, tag := range request.Tag {
		tid, err := filter.TagString(tag)

		if err != nil {
			return nil, err
		}

		tagWorker, exists := TagToWorkersMap[tid]

		if !exists {
			return nil, err
		}

		for _, worker := range tagWorker {
			tagWorkers[worker] = struct{}{}
		}
	}

	// check platform
	for _, platform := range request.Platform {
		pid, err := filter.PlatformString(platform)

		if err != nil {
			return nil, err
		}

		platformWorker, exists := PlatformToWorkerMap[pid]

		if !exists {
			return nil, err
		}

		platformWorkers[platformWorker] = struct{}{}
	}

	var (
		nodes        []common.Address
		err          error
		workers      []string
		needsWorker  = len(tagWorkers) > 0 || len(platformWorkers) > 0
		needsNetwork = len(request.Network) > 0
	)

	if len(tagWorkers) > 0 && len(platformWorkers) > 0 {
		workers = findCommonStr(lo.Keys(tagWorkers), lo.Keys(platformWorkers))
	} else if len(tagWorkers) > 0 {
		workers = lo.Keys(tagWorkers)
	} else if len(platformWorkers) > 0 {
		workers = lo.Keys(platformWorkers)
	}

	switch {
	case !needsWorker && needsNetwork:
		nodes, err = d.matchNetwork(ctx, request.Network)
	case needsWorker && !needsNetwork:
		nodes, err = d.matchWorker(ctx, workers)
	case needsWorker && needsNetwork:
		nodes, err = d.matchWorkerAndNetwork(ctx, workers, request.Network)
	default:
	}

	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// matchNetwork matches nodes based on the given network requests,
// and returns the addresses of nodes that match the requests.
func (d *Distributor) matchNetwork(ctx context.Context, requestNetworks []string) ([]common.Address, error) {
	nodes := make([]common.Address, 0)

	indexers, err := d.databaseClient.FindNodeIndexers(ctx, nil, requestNetworks, nil)

	if err != nil {
		return nil, err
	}

	nodeNetworkToWorkerMap := make(map[common.Address]map[string][]string)

	for _, indexer := range indexers {
		if _, exists := nodeNetworkToWorkerMap[indexer.Address]; !exists {
			nodeNetworkToWorkerMap[indexer.Address] = make(map[string][]string)
		}

		nodeNetworkToWorkerMap[indexer.Address][indexer.Network] = append(nodeNetworkToWorkerMap[indexer.Address][indexer.Network], indexer.Worker)
	}

	for address, networkToWorkersMap := range nodeNetworkToWorkerMap {
		if len(networkToWorkersMap) != len(requestNetworks) {
			// number of networks not match
			continue
		}

		flag := true

		// check every network workers
		for network, workers := range networkToWorkersMap {
			nid, _ := filter.NetworkString(network)

			requiredWorkers := NetworkToWorkersMap[nid]

			// number of workers not match
			if len(requiredWorkers) != len(workers) {
				flag = false

				break
			}

			workerMap := make(map[string]struct{}, len(workers))

			for _, worker := range workers {
				workerMap[worker] = struct{}{}
			}

			// check every worker
			for _, worker := range requiredWorkers {
				if _, exists := workerMap[worker]; !exists {
					flag = false

					break
				}
			}
		}

		if flag {
			nodes = append(nodes, address)
		}
	}

	return nodes, nil
}

// matchWorker matches nodes based on the given worker names,
// and returns the addresses of nodes that match the requests.
func (d *Distributor) matchWorker(ctx context.Context, workers []string) ([]common.Address, error) {
	nodes := make([]common.Address, 0)

	indexers, err := d.databaseClient.FindNodeIndexers(ctx, nil, nil, workers)

	if err != nil {
		return nil, err
	}

	nodeWorkerToNetworksMap := make(map[common.Address]map[string][]string)

	for _, indexer := range indexers {
		if _, exists := nodeWorkerToNetworksMap[indexer.Address]; !exists {
			nodeWorkerToNetworksMap[indexer.Address] = make(map[string][]string)
		}

		nodeWorkerToNetworksMap[indexer.Address][indexer.Worker] = append(nodeWorkerToNetworksMap[indexer.Address][indexer.Worker], indexer.Network)
	}

	for address, workerToNetworksMap := range nodeWorkerToNetworksMap {
		if len(nodeWorkerToNetworksMap) != len(workers) {
			// number of workers not match
			continue
		}

		flag := true

		// check every worker networks
		for worker, networks := range workerToNetworksMap {
			wid, _ := filter.NameString(worker)

			requiredNetworks := WorkerToNetworksMap[wid]

			// number of networks not match
			if len(requiredNetworks) != len(networks) {
				flag = false

				break
			}

			networkMap := make(map[string]struct{}, len(networks))

			for _, network := range networks {
				networkMap[network] = struct{}{}
			}

			// check every network
			for _, network := range requiredNetworks {
				if _, exists := networkMap[network]; !exists {
					flag = false

					break
				}
			}
		}

		if flag {
			nodes = append(nodes, address)
		}
	}

	return nodes, nil
}

// matchWorkerAndNetwork matches nodes based on both worker and network criteria.
// It takes a context, slices of workers and request networks as input parameters.
// It returns a slice of common node addresses and an error if any occurred.
func (d *Distributor) matchWorkerAndNetwork(ctx context.Context, workers, requestNetworks []string) ([]common.Address, error) {
	nodes := make([]common.Address, 0)

	indexers, err := d.databaseClient.FindNodeIndexers(ctx, nil, requestNetworks, workers)

	if err != nil {
		return nil, err
	}

	nodeWorkerToNetworksMap := make(map[common.Address]map[string][]string)

	for _, indexer := range indexers {
		if _, exists := nodeWorkerToNetworksMap[indexer.Address]; !exists {
			nodeWorkerToNetworksMap[indexer.Address] = make(map[string][]string)
		}

		nodeWorkerToNetworksMap[indexer.Address][indexer.Worker] = append(nodeWorkerToNetworksMap[indexer.Address][indexer.Worker], indexer.Network)
	}

	for address, workerToNetworksMap := range nodeWorkerToNetworksMap {
		flag := true

		// check every worker networks
		for worker, networks := range workerToNetworksMap {
			wid, _ := filter.NameString(worker)

			workerRequiredNetworks := WorkerToNetworksMap[wid]

			requiredNetworks := findCommonStr(workerRequiredNetworks, requestNetworks)

			networkMap := make(map[string]struct{}, len(networks))

			for _, network := range networks {
				networkMap[network] = struct{}{}
			}

			// check every network
			for _, network := range requiredNetworks {
				if _, exists := networkMap[network]; !exists {
					flag = false

					break
				}
			}
		}

		if flag {
			nodes = append(nodes, address)
		}
	}

	return nodes, nil
}

// findCommonStr finds common elements between two string slices.
// It takes two string slices as input parameters.
// It returns a slice containing common elements.
func findCommonStr(slice1, slice2 []string) []string {
	elementMap := make(map[string]struct{})

	var commonElements []string

	for _, v := range slice1 {
		elementMap[v] = struct{}{}
	}

	for _, v := range slice2 {
		if _, found := elementMap[v]; found {
			commonElements = append(commonElements, v)
			delete(elementMap, v)
		}
	}

	return commonElements
}

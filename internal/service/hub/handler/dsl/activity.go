package dsl

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/creasty/defaults"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/dsl"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/node/schema/worker/decentralized"
	"github.com/rss3-network/protocol-go/schema"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (d *DSL) GetActivity(c echo.Context) (err error) {
	var request dsl.ActivityRequest

	if err = c.Bind(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if err = defaults.Set(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if err = c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	activity, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestActivity, request, nil, nil)
	if err != nil {
		zap.L().Error("distribute activity request error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, activity)
}

func (d *DSL) GetAccountActivities(c echo.Context) (err error) {
	var request dsl.ActivitiesRequest

	if err = c.Bind(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if request.Type, err = parseParams(c.QueryParams(), request.Tag); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if err = defaults.Set(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if err = c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	workers, networks, err := validateCombinedParams(request.Tag, request.Network, request.Platform)
	if err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	// Resolve name to EVM address
	if !validEvmAddress(request.Account) {
		request.Account, err = d.nameService.Resolve(c.Request().Context(), request.Account)
		if err != nil {
			zap.L().Error("name service resolve error", zap.Error(err), zap.String("account", request.Account))

			return errorx.BadRequestError(c, err)
		}
	}

	activities, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestAccountActivities, request, workers, networks)
	if err != nil {
		zap.L().Error("distribute activities data error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

func (d *DSL) BatchGetAccountsActivities(c echo.Context) (err error) {
	var request dsl.AccountsActivitiesRequest

	if err = c.Bind(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if request.Type, err = parseParams(c.QueryParams(), request.Tag); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if err = defaults.Set(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if err = c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	workers, networks, err := validateCombinedParams(request.Tag, request.Network, request.Platform)
	if err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	activities, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestBatchAccountActivities, request, workers, networks)
	if err != nil {
		zap.L().Error("distribute batch activities data error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

func (d *DSL) GetNetworkActivities(c echo.Context) (err error) {
	var request dsl.NetworkActivitiesRequest

	if err = c.Bind(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if request.Type, err = parseParams(c.QueryParams(), request.Tag); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if err = defaults.Set(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if err = c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	workers, networks, err := validateCombinedParams(request.Tag, []string{request.Network}, request.Platform)
	if err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	activities, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestNetworkActivities, request, workers, networks)
	if err != nil {
		zap.L().Error("distribute network activities data error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

func (d *DSL) GetPlatformActivities(c echo.Context) (err error) {
	var request dsl.PlatformActivitiesRequest

	if err = c.Bind(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if request.Type, err = parseParams(c.QueryParams(), request.Tag); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if err = defaults.Set(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if err = c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	workers, networks, err := validateCombinedParams(request.Tag, request.Network, []string{request.Platform})
	if err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	activities, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestPlatformActivities, request, workers, networks)
	if err != nil {
		zap.L().Error("distribute platform activities data error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

// validEvmAddress checks if the address is a valid EVM address.
func validEvmAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}

// parseParams parses the type parameter and returns the corresponding types.
func parseParams(params url.Values, tags []string) ([]string, error) {
	if len(tags) == 0 {
		return nil, nil
	}

	types := make([]string, 0)

	for _, typeX := range params["type"] {
		for _, tagX := range tags {
			t, err := tag.TagString(tagX)

			if err != nil {
				continue
			}

			value, err := schema.ParseTypeFromString(t, typeX)

			if err != nil {
				continue
			}

			types = append(types, value.Name())

			break
		}

		if len(types) == 0 {
			return nil, fmt.Errorf("invalid type: %s", typeX)
		}
	}

	return types, nil
}

// validateCombinedParams validates the input tags, networks, and platforms and matches the workers.
func validateCombinedParams(inputTags, inputNetworks, inputPlatforms []string) ([]string, []string, error) {
	// Find network nodes that match the network requests.
	networks, networkWorks, err := getNetworks(inputNetworks)
	if err != nil {
		return nil, nil, err
	}

	// Find tag workers that match the tag requests.
	tagWorkers, err := getWorkersByTag(inputTags)
	if err != nil {
		return nil, nil, err
	}

	// Find platform workers that match the platform requests.
	platformWorkers, err := getWorkersByPlatform(inputPlatforms)
	if err != nil {
		return nil, nil, err
	}

	workers := combineWorkers(tagWorkers, networkWorks)
	// If no common workers are found between tag workers and network workers,
	// it indicates that tags and networks are not compatible.
	if len(workers) == 0 && (len(tagWorkers) > 0 || len(networkWorks) > 0) {
		return nil, nil, fmt.Errorf("no workers found for tags and networks")
	}

	workers = combineWorkers(networkWorks, platformWorkers)
	// If no common workers are found between network workers and platform workers,
	// it indicates that networks and platforms are not compatible.
	if len(workers) == 0 && (len(networkWorks) > 0 || len(platformWorkers) > 0) {
		return nil, nil, fmt.Errorf("no workers found for networks and platforms")
	}

	workers = combineWorkers(tagWorkers, platformWorkers)
	// If no common workers are found between tag workers and platform workers,
	// it indicates that tags and platforms are not compatible.
	if len(workers) == 0 && (len(tagWorkers) > 0 || len(platformWorkers) > 0) {
		return nil, nil, fmt.Errorf("no workers found for tags and platforms")
	}

	return lo.Keys(workers), networks, nil
}

type WorkerSet map[string]struct{}

// getNetworks returns a slice of networks based on the given network names.
func getNetworks(networks []string) ([]string, WorkerSet, error) {
	networkWorkers := make(WorkerSet)

	for i, n := range networks {
		nid, err := network.NetworkString(n)
		if err != nil {
			return nil, nil, err
		}

		networks[i] = nid.String()

		networkWorker, exists := model.NetworkToWorkersMap[nid.String()]
		if !exists {
			return nil, nil, fmt.Errorf("no workers found for network: %s", nid)
		}

		for _, w := range networkWorker {
			networkWorkers[w] = struct{}{}
		}
	}

	return networks, networkWorkers, nil
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
		pid, err := decentralized.PlatformString(platform)
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

// combineWorkers combines two worker sets.
func combineWorkers(workers1, workers2 WorkerSet) WorkerSet {
	if len(workers1) == 0 {
		return workers2
	}

	if len(workers2) == 0 {
		return workers1
	}

	commonWorkers := make(WorkerSet)

	for w := range workers1 {
		if _, exists := workers2[w]; exists {
			commonWorkers[w] = struct{}{}
		}
	}

	return commonWorkers
}

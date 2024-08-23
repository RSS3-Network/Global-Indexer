package dsl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/dsl"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/node/schema/worker/decentralized"
	"github.com/rss3-network/protocol-go/schema"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
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

	requestCounter.WithLabelValues("GetActivity").Inc()

	activity, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestActivity, request, c.QueryParams(), nil, nil)
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

	if request.Type, err = parseTypes(c.QueryParams()["type"], request.Tag); err != nil {
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
		resolvedName, err := d.getEVMAddress(c.Request().Context(), request.Account)
		if err == nil {
			request.Account = resolvedName
		}
	}

	incrementRequestCounter("GetAccountActivities", request.Network, request.Tag, request.Platform)

	activities, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestAccountActivities, request, c.QueryParams(), workers, networks)
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

	if request.Type, err = parseTypes(request.Type, request.Tag); err != nil {
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

	// Resolve names to EVM addresses
	if err = d.transformAccounts(c.Request().Context(), request.Accounts); err != nil {
		return errorx.BadRequestError(c, err)
	}

	request.Accounts = lo.Uniq(request.Accounts)

	incrementRequestCounter("BatchGetAccountsActivities", request.Network, request.Tag, request.Platform)

	activities, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestBatchAccountActivities, request, nil, workers, networks)
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

	if request.Type, err = parseTypes(c.QueryParams()["type"], request.Tag); err != nil {
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

	incrementRequestCounter("GetNetworkActivities", []string{request.Network}, request.Tag, request.Platform)

	activities, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestNetworkActivities, request, c.QueryParams(), workers, networks)
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

	if request.Type, err = parseTypes(c.QueryParams()["type"], request.Tag); err != nil {
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

	incrementRequestCounter("GetPlatformActivities", request.Network, request.Tag, []string{request.Platform})

	activities, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestPlatformActivities, request, c.QueryParams(), workers, networks)
	if err != nil {
		zap.L().Error("distribute platform activities data error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

func (d *DSL) transformAccounts(ctx context.Context, accounts []string) error {
	nsPool := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	for i := range accounts {
		i := i

		nsPool.Go(func(ctx context.Context) error {
			if !validEvmAddress(accounts[i]) {
				resolvedName, err := d.getEVMAddress(ctx, accounts[i])
				if err == nil {
					accounts[i] = resolvedName
				}
			}

			return nil
		})
	}

	return nsPool.Wait()
}

func (d *DSL) getEVMAddress(ctx context.Context, account string) (string, error) {
	var (
		err error

		key = buildNameServiceKey(account)
	)
	// Try to get the resolved address from Redis cache first
	if err = d.cacheClient.Get(ctx, key, &account); err == nil {
		return account, nil
	}

	// If not found in cache, resolve the name to an EVM address
	if errors.Is(err, redis.Nil) {
		// Cache miss, need to resolve the name
		account, err = d.nameService.Resolve(ctx, account)
		if err != nil {
			zap.L().Error("name service resolve error", zap.Error(err), zap.String("account", account))

			return "", err
		}

		// Cache the resolved address, with a TTL of 10 minutes
		// Fixme: The TTL should be configurable
		if err = d.cacheClient.Set(ctx, key, account, 10*time.Minute); err != nil {
			return "", err
		}

		return account, nil
	}

	return "", err
}

// buildNameServiceKey builds the key for the name service cache.
func buildNameServiceKey(account string) string {
	return fmt.Sprintf("name:service:%s", strings.ToLower(account))
}

// validEvmAddress checks if the address is a valid EVM address.
func validEvmAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}

// parseTypes parses the type parameter and returns the corresponding types.
func parseTypes(types []string, tags []string) ([]string, error) {
	if len(tags) == 0 {
		return nil, nil
	}

	var schemaTypes []string

	for _, typeX := range types {
		var (
			value schema.Type
			t     tag.Tag

			err error
		)

		for _, tagX := range tags {
			t, err = tag.TagString(tagX)

			if err != nil {
				return nil, err
			}

			value, err = schema.ParseTypeFromString(t, typeX)

			if err == nil {
				schemaTypes = append(schemaTypes, value.Name())
				break
			}
		}

		if err != nil {
			return nil, fmt.Errorf("invalid type: %s", typeX)
		}
	}

	return schemaTypes, nil
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
		return nil, nil, fmt.Errorf("no workers meet the conditions tags and networks")
	}

	workers = combineWorkers(networkWorks, platformWorkers)
	// If no common workers are found between network workers and platform workers,
	// it indicates that networks and platforms are not compatible.
	if len(workers) == 0 && (len(networkWorks) > 0 || len(platformWorkers) > 0) {
		return nil, nil, fmt.Errorf("no workers meet the conditions networks and platforms")
	}

	workers = combineWorkers(tagWorkers, platformWorkers)
	// If no common workers are found between tag workers and platform workers,
	// it indicates that tags and platforms are not compatible.
	if len(workers) == 0 && (len(tagWorkers) > 0 || len(platformWorkers) > 0) {
		return nil, nil, fmt.Errorf("no workers meet the conditions tags and platforms")
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

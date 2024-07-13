package dsl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

	activity, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestActivity, request)
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

	if err = validParams(request.Tag, request.Network, request.Platform); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	// Resolve name to EVM address
	if !validEvmAddress(request.Account) {
		request.Account, err = d.getEVMAddress(c.Request().Context(), request.Account)
		if err != nil {
			return errorx.BadRequestError(c, err)
		}
	}

	activities, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestAccountActivities, request)
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

	if err = validParams(request.Tag, request.Network, request.Platform); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	// Resolve names to EVM addresses
	if err = d.transformAccounts(c.Request().Context(), request.Accounts); err != nil {
		return errorx.BadRequestError(c, err)
	}

	request.Accounts = lo.Uniq(request.Accounts)

	activities, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestBatchAccountActivities, request)
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

	if err = validParams(request.Tag, []string{request.Network}, request.Platform); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	activities, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestNetworkActivities, request)
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

	if err = validParams(request.Tag, request.Network, []string{request.Platform}); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	activities, err := d.distributor.DistributeDecentralizedData(c.Request().Context(), model.DistributorRequestPlatformActivities, request)
	if err != nil {
		zap.L().Error("distribute platform activities data error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

func (d *DSL) transformAccounts(ctx context.Context, accounts []string) error {
	var err error

	nsPool := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError().WithMaxGoroutines(len(accounts))

	for i := range accounts {
		i := i

		nsPool.Go(func(ctx context.Context) error {
			if !validEvmAddress(accounts[i]) {
				accounts[i], err = d.getEVMAddress(ctx, accounts[i])
				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	if err = nsPool.Wait(); err != nil {
		return err
	}

	return nil
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

// validParams checks if the tags, networks, and platforms are valid.
func validParams(tags, networks, platforms []string) error {
	for _, tagX := range tags {
		if _, err := tag.TagString(tagX); err != nil {
			return err
		}
	}

	for _, networkX := range networks {
		if _, err := network.NetworkString(networkX); err != nil {
			return err
		}
	}

	for _, platform := range platforms {
		if _, err := decentralized.PlatformString(platform); err != nil {
			return err
		}
	}

	return nil
}

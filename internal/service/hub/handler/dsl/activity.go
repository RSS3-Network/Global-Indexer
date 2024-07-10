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
	"go.uber.org/zap"
)

func (d *DSL) GetActivity(c echo.Context) (err error) {
	var request dsl.ActivityRequest

	if err = c.Bind(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if err = c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	if err = defaults.Set(&request); err != nil {
		return errorx.BadRequestError(c, err)
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
		request.Account, err = d.nameService.Resolve(c.Request().Context(), request.Account)
		if err != nil {
			zap.L().Error("name service resolve error", zap.Error(err), zap.String("account", request.Account))

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

package dsl

import (
	"errors"
	"net/http"

	"github.com/creasty/defaults"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/dsl"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (d *DSL) GetFederatedActivity(c echo.Context) (err error) {
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

	requestCounter.WithLabelValues("GetFederatedActivity").Inc()

	activity, err := d.distributor.DistributeData(c.Request().Context(), model.DistributorRequestActivity, model.ComponentFederated, request, c.QueryParams(), nil, nil)
	if err != nil {
		if errors.Is(err, errorx.ErrNoNodesAvailable) {
			return errorx.BadRequestError(c, err)
		}

		zap.L().Error("distribute activity request error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, activity)
}

func (d *DSL) GetFederatedAccountActivities(c echo.Context) (err error) {
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

	incrementRequestCounter("GetFederatedAccountActivities", request.Network, request.Tag, request.Platform)

	activities, err := d.distributor.DistributeData(c.Request().Context(), model.DistributorRequestAccountActivities, model.ComponentFederated, request, c.QueryParams(), nil, nil)
	if err != nil {
		if errors.Is(err, errorx.ErrNoNodesAvailable) {
			return errorx.BadRequestError(c, err)
		}

		zap.L().Error("distribute activities data error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

func (d *DSL) BatchGetFederatedAccountsActivities(c echo.Context) (err error) {
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

	request.Accounts = lo.Uniq(request.Accounts)
	request.Network = lo.Uniq(request.Network)

	incrementRequestCounter("BatchGetFederatedAccountsActivities", request.Network, request.Tag, request.Platform)

	activities, err := d.distributor.DistributeData(c.Request().Context(), model.DistributorRequestBatchAccountActivities, model.ComponentFederated, request, nil, nil, request.Network)
	if err != nil {
		if errors.Is(err, errorx.ErrNoNodesAvailable) {
			return errorx.BadRequestError(c, err)
		}

		zap.L().Error("distribute batch activities data error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

func (d *DSL) GetFederatedNetworkActivities(c echo.Context) (err error) {
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

	incrementRequestCounter("GetFederatedNetworkActivities", []string{request.Network}, request.Tag, request.Platform)

	activities, err := d.distributor.DistributeData(c.Request().Context(), model.DistributorRequestNetworkActivities, model.ComponentFederated, request, c.QueryParams(), nil, []string{request.Network})
	if err != nil {
		if errors.Is(err, errorx.ErrNoNodesAvailable) {
			return errorx.BadRequestError(c, err)
		}

		zap.L().Error("distribute network activities data error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

func (d *DSL) GetFederatedPlatformActivities(c echo.Context) (err error) {
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

	incrementRequestCounter("GetFederatedPlatformActivities", request.Network, request.Tag, []string{request.Platform})

	activities, err := d.distributor.DistributeData(c.Request().Context(), model.DistributorRequestPlatformActivities, model.ComponentFederated, request, c.QueryParams(), nil, request.Network)
	if err != nil {
		if errors.Is(err, errorx.ErrNoNodesAvailable) {
			return errorx.BadRequestError(c, err)
		}

		zap.L().Error("distribute platform activities data error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

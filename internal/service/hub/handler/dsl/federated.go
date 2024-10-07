package dsl

import (
	"net/http"

	"github.com/creasty/defaults"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/dsl"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
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
		zap.L().Error("distribute activities data error", zap.Error(err))

		return errorx.InternalError(c)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

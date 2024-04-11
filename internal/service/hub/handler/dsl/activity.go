package dsl

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/creasty/defaults"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/hub/model/dsl"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/protocol-go/schema/filter"
)

func (d *DSL) GetActivity(c echo.Context) (err error) {
	var request dsl.ActivityRequest

	if err = c.Bind(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if err = c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, err)
	}

	if err = defaults.Set(&request); err != nil {
		return errorx.InternalError(c, err)
	}

	activity, err := d.Distributor.RouterActivityData(c.Request().Context(), request)
	if err != nil {
		return errorx.InternalError(c, err)
	}

	return c.JSONBlob(http.StatusOK, activity)
}

func (d *DSL) GetAccountActivities(c echo.Context) (err error) {
	var request dsl.AccountActivitiesRequest

	if err = c.Bind(&request); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if request.Type, err = d.parseParams(c.QueryParams(), request.Tag); err != nil {
		return errorx.BadRequestError(c, err)
	}

	if err = c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, err)
	}

	if err = defaults.Set(&request); err != nil {
		return errorx.InternalError(c, err)
	}

	if !d.validEvmAddress(request.Account) {
		request.Account, err = d.nameService.Resolve(c.Request().Context(), request.Account)
		if err != nil {
			return errorx.InternalError(c, err)
		}
	}

	activities, err := d.Distributor.RouterActivitiesData(c.Request().Context(), request)
	if err != nil {
		return errorx.InternalError(c, err)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

func (d *DSL) validEvmAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}

func (d *DSL) parseParams(params url.Values, tags []string) ([]string, error) {
	if len(tags) == 0 {
		return nil, nil
	}

	types := make([]string, 0)

	for _, typex := range params["type"] {
		var (
			value filter.Type
			err   error
		)

		for _, tag := range tags {
			t, err := filter.TagString(tag)
			if err == nil {
				value, err = filter.TypeString(t, typex)
				if err == nil {
					types = append(types, value.Name())

					break
				}
			}
		}

		if err != nil {
			return nil, fmt.Errorf("invalid type: %s", typex)
		}
	}

	return types, nil
}

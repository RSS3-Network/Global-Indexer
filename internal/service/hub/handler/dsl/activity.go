package dsl

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/creasty/defaults"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/dsl"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
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

	activity, err := d.distributor.RouteActivityRequest(c.Request().Context(), request)
	if err != nil {
		return errorx.InternalError(c, err)
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

	if err = c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, err)
	}

	if err = defaults.Set(&request); err != nil {
		return errorx.InternalError(c, err)
	}

	// Resolve name to EVM address
	if !validEvmAddress(request.Account) {
		request.Account, err = d.nameService.Resolve(c.Request().Context(), request.Account)
		if err != nil {
			return errorx.InternalError(c, err)
		}
	}

	activities, err := d.distributor.RouteActivitiesData(c.Request().Context(), request)
	if err != nil {
		return errorx.InternalError(c, err)
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
		for _, tag := range tags {
			t, err := filter.TagString(tag)

			if err != nil {
				continue
			}

			value, err := filter.TypeString(t, typeX)

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

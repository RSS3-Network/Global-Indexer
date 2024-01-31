package hub

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/creasty/defaults"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model/response"
	"github.com/rss3-network/protocol-go/schema/filter"
)

func (h *Hub) GetActivityHandler(c echo.Context) (err error) {
	request := model.ActivityRequest{}

	if err = c.Bind(&request); err != nil {
		return response.BadRequestError(c, err)
	}

	if err = c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, err)
	}

	if err = defaults.Set(&request); err != nil {
		return response.InternalError(c, err)
	}

	activity, err := h.routerActivityData(c.Request().Context(), request)
	if err != nil {
		return response.InternalError(c, err)
	}

	return c.JSONBlob(http.StatusOK, activity)
}

func (h *Hub) GetAccountActivitiesHandler(c echo.Context) (err error) {
	request := model.AccountActivitiesRequest{}

	if err = c.Bind(&request); err != nil {
		return response.BadRequestError(c, err)
	}

	if request.Type, err = h.parseParams(c.QueryParams(), request.Tag); err != nil {
		return response.BadRequestError(c, err)
	}

	if err = c.Validate(&request); err != nil {
		return response.ValidateFailedError(c, err)
	}

	if err = defaults.Set(&request); err != nil {
		return response.InternalError(c, err)
	}

	activities, err := h.routerActivitiesData(c.Request().Context(), request)
	if err != nil {
		return response.InternalError(c, err)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

func (h *Hub) parseParams(params url.Values, tags []string) ([]string, error) {
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

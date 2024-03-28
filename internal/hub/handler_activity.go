package hub

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/creasty/defaults"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model/response"
	"github.com/rss3-network/protocol-go/schema/filter"
)

func (h *Hub) GetRSSHub(c echo.Context) error {
	path := c.Param("*")
	query := c.Request().URL.RawQuery

	data, err := h.routerRSSHubData(c.Request().Context(), path, query)

	if err != nil {
		return response.InternalError(c, err)
	}

	return c.JSONBlob(http.StatusOK, data)
}

type ActivityRequest struct {
	ID          string `param:"id"`
	ActionLimit int    `query:"action_limit" default:"10" min:"1" max:"20"`
	ActionPage  int    `query:"action_page" default:"1" min:"1"`
}

func (h *Hub) GetActivity(c echo.Context) (err error) {
	var request ActivityRequest

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

type AccountActivitiesRequest struct {
	Account        string   `param:"account"`
	Limit          *int     `query:"limit" default:"100" min:"1" max:"100"`
	ActionLimit    *int     `query:"action_limit" default:"10" min:"1" max:"20"`
	Cursor         *string  `query:"cursor"`
	SinceTimestamp *uint64  `query:"since_timestamp"`
	UntilTimestamp *uint64  `query:"until_timestamp"`
	Status         *bool    `query:"success"`
	Direction      *string  `query:"direction"`
	Network        []string `query:"network"`
	Tag            []string `query:"tag"`
	Type           []string `query:"-"`
	Platform       []string `query:"platform"`
}

func (h *Hub) GetAccountActivities(c echo.Context) (err error) {
	var request AccountActivitiesRequest

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

	if !h.validEvmAddress(request.Account) {
		request.Account, err = h.nameService.Resolve(c.Request().Context(), request.Account)
		if err != nil {
			return response.InternalError(c, err)
		}
	}

	activities, err := h.routerActivitiesData(c.Request().Context(), request)
	if err != nil {
		return response.InternalError(c, err)
	}

	return c.JSONBlob(http.StatusOK, activities)
}

func (h *Hub) validEvmAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
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

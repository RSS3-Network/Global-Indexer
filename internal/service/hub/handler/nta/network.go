package nta

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
)

// GetNetworks returns all networks supported by the DSL.
func (n *NTA) GetNetworks(c echo.Context) error {
	endpoint, err := n.getNodeEndpoint(c)
	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("get node endpoint: %w", err))
	}

	return n.fetchResponse(c, fmt.Sprintf("%s/networks", endpoint))
}

// GetWorkersByNetwork returns all available workers on a specific Network.
func (n *NTA) GetWorkersByNetwork(c echo.Context) error {
	var request nta.NetworkRequest

	if err := n.bindAndValidateRequest(c, &request); err != nil {
		return err
	}

	endpoint, err := n.getNodeEndpoint(c)
	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("get node endpoint: %w", err))
	}

	return n.fetchResponse(c, fmt.Sprintf("%s/networks/%s/list_workers", endpoint, request.Network))
}

// GetWorkerDetail returns a worker's detail and possible configuration.
func (n *NTA) GetWorkerDetail(c echo.Context) error {
	var request nta.WorkerRequest

	if err := n.bindAndValidateRequest(c, &request); err != nil {
		return err
	}

	endpoint, err := n.getNodeEndpoint(c)

	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("get node endpoint: %w", err))
	}

	return n.fetchResponse(c, fmt.Sprintf("%s/networks/%s/workers/%s", endpoint, request.Network, request.Worker))
}

// GetEndpointConfig returns possible configurations for an endpoint.
func (n *NTA) GetEndpointConfig(c echo.Context) error {
	endpoint, err := n.getNodeEndpoint(c)

	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("get node endpoint: %w", err))
	}

	return n.fetchResponse(c, fmt.Sprintf("%s/networks/get_endpoint_config", endpoint))
}

// bindAndValidateRequest binds and validates the request.
func (n *NTA) bindAndValidateRequest(c echo.Context, request interface{}) error {
	if err := c.Bind(request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validation failed: %w", err))
	}

	return nil
}

// getNodeEndpoint returns the endpoint of the node with the highest points.
func (n *NTA) getNodeEndpoint(c echo.Context) (string, error) {
	stats, err := n.databaseClient.FindNodeStats(c.Request().Context(), &schema.StatQuery{
		ValidRequest: lo.ToPtr(model.DemotionCountBeforeSlashing),
		Limit:        lo.ToPtr(model.RequiredQualifiedNodeCount),
		PointsOrder:  lo.ToPtr("DESC"),
	})

	if err != nil {
		return "", err
	}

	if len(stats) == 0 {
		return "", fmt.Errorf("no available nodes")
	}

	return stats[0].Endpoint, nil
}

// fetchResponse fetches data from the provided endpoint.
func (n *NTA) fetchResponse(c echo.Context, url string) error {
	response, err := n.httpClient.Fetch(c.Request().Context(), url)
	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("network: %w", err))
	}

	data, err := io.ReadAll(response)

	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("read response: %w", err))
	}

	return c.JSONBlob(http.StatusOK, data)
}

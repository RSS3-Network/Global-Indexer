package nta

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

	return n.fetchResponse(c, fmt.Sprintf("%s/networks/%s/list_workers", endpoint, request.NetworkName))
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

	return n.fetchResponse(c, fmt.Sprintf("%s/networks/%s/workers/%s", endpoint, request.NetworkName, request.WorkerName))
}

// GetEndpointConfig returns possible configurations for an endpoint.
func (n *NTA) GetEndpointConfig(c echo.Context) error {
	endpoint, err := n.getNodeEndpoint(c)

	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("get node endpoint: %w", err))
	}

	return n.fetchResponse(c, fmt.Sprintf("%s/networks/endpoint_config", endpoint))
}

// GetAssets returns all assets supported by the DSL.
func (n *NTA) GetAssets(c echo.Context) error {
	// Get parameters for the current epoch from networkParams
	params, err := n.networkParamsContract.GetParams(&bind.CallOpts{}, math.MaxUint64)

	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("failed to get params for epoch %w", err))
	}

	var networkParam nta.NetworkParamsData
	if err = json.Unmarshal([]byte(params), &networkParam); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("failed to unmarshal network params %w", err))
	}

	return c.JSON(http.StatusOK, nta.Response{Data: struct {
		Networks map[string]nta.Asset `json:"networks"`
		Workers  map[string]nta.Asset `json:"workers"`
	}{
		Networks: networkParam.NetworkAssets,
		Workers:  networkParam.WorkerAssets,
	}})
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
		IsFullNode:   lo.ToPtr(true),
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
	response, err := n.httpClient.FetchWithMethod(c.Request().Context(), http.MethodGet, url, "", nil)
	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("network: %w", err))
	}

	data, err := io.ReadAll(response)

	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("read response: %w", err))
	}

	return c.JSONBlob(http.StatusOK, data)
}

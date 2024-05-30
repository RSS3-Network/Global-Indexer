package nta

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
)

// GetNetworks returns all networks supported by the DSL.
func (n *NTA) GetNetworks(c echo.Context) error {
	var request nta.NodeRequest

	if err := n.bindAndValidateRequest(c, &request); err != nil {
		return err
	}

	endpoint, err := n.getNodeEndpoint(c, request.Address)
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

	endpoint, err := n.getNodeEndpoint(c, request.Address)
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

	endpoint, err := n.getNodeEndpoint(c, request.Address)

	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("get node endpoint: %w", err))
	}

	return n.fetchResponse(c, fmt.Sprintf("%s/networks/%s/workers/%s", endpoint, request.NetworkName, request.WorkerName))
}

// GetEndpointConfig returns possible configurations for an endpoint.
func (n *NTA) GetEndpointConfig(c echo.Context) error {
	var request nta.NodeRequest

	if err := n.bindAndValidateRequest(c, &request); err != nil {
		return err
	}

	endpoint, err := n.getNodeEndpoint(c, request.Address)

	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("get node endpoint: %w", err))
	}

	return n.fetchResponse(c, fmt.Sprintf("%s/networks/endpoint_config", endpoint))
}

// GetNodeVersion returns the version of the node.
func (n *NTA) GetNodeVersion(c echo.Context) error {
	var request nta.NodeRequest

	if err := n.bindAndValidateRequest(c, &request); err != nil {
		return err
	}

	endpoint, err := n.getNodeEndpoint(c, request.Address)

	if err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("get node endpoint: %w", err))
	}

	return n.fetchResponse(c, fmt.Sprintf("%s/version", endpoint))
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

// getNodeEndpoint returns the endpoint of the node.
func (n *NTA) getNodeEndpoint(c echo.Context, nodeAddress common.Address) (string, error) {
	node, err := n.databaseClient.FindNode(c.Request().Context(), nodeAddress)

	if err != nil {
		return "", err
	}

	if node == nil {
		return "", fmt.Errorf("no available nodes")
	}

	return node.Endpoint, nil
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

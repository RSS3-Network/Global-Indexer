package nta

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
)

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

// GetNetworkConfig returns the network configuration for the current epoch.
func (n *NTA) GetNetworkConfig(c echo.Context) error {
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
		RSSConfig           any `json:"rss"`
		DecentralizedConfig any `json:"decentralized"`
		FederatedConfig     any `json:"federated"`
	}{
		RSSConfig:           networkParam.NetworkConfig["rss"],
		DecentralizedConfig: networkParam.NetworkConfig["decentralized"],
		FederatedConfig:     networkParam.NetworkConfig["federated"],
	}})
}

package nta

import (
	"context"
	"math/big"
	"net/url"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/common/geolite2"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/database"
)

type NTA struct {
	databaseClient          database.Client
	stakingContract         *l2.StakingV2MulticallClient
	networkParamsContract   *l2.NetworkParams
	contractGovernanceToken *bindings.GovernanceToken
	geoLite2                *geolite2.Client
	cacheClient             cache.Client
	httpClient              httputil.Client
}

var MinDeposit = new(big.Int).Mul(big.NewInt(10000), big.NewInt(1e18))

func (n *NTA) baseURL(c echo.Context) url.URL {
	return url.URL{
		Scheme: c.Scheme(),
		Host:   c.Request().Host,
	}
}

func NewNTA(_ context.Context, databaseClient database.Client, stakingContract *l2.StakingV2MulticallClient, networkParamsContract *l2.NetworkParams, contractGovernanceToken *bindings.GovernanceToken, geoLite2 *geolite2.Client, cacheClient cache.Client, httpClient httputil.Client) *NTA {
	return &NTA{
		databaseClient:          databaseClient,
		stakingContract:         stakingContract,
		networkParamsContract:   networkParamsContract,
		contractGovernanceToken: contractGovernanceToken,
		geoLite2:                geoLite2,
		cacheClient:             cacheClient,
		httpClient:              httpClient,
	}
}

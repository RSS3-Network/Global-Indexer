package nta

import (
	"context"
	"math/big"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/common/geolite2"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/contract/l2"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/database"
)

type NTA struct {
	databaseClient        database.Client
	stakingContract       *stakingv2.Staking
	networkParamsContract *l2.NetworkParams
	geoLite2              *geolite2.Client
	cacheClient           cache.Client
	httpClient            httputil.Client
}

var MinDeposit = new(big.Int).Mul(big.NewInt(10000), big.NewInt(1e18))

func (n *NTA) baseURL(c echo.Context) url.URL {
	return url.URL{
		Scheme: c.Scheme(),
		Host:   c.Request().Host,
	}
}

func NewNTA(_ context.Context, databaseClient database.Client, stakingContract *stakingv2.Staking, networkParamsContract *l2.NetworkParams, geoLite2 *geolite2.Client, cacheClient cache.Client, httpClient httputil.Client) *NTA {
	return &NTA{
		databaseClient:        databaseClient,
		stakingContract:       stakingContract,
		networkParamsContract: networkParamsContract,
		geoLite2:              geoLite2,
		cacheClient:           cacheClient,
		httpClient:            httpClient,
	}
}

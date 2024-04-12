package nta

import (
	"context"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/common/geolite2"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/database"
)

type NTA struct {
	databaseClient  database.Client
	stakingContract *l2.Staking
	geoLite2        *geolite2.Client
	cacheClient     cache.Client
}

func (n *NTA) baseURL(c echo.Context) url.URL {
	return url.URL{
		Scheme: c.Scheme(),
		Host:   c.Request().Host,
	}
}

func NewNTA(_ context.Context, databaseClient database.Client, stakingContract *l2.Staking, geoLite2 *geolite2.Client, cacheClient cache.Client) *NTA {
	return &NTA{
		databaseClient:  databaseClient,
		stakingContract: stakingContract,
		geoLite2:        geoLite2,
		cacheClient:     cacheClient,
	}
}

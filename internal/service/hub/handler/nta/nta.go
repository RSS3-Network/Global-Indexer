package nta

import (
	"context"
	"math/big"
	"net/url"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/common/geolite2"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/config"
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
	erc20TokenMap           map[common.Address]*bindings.GovernanceToken
	configFile              *config.File
	chainL1ID               uint64
	chainL2ID               uint64
}

var MinDeposit = new(big.Int).Mul(big.NewInt(10000), big.NewInt(1e18))

func (n *NTA) baseURL(c echo.Context) url.URL {
	return url.URL{
		Scheme: c.Scheme(),
		Host:   c.Request().Host,
	}
}

func NewNTA(_ context.Context, configFile *config.File, databaseClient database.Client, stakingContract *l2.StakingV2MulticallClient, networkParamsContract *l2.NetworkParams, contractGovernanceToken *bindings.GovernanceToken, geoLite2 *geolite2.Client, cacheClient cache.Client, httpClient httputil.Client, erc20TokenMap map[common.Address]*bindings.GovernanceToken, chainL1ID, chainL2ID uint64) *NTA {
	return &NTA{
		databaseClient:          databaseClient,
		stakingContract:         stakingContract,
		networkParamsContract:   networkParamsContract,
		contractGovernanceToken: contractGovernanceToken,
		geoLite2:                geoLite2,
		cacheClient:             cacheClient,
		httpClient:              httpClient,
		erc20TokenMap:           erc20TokenMap,
		configFile:              configFile,
		chainL1ID:               chainL1ID,
		chainL2ID:               chainL2ID,
	}
}

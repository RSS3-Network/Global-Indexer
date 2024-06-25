package hub

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/geolite2"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/contract/l2"
	stakingv1 "github.com/rss3-network/global-indexer/contract/l2/staking/v1"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/client/ethereum"
	"github.com/rss3-network/global-indexer/internal/config/flag"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/nameresolver"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/nta"
	"github.com/spf13/viper"
)

type Hub struct {
	dsl *dsl.DSL
	nta *nta.NTA
}

var _ echo.Validator = (*Validator)(nil)

var defaultValidator = &Validator{
	validate: validator.New(),
}

type Validator struct {
	validate *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func NewHub(ctx context.Context, databaseClient database.Client, redisClient *redis.Client, ethereumMultiChainClient *ethereum.MultiChainClient, geoLite2 *geolite2.Client, nameService *nameresolver.NameResolver, httpClient httputil.Client) (*Hub, error) {
	chainID := viper.GetUint64(flag.KeyChainIDL2)

	ethereumClient, err := ethereumMultiChainClient.Get(chainID)
	if err != nil {
		return nil, fmt.Errorf("get ethereum client: %w", err)
	}

	contractAddresses := l2.ContractMap[chainID]
	if contractAddresses == nil {
		return nil, fmt.Errorf("contract address not found for chain id: %d", chainID)
	}

	stakingContract, err := stakingv1.NewStaking(contractAddresses.AddressStakingProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new staking contract: %w", err)
	}

	cacheClient := cache.New(redisClient)

	dsl, err := dsl.NewDSL(ctx, databaseClient, cacheClient, nameService, stakingContract, httpClient)
	if err != nil {
		return nil, fmt.Errorf("new dsl: %w", err)
	}

	return &Hub{
		dsl: dsl,
		nta: nta.NewNTA(ctx, databaseClient, stakingContract, geoLite2, cacheClient, httpClient),
	}, nil
}

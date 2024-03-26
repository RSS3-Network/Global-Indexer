package hub

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/geolite2"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cache"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/nameresolver"
	"github.com/redis/go-redis/v9"
)

type Hub struct {
	databaseClient  database.Client
	geoLite2        *geolite2.Client
	cacheClient     cache.Client
	stakingContract *l2.Staking
	httpClient      *http.Client
	nameService     *nameresolver.NameResolver
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

func NewHub(ctx context.Context, databaseClient database.Client, ethereumClient *ethclient.Client, redisClient *redis.Client, geoLite2 *geolite2.Client, nameService *nameresolver.NameResolver) (*Hub, error) {
	chainID, err := ethereumClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}

	contractAddresses := l2.ContractMap[chainID.Uint64()]
	if contractAddresses == nil {
		return nil, fmt.Errorf("contract address not found for chain id: %d", chainID.Uint64())
	}

	stakingContract, err := l2.NewStaking(contractAddresses.AddressStakingProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new staking contract: %w", err)
	}

	return &Hub{
		databaseClient:  databaseClient,
		geoLite2:        geoLite2,
		cacheClient:     cache.New(redisClient),
		stakingContract: stakingContract,
		httpClient:      http.DefaultClient,
		nameService:     nameService,
	}, nil
}

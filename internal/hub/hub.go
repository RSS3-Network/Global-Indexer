package hub

import (
	"context"
	"fmt"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/geolite2"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cache"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/redis/go-redis/v9"
)

type Hub struct {
	databaseClient  database.Client
	geoLite2        *geolite2.Client
	cacheClient     cache.Client
	stakingContract *l2.Staking
	httpClient      *http.Client
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

func NewHub(_ context.Context, databaseClient database.Client, ethereumClient *ethclient.Client, redisClient *redis.Client, geoLite2 *geolite2.Client) (*Hub, error) {
	stakingContract, err := l2.NewStaking(l2.AddressStakingProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new staking contract: %w", err)
	}

	return &Hub{
		databaseClient:  databaseClient,
		geoLite2:        geoLite2,
		cacheClient:     cache.New(redisClient),
		stakingContract: stakingContract,
		httpClient:      http.DefaultClient,
	}, nil
}

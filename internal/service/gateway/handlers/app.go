package handlers

import (
	apisixHTTPAPI "github.com/naturalselectionlabs/rss3-global-indexer/internal/apisix/httpapi"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/jwt"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/siwe"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ oapi.ServerInterface = (*App)(nil)

type App struct {
	apisixHTTPAPIClient *apisixHTTPAPI.Client
	redisClient         *redis.Client
	databaseClient      *gorm.DB
	jwtClient           *jwt.JWT
	siweClient          *siwe.SIWE
}

func NewApp(apiService *apisixHTTPAPI.Client, redis *redis.Client, databaseClient *gorm.DB, jwtClient *jwt.JWT, siweClient *siwe.SIWE) (*App, error) {
	return &App{
		apisixHTTPAPIClient: apiService,
		redisClient:         redis,
		databaseClient:      databaseClient,
		jwtClient:           jwtClient,
		siweClient:          siweClient,
	}, nil
}

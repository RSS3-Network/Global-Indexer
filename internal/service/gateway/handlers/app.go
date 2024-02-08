package handlers

import (
	apisixHTTPAPI "github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/apisix/httpapi"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ oapi.ServerInterface = (*App)(nil)

type App struct {
	siweDomain       string
	apiSixAPIService *apisixHTTPAPI.HTTPAPIService
	redisClient      *redis.Client
	databaseClient   *gorm.DB
}

func NewApp(siweDomain string, apiService *apisixHTTPAPI.HTTPAPIService, redis *redis.Client, databaseClient *gorm.DB) *App {
	return &App{
		siweDomain:       siweDomain,
		apiSixAPIService: apiService,
		redisClient:      redis,
		databaseClient:   databaseClient,
	}
}

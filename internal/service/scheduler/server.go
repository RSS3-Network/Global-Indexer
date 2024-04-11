package scheduler

import (
	"fmt"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/client/ethereum"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config/flag"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/detector"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/integrator"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/snapshot"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler/taxer"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewServer(databaseClient database.Client, redis *redis.Client, ethereumMultiChainClient *ethereum.MultiChainClient, config *config.File) (service.Server, error) {
	ethereumClient, err := ethereumMultiChainClient.Get(viper.GetUint64(flag.KeyChainIDL2))
	if err != nil {
		return nil, fmt.Errorf("get ethereum client: %w", err)
	}

	switch server := viper.GetString(flag.KeyServer); server {
	case detector.Name:
		return detector.New(databaseClient, redis)
	case integrator.Name:
		return integrator.New(databaseClient, redis, ethereumClient)
	case snapshot.Name:
		return snapshot.New(databaseClient, redis, ethereumClient)
	case taxer.Name:
		return taxer.New(databaseClient, redis, ethereumClient, config)
	default:
		return nil, fmt.Errorf("unknown scheduler server: %s", server)
	}
}

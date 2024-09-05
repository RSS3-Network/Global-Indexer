package scheduler

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/internal/client/ethereum"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/config/flag"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service"
	"github.com/rss3-network/global-indexer/internal/service/scheduler/detector"
	"github.com/rss3-network/global-indexer/internal/service/scheduler/enforcer"
	"github.com/rss3-network/global-indexer/internal/service/scheduler/snapshot"
	"github.com/rss3-network/global-indexer/internal/service/scheduler/taxer"
	"github.com/spf13/viper"
)

// NewServer creates a new scheduler server that executes cron jobs.
func NewServer(databaseClient database.Client, redis *redis.Client, ethereumMultiChainClient *ethereum.MultiChainClient, httpClient httputil.Client, config *config.File, txManager *txmgr.SimpleTxManager) (service.Server, error) {
	ethereumClient, err := ethereumMultiChainClient.Get(viper.GetUint64(flag.KeyChainIDL2))
	if err != nil {
		return nil, fmt.Errorf("get ethereum client: %w", err)
	}

	switch server := viper.GetString(flag.KeyServer); server {
	case detector.Name:
		return detector.New(databaseClient, redis)
	case enforcer.Name:
		return enforcer.New(databaseClient, redis, ethereumClient, httpClient)
	case snapshot.Name:
		return snapshot.New(databaseClient, redis, ethereumClient)
	case taxer.Name:
		return taxer.New(databaseClient, redis, ethereumClient, config, txManager)
	default:
		return nil, fmt.Errorf("unknown scheduler server: %s", server)
	}
}

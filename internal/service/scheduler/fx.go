package scheduler

import (
	"github.com/rss3-network/global-indexer/internal/provider"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(provider.ProvideDatabaseClient),
	fx.Provide(provider.ProvideRedisClient),
	fx.Provide(provider.ProvideEthereumMultiChainClient),
	fx.Provide(provider.ProvideHTTPClient),
	fx.Provide(provider.ProvideTxManager),
)

package hub

import (
	"github.com/rss3-network/global-indexer/internal/provider"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(provider.ProvideDatabaseClient),
	fx.Provide(provider.ProvideRedisClient),
	fx.Provide(provider.ProvideEthereumMultiChainClient),
	fx.Provide(provider.ProvideGeoIP2),
	fx.Provide(provider.ProvideNameResolver),
)

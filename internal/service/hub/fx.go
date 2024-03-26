package hub

import (
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/provider"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(provider.ProvideDatabaseClient),
	fx.Provide(provider.ProvideRedisClient),
	fx.Provide(provider.ProvideEthereumMultiChainClient),
	fx.Provide(provider.ProvideGeoIP2),
	fx.Provide(provider.ProvideNameResolver),
)

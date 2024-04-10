package provider

import (
	"context"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/nameresolver"
)

func ProvideNameResolver(configFile *config.File) (*nameresolver.NameResolver, error) {
	return nameresolver.NewNameResolver(context.TODO(), configFile.RPC.RPCNetwork)
}

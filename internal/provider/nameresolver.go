package provider

import (
	"context"

	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/nameresolver"
)

func ProvideNameResolver(configFile *config.File) (*nameresolver.NameResolver, error) {
	return nameresolver.NewNameResolver(context.TODO(), configFile.RPC.RPCNetwork)
}

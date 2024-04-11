package provider

import (
	"context"

	"github.com/rss3-network/global-indexer/internal/client/ethereum"
	"github.com/rss3-network/global-indexer/internal/config"
)

func ProvideEthereumMultiChainClient(configFile *config.File) (*ethereum.MultiChainClient, error) {
	endpoints := []string{
		configFile.RSS3Chain.EndpointL1,
		configFile.RSS3Chain.EndpointL2,
	}

	return ethereum.Dial(context.TODO(), endpoints)
}

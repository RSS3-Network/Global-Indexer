package provider

import (
	"context"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/client/ethereum"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
)

func ProvideEthereumMultiChainClient(configFile *config.File) (*ethereum.MultiChainClient, error) {
	endpoints := []string{
		configFile.RSS3Chain.EndpointL1,
		configFile.RSS3Chain.EndpointL2,
	}

	return ethereum.Dial(context.TODO(), endpoints)
}

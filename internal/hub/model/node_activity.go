package model

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/form/v4"
	"github.com/rss3-network/protocol-go/schema/filter"
)

var (
	RssNodeCacheKey  = "nodes:rss"
	FullNodeCacheKey = "nodes:full"

	MessageNodeDataFailed = "failed to request node data "

	DefaultNodeCount   = 3
	DefaultSlashCount  = 4
	DefaultVerifyCount = 3
)

type Cache struct {
	Address  string `json:"address"`
	Endpoint string `json:"endpoint"`
}

// WorkerToNetworksMap Supplement the conditions for a full node based on the configuration file.
// https://github.com/NaturalSelectionLabs/RSS3-Node/blob/develop/deploy/config.development.yaml
var WorkerToNetworksMap = map[filter.Name][]string{
	filter.Fallback: {
		filter.NetworkEthereum.String(),
	},
	filter.Mirror: {
		filter.NetworkArweave.String(),
	},
	filter.Farcaster: {
		filter.NetworkFarcaster.String(),
	},
	filter.RSS3: {
		filter.NetworkEthereum.String(),
	},
	filter.Paragraph: {
		filter.NetworkArweave.String(),
	},
	filter.OpenSea: {
		filter.NetworkEthereum.String(),
	},
	filter.Uniswap: {
		filter.NetworkEthereum.String(),
	},
	filter.Optimism: {
		filter.NetworkEthereum.String(),
	},
	filter.Aavegotchi: {
		filter.NetworkPolygon.String(),
	},
	filter.Lens: {
		filter.NetworkPolygon.String(),
	},
	filter.Looksrare: {
		filter.NetworkEthereum.String(),
	},
	filter.Matters: {
		filter.NetworkPolygon.String(),
	},
	filter.Momoka: {
		filter.NetworkArweave.String(),
	},
	filter.Highlight: {
		filter.NetworkEthereum.String(),
	},
	filter.Aave: {
		filter.NetworkPolygon.String(),
		filter.NetworkEthereum.String(),
		filter.NetworkAvalanche.String(),
		filter.NetworkBase.String(),
		filter.NetworkOptimism.String(),
		filter.NetworkArbitrum.String(),
		filter.NetworkFantom.String(),
	},
	filter.IQWiki: {
		filter.NetworkPolygon.String(),
	},
	filter.Lido: {
		filter.NetworkEthereum.String(),
	},
	filter.Crossbell: {
		filter.NetworkCrossbell.String(),
	},
	filter.ENS: {
		filter.NetworkEthereum.String(),
	},
}

var NetworkToWorkersMap = map[filter.Network][]string{
	filter.NetworkEthereum: {
		filter.Fallback.String(),
		filter.RSS3.String(),
		filter.OpenSea.String(),
		filter.Uniswap.String(),
		filter.Optimism.String(),
		filter.Looksrare.String(),
		filter.Highlight.String(),
		filter.Aave.String(),
		filter.Lido.String(),
		filter.ENS.String(),
	},
	filter.NetworkArweave: {
		filter.Mirror.String(),
		filter.Paragraph.String(),
		filter.Momoka.String(),
	},
	filter.NetworkFarcaster: {
		filter.Farcaster.String(),
	},
	filter.NetworkPolygon: {
		filter.Aavegotchi.String(),
		filter.Lens.String(),
		filter.Matters.String(),
		filter.Aave.String(),
		filter.IQWiki.String(),
	},
	filter.NetworkCrossbell: {
		filter.Crossbell.String(),
	},
	filter.NetworkAvalanche: {
		filter.Aave.String(),
	},
	filter.NetworkBase: {
		filter.Aave.String(),
	},
	filter.NetworkOptimism: {
		filter.Aave.String(),
	},
	filter.NetworkArbitrum: {
		filter.Aave.String(),
	},
	filter.NetworkFantom: {
		filter.Aave.String(),
	},
}

var PlatformToWorkerMap = map[filter.Platform]string{
	filter.PlatformRSS3:       filter.RSS3.String(),
	filter.PlatformMirror:     filter.Mirror.String(),
	filter.PlatformFarcaster:  filter.Farcaster.String(),
	filter.PlatformParagraph:  filter.Paragraph.String(),
	filter.PlatformOpenSea:    filter.OpenSea.String(),
	filter.PlatformUniswap:    filter.Uniswap.String(),
	filter.PlatformOptimism:   filter.Optimism.String(),
	filter.PlatformAavegotchi: filter.Aavegotchi.String(),
	filter.PlatformLens:       filter.Lens.String(),
	filter.PlatformLooksRare:  filter.Looksrare.String(),
	filter.PlatformMatters:    filter.Matters.String(),
	filter.PlatformMomoka:     filter.Momoka.String(),
	filter.PlatformHighlight:  filter.Highlight.String(),
	filter.PlatformAAVE:       filter.Aave.String(),
	filter.PlatformIQWiki:     filter.IQWiki.String(),
	filter.PlatformLido:       filter.Lido.String(),
	filter.PlatformCrossbell:  filter.Crossbell.String(),
	filter.PlatformENS:        filter.ENS.String(),
}

var TagToWorkersMap = map[filter.Tag][]string{
	filter.TagTransaction: {
		filter.Optimism.String(),
	},
	filter.TagCollectible: {
		filter.OpenSea.String(),
		filter.ENS.String(),
		filter.Highlight.String(),
		filter.Lido.String(),
		filter.Looksrare.String(),
	},
	filter.TagExchange: {
		filter.RSS3.String(),
		filter.Uniswap.String(),
		filter.Aave.String(),
		filter.Lido.String(),
	},
	filter.TagSocial: {
		filter.Farcaster.String(),
		filter.Mirror.String(),
		filter.Lens.String(),
		filter.Paragraph.String(),
		filter.Crossbell.String(),
		filter.ENS.String(),
		filter.IQWiki.String(),
		filter.Matters.String(),
		filter.Momoka.String(),
	},
	filter.TagMetaverse: {
		filter.Aavegotchi.String(),
	},
}

func BuildPath(path string, query any, nodes []Cache) (map[common.Address]string, error) {
	if query != nil {
		values, err := form.NewEncoder().Encode(query)

		if err != nil {
			return nil, fmt.Errorf("build params %w", err)
		}

		path = fmt.Sprintf("%s?%s", path, values.Encode())
	}

	urls := make(map[common.Address]string, len(nodes))

	for _, node := range nodes {
		fullURL, err := url.JoinPath(node.Endpoint, path)
		if err != nil {
			return nil, fmt.Errorf("failed to join path for node %s: %w", node.Address, err)
		}

		decodedURL, err := url.QueryUnescape(fullURL)
		if err != nil {
			return nil, fmt.Errorf("failed to unescape url for node %s: %w", node.Address, err)
		}

		urls[common.HexToAddress(node.Address)] = decodedURL
	}

	return urls, nil
}

type DataResponse struct {
	Address common.Address
	Data    []byte
	// Valid indicates whether the data is non-null.
	Valid          bool
	Err            error
	Request        int
	InvalidRequest int
}

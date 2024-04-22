package model

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/rss3-network/protocol-go/schema/metadata"
)

var (
	RssNodeCacheKey  = "nodes:rss"
	FullNodeCacheKey = "nodes:full"

	// RequiredQualifiedNodeCount the required number of qualified Nodes
	RequiredQualifiedNodeCount = 3
	// RequiredVerificationCount the required number of verifications before a request is considered valid
	RequiredVerificationCount = 3
	// DemotionCountBeforeSlashing the number of demotions that trigger a slashing
	DemotionCountBeforeSlashing = 4

	// MutablePlatformMap is a map of mutable platforms which should be excluded from the data comparison.
	MutablePlatformMap = map[string]struct{}{
		filter.PlatformFarcaster.String(): {},
	}
)

// NodeEndpointCache represents a cache of a Node.
type NodeEndpointCache struct {
	Address  string `json:"address"`
	Endpoint string `json:"endpoint"`
}

// DataResponse represents the response returned by a Node.
// It is also used to store the verification result.
type DataResponse struct {
	Address  common.Address
	Endpoint string
	Data     []byte
	// A valid response must be non-null and non-error
	Valid bool
	Err   error
	// ValidPoint is the points given to the response
	ValidPoint int
	// InvalidPoint is the points given to the response when it is invalid
	InvalidPoint int
}

type ErrResponse struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

// ActivityResponse represents a single Activity in a response being returned to the requester.
type ActivityResponse struct {
	Data *Activity `json:"data"`
}

// ActivitiesResponse represents a list of Activity in a response being returned to the requester.
type ActivitiesResponse struct {
	Data []*Activity `json:"data"`
	Meta *MetaCursor `json:"meta,omitempty"`
}

type MetaCursor struct {
	Cursor string `json:"cursor"`
}

// Activity represents an activity.
type Activity struct {
	ID       string    `json:"id"`
	Owner    string    `json:"owner,omitempty"`
	Network  string    `json:"network"`
	Index    uint      `json:"index"`
	From     string    `json:"from"`
	To       string    `json:"to"`
	Tag      string    `json:"tag"`
	Type     string    `json:"type"`
	Platform string    `json:"platform,omitempty"`
	Actions  []*Action `json:"actions"`
}

// Action represents an action within an Activity.
type Action struct {
	Tag         string            `json:"tag"`
	Type        string            `json:"type"`
	Platform    string            `json:"platform,omitempty"`
	From        string            `json:"from"`
	To          string            `json:"to"`
	Metadata    metadata.Metadata `json:"metadata"`
	RelatedURLs []string          `json:"related_urls,omitempty"`
}

type Actions []*Action

var _ json.Unmarshaler = (*Action)(nil)

func (a *Action) UnmarshalJSON(bytes []byte) error {
	type ActionAlias Action

	type action struct {
		ActionAlias

		MetadataX json.RawMessage `json:"metadata"`
	}

	var temp action

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("unmarshal action: %w", err)
	}

	tag, err := filter.TagString(temp.Tag)
	if err != nil {
		return fmt.Errorf("invalid action tag: %w", err)
	}

	typeX, err := filter.TypeString(tag, temp.Type)
	if err != nil {
		return fmt.Errorf("invalid action type: %w", err)
	}

	temp.Metadata, err = metadata.Unmarshal(typeX, temp.MetadataX)
	if err != nil {
		return fmt.Errorf("invalid action metadata: %w", err)
	}

	*a = Action(temp.ActionAlias)

	return nil
}

// WorkerToNetworksMap Supplement the conditions for a full node based on the configuration file.
// https://github.com/RSS3-Network/Node/blob/develop/deploy/config.yaml
var WorkerToNetworksMap = map[filter.Name][]string{
	filter.Fallback: {
		filter.NetworkEthereum.String(),
		filter.NetworkVSL.String(),
		filter.NetworkSatoshiVM.String(),
		filter.Optimism.String(),
		filter.NetworkPolygon.String(),
		filter.NetworkArbitrum.String(),
		filter.NetworkBase.String(),
		filter.NetworkBinanceSmartChain.String(),
		filter.NetworkGnosis.String(),
		filter.NetworkLinea.String(),
	},
	filter.RSS3: {
		filter.NetworkEthereum.String(),
		filter.NetworkVSL.String(),
	},
	filter.Lens: {
		filter.NetworkPolygon.String(),
	},
	filter.OpenSea: {
		filter.NetworkEthereum.String(),
	},
	filter.Uniswap: {
		filter.NetworkEthereum.String(),
		filter.NetworkSatoshiVM.String(),
		filter.NetworkLinea.String(),
	},
	filter.Optimism: {
		filter.NetworkEthereum.String(),
	},
	filter.Aavegotchi: {
		filter.NetworkPolygon.String(),
	},
	filter.Highlight: {
		filter.NetworkEthereum.String(),
		filter.NetworkArbitrum.String(),
		filter.NetworkPolygon.String(),
		filter.NetworkOptimism.String(),
	},
	filter.Crossbell: {
		filter.NetworkCrossbell.String(),
	},
	filter.Farcaster: {
		filter.NetworkFarcaster.String(),
	},
	filter.Mirror: {
		filter.NetworkArweave.String(),
	},
	filter.Paragraph: {
		filter.NetworkArweave.String(),
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
	filter.Aave: {
		filter.NetworkPolygon.String(),
		filter.NetworkEthereum.String(),
		filter.NetworkAvalanche.String(),
		filter.NetworkBase.String(),
		filter.NetworkOptimism.String(),
		filter.NetworkArbitrum.String(),
	},
	filter.IQWiki: {
		filter.NetworkPolygon.String(),
	},
	filter.Lido: {
		filter.NetworkEthereum.String(),
	},
	filter.ENS: {
		filter.NetworkEthereum.String(),
	},
	filter.Oneinch: {
		filter.NetworkEthereum.String(),
	},
	filter.KiwiStand: {
		filter.NetworkOptimism.String(),
	},
	filter.SAVM: {
		filter.NetworkSatoshiVM.String(),
	},
	filter.VSL: {
		filter.NetworkVSL.String(),
	},
	filter.Stargate: {
		filter.NetworkEthereum.String(),
		filter.NetworkArbitrum.String(),
		filter.NetworkLinea.String(),
		filter.NetworkBinanceSmartChain.String(),
		filter.NetworkBase.String(),
		filter.NetworkOptimism.String(),
		filter.NetworkPolygon.String(),
		filter.NetworkAvalanche.String(),
	},
	filter.Curve: {
		filter.NetworkEthereum.String(),
		filter.NetworkArbitrum.String(),
		filter.NetworkAvalanche.String(),
		filter.NetworkGnosis.String(),
		filter.NetworkOptimism.String(),
		filter.NetworkPolygon.String(),
	},
}

// NetworkToWorkersMap is a map of network to workers.
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
		filter.Oneinch.String(),
		filter.Stargate.String(),
		filter.Curve.String(),
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
		filter.Fallback.String(),
		filter.Highlight.String(),
		filter.Stargate.String(),
		filter.Curve.String(),
	},
	filter.NetworkCrossbell: {
		filter.Crossbell.String(),
	},
	filter.NetworkAvalanche: {
		filter.Aave.String(),
		filter.Stargate.String(),
		filter.Curve.String(),
	},
	filter.NetworkBase: {
		filter.Aave.String(),
		filter.Fallback.String(),
		filter.Stargate.String(),
	},
	filter.NetworkOptimism: {
		filter.Aave.String(),
		filter.Fallback.String(),
		filter.Highlight.String(),
		filter.KiwiStand.String(),
		filter.Stargate.String(),
		filter.Curve.String(),
	},
	filter.NetworkArbitrum: {
		filter.Aave.String(),
		filter.Fallback.String(),
		filter.Highlight.String(),
		filter.Stargate.String(),
		filter.Curve.String(),
	},
	filter.NetworkVSL: {
		filter.Fallback.String(),
		filter.RSS3.String(),
		filter.VSL.String(),
	},
	filter.NetworkSatoshiVM: {
		filter.Fallback.String(),
		filter.Uniswap.String(),
		filter.SAVM.String(),
	},
	filter.NetworkBinanceSmartChain: {
		filter.Fallback.String(),
		filter.Stargate.String(),
	},
	filter.NetworkGnosis: {
		filter.Fallback.String(),
		filter.Curve.String(),
	},
	filter.NetworkLinea: {
		filter.Fallback.String(),
		filter.Uniswap.String(),
		filter.Stargate.String(),
	},
}

// PlatformToWorkerMap is a map of platform to worker.
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
	filter.Platform1inch:      filter.Oneinch.String(),
	filter.PlatformKiwiStand:  filter.KiwiStand.String(),
	filter.PlatformSAVM:       filter.SAVM.String(),
	filter.PlatformVSL:        filter.VSL.String(),
	filter.PlatformStargate:   filter.Stargate.String(),
	filter.PlatformCurve:      filter.Curve.String(),
}

// TagToWorkersMap is a map of tag to workers.
var TagToWorkersMap = map[filter.Tag][]string{
	filter.TagTransaction: {
		filter.Optimism.String(),
		filter.SAVM.String(),
		filter.VSL.String(),
		filter.Stargate.String(),
	},
	filter.TagCollectible: {
		filter.OpenSea.String(),
		filter.ENS.String(),
		filter.Highlight.String(),
		filter.Lido.String(),
		filter.Looksrare.String(),
		filter.KiwiStand.String(),
	},
	filter.TagExchange: {
		filter.RSS3.String(),
		filter.Uniswap.String(),
		filter.Aave.String(),
		filter.Lido.String(),
		filter.Oneinch.String(),
		filter.Curve.String(),
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

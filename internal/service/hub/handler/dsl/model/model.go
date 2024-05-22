package model

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/node/schema/worker"
	"github.com/rss3-network/protocol-go/schema"
	"github.com/rss3-network/protocol-go/schema/metadata"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
)

var (
	RssNodeCacheKey  = "nodes:rss"
	FullNodeCacheKey = "nodes:full"

	SubscribeNodeCacheKey = "epoch"

	// RequiredQualifiedNodeCount the required number of qualified Nodes
	RequiredQualifiedNodeCount = 3
	// RequiredVerificationCount the required number of verifications before a request is considered valid
	RequiredVerificationCount = 3
	// DemotionCountBeforeSlashing the number of demotions that trigger a slashing
	DemotionCountBeforeSlashing = 4

	// MutablePlatformMap is a map of mutable platforms which should be excluded from the data comparison.
	MutablePlatformMap = map[string]struct{}{
		worker.PlatformFarcaster.String(): {},
	}
)

// NodeEndpointCache stores the elements in the heap.
type NodeEndpointCache struct {
	Address      string `json:"address"`
	Endpoint     string `json:"endpoint"`
	Score        float64
	InvalidCount int64
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

	tagX, err := tag.TagString(temp.Tag)
	if err != nil {
		return fmt.Errorf("invalid action tag: %w", err)
	}

	typeX, err := schema.ParseTypeFromString(tagX, temp.Type)
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
var WorkerToNetworksMap = map[worker.Worker][]string{
	worker.Aave: {
		network.Arbitrum.String(),
		network.Avalanche.String(),
		network.Base.String(),
		network.Ethereum.String(),
		network.Optimism.String(),
		network.Polygon.String(),
	},
	worker.Aavegotchi: {
		network.Polygon.String(),
	},
	worker.Core: {
		network.Arbitrum.String(),
		network.Base.String(),
		network.BinanceSmartChain.String(),
		network.Ethereum.String(),
		network.Farcaster.String(),
		network.Gnosis.String(),
		network.Linea.String(),
		network.Optimism.String(),
		network.Polygon.String(),
		network.SatoshiVM.String(),
		network.VSL.String(),
	},
	worker.Crossbell: {
		network.Crossbell.String(),
	},
	worker.Curve: {
		network.Arbitrum.String(),
		network.Avalanche.String(),
		network.Ethereum.String(),
		network.Gnosis.String(),
		network.Optimism.String(),
		network.Polygon.String(),
	},
	worker.ENS: {
		network.Ethereum.String(),
	},
	worker.Highlight: {
		network.Arbitrum.String(),
		network.Ethereum.String(),
		network.Optimism.String(),
		network.Polygon.String(),
	},
	worker.IQWiki: {
		network.Polygon.String(),
	},
	worker.KiwiStand: {
		network.Optimism.String(),
	},
	worker.Lens: {
		network.Polygon.String(),
	},
	worker.Lido: {
		network.Ethereum.String(),
	},
	worker.Looksrare: {
		network.Ethereum.String(),
	},
	worker.Matters: {
		network.Optimism.String(),
	},
	worker.Mirror: {
		network.Arweave.String(),
	},
	worker.Momoka: {
		network.Arweave.String(),
	},
	worker.Oneinch: {
		network.Ethereum.String(),
	},
	worker.OpenSea: {
		network.Ethereum.String(),
	},
	worker.Optimism: {
		network.Ethereum.String(),
	},
	worker.Paragraph: {
		network.Arweave.String(),
	},
	worker.RSS3: {
		network.Ethereum.String(),
		network.VSL.String(),
	},
	worker.SAVM: {
		network.SatoshiVM.String(),
	},
	worker.Stargate: {
		network.Arbitrum.String(),
		network.Avalanche.String(),
		network.Base.String(),
		network.BinanceSmartChain.String(),
		network.Ethereum.String(),
		network.Linea.String(),
		network.Optimism.String(),
		network.Polygon.String(),
	},
	worker.Uniswap: {
		network.Ethereum.String(),
		network.Linea.String(),
		network.SatoshiVM.String(),
	},
	worker.VSL: {
		network.VSL.String(),
	},
}

// NetworkToWorkersMap is a map of network to workers.
var NetworkToWorkersMap = map[network.Network][]string{
	network.Arbitrum: {
		worker.Aave.String(),
		worker.Core.String(),
		worker.Curve.String(),
		worker.Highlight.String(),
		worker.Stargate.String(),
	},
	network.Arweave: {
		worker.Mirror.String(),
		worker.Momoka.String(),
		worker.Paragraph.String(),
	},
	network.Avalanche: {
		worker.Aave.String(),
		worker.Curve.String(),
		worker.Stargate.String(),
	},
	network.Base: {
		worker.Aave.String(),
		worker.Core.String(),
		worker.Stargate.String(),
	},
	network.BinanceSmartChain: {
		worker.Core.String(),
		worker.Stargate.String(),
	},
	network.Crossbell: {
		worker.Crossbell.String(),
	},
	network.Ethereum: {
		worker.Aave.String(),
		worker.Core.String(),
		worker.Curve.String(),
		worker.ENS.String(),
		worker.Highlight.String(),
		worker.Lido.String(),
		worker.Looksrare.String(),
		worker.Oneinch.String(),
		worker.OpenSea.String(),
		worker.Optimism.String(),
		worker.RSS3.String(),
		worker.Stargate.String(),
		worker.Uniswap.String(),
	},
	network.Farcaster: {
		worker.Core.String(),
	},
	network.Gnosis: {
		worker.Core.String(),
		worker.Curve.String(),
	},
	network.Linea: {
		worker.Core.String(),
		worker.Uniswap.String(),
		worker.Stargate.String(),
	},
	network.Optimism: {
		worker.Aave.String(),
		worker.Core.String(),
		worker.Curve.String(),
		worker.Highlight.String(),
		worker.KiwiStand.String(),
		worker.Matters.String(),
		worker.Stargate.String(),
	},
	network.Polygon: {
		worker.Aave.String(),
		worker.Aavegotchi.String(),
		worker.Core.String(),
		worker.Curve.String(),
		worker.Highlight.String(),
		worker.Lens.String(),
		worker.IQWiki.String(),
		worker.Stargate.String(),
	},
	network.SatoshiVM: {
		worker.Core.String(),
		worker.Uniswap.String(),
		worker.SAVM.String(),
	},
	network.VSL: {
		worker.Core.String(),
		worker.RSS3.String(),
		worker.VSL.String(),
	},
}

// PlatformToWorkersMap is a map of platform to workers.
var PlatformToWorkersMap = map[worker.Platform][]string{
	worker.Platform1Inch:      {worker.Oneinch.String()},
	worker.PlatformAAVE:       {worker.Aave.String()},
	worker.PlatformAavegotchi: {worker.Aavegotchi.String()},
	worker.PlatformCrossbell:  {worker.Crossbell.String()},
	worker.PlatformCurve:      {worker.Curve.String()},
	worker.PlatformENS:        {worker.ENS.String()},
	worker.PlatformFarcaster:  {worker.Core.String()},
	worker.PlatformHighlight:  {worker.Highlight.String()},
	worker.PlatformIQWiki:     {worker.IQWiki.String()},
	worker.PlatformKiwiStand:  {worker.KiwiStand.String()},
	worker.PlatformLens:       {worker.Lens.String(), worker.Momoka.String()},
	worker.PlatformLido:       {worker.Lido.String()},
	worker.PlatformLooksRare:  {worker.Looksrare.String()},
	worker.PlatformMatters:    {worker.Matters.String()},
	worker.PlatformMirror:     {worker.Mirror.String()},
	worker.PlatformOpenSea:    {worker.OpenSea.String()},
	worker.PlatformOptimism:   {worker.Optimism.String()},
	worker.PlatformParagraph:  {worker.Paragraph.String()},
	worker.PlatformRSS3:       {worker.RSS3.String()},
	worker.PlatformSAVM:       {worker.SAVM.String()},
	worker.PlatformStargate:   {worker.Stargate.String()},
	worker.PlatformUniswap:    {worker.Uniswap.String()},
	worker.PlatformVSL:        {worker.VSL.String()},
}

// TagToWorkersMap is a map of tag to workers.
var TagToWorkersMap = map[tag.Tag][]string{
	tag.Collectible: {
		worker.ENS.String(),
		worker.Highlight.String(),
		worker.KiwiStand.String(),
		worker.Lido.String(),
		worker.Looksrare.String(),
		worker.OpenSea.String(),
	},
	tag.Exchange: {
		worker.Aave.String(),
		worker.Curve.String(),
		worker.Lido.String(),
		worker.Oneinch.String(),
		worker.RSS3.String(),
		worker.Uniswap.String(),
	},
	tag.Metaverse: {
		worker.Aavegotchi.String(),
	},
	tag.Social: {
		worker.Core.String(), // farcaster core
		worker.Crossbell.String(),
		worker.ENS.String(),
		worker.IQWiki.String(),
		worker.Lens.String(),
		worker.Matters.String(),
		worker.Mirror.String(),
		worker.Momoka.String(),
		worker.Paragraph.String(),
	},
	tag.Transaction: {
		worker.Optimism.String(),
		worker.SAVM.String(),
		worker.Stargate.String(),
		worker.VSL.String(),
	},
}

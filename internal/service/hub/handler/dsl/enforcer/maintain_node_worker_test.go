package enforcer

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/rss3-network/node/schema/worker"
	"github.com/rss3-network/node/schema/worker/decentralized"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	workerStatusNode1 = `{"data":[{"network":"optimism","worker":"kiwistand","tags":["collectible","transaction","social"],"platform":"KiwiStand","status":"Ready"},{"network":"ethereum","worker":"uniswap","tags":["exchange","transaction"],"platform":"Uniswap","status":"Ready"},{"network":"savm","worker":"uniswap","tags":["exchange","transaction"],"platform":"Uniswap","status":"Ready"},{"network":"crossbell","worker":"core","tags":null,"platform":"Unknown","status":"Indexing"},{"network":"gnosis","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready"},{"network":"base","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready"},{"network":"optimism","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready"},{"network":"arweave","worker":"momoka","tags":["social"],"platform":"Lens","status":"Unhealthy"},{"network":"arbitrum","worker":"core","tags":null,"platform":"Unknown","status":"Indexing"},{"network":"ethereum","worker":"looksrare","tags":["collectible"],"platform":"LooksRare","status":"Ready"},{"network":"ethereum","worker":"highlight","tags":["collectible","transaction"],"platform":"Highlight","status":"Ready"},{"network":"farcaster","worker":"core","tags":null,"platform":"Unknown","status":"Indexing"},{"network":"ethereum","worker":"1inch","tags":["exchange"],"platform":"1inch","status":"Ready"},{"network":"optimism","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready"},{"network":"arbitrum","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Indexing"},{"network":"binance-smart-chain","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready"},{"network":"vsl","worker":"rss3","tags":["exchange","collectible"],"platform":"RSS3","status":"Ready"},{"network":"polygon","worker":"aavegotchi","tags":["metaverse"],"platform":"Aavegotchi","status":"Ready"},{"network":"polygon","worker":"highlight","tags":["collectible","transaction"],"platform":"Highlight","status":"Ready"},{"network":"savm","worker":"core","tags":null,"platform":"Unknown","status":"Ready"},{"network":"savm","worker":"savm","tags":["transaction"],"platform":"SAVM","status":"Ready"},{"network":"ethereum","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready"},{"network":"linea","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Unhealthy"},{"network":"arbitrum","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready"},{"network":"arweave","worker":"mirror","tags":["social"],"platform":"Mirror","status":"Unhealthy"},{"network":"gnosis","worker":"core","tags":null,"platform":"Unknown","status":"Indexing"},{"network":"ethereum","worker":"vsl","tags":["transaction"],"platform":"VSL","status":"Ready"},{"network":"ethereum","worker":"core","tags":null,"platform":"Unknown","status":"Ready"},{"network":"arweave","worker":"paragraph","tags":["social"],"platform":"Paragraph","status":"Unhealthy"},{"network":"linea","worker":"core","tags":null,"platform":"Unknown","status":"Ready"},{"network":"optimism","worker":"matters","tags":["social"],"platform":"Matters","status":"Indexing"},{"network":"ethereum","worker":"opensea","tags":["collectible"],"platform":"OpenSea","status":"Ready"},{"network":"base","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready"},{"network":"polygon","worker":"iqwiki","tags":["social"],"platform":"IQWiki","status":"Ready"},{"network":"ethereum","worker":"lido","tags":["exchange","transaction","collectible"],"platform":"Lido","status":"Ready"},{"network":"arbitrum","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready"},{"network":"polygon","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready"},{"network":"optimism","worker":"highlight","tags":["collectible","transaction"],"platform":"Highlight","status":"Ready"},{"network":"polygon","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready"},{"network":"optimism","worker":"core","tags":null,"platform":"Unknown","status":"Ready"},{"network":"ethereum","worker":"optimism","tags":["transaction"],"platform":"Optimism","status":"Ready"},{"network":"crossbell","worker":"crossbell","tags":["social"],"platform":"Crossbell","status":"Indexing"},{"network":"arbitrum","worker":"highlight","tags":["collectible","transaction"],"platform":"Highlight","status":"Indexing"},{"network":"polygon","worker":"core","tags":null,"platform":"Unknown","status":"Ready"},{"network":"polygon","worker":"lens","tags":["social"],"platform":"Lens","status":"Unhealthy"},{"network":"avax","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready"},{"network":"linea","worker":"uniswap","tags":["exchange","transaction"],"platform":"Uniswap","status":"Ready"},{"network":"avax","worker":"core","tags":null,"platform":"Unknown","status":"Indexing"},{"network":"polygon","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready"},{"network":"ethereum","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready"},{"network":"base","worker":"core","tags":null,"platform":"Unknown","status":"Indexing"},{"network":"avax","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready"},{"network":"vsl","worker":"core","tags":null,"platform":"Unknown","status":"Ready"},{"network":"ethereum","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready"},{"network":"avax","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready"},{"network":"ethereum","worker":"rss3","tags":["exchange","collectible"],"platform":"RSS3","status":"Ready"},{"network":"binance-smart-chain","worker":"core","tags":null,"platform":"Unknown","status":"Ready"},{"network":"optimism","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready"}]}`
	workerStatusNode2 = `{"data":[{"network":"farcaster","worker":"core","tags":null,"platform":"Unknown","status":"Ready"}]}`
	workerStatusNode3 = `{"data":[{"network":"polygon","worker":"lens","tags":["social"],"platform":"Lens","status":"Ready"},{"network":"arweave","worker":"momoka","tags":["social"],"platform":"Lens","status":"Ready"}]}`

	workerInfoNode1 = []*WorkerInfo{
		{Network: network.Optimism, Worker: decentralized.KiwiStand, Tags: []tag.Tag{tag.Collectible, tag.Transaction, tag.Social}, Platform: decentralized.PlatformKiwiStand, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: decentralized.Uniswap, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: decentralized.PlatformUniswap, Status: worker.StatusReady},
		{Network: network.SatoshiVM, Worker: decentralized.Uniswap, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: decentralized.PlatformUniswap, Status: worker.StatusReady},
		{Network: network.Crossbell, Worker: decentralized.Core, Tags: nil, Platform: decentralized.PlatformUnknown, Status: worker.StatusIndexing},
		{Network: network.Gnosis, Worker: decentralized.Curve, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: decentralized.PlatformCurve, Status: worker.StatusReady},
		{Network: network.Base, Worker: decentralized.Aave, Tags: []tag.Tag{tag.Exchange}, Platform: decentralized.PlatformAAVE, Status: worker.StatusReady},
		{Network: network.Optimism, Worker: decentralized.Aave, Tags: []tag.Tag{tag.Exchange}, Platform: decentralized.PlatformAAVE, Status: worker.StatusReady},
		{Network: network.Arweave, Worker: decentralized.Momoka, Tags: []tag.Tag{tag.Social}, Platform: decentralized.PlatformLens, Status: worker.StatusUnhealthy},
		{Network: network.Arbitrum, Worker: decentralized.Core, Tags: nil, Platform: decentralized.PlatformUnknown, Status: worker.StatusIndexing},
		{Network: network.Ethereum, Worker: decentralized.Looksrare, Tags: []tag.Tag{tag.Collectible}, Platform: decentralized.PlatformLooksRare, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: decentralized.Highlight, Tags: []tag.Tag{tag.Collectible, tag.Transaction}, Platform: decentralized.PlatformHighlight, Status: worker.StatusReady},
		{Network: network.Farcaster, Worker: decentralized.Core, Tags: []tag.Tag{tag.Social}, Platform: decentralized.PlatformFarcaster, Status: worker.StatusIndexing},
		{Network: network.Ethereum, Worker: decentralized.Oneinch, Tags: []tag.Tag{tag.Exchange}, Platform: decentralized.Platform1Inch, Status: worker.StatusReady},
		{Network: network.Optimism, Worker: decentralized.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: decentralized.PlatformStargate, Status: worker.StatusReady},
		{Network: network.Arbitrum, Worker: decentralized.Aave, Tags: []tag.Tag{tag.Exchange}, Platform: decentralized.PlatformAAVE, Status: worker.StatusIndexing},
		{Network: network.BinanceSmartChain, Worker: decentralized.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: decentralized.PlatformStargate, Status: worker.StatusReady},
		{Network: network.VSL, Worker: decentralized.RSS3, Tags: []tag.Tag{tag.Exchange, tag.Collectible}, Platform: decentralized.PlatformRSS3, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: decentralized.Aavegotchi, Tags: []tag.Tag{tag.Metaverse}, Platform: decentralized.PlatformAavegotchi, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: decentralized.Highlight, Tags: []tag.Tag{tag.Collectible, tag.Transaction}, Platform: decentralized.PlatformHighlight, Status: worker.StatusReady},
		{Network: network.SatoshiVM, Worker: decentralized.Core, Tags: nil, Platform: decentralized.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.SatoshiVM, Worker: decentralized.SAVM, Tags: []tag.Tag{tag.Transaction}, Platform: decentralized.PlatformSAVM, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: decentralized.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: decentralized.PlatformStargate, Status: worker.StatusReady},
		{Network: network.Linea, Worker: decentralized.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: decentralized.PlatformStargate, Status: worker.StatusUnhealthy},
		{Network: network.Arbitrum, Worker: decentralized.Curve, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: decentralized.PlatformCurve, Status: worker.StatusReady},
		{Network: network.Arweave, Worker: decentralized.Mirror, Tags: []tag.Tag{tag.Social}, Platform: decentralized.PlatformMirror, Status: worker.StatusUnhealthy},
		{Network: network.Gnosis, Worker: decentralized.Core, Tags: nil, Platform: decentralized.PlatformUnknown, Status: worker.StatusIndexing},
		{Network: network.Ethereum, Worker: decentralized.VSL, Tags: []tag.Tag{tag.Transaction}, Platform: decentralized.PlatformVSL, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: decentralized.Core, Tags: nil, Platform: decentralized.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.Arweave, Worker: decentralized.Paragraph, Tags: []tag.Tag{tag.Social}, Platform: decentralized.PlatformParagraph, Status: worker.StatusUnhealthy},
		{Network: network.Linea, Worker: decentralized.Core, Tags: nil, Platform: decentralized.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.Optimism, Worker: decentralized.Matters, Tags: []tag.Tag{tag.Social}, Platform: decentralized.PlatformMatters, Status: worker.StatusIndexing},
		{Network: network.Ethereum, Worker: decentralized.OpenSea, Tags: []tag.Tag{tag.Collectible}, Platform: decentralized.PlatformOpenSea, Status: worker.StatusReady},
		{Network: network.Base, Worker: decentralized.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: decentralized.PlatformStargate, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: decentralized.IQWiki, Tags: []tag.Tag{tag.Social}, Platform: decentralized.PlatformIQWiki, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: decentralized.Lido, Tags: []tag.Tag{tag.Exchange, tag.Transaction, tag.Collectible}, Platform: decentralized.PlatformLido, Status: worker.StatusReady},
		{Network: network.Arbitrum, Worker: decentralized.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: decentralized.PlatformStargate, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: decentralized.Aave, Tags: []tag.Tag{tag.Exchange}, Platform: decentralized.PlatformAAVE, Status: worker.StatusReady},
		{Network: network.Optimism, Worker: decentralized.Highlight, Tags: []tag.Tag{tag.Collectible, tag.Transaction}, Platform: decentralized.PlatformHighlight, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: decentralized.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: decentralized.PlatformStargate, Status: worker.StatusReady},
		{Network: network.Optimism, Worker: decentralized.Core, Tags: nil, Platform: decentralized.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: decentralized.Optimism, Tags: []tag.Tag{tag.Transaction}, Platform: decentralized.PlatformOptimism, Status: worker.StatusReady},
		{Network: network.Crossbell, Worker: decentralized.Crossbell, Tags: []tag.Tag{tag.Social}, Platform: decentralized.PlatformCrossbell, Status: worker.StatusIndexing},
		{Network: network.Arbitrum, Worker: decentralized.Highlight, Tags: []tag.Tag{tag.Collectible, tag.Transaction}, Platform: decentralized.PlatformHighlight, Status: worker.StatusIndexing},
		{Network: network.Polygon, Worker: decentralized.Core, Tags: nil, Platform: decentralized.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: decentralized.Lens, Tags: []tag.Tag{tag.Social}, Platform: decentralized.PlatformLens, Status: worker.StatusUnhealthy},
		{Network: network.Avalanche, Worker: decentralized.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: decentralized.PlatformStargate, Status: worker.StatusReady},
		{Network: network.Linea, Worker: decentralized.Uniswap, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: decentralized.PlatformUniswap, Status: worker.StatusReady},
		{Network: network.Avalanche, Worker: decentralized.Core, Tags: nil, Platform: decentralized.PlatformUnknown, Status: worker.StatusIndexing},
		{Network: network.Polygon, Worker: decentralized.Curve, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: decentralized.PlatformCurve, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: decentralized.Curve, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: decentralized.PlatformCurve, Status: worker.StatusReady},
		{Network: network.Base, Worker: decentralized.Core, Tags: nil, Platform: decentralized.PlatformUnknown, Status: worker.StatusIndexing},
		{Network: network.Avalanche, Worker: decentralized.Curve, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: decentralized.PlatformCurve, Status: worker.StatusReady},
		{Network: network.VSL, Worker: decentralized.Core, Tags: nil, Platform: decentralized.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: decentralized.Aave, Tags: []tag.Tag{tag.Exchange}, Platform: decentralized.PlatformAAVE, Status: worker.StatusReady},
		{Network: network.Avalanche, Worker: decentralized.Aave, Tags: []tag.Tag{tag.Exchange}, Platform: decentralized.PlatformAAVE, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: decentralized.RSS3, Tags: []tag.Tag{tag.Exchange, tag.Collectible}, Platform: decentralized.PlatformRSS3, Status: worker.StatusReady},
		{Network: network.BinanceSmartChain, Worker: decentralized.Core, Tags: nil, Platform: decentralized.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.Optimism, Worker: decentralized.Curve, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: decentralized.PlatformCurve, Status: worker.StatusReady},
	}
	workerInfoNode2 = []*WorkerInfo{
		{Network: network.Farcaster, Worker: decentralized.Core, Tags: []tag.Tag{tag.Social}, Platform: decentralized.PlatformFarcaster, Status: worker.StatusReady},
	}
	workerInfoNode3 = []*WorkerInfo{
		{Network: network.Arweave, Worker: decentralized.Momoka, Tags: []tag.Tag{tag.Social}, Platform: decentralized.PlatformLens, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: decentralized.Lens, Tags: []tag.Tag{tag.Social}, Platform: decentralized.PlatformLens, Status: worker.StatusReady},
	}
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Fetch(ctx context.Context, endpoint string) (io.ReadCloser, error) {
	args := m.Called(ctx, endpoint)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func TestGetNodeWorkerStatus(t *testing.T) {
	t.Parallel()

	mockClient := new(MockHTTPClient)
	mockClient.On("Fetch", mock.Anything, "http://localhost:8080/worker_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNode1))), nil)

	enforcer := &SimpleEnforcer{httpClient: mockClient}
	response, err := enforcer.getNodeWorkerStatus(context.Background(), "http://localhost:8080")

	assert.NoError(t, err)
	assert.Equal(t, len(workerInfoNode1), len(response.Data))
}

func TestGenerateMaps(t *testing.T) {
	t.Parallel()

	mockClient := new(MockHTTPClient)
	mockClient.On("Fetch", mock.Anything, "http://localhost:8080/worker_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNode1))), nil)
	mockClient.On("Fetch", mock.Anything, "http://localhost:8081/worker_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNode2))), nil)
	mockClient.On("Fetch", mock.Anything, "http://localhost:8082/worker_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNode3))), nil)

	enforcer := &SimpleEnforcer{httpClient: mockClient}
	stats := []*schema.Stat{
		{
			Address:  common.Address{1},
			Endpoint: "http://localhost:8080",
		},
		{
			Address:  common.Address{2},
			Endpoint: "http://localhost:8081",
		},
		{
			Address:  common.Address{3},
			Endpoint: "http://localhost:8082",
		},
	}

	nodeToWorkersMap, fullNodeWorkerToNetworksMap, networkToWorkersMap, platformToWorkersMap, tagToWorkersMap := enforcer.generateMaps(context.Background(), stats)

	expectedNodeToWorkersMap := map[common.Address][]*WorkerInfo{
		common.Address{1}: workerInfoNode1,
		common.Address{2}: workerInfoNode2,
		common.Address{3}: workerInfoNode3,
	}
	expectedFullNodeWorkerToNetworksMap := map[decentralized.Worker]map[network.Network]struct{}{
		decentralized.Aave: {
			network.Avalanche: {},
			network.Base:      {},
			network.Ethereum:  {},
			network.Optimism:  {},
			network.Polygon:   {},
		},
		decentralized.Aavegotchi: {
			network.Polygon: {},
		},
		decentralized.Core: {
			network.BinanceSmartChain: {},
			network.Ethereum:          {},
			network.Linea:             {},
			network.Optimism:          {},
			network.Polygon:           {},
			network.SatoshiVM:         {},
			network.VSL:               {},
			network.Farcaster:         {},
		},
		decentralized.Curve: {
			network.Arbitrum:  {},
			network.Avalanche: {},
			network.Ethereum:  {},
			network.Gnosis:    {},
			network.Optimism:  {},
			network.Polygon:   {},
		},
		decentralized.Highlight: {
			network.Ethereum: {},
			network.Optimism: {},
			network.Polygon:  {},
		},
		decentralized.IQWiki: {
			network.Polygon: {},
		},
		decentralized.KiwiStand: {
			network.Optimism: {},
		},
		decentralized.Lido: {
			network.Ethereum: {},
		},
		decentralized.Looksrare: {
			network.Ethereum: {},
		},
		decentralized.Oneinch: {
			network.Ethereum: {},
		},
		decentralized.OpenSea: {
			network.Ethereum: {},
		},
		decentralized.Optimism: {
			network.Ethereum: {},
		},
		decentralized.RSS3: {
			network.Ethereum: {},
			network.VSL:      {},
		},
		decentralized.SAVM: {
			network.SatoshiVM: {},
		},
		decentralized.Stargate: {
			network.Arbitrum:          {},
			network.Avalanche:         {},
			network.Base:              {},
			network.BinanceSmartChain: {},
			network.Ethereum:          {},
			network.Optimism:          {},
			network.Polygon:           {},
		},
		decentralized.Uniswap: {
			network.Ethereum:  {},
			network.Linea:     {},
			network.SatoshiVM: {},
		},
		decentralized.VSL: {
			network.Ethereum: {},
		},
		decentralized.Momoka: {
			network.Arweave: {},
		},
		decentralized.Lens: {
			network.Polygon: {},
		},
	}
	expectedNetworkToWorkersMap := map[network.Network]map[decentralized.Worker]struct{}{
		network.Arbitrum: {
			decentralized.Curve:    {},
			decentralized.Stargate: {},
		},
		network.Avalanche: {
			decentralized.Aave:     {},
			decentralized.Curve:    {},
			decentralized.Stargate: {},
		},
		network.Arweave: {
			decentralized.Momoka: {},
		},
		network.Base: {
			decentralized.Aave:     {},
			decentralized.Stargate: {},
		},
		network.BinanceSmartChain: {
			decentralized.Core:     {},
			decentralized.Stargate: {},
		},
		network.Ethereum: {
			decentralized.Aave:      {},
			decentralized.Core:      {},
			decentralized.Curve:     {},
			decentralized.Highlight: {},
			decentralized.Lido:      {},
			decentralized.Looksrare: {},
			decentralized.Oneinch:   {},
			decentralized.OpenSea:   {},
			decentralized.Optimism:  {},
			decentralized.RSS3:      {},
			decentralized.Stargate:  {},
			decentralized.Uniswap:   {},
			decentralized.VSL:       {},
		},
		network.Gnosis: {
			decentralized.Curve: {},
		},
		network.Linea: {
			decentralized.Core:    {},
			decentralized.Uniswap: {},
		},
		network.Optimism: {
			decentralized.Aave:      {},
			decentralized.Core:      {},
			decentralized.Curve:     {},
			decentralized.Highlight: {},
			decentralized.KiwiStand: {},
			decentralized.Stargate:  {},
		},
		network.Polygon: {
			decentralized.Aave:       {},
			decentralized.Aavegotchi: {},
			decentralized.Core:       {},
			decentralized.Curve:      {},
			decentralized.Highlight:  {},
			decentralized.IQWiki:     {},
			decentralized.Stargate:   {},
			decentralized.Lens:       {},
		},
		network.SatoshiVM: {
			decentralized.Core:    {},
			decentralized.Uniswap: {},
			decentralized.SAVM:    {},
		},
		network.VSL: {
			decentralized.Core: {},
			decentralized.RSS3: {},
		},

		network.Farcaster: {
			decentralized.Core: {},
		},
	}
	expectedPlatformToWorkersMap := map[decentralized.Platform]map[decentralized.Worker]struct{}{
		decentralized.Platform1Inch:      {decentralized.Oneinch: {}},
		decentralized.PlatformAAVE:       {decentralized.Aave: {}},
		decentralized.PlatformAavegotchi: {decentralized.Aavegotchi: {}},
		decentralized.PlatformCurve:      {decentralized.Curve: {}},
		decentralized.PlatformHighlight:  {decentralized.Highlight: {}},
		decentralized.PlatformIQWiki:     {decentralized.IQWiki: {}},
		decentralized.PlatformKiwiStand:  {decentralized.KiwiStand: {}},
		decentralized.PlatformLido:       {decentralized.Lido: {}},
		decentralized.PlatformLooksRare:  {decentralized.Looksrare: {}},
		decentralized.PlatformOpenSea:    {decentralized.OpenSea: {}},
		decentralized.PlatformOptimism:   {decentralized.Optimism: {}},
		decentralized.PlatformRSS3:       {decentralized.RSS3: {}},
		decentralized.PlatformSAVM:       {decentralized.SAVM: {}},
		decentralized.PlatformStargate:   {decentralized.Stargate: {}},
		decentralized.PlatformUniswap:    {decentralized.Uniswap: {}},
		decentralized.PlatformVSL:        {decentralized.VSL: {}},
		decentralized.PlatformFarcaster:  {decentralized.Core: {}},
		decentralized.PlatformLens:       {decentralized.Lens: {}, decentralized.Momoka: {}},
	}
	expectedTagToWorkersMap := map[tag.Tag]map[decentralized.Worker]struct{}{
		tag.Collectible: {
			decentralized.Highlight: {},
			decentralized.KiwiStand: {},
			decentralized.Lido:      {},
			decentralized.Looksrare: {},
			decentralized.OpenSea:   {},
			decentralized.RSS3:      {},
		},
		tag.Exchange: {
			decentralized.Aave:    {},
			decentralized.Curve:   {},
			decentralized.Lido:    {},
			decentralized.Oneinch: {},
			decentralized.RSS3:    {},
			decentralized.Uniswap: {},
		},
		tag.Metaverse: {
			decentralized.Aavegotchi: {},
		},
		tag.Social: {
			decentralized.KiwiStand: {},
			decentralized.IQWiki:    {},
			decentralized.Core:      {},
			decentralized.Momoka:    {},
			decentralized.Lens:      {},
		},
		tag.Transaction: {
			decentralized.KiwiStand: {},
			decentralized.Uniswap:   {},
			decentralized.Curve:     {},
			decentralized.Optimism:  {},
			decentralized.Highlight: {},
			decentralized.Lido:      {},
			decentralized.SAVM:      {},
			decentralized.Stargate:  {},
			decentralized.VSL:       {},
		},
	}

	assert.Equal(t, expectedNodeToWorkersMap[common.Address{1}], nodeToWorkersMap[common.Address{1}])
	assert.Equal(t, expectedFullNodeWorkerToNetworksMap, fullNodeWorkerToNetworksMap)
	assert.Equal(t, expectedNetworkToWorkersMap, networkToWorkersMap)
	assert.Equal(t, expectedPlatformToWorkersMap, platformToWorkersMap)
	assert.Equal(t, expectedTagToWorkersMap, tagToWorkersMap)
}

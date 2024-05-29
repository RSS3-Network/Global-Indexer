package enforcer

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/rss3-network/node/schema/worker"
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
		{Network: network.Optimism, Worker: worker.KiwiStand, Tags: []tag.Tag{tag.Collectible, tag.Transaction, tag.Social}, Platform: worker.PlatformKiwiStand, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: worker.Uniswap, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: worker.PlatformUniswap, Status: worker.StatusReady},
		{Network: network.SatoshiVM, Worker: worker.Uniswap, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: worker.PlatformUniswap, Status: worker.StatusReady},
		{Network: network.Crossbell, Worker: worker.Core, Tags: nil, Platform: worker.PlatformUnknown, Status: worker.StatusIndexing},
		{Network: network.Gnosis, Worker: worker.Curve, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: worker.PlatformCurve, Status: worker.StatusReady},
		{Network: network.Base, Worker: worker.Aave, Tags: []tag.Tag{tag.Exchange}, Platform: worker.PlatformAAVE, Status: worker.StatusReady},
		{Network: network.Optimism, Worker: worker.Aave, Tags: []tag.Tag{tag.Exchange}, Platform: worker.PlatformAAVE, Status: worker.StatusReady},
		{Network: network.Arweave, Worker: worker.Momoka, Tags: []tag.Tag{tag.Social}, Platform: worker.PlatformLens, Status: worker.StatusUnhealthy},
		{Network: network.Arbitrum, Worker: worker.Core, Tags: nil, Platform: worker.PlatformUnknown, Status: worker.StatusIndexing},
		{Network: network.Ethereum, Worker: worker.Looksrare, Tags: []tag.Tag{tag.Collectible}, Platform: worker.PlatformLooksRare, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: worker.Highlight, Tags: []tag.Tag{tag.Collectible, tag.Transaction}, Platform: worker.PlatformHighlight, Status: worker.StatusReady},
		{Network: network.Farcaster, Worker: worker.Core, Tags: []tag.Tag{tag.Social}, Platform: worker.PlatformFarcaster, Status: worker.StatusIndexing},
		{Network: network.Ethereum, Worker: worker.Oneinch, Tags: []tag.Tag{tag.Exchange}, Platform: worker.Platform1Inch, Status: worker.StatusReady},
		{Network: network.Optimism, Worker: worker.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: worker.PlatformStargate, Status: worker.StatusReady},
		{Network: network.Arbitrum, Worker: worker.Aave, Tags: []tag.Tag{tag.Exchange}, Platform: worker.PlatformAAVE, Status: worker.StatusIndexing},
		{Network: network.BinanceSmartChain, Worker: worker.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: worker.PlatformStargate, Status: worker.StatusReady},
		{Network: network.VSL, Worker: worker.RSS3, Tags: []tag.Tag{tag.Exchange, tag.Collectible}, Platform: worker.PlatformRSS3, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: worker.Aavegotchi, Tags: []tag.Tag{tag.Metaverse}, Platform: worker.PlatformAavegotchi, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: worker.Highlight, Tags: []tag.Tag{tag.Collectible, tag.Transaction}, Platform: worker.PlatformHighlight, Status: worker.StatusReady},
		{Network: network.SatoshiVM, Worker: worker.Core, Tags: nil, Platform: worker.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.SatoshiVM, Worker: worker.SAVM, Tags: []tag.Tag{tag.Transaction}, Platform: worker.PlatformSAVM, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: worker.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: worker.PlatformStargate, Status: worker.StatusReady},
		{Network: network.Linea, Worker: worker.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: worker.PlatformStargate, Status: worker.StatusUnhealthy},
		{Network: network.Arbitrum, Worker: worker.Curve, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: worker.PlatformCurve, Status: worker.StatusReady},
		{Network: network.Arweave, Worker: worker.Mirror, Tags: []tag.Tag{tag.Social}, Platform: worker.PlatformMirror, Status: worker.StatusUnhealthy},
		{Network: network.Gnosis, Worker: worker.Core, Tags: nil, Platform: worker.PlatformUnknown, Status: worker.StatusIndexing},
		{Network: network.Ethereum, Worker: worker.VSL, Tags: []tag.Tag{tag.Transaction}, Platform: worker.PlatformVSL, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: worker.Core, Tags: nil, Platform: worker.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.Arweave, Worker: worker.Paragraph, Tags: []tag.Tag{tag.Social}, Platform: worker.PlatformParagraph, Status: worker.StatusUnhealthy},
		{Network: network.Linea, Worker: worker.Core, Tags: nil, Platform: worker.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.Optimism, Worker: worker.Matters, Tags: []tag.Tag{tag.Social}, Platform: worker.PlatformMatters, Status: worker.StatusIndexing},
		{Network: network.Ethereum, Worker: worker.OpenSea, Tags: []tag.Tag{tag.Collectible}, Platform: worker.PlatformOpenSea, Status: worker.StatusReady},
		{Network: network.Base, Worker: worker.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: worker.PlatformStargate, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: worker.IQWiki, Tags: []tag.Tag{tag.Social}, Platform: worker.PlatformIQWiki, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: worker.Lido, Tags: []tag.Tag{tag.Exchange, tag.Transaction, tag.Collectible}, Platform: worker.PlatformLido, Status: worker.StatusReady},
		{Network: network.Arbitrum, Worker: worker.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: worker.PlatformStargate, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: worker.Aave, Tags: []tag.Tag{tag.Exchange}, Platform: worker.PlatformAAVE, Status: worker.StatusReady},
		{Network: network.Optimism, Worker: worker.Highlight, Tags: []tag.Tag{tag.Collectible, tag.Transaction}, Platform: worker.PlatformHighlight, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: worker.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: worker.PlatformStargate, Status: worker.StatusReady},
		{Network: network.Optimism, Worker: worker.Core, Tags: nil, Platform: worker.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: worker.Optimism, Tags: []tag.Tag{tag.Transaction}, Platform: worker.PlatformOptimism, Status: worker.StatusReady},
		{Network: network.Crossbell, Worker: worker.Crossbell, Tags: []tag.Tag{tag.Social}, Platform: worker.PlatformCrossbell, Status: worker.StatusIndexing},
		{Network: network.Arbitrum, Worker: worker.Highlight, Tags: []tag.Tag{tag.Collectible, tag.Transaction}, Platform: worker.PlatformHighlight, Status: worker.StatusIndexing},
		{Network: network.Polygon, Worker: worker.Core, Tags: nil, Platform: worker.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: worker.Lens, Tags: []tag.Tag{tag.Social}, Platform: worker.PlatformLens, Status: worker.StatusUnhealthy},
		{Network: network.Avalanche, Worker: worker.Stargate, Tags: []tag.Tag{tag.Transaction}, Platform: worker.PlatformStargate, Status: worker.StatusReady},
		{Network: network.Linea, Worker: worker.Uniswap, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: worker.PlatformUniswap, Status: worker.StatusReady},
		{Network: network.Avalanche, Worker: worker.Core, Tags: nil, Platform: worker.PlatformUnknown, Status: worker.StatusIndexing},
		{Network: network.Polygon, Worker: worker.Curve, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: worker.PlatformCurve, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: worker.Curve, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: worker.PlatformCurve, Status: worker.StatusReady},
		{Network: network.Base, Worker: worker.Core, Tags: nil, Platform: worker.PlatformUnknown, Status: worker.StatusIndexing},
		{Network: network.Avalanche, Worker: worker.Curve, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: worker.PlatformCurve, Status: worker.StatusReady},
		{Network: network.VSL, Worker: worker.Core, Tags: nil, Platform: worker.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: worker.Aave, Tags: []tag.Tag{tag.Exchange}, Platform: worker.PlatformAAVE, Status: worker.StatusReady},
		{Network: network.Avalanche, Worker: worker.Aave, Tags: []tag.Tag{tag.Exchange}, Platform: worker.PlatformAAVE, Status: worker.StatusReady},
		{Network: network.Ethereum, Worker: worker.RSS3, Tags: []tag.Tag{tag.Exchange, tag.Collectible}, Platform: worker.PlatformRSS3, Status: worker.StatusReady},
		{Network: network.BinanceSmartChain, Worker: worker.Core, Tags: nil, Platform: worker.PlatformUnknown, Status: worker.StatusReady},
		{Network: network.Optimism, Worker: worker.Curve, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Platform: worker.PlatformCurve, Status: worker.StatusReady},
	}
	workerInfoNode2 = []*WorkerInfo{
		{Network: network.Farcaster, Worker: worker.Core, Tags: []tag.Tag{tag.Social}, Platform: worker.PlatformFarcaster, Status: worker.StatusReady},
	}
	workerInfoNode3 = []*WorkerInfo{
		{Network: network.Arweave, Worker: worker.Momoka, Tags: []tag.Tag{tag.Social}, Platform: worker.PlatformLens, Status: worker.StatusReady},
		{Network: network.Polygon, Worker: worker.Lens, Tags: []tag.Tag{tag.Social}, Platform: worker.PlatformLens, Status: worker.StatusReady},
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
	expectedFullNodeWorkerToNetworksMap := map[worker.Worker]map[network.Network]struct{}{
		worker.Aave: {
			network.Avalanche: {},
			network.Base:      {},
			network.Ethereum:  {},
			network.Optimism:  {},
			network.Polygon:   {},
		},
		worker.Aavegotchi: {
			network.Polygon: {},
		},
		worker.Core: {
			network.BinanceSmartChain: {},
			network.Ethereum:          {},
			network.Linea:             {},
			network.Optimism:          {},
			network.Polygon:           {},
			network.SatoshiVM:         {},
			network.VSL:               {},
			network.Farcaster:         {},
		},
		worker.Curve: {
			network.Arbitrum:  {},
			network.Avalanche: {},
			network.Ethereum:  {},
			network.Gnosis:    {},
			network.Optimism:  {},
			network.Polygon:   {},
		},
		worker.Highlight: {
			network.Ethereum: {},
			network.Optimism: {},
			network.Polygon:  {},
		},
		worker.IQWiki: {
			network.Polygon: {},
		},
		worker.KiwiStand: {
			network.Optimism: {},
		},
		worker.Lido: {
			network.Ethereum: {},
		},
		worker.Looksrare: {
			network.Ethereum: {},
		},
		worker.Oneinch: {
			network.Ethereum: {},
		},
		worker.OpenSea: {
			network.Ethereum: {},
		},
		worker.Optimism: {
			network.Ethereum: {},
		},
		worker.RSS3: {
			network.Ethereum: {},
			network.VSL:      {},
		},
		worker.SAVM: {
			network.SatoshiVM: {},
		},
		worker.Stargate: {
			network.Arbitrum:          {},
			network.Avalanche:         {},
			network.Base:              {},
			network.BinanceSmartChain: {},
			network.Ethereum:          {},
			network.Optimism:          {},
			network.Polygon:           {},
		},
		worker.Uniswap: {
			network.Ethereum:  {},
			network.Linea:     {},
			network.SatoshiVM: {},
		},
		worker.VSL: {
			network.Ethereum: {},
		},
		worker.Momoka: {
			network.Arweave: {},
		},
		worker.Lens: {
			network.Polygon: {},
		},
	}
	expectedNetworkToWorkersMap := map[network.Network]map[worker.Worker]struct{}{
		network.Arbitrum: {
			worker.Curve:    {},
			worker.Stargate: {},
		},
		network.Avalanche: {
			worker.Aave:     {},
			worker.Curve:    {},
			worker.Stargate: {},
		},
		network.Arweave: {
			worker.Momoka: {},
		},
		network.Base: {
			worker.Aave:     {},
			worker.Stargate: {},
		},
		network.BinanceSmartChain: {
			worker.Core:     {},
			worker.Stargate: {},
		},
		network.Ethereum: {
			worker.Aave:      {},
			worker.Core:      {},
			worker.Curve:     {},
			worker.Highlight: {},
			worker.Lido:      {},
			worker.Looksrare: {},
			worker.Oneinch:   {},
			worker.OpenSea:   {},
			worker.Optimism:  {},
			worker.RSS3:      {},
			worker.Stargate:  {},
			worker.Uniswap:   {},
			worker.VSL:       {},
		},
		network.Gnosis: {
			worker.Curve: {},
		},
		network.Linea: {
			worker.Core:    {},
			worker.Uniswap: {},
		},
		network.Optimism: {
			worker.Aave:      {},
			worker.Core:      {},
			worker.Curve:     {},
			worker.Highlight: {},
			worker.KiwiStand: {},
			worker.Stargate:  {},
		},
		network.Polygon: {
			worker.Aave:       {},
			worker.Aavegotchi: {},
			worker.Core:       {},
			worker.Curve:      {},
			worker.Highlight:  {},
			worker.IQWiki:     {},
			worker.Stargate:   {},
			worker.Lens:       {},
		},
		network.SatoshiVM: {
			worker.Core:    {},
			worker.Uniswap: {},
			worker.SAVM:    {},
		},
		network.VSL: {
			worker.Core: {},
			worker.RSS3: {},
		},

		network.Farcaster: {
			worker.Core: {},
		},
	}
	expectedPlatformToWorkersMap := map[worker.Platform]map[worker.Worker]struct{}{
		worker.Platform1Inch:      {worker.Oneinch: {}},
		worker.PlatformAAVE:       {worker.Aave: {}},
		worker.PlatformAavegotchi: {worker.Aavegotchi: {}},
		worker.PlatformCurve:      {worker.Curve: {}},
		worker.PlatformHighlight:  {worker.Highlight: {}},
		worker.PlatformIQWiki:     {worker.IQWiki: {}},
		worker.PlatformKiwiStand:  {worker.KiwiStand: {}},
		worker.PlatformLido:       {worker.Lido: {}},
		worker.PlatformLooksRare:  {worker.Looksrare: {}},
		worker.PlatformOpenSea:    {worker.OpenSea: {}},
		worker.PlatformOptimism:   {worker.Optimism: {}},
		worker.PlatformRSS3:       {worker.RSS3: {}},
		worker.PlatformSAVM:       {worker.SAVM: {}},
		worker.PlatformStargate:   {worker.Stargate: {}},
		worker.PlatformUniswap:    {worker.Uniswap: {}},
		worker.PlatformVSL:        {worker.VSL: {}},
		worker.PlatformFarcaster:  {worker.Core: {}},
		worker.PlatformLens:       {worker.Lens: {}, worker.Momoka: {}},
	}
	expectedTagToWorkersMap := map[tag.Tag]map[worker.Worker]struct{}{
		tag.Collectible: {
			worker.Highlight: {},
			worker.KiwiStand: {},
			worker.Lido:      {},
			worker.Looksrare: {},
			worker.OpenSea:   {},
			worker.RSS3:      {},
		},
		tag.Exchange: {
			worker.Aave:    {},
			worker.Curve:   {},
			worker.Lido:    {},
			worker.Oneinch: {},
			worker.RSS3:    {},
			worker.Uniswap: {},
		},
		tag.Metaverse: {
			worker.Aavegotchi: {},
		},
		tag.Social: {
			worker.KiwiStand: {},
			worker.IQWiki:    {},
			worker.Core:      {},
			worker.Momoka:    {},
			worker.Lens:      {},
		},
		tag.Transaction: {
			worker.KiwiStand: {},
			worker.Uniswap:   {},
			worker.Curve:     {},
			worker.Optimism:  {},
			worker.Highlight: {},
			worker.Lido:      {},
			worker.SAVM:      {},
			worker.Stargate:  {},
			worker.VSL:       {},
		},
	}

	assert.Equal(t, expectedNodeToWorkersMap[common.Address{1}], nodeToWorkersMap[common.Address{1}])
	assert.Equal(t, expectedFullNodeWorkerToNetworksMap, fullNodeWorkerToNetworksMap)
	assert.Equal(t, expectedNetworkToWorkersMap, networkToWorkersMap)
	assert.Equal(t, expectedPlatformToWorkersMap, platformToWorkersMap)
	assert.Equal(t, expectedTagToWorkersMap, tagToWorkersMap)
}

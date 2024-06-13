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
	workerStatusNode1 = `{"data":{"decentralized":[{"network":"optimism","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready","remote_state":121308330,"indexed_state":121308321},{"network":"arbitrum","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":221155623,"indexed_state":221155361},{"network":"ethereum","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"polygon","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready","remote_state":58079616,"indexed_state":58079615},{"network":"ethereum","worker":"looksrare","tags":["collectible"],"platform":"LooksRare","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"optimism","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":121308330,"indexed_state":121308327},{"network":"optimism","worker":"matters","tags":["social"],"platform":"Matters","status":"Ready","remote_state":121308330,"indexed_state":121308326},{"network":"avax","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready","remote_state":46634317,"indexed_state":46634312},{"network":"ethereum","worker":"optimism","tags":["transaction"],"platform":"Optimism","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"polygon","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":58079615,"indexed_state":58079613},{"network":"arbitrum","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready","remote_state":221155623,"indexed_state":221155055},{"network":"arweave","worker":"paragraph","tags":["social"],"platform":"Paragraph","status":"Ready","remote_state":1443533,"indexed_state":1443532},{"network":"polygon","worker":"highlight","tags":["collectible","transaction"],"platform":"Highlight","status":"Ready","remote_state":58079615,"indexed_state":58079611},{"network":"vsl","worker":"rss3","tags":["exchange","collectible"],"platform":"RSS3","status":"Ready","remote_state":4178459,"indexed_state":4178455},{"network":"savm","worker":"uniswap","tags":["exchange","transaction"],"platform":"Uniswap","status":"Ready","remote_state":2660829,"indexed_state":2660828},{"network":"avax","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":46634317,"indexed_state":46634309},{"network":"gnosis","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":34430243,"indexed_state":34430241},{"network":"polygon","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":58079616,"indexed_state":58079614},{"network":"optimism","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready","remote_state":121308330,"indexed_state":121308323},{"network":"polygon","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready","remote_state":58079616,"indexed_state":58079615},{"network":"ethereum","worker":"1inch","tags":["exchange"],"platform":"1inch","status":"Ready","remote_state":20077417,"indexed_state":20077417},{"network":"ethereum","worker":"rss3","tags":["exchange","collectible"],"platform":"RSS3","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"binance-smart-chain","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":39554478,"indexed_state":39554472},{"network":"ethereum","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":20077417,"indexed_state":20077417},{"network":"crossbell","worker":"crossbell","tags":["social"],"platform":"Crossbell","status":"Ready","remote_state":67839695,"indexed_state":67839680},{"network":"arbitrum","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready","remote_state":221155623,"indexed_state":218087588},{"network":"polygon","worker":"aavegotchi","tags":["metaverse"],"platform":"Aavegotchi","status":"Ready","remote_state":58079616,"indexed_state":58079611},{"network":"farcaster","worker":"core","tags":null,"platform":"Unknown","status":"Indexing","remote_state":1718215438006,"indexed_state":1718215435040},{"network":"ethereum","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"ethereum","worker":"vsl","tags":["transaction"],"platform":"VSL","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"ethereum","worker":"lido","tags":["exchange","transaction","collectible"],"platform":"Lido","status":"Ready","remote_state":20077417,"indexed_state":20077417},{"network":"avax","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":46634317,"indexed_state":41648157},{"network":"optimism","worker":"kiwistand","tags":["collectible","transaction","social"],"platform":"KiwiStand","status":"Ready","remote_state":121308330,"indexed_state":121308327},{"network":"base","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready","remote_state":15713044,"indexed_state":15713040},{"network":"base","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":15713044,"indexed_state":15713040},{"network":"arweave","worker":"momoka","tags":["social"],"platform":"Lens","status":"Indexing","remote_state":1718215438171,"indexed_state":1703358338044},{"network":"arbitrum","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":221155623,"indexed_state":221155224},{"network":"base","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":15713044,"indexed_state":15713034},{"network":"optimism","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":121308330,"indexed_state":121308328},{"network":"binance-smart-chain","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":39554477,"indexed_state":39554474},{"network":"gnosis","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready","remote_state":34430243,"indexed_state":34430242},{"network":"linea","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":5407344,"indexed_state":5407338},{"network":"ethereum","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":20077417,"indexed_state":20077417},{"network":"savm","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":2660829,"indexed_state":2660825},{"network":"ethereum","worker":"opensea","tags":["collectible"],"platform":"OpenSea","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"arbitrum","worker":"highlight","tags":["collectible","transaction"],"platform":"Highlight","status":"Ready","remote_state":221155623,"indexed_state":218991511},{"network":"linea","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":5407344,"indexed_state":5407343},{"network":"polygon","worker":"lens","tags":["social"],"platform":"Lens","status":"Ready","remote_state":58079616,"indexed_state":58079613},{"network":"linea","worker":"uniswap","tags":["exchange","transaction"],"platform":"Uniswap","status":"Ready","remote_state":5407344,"indexed_state":5407340},{"network":"avax","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready","remote_state":46634317,"indexed_state":46634310},{"network":"savm","worker":"savm","tags":["transaction"],"platform":"SAVM","status":"Ready","remote_state":2660829,"indexed_state":2660826},{"network":"ethereum","worker":"highlight","tags":["collectible","transaction"],"platform":"Highlight","status":"Ready","remote_state":20077417,"indexed_state":20077417},{"network":"ethereum","worker":"uniswap","tags":["exchange","transaction"],"platform":"Uniswap","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"vsl","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":4178459,"indexed_state":4178459},{"network":"polygon","worker":"iqwiki","tags":["social"],"platform":"IQWiki","status":"Ready","remote_state":58079616,"indexed_state":58079614},{"network":"optimism","worker":"highlight","tags":["collectible","transaction"],"platform":"Highlight","status":"Ready","remote_state":121308330,"indexed_state":121308327},{"network":"crossbell","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":67839695,"indexed_state":67839517},{"network":"arweave","worker":"mirror","tags":["social"],"platform":"Mirror","status":"Ready","remote_state":1443533,"indexed_state":1443532}],"rss":[],"federated":null}}`
	workerStatusNode2 = `{"data":{"decentralized":[{"network":"farcaster","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":1718215438006,"indexed_state":1718215435040}],"rss":[],"federated":null}}`
	workerStatusNode3 = `{"data":{"decentralized":[{"network":"arweave","worker":"momoka","tags":["social"],"platform":"Lens","status":"Ready","remote_state":1718215438171,"indexed_state":1703358338044},{"network":"polygon","worker":"lens","tags":["social"],"platform":"Lens","status":"Ready","remote_state":58079616,"indexed_state":58079613}],"rss":[],"federated":null}}`

	workerInfoNode1 = []*DecentralizedWorkerInfo{
		{WorkerInfo: WorkerInfo{Network: network.Optimism, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange, tag.Transaction}}, Worker: decentralized.Curve, Platform: decentralized.PlatformCurve},
		{WorkerInfo: WorkerInfo{Network: network.Arbitrum, Status: worker.StatusReady, Tags: nil}, Worker: decentralized.Core, Platform: decentralized.PlatformUnknown},
		{WorkerInfo: WorkerInfo{Network: network.Ethereum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange, tag.Transaction}}, Worker: decentralized.Curve, Platform: decentralized.PlatformCurve},
		{WorkerInfo: WorkerInfo{Network: network.Polygon, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange, tag.Transaction}}, Worker: decentralized.Curve, Platform: decentralized.PlatformCurve},
		{WorkerInfo: WorkerInfo{Network: network.Ethereum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Collectible}}, Worker: decentralized.Looksrare, Platform: decentralized.PlatformLooksRare},
		{WorkerInfo: WorkerInfo{Network: network.Optimism, Status: worker.StatusReady, Tags: nil}, Worker: decentralized.Core, Platform: decentralized.PlatformUnknown},
		{WorkerInfo: WorkerInfo{Network: network.Optimism, Status: worker.StatusReady, Tags: []tag.Tag{tag.Social}}, Worker: decentralized.Matters, Platform: decentralized.PlatformMatters},
		{WorkerInfo: WorkerInfo{Network: network.Avalanche, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange}}, Worker: decentralized.Aave, Platform: decentralized.PlatformAAVE},
		{WorkerInfo: WorkerInfo{Network: network.Ethereum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Transaction}}, Worker: decentralized.Optimism, Platform: decentralized.PlatformOptimism},
		{WorkerInfo: WorkerInfo{Network: network.Polygon, Status: worker.StatusReady, Tags: []tag.Tag{tag.Transaction}}, Worker: decentralized.Stargate, Platform: decentralized.PlatformStargate},
		{WorkerInfo: WorkerInfo{Network: network.Arbitrum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange, tag.Transaction}}, Worker: decentralized.Curve, Platform: decentralized.PlatformCurve},
		{WorkerInfo: WorkerInfo{Network: network.Arweave, Status: worker.StatusReady, Tags: []tag.Tag{tag.Social}}, Worker: decentralized.Paragraph, Platform: decentralized.PlatformParagraph},
		{WorkerInfo: WorkerInfo{Network: network.Polygon, Status: worker.StatusReady, Tags: []tag.Tag{tag.Collectible, tag.Transaction}}, Worker: decentralized.Highlight, Platform: decentralized.PlatformHighlight},
		{WorkerInfo: WorkerInfo{Network: network.VSL, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange, tag.Collectible}}, Worker: decentralized.RSS3, Platform: decentralized.PlatformRSS3},
		{WorkerInfo: WorkerInfo{Network: network.SatoshiVM, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange, tag.Transaction}}, Worker: decentralized.Uniswap, Platform: decentralized.PlatformUniswap},
		{WorkerInfo: WorkerInfo{Network: network.Avalanche, Status: worker.StatusReady, Tags: []tag.Tag{tag.Transaction}}, Worker: decentralized.Stargate, Platform: decentralized.PlatformStargate},
		{WorkerInfo: WorkerInfo{Network: network.Gnosis, Status: worker.StatusReady, Tags: nil}, Worker: decentralized.Core, Platform: decentralized.PlatformUnknown},
		{WorkerInfo: WorkerInfo{Network: network.Polygon, Status: worker.StatusReady, Tags: nil}, Worker: decentralized.Core, Platform: decentralized.PlatformUnknown},
		{WorkerInfo: WorkerInfo{Network: network.Optimism, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange}}, Worker: decentralized.Aave, Platform: decentralized.PlatformAAVE},
		{WorkerInfo: WorkerInfo{Network: network.Polygon, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange}}, Worker: decentralized.Aave, Platform: decentralized.PlatformAAVE},
		{WorkerInfo: WorkerInfo{Network: network.Ethereum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange}}, Worker: decentralized.Oneinch, Platform: decentralized.Platform1Inch},
		{WorkerInfo: WorkerInfo{Network: network.Ethereum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange, tag.Collectible}}, Worker: decentralized.RSS3, Platform: decentralized.PlatformRSS3},
		{WorkerInfo: WorkerInfo{Network: network.BinanceSmartChain, Status: worker.StatusReady, Tags: []tag.Tag{tag.Transaction}}, Worker: decentralized.Stargate, Platform: decentralized.PlatformStargate},
		{WorkerInfo: WorkerInfo{Network: network.Ethereum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Transaction}}, Worker: decentralized.Stargate, Platform: decentralized.PlatformStargate},
		{WorkerInfo: WorkerInfo{Network: network.Crossbell, Status: worker.StatusReady, Tags: []tag.Tag{tag.Social}}, Worker: decentralized.Crossbell, Platform: decentralized.PlatformCrossbell},
		{WorkerInfo: WorkerInfo{Network: network.Arbitrum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange}}, Worker: decentralized.Aave, Platform: decentralized.PlatformAAVE},
		{WorkerInfo: WorkerInfo{Network: network.Polygon, Status: worker.StatusReady, Tags: []tag.Tag{tag.Metaverse}}, Worker: decentralized.Aavegotchi, Platform: decentralized.PlatformAavegotchi},
		{WorkerInfo: WorkerInfo{Network: network.Farcaster, Status: worker.StatusIndexing, Tags: []tag.Tag{tag.Social}}, Worker: decentralized.Core, Platform: decentralized.PlatformFarcaster},
		{WorkerInfo: WorkerInfo{Network: network.Ethereum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange}}, Worker: decentralized.Aave, Platform: decentralized.PlatformAAVE},
		{WorkerInfo: WorkerInfo{Network: network.Ethereum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Transaction}}, Worker: decentralized.VSL, Platform: decentralized.PlatformVSL},
		{WorkerInfo: WorkerInfo{Network: network.Ethereum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange, tag.Transaction, tag.Collectible}}, Worker: decentralized.Lido, Platform: decentralized.PlatformLido},
		{WorkerInfo: WorkerInfo{Network: network.Avalanche, Status: worker.StatusReady, Tags: nil}, Worker: decentralized.Core, Platform: decentralized.PlatformUnknown},
		{WorkerInfo: WorkerInfo{Network: network.Optimism, Tags: []tag.Tag{tag.Collectible, tag.Transaction, tag.Social}, Status: worker.StatusReady}, Worker: decentralized.KiwiStand, Platform: decentralized.PlatformKiwiStand},
		{WorkerInfo: WorkerInfo{Network: network.Base, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange}}, Worker: decentralized.Aave, Platform: decentralized.PlatformAAVE},
		{WorkerInfo: WorkerInfo{Network: network.Base, Status: worker.StatusReady, Tags: []tag.Tag{tag.Transaction}}, Worker: decentralized.Stargate, Platform: decentralized.PlatformStargate},
		{WorkerInfo: WorkerInfo{Network: network.Arweave, Status: worker.StatusIndexing, Tags: []tag.Tag{tag.Social}}, Worker: decentralized.Momoka, Platform: decentralized.PlatformLens},
		{WorkerInfo: WorkerInfo{Network: network.Arbitrum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Transaction}}, Worker: decentralized.Stargate, Platform: decentralized.PlatformStargate},
		{WorkerInfo: WorkerInfo{Network: network.Base, Status: worker.StatusReady, Tags: nil}, Worker: decentralized.Core, Platform: decentralized.PlatformUnknown},
		{WorkerInfo: WorkerInfo{Network: network.Optimism, Status: worker.StatusReady, Tags: []tag.Tag{tag.Transaction}}, Worker: decentralized.Stargate, Platform: decentralized.PlatformStargate},
		{WorkerInfo: WorkerInfo{Network: network.BinanceSmartChain, Status: worker.StatusReady, Tags: nil}, Worker: decentralized.Core, Platform: decentralized.PlatformUnknown},
		{WorkerInfo: WorkerInfo{Network: network.Gnosis, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange, tag.Transaction}}, Worker: decentralized.Curve, Platform: decentralized.PlatformCurve},
		{WorkerInfo: WorkerInfo{Network: network.Linea, Status: worker.StatusReady, Tags: nil}, Worker: decentralized.Core, Platform: decentralized.PlatformUnknown},
		{WorkerInfo: WorkerInfo{Network: network.Ethereum, Status: worker.StatusReady, Tags: nil}, Worker: decentralized.Core, Platform: decentralized.PlatformUnknown},
		{WorkerInfo: WorkerInfo{Network: network.SatoshiVM, Status: worker.StatusReady, Tags: nil}, Worker: decentralized.Core, Platform: decentralized.PlatformUnknown},
		{WorkerInfo: WorkerInfo{Network: network.Ethereum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Collectible}}, Worker: decentralized.OpenSea, Platform: decentralized.PlatformOpenSea},
		{WorkerInfo: WorkerInfo{Network: network.Arbitrum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Collectible, tag.Transaction}}, Worker: decentralized.Highlight, Platform: decentralized.PlatformHighlight},
		{WorkerInfo: WorkerInfo{Network: network.Linea, Status: worker.StatusReady, Tags: []tag.Tag{tag.Transaction}}, Worker: decentralized.Stargate, Platform: decentralized.PlatformStargate},
		{WorkerInfo: WorkerInfo{Network: network.Polygon, Status: worker.StatusReady, Tags: []tag.Tag{tag.Social}}, Worker: decentralized.Lens, Platform: decentralized.PlatformLens},
		{WorkerInfo: WorkerInfo{Network: network.Linea, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange, tag.Transaction}}, Worker: decentralized.Uniswap, Platform: decentralized.PlatformUniswap},
		{WorkerInfo: WorkerInfo{Network: network.Avalanche, Status: worker.StatusReady, Tags: []tag.Tag{tag.Exchange, tag.Transaction}}, Worker: decentralized.Curve, Platform: decentralized.PlatformCurve},
		{WorkerInfo: WorkerInfo{Network: network.SatoshiVM, Status: worker.StatusReady, Tags: []tag.Tag{tag.Transaction}}, Worker: decentralized.SAVM, Platform: decentralized.PlatformSAVM},
		{WorkerInfo: WorkerInfo{Network: network.Ethereum, Status: worker.StatusReady, Tags: []tag.Tag{tag.Collectible, tag.Transaction}}, Worker: decentralized.Highlight, Platform: decentralized.PlatformHighlight},
		{WorkerInfo: WorkerInfo{Network: network.Ethereum, Tags: []tag.Tag{tag.Exchange, tag.Transaction}, Status: worker.StatusReady}, Worker: decentralized.Uniswap, Platform: decentralized.PlatformUniswap},
		{WorkerInfo: WorkerInfo{Network: network.VSL, Status: worker.StatusReady, Tags: nil}, Worker: decentralized.Core, Platform: decentralized.PlatformUnknown},
		{WorkerInfo: WorkerInfo{Network: network.Polygon, Status: worker.StatusReady, Tags: []tag.Tag{tag.Social}}, Worker: decentralized.IQWiki, Platform: decentralized.PlatformIQWiki},
		{WorkerInfo: WorkerInfo{Network: network.Optimism, Status: worker.StatusReady, Tags: []tag.Tag{tag.Collectible, tag.Transaction}}, Worker: decentralized.Highlight, Platform: decentralized.PlatformHighlight},
		{WorkerInfo: WorkerInfo{Network: network.Crossbell, Status: worker.StatusReady, Tags: nil}, Worker: decentralized.Core, Platform: decentralized.PlatformUnknown},
		{WorkerInfo: WorkerInfo{Network: network.Arweave, Status: worker.StatusReady, Tags: []tag.Tag{tag.Social}}, Worker: decentralized.Mirror, Platform: decentralized.PlatformMirror},
	}
	workerInfoNode2 = []*DecentralizedWorkerInfo{
		{WorkerInfo: WorkerInfo{Network: network.Farcaster, Status: worker.StatusReady, Tags: []tag.Tag{tag.Social}}, Worker: decentralized.Core, Platform: decentralized.PlatformFarcaster},
	}
	workerInfoNode3 = []*DecentralizedWorkerInfo{
		{WorkerInfo: WorkerInfo{Network: network.Arweave, Status: worker.StatusReady, Tags: []tag.Tag{tag.Social}}, Worker: decentralized.Momoka, Platform: decentralized.PlatformLens},
		{WorkerInfo: WorkerInfo{Network: network.Polygon, Status: worker.StatusReady, Tags: []tag.Tag{tag.Social}}, Worker: decentralized.Lens, Platform: decentralized.PlatformLens},
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
	mockClient.On("Fetch", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNode1))), nil)

	enforcer := &SimpleEnforcer{httpClient: mockClient}
	response, err := enforcer.getNodeWorkerStatus(context.Background(), "http://localhost:8080")

	assert.NoError(t, err)
	assert.Equal(t, workerInfoNode1, response.Data.Decentralized)
}

func TestGenerateMaps(t *testing.T) {
	t.Parallel()

	mockClient := new(MockHTTPClient)
	mockClient.On("Fetch", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNode1))), nil)
	mockClient.On("Fetch", mock.Anything, "http://localhost:8081/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNode2))), nil)
	mockClient.On("Fetch", mock.Anything, "http://localhost:8082/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNode3))), nil)

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

	expectedNodeToWorkersMap := map[common.Address]*ComponentInfo{
		common.Address{1}: {
			Decentralized: workerInfoNode1,
		},
		common.Address{2}: {
			Decentralized: workerInfoNode2,
		},
		common.Address{3}: {
			Decentralized: workerInfoNode3,
		},
	}
	assert.Equal(t, expectedNodeToWorkersMap[common.Address{1}].Decentralized, nodeToWorkersMap[common.Address{1}].Decentralized)

	expectedFullNodeWorkerToNetworksMap := map[decentralized.Worker]map[network.Network]struct{}{
		decentralized.Aave: {
			network.Avalanche: {},
			network.Arbitrum:  {},
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
			network.Crossbell:         {},
			network.Arbitrum:          {},
			network.Gnosis:            {},
			network.Avalanche:         {},
			network.Base:              {},
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
			network.Arbitrum: {},
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
			network.Linea:             {},
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
		decentralized.Mirror: {
			network.Arweave: {},
		},
		decentralized.Paragraph: {
			network.Arweave: {},
		},
		decentralized.Matters: {
			network.Optimism: {},
		},
		decentralized.Crossbell: {
			network.Crossbell: {},
		},
	}

	assert.Equal(t, expectedFullNodeWorkerToNetworksMap, fullNodeWorkerToNetworksMap)

	expectedNetworkToWorkersMap := map[network.Network]map[decentralized.Worker]struct{}{
		network.Arbitrum: {
			decentralized.Curve:     {},
			decentralized.Stargate:  {},
			decentralized.Core:      {},
			decentralized.Aave:      {},
			decentralized.Highlight: {},
		},
		network.Avalanche: {
			decentralized.Aave:     {},
			decentralized.Curve:    {},
			decentralized.Stargate: {},
			decentralized.Core:     {},
		},
		network.Arweave: {
			decentralized.Momoka:    {},
			decentralized.Mirror:    {},
			decentralized.Paragraph: {},
		},
		network.Base: {
			decentralized.Aave:     {},
			decentralized.Stargate: {},
			decentralized.Core:     {},
		},
		network.BinanceSmartChain: {
			decentralized.Core:     {},
			decentralized.Stargate: {},
		},
		network.Crossbell: {
			decentralized.Crossbell: {},
			decentralized.Core:      {},
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
			decentralized.Core:  {},
		},
		network.Linea: {
			decentralized.Core:     {},
			decentralized.Stargate: {},
			decentralized.Uniswap:  {},
		},
		network.Optimism: {
			decentralized.Aave:      {},
			decentralized.Core:      {},
			decentralized.Curve:     {},
			decentralized.Highlight: {},
			decentralized.KiwiStand: {},
			decentralized.Stargate:  {},
			decentralized.Matters:   {},
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

	assert.Equal(t, expectedNetworkToWorkersMap, networkToWorkersMap)

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
		decentralized.PlatformMirror:     {decentralized.Mirror: {}},
		decentralized.PlatformParagraph:  {decentralized.Paragraph: {}},
		decentralized.PlatformMatters:    {decentralized.Matters: {}},
		decentralized.PlatformCrossbell:  {decentralized.Crossbell: {}},
	}

	assert.Equal(t, expectedPlatformToWorkersMap, platformToWorkersMap)

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
			decentralized.Mirror:    {},
			decentralized.Paragraph: {},
			decentralized.Crossbell: {},
			decentralized.Matters:   {},
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

	assert.Equal(t, expectedTagToWorkersMap, tagToWorkersMap)
}

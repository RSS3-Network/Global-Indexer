package enforcer

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/rss3-network/node/schema/worker"
	"github.com/rss3-network/node/schema/worker/decentralized"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	workerStatusNode1 = `{"data":{"decentralized":[{"network":"optimism","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready","remote_state":121308330,"indexed_state":121308321},{"network":"arbitrum","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":221155623,"indexed_state":221155361},{"network":"ethereum","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"polygon","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready","remote_state":58079616,"indexed_state":58079615},{"network":"ethereum","worker":"looksrare","tags":["collectible"],"platform":"LooksRare","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"optimism","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":121308330,"indexed_state":121308327},{"network":"optimism","worker":"matters","tags":["social"],"platform":"Matters","status":"Ready","remote_state":121308330,"indexed_state":121308326},{"network":"avax","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready","remote_state":46634317,"indexed_state":46634312},{"network":"ethereum","worker":"optimism","tags":["transaction"],"platform":"Optimism","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"polygon","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":58079615,"indexed_state":58079613},{"network":"arbitrum","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready","remote_state":221155623,"indexed_state":221155055},{"network":"arweave","worker":"paragraph","tags":["social"],"platform":"Paragraph","status":"Ready","remote_state":1443533,"indexed_state":1443532},{"network":"polygon","worker":"highlight","tags":["collectible","transaction"],"platform":"Highlight","status":"Ready","remote_state":58079615,"indexed_state":58079611},{"network":"savm","worker":"uniswap","tags":["exchange","transaction"],"platform":"Uniswap","status":"Ready","remote_state":2660829,"indexed_state":2660828},{"network":"avax","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":46634317,"indexed_state":46634309},{"network":"gnosis","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":34430243,"indexed_state":34430241},{"network":"polygon","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":58079616,"indexed_state":58079614},{"network":"optimism","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready","remote_state":121308330,"indexed_state":121308323},{"network":"polygon","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready","remote_state":58079616,"indexed_state":58079615},{"network":"ethereum","worker":"1inch","tags":["exchange"],"platform":"1inch","status":"Ready","remote_state":20077417,"indexed_state":20077417},{"network":"ethereum","worker":"rss3","tags":["exchange","collectible"],"platform":"RSS3","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"binance-smart-chain","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":39554478,"indexed_state":39554472},{"network":"ethereum","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":20077417,"indexed_state":20077417},{"network":"crossbell","worker":"crossbell","tags":["social"],"platform":"Crossbell","status":"Ready","remote_state":67839695,"indexed_state":67839680},{"network":"arbitrum","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready","remote_state":221155623,"indexed_state":218087588},{"network":"polygon","worker":"aavegotchi","tags":["metaverse"],"platform":"Aavegotchi","status":"Ready","remote_state":58079616,"indexed_state":58079611},{"network":"farcaster","worker":"core","tags":null,"platform":"Unknown","status":"Indexing","remote_state":1718215438006,"indexed_state":1718215435040},{"network":"ethereum","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"ethereum","worker":"lido","tags":["exchange","transaction","collectible"],"platform":"Lido","status":"Ready","remote_state":20077417,"indexed_state":20077417},{"network":"avax","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":46634317,"indexed_state":41648157},{"network":"optimism","worker":"kiwistand","tags":["collectible","transaction","social"],"platform":"KiwiStand","status":"Ready","remote_state":121308330,"indexed_state":121308327},{"network":"base","worker":"aave","tags":["exchange"],"platform":"AAVE","status":"Ready","remote_state":15713044,"indexed_state":15713040},{"network":"base","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":15713044,"indexed_state":15713040},{"network":"arweave","worker":"momoka","tags":["social"],"platform":"Lens","status":"Indexing","remote_state":1718215438171,"indexed_state":1703358338044},{"network":"arbitrum","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":221155623,"indexed_state":221155224},{"network":"base","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":15713044,"indexed_state":15713034},{"network":"optimism","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":121308330,"indexed_state":121308328},{"network":"binance-smart-chain","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":39554477,"indexed_state":39554474},{"network":"gnosis","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready","remote_state":34430243,"indexed_state":34430242},{"network":"linea","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":5407344,"indexed_state":5407338},{"network":"ethereum","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":20077417,"indexed_state":20077417},{"network":"savm","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":2660829,"indexed_state":2660825},{"network":"ethereum","worker":"opensea","tags":["collectible"],"platform":"OpenSea","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"arbitrum","worker":"highlight","tags":["collectible","transaction"],"platform":"Highlight","status":"Ready","remote_state":221155623,"indexed_state":218991511},{"network":"linea","worker":"stargate","tags":["transaction"],"platform":"Stargate","status":"Ready","remote_state":5407344,"indexed_state":5407343},{"network":"polygon","worker":"lens","tags":["social"],"platform":"Lens","status":"Ready","remote_state":58079616,"indexed_state":58079613},{"network":"linea","worker":"uniswap","tags":["exchange","transaction"],"platform":"Uniswap","status":"Ready","remote_state":5407344,"indexed_state":5407340},{"network":"avax","worker":"curve","tags":["exchange","transaction"],"platform":"Curve","status":"Ready","remote_state":46634317,"indexed_state":46634310},{"network":"savm","worker":"savm","tags":["transaction"],"platform":"SAVM","status":"Ready","remote_state":2660829,"indexed_state":2660826},{"network":"ethereum","worker":"highlight","tags":["collectible","transaction"],"platform":"Highlight","status":"Ready","remote_state":20077417,"indexed_state":20077417},{"network":"ethereum","worker":"uniswap","tags":["exchange","transaction"],"platform":"Uniswap","status":"Ready","remote_state":20077417,"indexed_state":20077416},{"network":"vsl","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":4178459,"indexed_state":4178459},{"network":"polygon","worker":"iqwiki","tags":["social"],"platform":"IQWiki","status":"Ready","remote_state":58079616,"indexed_state":58079614},{"network":"optimism","worker":"highlight","tags":["collectible","transaction"],"platform":"Highlight","status":"Ready","remote_state":121308330,"indexed_state":121308327},{"network":"crossbell","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":67839695,"indexed_state":67839517},{"network":"arweave","worker":"mirror","tags":["social"],"platform":"Mirror","status":"Ready","remote_state":1443533,"indexed_state":1443532}],"rss":null,"federated":null}}`
	workerStatusNode2 = `{"data":{"decentralized":[{"network":"farcaster","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":1718215438006,"indexed_state":1718215435040}],"rss":null,"federated":null}}`
	workerStatusNode3 = `{"data":{"decentralized":[{"network":"arweave","worker":"momoka","tags":["social"],"platform":"Lens","status":"Ready","remote_state":1718215438171,"indexed_state":1703358338044},{"network":"arweave","worker":"momoka","tags":["social"],"platform":"Lens","status":"Ready","remote_state":718215438171,"indexed_state":703358338044},{"network":"polygon","worker":"lens","tags":["social"],"platform":"Lens","status":"Ready","remote_state":58079616,"indexed_state":58079613},{"network":"polygon","worker":"lens","tags":["social"],"platform":"Lens","status":"Ready","remote_state":48079616,"indexed_state":48079613},{"network":"polygon","worker":"lens","tags":["social"],"platform":"Lens","status":"Indexing","remote_state":8079616,"indexed_state":8079613}],"rss":null,"federated":null}}`

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
		{WorkerInfo: WorkerInfo{Network: network.Polygon, Status: worker.StatusIndexing, Tags: []tag.Tag{tag.Social}}, Worker: decentralized.Lens, Platform: decentralized.PlatformLens},
	}
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) FetchWithMethod(ctx context.Context, _, endpoint string, _ string, _ io.Reader) (io.ReadCloser, error) {
	args := m.Called(ctx, endpoint)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func Test_GetNodeWorkerStatus(t *testing.T) {
	t.Parallel()

	mockClient := new(MockHTTPClient)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNode1))), nil)

	enforcer := &SimpleEnforcer{httpClient: mockClient}
	response, err := enforcer.getNodeWorkerStatus(context.Background(), "http://localhost:8080", "")

	assert.NoError(t, err)
	assert.Equal(t, workerInfoNode1, response.Data.Decentralized)
}

func Test_GenerateMaps(t *testing.T) {
	t.Parallel()

	mockClient := new(MockHTTPClient)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNode1))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8081/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNode2))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8082/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNode3))), nil)

	enforcer := &SimpleEnforcer{
		httpClient: mockClient,
	}
	stats := []*schema.Stat{
		{
			Address:  common.Address{1},
			Endpoint: "http://localhost:8080",
			Status:   schema.NodeStatusOnline,
			Version:  "v1.0.0",
		},
		{
			Address:  common.Address{2},
			Endpoint: "http://localhost:8081",
			Status:   schema.NodeStatusOnline,
			Version:  "v1.0.0",
		},
		{
			Address:  common.Address{3},
			Endpoint: "http://localhost:8082",
			Status:   schema.NodeStatusOnline,
			Version:  "v1.0.0",
		},
	}

	nodeToWorkersMap, fullNodeWorkerToNetworksMap, networkToWorkersMap, platformToWorkersMap, tagToWorkersMap := enforcer.generateMaps(context.Background(), stats, "v1.0.0")

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

	assert.Equal(t, lo.SliceToMap(expectedNodeToWorkersMap[common.Address{3}].Decentralized, func(workerInfo *DecentralizedWorkerInfo) (string, worker.Status) {
		return workerInfo.Network.String() + workerInfo.Worker.String(), workerInfo.Status
	}), lo.SliceToMap(nodeToWorkersMap[common.Address{3}].Decentralized, func(workerInfo *DecentralizedWorkerInfo) (string, worker.Status) {
		return workerInfo.Network.String() + workerInfo.Worker.String(), workerInfo.Status
	}))

	expectedFullNodeWorkerToNetworksMap := map[string]map[string]struct{}{
		decentralized.Aave.String(): {
			network.Avalanche.String(): {},
			network.Arbitrum.String():  {},
			network.Base.String():      {},
			network.Ethereum.String():  {},
			network.Optimism.String():  {},
			network.Polygon.String():   {},
		},
		decentralized.Aavegotchi.String(): {
			network.Polygon.String(): {},
		},
		decentralized.Core.String(): {
			network.BinanceSmartChain.String(): {},
			network.Ethereum.String():          {},
			network.Linea.String():             {},
			network.Optimism.String():          {},
			network.Polygon.String():           {},
			network.SatoshiVM.String():         {},
			network.VSL.String():               {},
			//network.Farcaster.String():         {},
			network.Crossbell.String(): {},
			network.Arbitrum.String():  {},
			network.Gnosis.String():    {},
			network.Avalanche.String(): {},
			network.Base.String():      {},
		},
		decentralized.Curve.String(): {
			network.Arbitrum.String():  {},
			network.Avalanche.String(): {},
			network.Ethereum.String():  {},
			network.Gnosis.String():    {},
			network.Optimism.String():  {},
			network.Polygon.String():   {},
		},
		"farcaster": {
			network.Farcaster.String(): {},
		},
		decentralized.Highlight.String(): {
			network.Ethereum.String(): {},
			network.Optimism.String(): {},
			network.Polygon.String():  {},
			network.Arbitrum.String(): {},
		},
		decentralized.IQWiki.String(): {
			network.Polygon.String(): {},
		},
		decentralized.KiwiStand.String(): {
			network.Optimism.String(): {},
		},
		decentralized.Lido.String(): {
			network.Ethereum.String(): {},
		},
		decentralized.Looksrare.String(): {
			network.Ethereum.String(): {},
		},
		decentralized.Oneinch.String(): {
			network.Ethereum.String(): {},
		},
		decentralized.OpenSea.String(): {
			network.Ethereum.String(): {},
		},
		decentralized.Optimism.String(): {
			network.Ethereum.String(): {},
		},
		decentralized.RSS3.String(): {
			network.Ethereum.String(): {},
		},
		decentralized.SAVM.String(): {
			network.SatoshiVM.String(): {},
		},
		decentralized.Stargate.String(): {
			network.Arbitrum.String():          {},
			network.Avalanche.String():         {},
			network.Base.String():              {},
			network.BinanceSmartChain.String(): {},
			network.Ethereum.String():          {},
			network.Optimism.String():          {},
			network.Polygon.String():           {},
			network.Linea.String():             {},
		},
		decentralized.Uniswap.String(): {
			network.Ethereum.String():  {},
			network.Linea.String():     {},
			network.SatoshiVM.String(): {},
		},
		decentralized.Momoka.String(): {
			network.Arweave.String(): {},
		},
		decentralized.Lens.String(): {
			network.Polygon.String(): {},
		},
		decentralized.Mirror.String(): {
			network.Arweave.String(): {},
		},
		decentralized.Paragraph.String(): {
			network.Arweave.String(): {},
		},
		decentralized.Matters.String(): {
			network.Optimism.String(): {},
		},
		decentralized.Crossbell.String(): {
			network.Crossbell.String(): {},
		},
	}

	assert.Equal(t, expectedFullNodeWorkerToNetworksMap, fullNodeWorkerToNetworksMap)

	expectedNetworkToWorkersMap := map[string]map[string]struct{}{
		network.Arbitrum.String(): {
			decentralized.Curve.String():     {},
			decentralized.Stargate.String():  {},
			decentralized.Core.String():      {},
			decentralized.Aave.String():      {},
			decentralized.Highlight.String(): {},
		},
		network.Avalanche.String(): {
			decentralized.Aave.String():     {},
			decentralized.Curve.String():    {},
			decentralized.Stargate.String(): {},
			decentralized.Core.String():     {},
		},
		network.Arweave.String(): {
			decentralized.Momoka.String():    {},
			decentralized.Mirror.String():    {},
			decentralized.Paragraph.String(): {},
		},
		network.Base.String(): {
			decentralized.Aave.String():     {},
			decentralized.Stargate.String(): {},
			decentralized.Core.String():     {},
		},
		network.BinanceSmartChain.String(): {
			decentralized.Core.String():     {},
			decentralized.Stargate.String(): {},
		},
		network.Crossbell.String(): {
			decentralized.Crossbell.String(): {},
			decentralized.Core.String():      {},
		},
		network.Ethereum.String(): {
			decentralized.Aave.String():      {},
			decentralized.Core.String():      {},
			decentralized.Curve.String():     {},
			decentralized.Highlight.String(): {},
			decentralized.Lido.String():      {},
			decentralized.Looksrare.String(): {},
			decentralized.Oneinch.String():   {},
			decentralized.OpenSea.String():   {},
			decentralized.Optimism.String():  {},
			decentralized.RSS3.String():      {},
			decentralized.Stargate.String():  {},
			decentralized.Uniswap.String():   {},
		},
		network.Gnosis.String(): {
			decentralized.Curve.String(): {},
			decentralized.Core.String():  {},
		},
		network.Linea.String(): {
			decentralized.Core.String():     {},
			decentralized.Stargate.String(): {},
			decentralized.Uniswap.String():  {},
		},
		network.Optimism.String(): {
			decentralized.Aave.String():      {},
			decentralized.Core.String():      {},
			decentralized.Curve.String():     {},
			decentralized.Highlight.String(): {},
			decentralized.KiwiStand.String(): {},
			decentralized.Stargate.String():  {},
			decentralized.Matters.String():   {},
		},
		network.Polygon.String(): {
			decentralized.Aave.String():       {},
			decentralized.Aavegotchi.String(): {},
			decentralized.Core.String():       {},
			decentralized.Curve.String():      {},
			decentralized.Highlight.String():  {},
			decentralized.IQWiki.String():     {},
			decentralized.Stargate.String():   {},
			decentralized.Lens.String():       {},
		},
		network.SatoshiVM.String(): {
			decentralized.Core.String():    {},
			decentralized.Uniswap.String(): {},
			decentralized.SAVM.String():    {},
		},
		network.VSL.String(): {
			decentralized.Core.String(): {},
		},

		network.Farcaster.String(): {
			"farcaster": {},
		},
	}

	assert.Equal(t, expectedNetworkToWorkersMap, networkToWorkersMap)

	expectedPlatformToWorkersMap := map[string]map[string]struct{}{
		decentralized.Platform1Inch.String():      {decentralized.Oneinch.String(): {}},
		decentralized.PlatformAAVE.String():       {decentralized.Aave.String(): {}},
		decentralized.PlatformAavegotchi.String(): {decentralized.Aavegotchi.String(): {}},
		decentralized.PlatformCurve.String():      {decentralized.Curve.String(): {}},
		decentralized.PlatformHighlight.String():  {decentralized.Highlight.String(): {}},
		decentralized.PlatformIQWiki.String():     {decentralized.IQWiki.String(): {}},
		decentralized.PlatformKiwiStand.String():  {decentralized.KiwiStand.String(): {}},
		decentralized.PlatformLido.String():       {decentralized.Lido.String(): {}},
		decentralized.PlatformLooksRare.String():  {decentralized.Looksrare.String(): {}},
		decentralized.PlatformOpenSea.String():    {decentralized.OpenSea.String(): {}},
		decentralized.PlatformOptimism.String():   {decentralized.Optimism.String(): {}},
		decentralized.PlatformRSS3.String():       {decentralized.RSS3.String(): {}},
		decentralized.PlatformSAVM.String():       {decentralized.SAVM.String(): {}},
		decentralized.PlatformStargate.String():   {decentralized.Stargate.String(): {}},
		decentralized.PlatformUniswap.String():    {decentralized.Uniswap.String(): {}},
		decentralized.PlatformFarcaster.String():  {"farcaster": {}},
		decentralized.PlatformLens.String():       {decentralized.Lens.String(): {}, decentralized.Momoka.String(): {}},
		decentralized.PlatformMirror.String():     {decentralized.Mirror.String(): {}},
		decentralized.PlatformParagraph.String():  {decentralized.Paragraph.String(): {}},
		decentralized.PlatformMatters.String():    {decentralized.Matters.String(): {}},
		decentralized.PlatformCrossbell.String():  {decentralized.Crossbell.String(): {}},
	}

	assert.Equal(t, expectedPlatformToWorkersMap, platformToWorkersMap)

	expectedTagToWorkersMap := map[string]map[string]struct{}{
		tag.Collectible.String(): {
			decentralized.Highlight.String(): {},
			decentralized.KiwiStand.String(): {},
			decentralized.Lido.String():      {},
			decentralized.Looksrare.String(): {},
			decentralized.OpenSea.String():   {},
			decentralized.RSS3.String():      {},
		},
		tag.Exchange.String(): {
			decentralized.Aave.String():    {},
			decentralized.Curve.String():   {},
			decentralized.Lido.String():    {},
			decentralized.Oneinch.String(): {},
			decentralized.RSS3.String():    {},
			decentralized.Uniswap.String(): {},
		},
		tag.Metaverse.String(): {
			decentralized.Aavegotchi.String(): {},
		},
		tag.Social.String(): {
			decentralized.KiwiStand.String(): {},
			decentralized.IQWiki.String():    {},
			"farcaster":                      {},
			decentralized.Momoka.String():    {},
			decentralized.Lens.String():      {},
			decentralized.Mirror.String():    {},
			decentralized.Paragraph.String(): {},
			decentralized.Crossbell.String(): {},
			decentralized.Matters.String():   {},
		},
		tag.Transaction.String(): {
			decentralized.KiwiStand.String(): {},
			decentralized.Uniswap.String():   {},
			decentralized.Curve.String():     {},
			decentralized.Optimism.String():  {},
			decentralized.Highlight.String(): {},
			decentralized.Lido.String():      {},
			decentralized.SAVM.String():      {},
			decentralized.Stargate.String():  {},
		},
	}

	assert.Equal(t, expectedTagToWorkersMap, tagToWorkersMap)
}

func Test_GenerateMapsNodeStatus(t *testing.T) {
	t.Parallel()

	mockClient := new(MockHTTPClient)

	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeOnline))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8081/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeUnhealthy))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8082/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeIndexing))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8083/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeIndexing))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8084/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(``))), errors.New("workers_status"))
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8085/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeRssOnline))), nil)

	enforcer := &SimpleEnforcer{
		httpClient: mockClient,
	}
	stats := []*schema.Stat{
		{
			Address:  common.Address{0},
			Endpoint: "http://localhost:8080",
			Status:   schema.NodeStatusOnline,
			HearBeat: schema.NodeStatusOnline,
			Version:  "v1.0.0",
		},
		{
			Address:  common.Address{1},
			Endpoint: "http://localhost:8081",
			Status:   schema.NodeStatusOnline,
			HearBeat: schema.NodeStatusOnline,
			Version:  "v1.0.0",
		},
		{
			Address:  common.Address{2},
			Endpoint: "http://localhost:8082",
			Status:   schema.NodeStatusOnline,
			HearBeat: schema.NodeStatusOnline,
			Version:  "v1.0.0",
		},
		{
			Address:  common.Address{3},
			Endpoint: "http://localhost:8083",
			Status:   schema.NodeStatusOnline,
			HearBeat: schema.NodeStatusOnline,
			Version:  "v0.1.0",
		},
		{
			Address:  common.Address{4},
			Endpoint: "http://localhost:8084",
			Status:   schema.NodeStatusRegistered,
			HearBeat: schema.NodeStatusOnline,
			Version:  "v1.0.0",
		},
		{
			Address:  common.Address{5},
			Endpoint: "http://localhost:8080",
			Status:   schema.NodeStatusExiting,
			HearBeat: schema.NodeStatusOnline,
			Version:  "v1.0.0",
		},
		{
			Address:  common.Address{6},
			Endpoint: "http://localhost:8084",
			Status:   schema.NodeStatusRegistered,
			HearBeat: schema.NodeStatusOffline,
			Version:  "v1.0.0",
		},
		{
			Address:  common.Address{7},
			Endpoint: "http://localhost:8085",
			Status:   schema.NodeStatusInitializing,
			HearBeat: schema.NodeStatusOnline,
			Version:  "v1.0.0",
		},
	}

	_, _, _, _, _ = enforcer.generateMaps(context.Background(), stats, "v1.0.0")

	assert.Equal(t, schema.NodeStatusOnline, stats[0].Status)
	assert.Equal(t, schema.NodeStatusRegistered, stats[1].Status)
	assert.Equal(t, schema.NodeStatusInitializing, stats[2].Status)
	assert.Equal(t, schema.NodeStatusOutdated, stats[3].Status)
	assert.Equal(t, schema.NodeStatusInitializing, stats[4].Status)
	assert.Equal(t, schema.NodeStatusExited, stats[5].Status)
	assert.Equal(t, schema.NodeStatusOffline, stats[6].Status)
	assert.Equal(t, schema.NodeStatusOnline, stats[7].Status)
}

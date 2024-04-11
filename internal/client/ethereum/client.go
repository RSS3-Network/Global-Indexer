package ethereum

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sourcegraph/conc/pool"
)

type MultiChainClient struct {
	chainMap map[uint64]*ethclient.Client
	locker   sync.RWMutex
}

func (m *MultiChainClient) Put(chainID uint64, ethereumClient *ethclient.Client) {
	m.locker.Lock()
	defer m.locker.Unlock()

	m.chainMap[chainID] = ethereumClient
}

func (m *MultiChainClient) Get(chainID uint64) (*ethclient.Client, error) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	ethereumClient, found := m.chainMap[chainID]

	if !found {
		return nil, fmt.Errorf("client with chain id %d not found", chainID)
	}

	return ethereumClient, nil
}

func Dial(ctx context.Context, endpoints []string) (*MultiChainClient, error) {
	client := MultiChainClient{
		chainMap: make(map[uint64]*ethclient.Client),
	}

	contextPool := pool.New().WithContext(ctx).WithFirstError().WithCancelOnError()

	for _, endpoint := range endpoints {
		endpoint := endpoint

		contextPool.Go(func(ctx context.Context) error {
			ethereumClient, err := ethclient.DialContext(ctx, endpoint)
			if err != nil {
				return fmt.Errorf("dial to endpoint: %w", err)
			}

			chainID, err := ethereumClient.ChainID(ctx)
			if err != nil {
				return fmt.Errorf("get chain id: %w", err)
			}

			client.Put(chainID.Uint64(), ethereumClient)

			return nil
		})
	}

	if err := contextPool.Wait(); err != nil {
		return nil, err
	}

	return &client, nil
}

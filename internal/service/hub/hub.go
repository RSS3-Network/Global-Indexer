package hub

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/geolite2"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/contract/l1"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/cache"
	"github.com/rss3-network/global-indexer/internal/client/ethereum"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/config/flag"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/nameresolver"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/nta"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type Hub struct {
	dsl *dsl.DSL
	nta *nta.NTA
}

var _ echo.Validator = (*Validator)(nil)

var defaultValidator = &Validator{
	validate: validator.New(),
}

type Validator struct {
	validate *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func NewHub(ctx context.Context, databaseClient database.Client, redisClient *redis.Client, ethereumMultiChainClient *ethereum.MultiChainClient, geoLite2 *geolite2.Client, nameService *nameresolver.NameResolver, httpClient httputil.Client, txManager *txmgr.SimpleTxManager, config *config.File) (*Hub, error) {
	chainL2ID := viper.GetUint64(flag.KeyChainIDL2)

	ethereumClient, err := ethereumMultiChainClient.Get(chainL2ID)
	if err != nil {
		return nil, fmt.Errorf("get ethereum client: %w", err)
	}

	stakingV2MulticallClient, err := l2.NewStakingV2MulticallClient(chainL2ID, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new staking v2 multicall client: %w", err)
	}

	contractAddresses := l2.ContractMap[chainL2ID]
	if contractAddresses == nil {
		return nil, fmt.Errorf("contract address not found for chain id: %d", chainL2ID)
	}

	networkParamsContract, err := l2.NewNetworkParams(contractAddresses.AddressNetworkParamsProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new network contract: %w", err)
	}

	cacheClient := cache.New(redisClient)

	dslService, err := dsl.NewDSL(ctx, databaseClient, cacheClient, nameService, stakingV2MulticallClient, networkParamsContract, httpClient, txManager, config.Settler, new(big.Int).SetUint64(chainL2ID))
	if err != nil {
		return nil, fmt.Errorf("new dsl: %w", err)
	}

	contractGovernanceToken := lo.Must(bindings.NewGovernanceToken(l2.AddressGovernanceTokenProxy, ethereumClient))

	chainL1ID := viper.GetUint64(flag.KeyChainIDL1)

	ethereumL1Client, err := ethereumMultiChainClient.Get(chainL1ID)
	if err != nil {
		return nil, fmt.Errorf("get ethereum l1 client: %w", err)
	}

	erc20TokenMap := map[common.Address]*bindings.GovernanceToken{
		l1.ContractMap[chainL1ID].AddressGovernanceTokenProxy: lo.Must(bindings.NewGovernanceToken(l1.ContractMap[chainL1ID].AddressGovernanceTokenProxy, ethereumL1Client)),
		l1.ContractMap[chainL1ID].AddressUSDCToken:            lo.Must(bindings.NewGovernanceToken(l1.ContractMap[chainL1ID].AddressUSDCToken, ethereumL1Client)),
		l1.ContractMap[chainL1ID].AddressUSDTToken:            lo.Must(bindings.NewGovernanceToken(l1.ContractMap[chainL1ID].AddressUSDTToken, ethereumL1Client)),
		l1.ContractMap[chainL1ID].AddressWETHToken:            lo.Must(bindings.NewGovernanceToken(l1.ContractMap[chainL1ID].AddressWETHToken, ethereumL1Client)),
		l2.ContractMap[chainL2ID].AddressPowerToken:           lo.Must(bindings.NewGovernanceToken(l2.ContractMap[chainL2ID].AddressPowerToken, ethereumClient)),
	}

	return &Hub{
		dsl: dslService,
		nta: nta.NewNTA(ctx, config, databaseClient, stakingV2MulticallClient, networkParamsContract, contractGovernanceToken, geoLite2, cacheClient, httpClient, erc20TokenMap, chainL1ID, chainL2ID),
	}, nil
}

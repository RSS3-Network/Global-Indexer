package integrator

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cache"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cronjob"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
)

var _ service.Server = (*server)(nil)

var Name = "sort"

type server struct {
	cronJob            *cronjob.CronJob
	stakingContract    *l2.Staking
	settlementContract *l2.Settlement
	databaseClient     database.Client
	cacheClient        cache.Client
}

func (s *server) Spec() string {
	return "0 */10 * * * *"
}

func (s *server) Run(ctx context.Context) error {
	err := s.cronJob.AddFunc(ctx, s.Spec(), func() {
		if err := s.sortNodes(ctx); err != nil {
			zap.L().Error("sort nodes error", zap.Error(err))
			return
		}
	})
	if err != nil {
		return fmt.Errorf("add detector cron job: %w", err)
	}

	s.cronJob.Start()
	defer s.cronJob.Stop()

	stopchan := make(chan os.Signal, 1)

	signal.Notify(stopchan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-stopchan

	return nil
}

func (s *server) sortNodes(ctx context.Context) error {
	var limit = 100

	epoch, err := s.settlementContract.CurrentEpoch(&bind.CallOpts{})

	if err != nil {
		return fmt.Errorf("get current epoch: %w", err)
	}

	query := &schema.StatQuery{
		Limit: &limit,
	}

	for first := true; query.Cursor != nil || first; first = false {
		stats, err := s.databaseClient.FindNodeStats(ctx, query)

		if err != nil {
			return err
		}

		statsPool := pool.New().
			WithContext(ctx).
			WithCancelOnError().
			WithFirstError()

		for _, stat := range stats {
			stat := stat

			statsPool.Go(func(_ context.Context) error {
				if err = s.updateNodeEpochStats(stat, epoch.Int64()); err != nil {
					return err
				}

				return s.updateNodePoints(stat)
			})
		}

		if err := statsPool.Wait(); err != nil {
			return fmt.Errorf("wait stats pool: %w", err)
		}

		if err = s.databaseClient.SaveNodeStats(ctx, stats); err != nil {
			return err
		}

		if len(stats) == 0 {
			break
		}

		lastStat, _ := lo.Last(stats)
		query.Cursor = lo.ToPtr(lastStat.Address.String())
	}

	return s.updateNodeCache(ctx)
}

func (s *server) updateNodeCache(ctx context.Context) error {
	rssNodes, err := s.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		IsRssNode:    lo.ToPtr(true),
		ValidRequest: lo.ToPtr(model.DefaultSlashCount),
		Limit:        lo.ToPtr(model.DefaultNodeCount),
	})

	if err != nil {
		return err
	}

	if err = s.setNodeCache(ctx, model.RssNodeCacheKey, rssNodes); err != nil {
		return err
	}

	fullNodes, err := s.databaseClient.FindNodeStats(ctx, &schema.StatQuery{
		IsFullNode:   lo.ToPtr(true),
		ValidRequest: lo.ToPtr(model.DefaultSlashCount),
		Limit:        lo.ToPtr(model.DefaultNodeCount),
	})
	if err != nil {
		return err
	}

	return s.setNodeCache(ctx, model.FullNodeCacheKey, fullNodes)
}

func (s *server) updateNodeEpochStats(stat *schema.Stat, epoch int64) error {
	nodeInfo, err := s.stakingContract.GetNode(&bind.CallOpts{}, stat.Address)

	if err != nil {
		return fmt.Errorf("get node info: %s,%w", stat.Address.String(), err)
	}

	stat.Staking = float64(nodeInfo.StakingPoolTokens.Uint64())

	if epoch != stat.Epoch {
		stat.EpochRequest = 0
		stat.EpochInvalidRequest = 0
		stat.Epoch = epoch
	}

	return nil
}

func (s *server) updateNodePoints(stat *schema.Stat) error {
	node, err := s.databaseClient.FindNode(context.Background(), stat.Address)

	if err != nil {
		return fmt.Errorf("find node: %s, %w", stat.Address.String(), err)
	}

	if node.Status == schema.NodeStatusOffline {
		stat.ResetAt = time.Now()
		stat.EpochInvalidRequest = int64(model.DefaultSlashCount)

		return nil
	}

	s.calcPoints(stat)

	return nil
}

func (s *server) setNodeCache(ctx context.Context, key string, stats []*schema.Stat) error {
	nodesCache := lo.Map(stats, func(n *schema.Stat, _ int) model.Cache {
		return model.Cache{Address: n.Address.String(), Endpoint: n.Endpoint}
	})

	if err := s.cacheClient.Set(ctx, key, nodesCache); err != nil {
		return fmt.Errorf("set nodes to cache: %s, %w", key, err)
	}

	return nil
}

func (s *server) calcPoints(stat *schema.Stat) {
	// staking pool tokens
	stat.Points = math.Min(math.Log2(stat.Staking/100000+1), 0.2)

	// public good
	stat.Points += float64(lo.Ternary(stat.IsPublicGood, 0, 1))

	// running time
	stat.Points += math.Min(math.Ceil(time.Since(stat.ResetAt).Hours()/18)/120, 0.3)

	// total requests
	stat.Points += math.Min(math.Log(float64(stat.TotalRequest)/100000+1)/math.Log(100), 0.3)

	// epoch requests
	stat.Points += math.Min(math.Log(float64(stat.EpochRequest)/1000000+1)/math.Log(5000), 1)

	// network count
	stat.Points += 0.1*float64(stat.DecentralizedNetwork+stat.FederatedNetwork) + 0.3*float64(lo.Ternary(stat.IsRssNode, 1, 0))

	// indexer count
	stat.Points += math.Min(float64(stat.Indexer)*0.05, 0.2)

	// epoch failure requests
	stat.Points -= 0.5 * float64(stat.EpochInvalidRequest)
}

func New(databaseClient database.Client, redis *redis.Client, config *config.File) (service.Server, error) {
	ethereumClient, err := ethclient.Dial(config.RSS3Chain.EndpointL2)
	if err != nil {
		return nil, fmt.Errorf("dial ethereum client: %w", err)
	}

	chainID, err := ethereumClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}

	contractAddresses := l2.ContractMap[chainID.Uint64()]
	if contractAddresses == nil {
		return nil, fmt.Errorf("contract address not found for chain id: %s", chainID.String())
	}

	stakingContract, err := l2.NewStaking(contractAddresses.AddressStakingProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new staking contract: %w", err)
	}

	settlementContract, err := l2.NewSettlement(contractAddresses.AddressSettlementProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new settlement contract: %w", err)
	}

	instance := server{
		databaseClient:     databaseClient,
		cacheClient:        cache.New(redis),
		stakingContract:    stakingContract,
		settlementContract: settlementContract,
		cronJob:            cronjob.New(redis, Name, 10*time.Second),
	}

	return &instance, nil
}

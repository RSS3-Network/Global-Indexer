package job

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cache"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/provider/node"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
)

var _ Job = (*SortNodesJob)(nil)

type SortNodesJob struct {
	databaseClient  database.Client
	stakingContract *l2.Staking
}

func (job *SortNodesJob) Name() string {
	return "sort_node_job"
}

func (job *SortNodesJob) Spec() string {
	return "* */5 * * * *" // 5 minutes
}

func (job *SortNodesJob) Timeout() time.Duration {
	return 10 * time.Minute
}

func (job *SortNodesJob) Run() error {
	return job.sortNodesTask()
}

func NewSortNodesJob(databaseClient database.Client, stakingContract *l2.Staking) Job {
	return &SortNodesJob{
		databaseClient:  databaseClient,
		stakingContract: stakingContract,
	}
}

func (job *SortNodesJob) sortNodesTask() error {
	var (
		stats []*schema.Stat

		err error
	)

	ctx := context.Background()

	stats, err = job.databaseClient.FindNodeStats(ctx, []common.Address{})

	if err != nil {
		return err
	}

	// TODO: parallel
	for _, stat := range stats {
		if err = job.updateNodeEpochStats(stat); err != nil {
			return err
		}

		calcPoints(stat)

		if err = job.databaseClient.SaveNodeStat(ctx, stat); err != nil {
			return err
		}
	}

	// Update node cache.
	rssNodes, err := job.databaseClient.FindNodeStatsByType(ctx, nil, lo.ToPtr(true), 3)

	if err != nil {
		return err
	}

	if err = setNodeCache(ctx, node.RssNodeCacheKey, rssNodes); err != nil {
		return err
	}

	fullNodes, err := job.databaseClient.FindNodeStatsByType(ctx, lo.ToPtr(true), nil, 3)

	if err != nil {
		return err
	}

	if err = setNodeCache(ctx, node.FullNodeCacheKey, fullNodes); err != nil {
		return err
	}

	return nil
}

func (job *SortNodesJob) updateNodeEpochStats(stat *schema.Stat) error {
	nodeInfo, err := job.stakingContract.GetNode(&bind.CallOpts{}, stat.Address)

	if err != nil {
		return fmt.Errorf("get node info: %s,%w", stat.Address.String(), err)
	}

	stat.Staking = float64(nodeInfo.StakingPoolTokens.Uint64())
	stat.EpochRequest = 0
	stat.EpochInvalidRequest = 0

	return nil
}

func setNodeCache(ctx context.Context, key string, stats []*schema.Stat) error {
	nodesCache := lo.Map(stats, func(n *schema.Stat, _ int) node.Cache {
		return node.Cache{Address: n.Address.String(), Endpoint: n.Endpoint}
	})

	if err := cache.Set(ctx, key, nodesCache); err != nil {
		return fmt.Errorf("set nodes to cache: %s, %w", key, err)
	}

	return nil
}

// calculation rule https://docs.google.com/spreadsheets/d/1N7zEwUooiOjCIHzhoHuf8aM_lbF5bS0ZC-4luxc2qNU/edit?pli=1#gid=0
func calcPoints(stat *schema.Stat) {
	// staking pool tokens
	stat.Points = math.Min(math.Log2(stat.Staking/100000)+1, 0.2)

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

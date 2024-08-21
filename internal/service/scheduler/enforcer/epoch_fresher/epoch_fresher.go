package epochfresher

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/contract/l2"
	stakingv2 "github.com/rss3-network/global-indexer/contract/l2/staking/v2"
	"github.com/rss3-network/global-indexer/internal/cronjob"
	"github.com/rss3-network/global-indexer/internal/service"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/enforcer"
	"go.uber.org/zap"
)

var _ service.Server = (*server)(nil)

var Name = "epoch_fresher"

type server struct {
	cronJob                   *cronjob.CronJob
	blockNumber               uint64
	simpleEnforcer            *enforcer.SimpleEnforcer
	stakingContract           *stakingv2.Staking
	settlementContract        *l2.Settlement
	ethereumClient            *ethclient.Client
	settlementContractAddress common.Address
}

func (s *server) Name() string {
	return Name
}

func (s *server) Run(ctx context.Context) error {
	retryableFunc := func() error {
		for {
			err := s.process(ctx)
			if err != nil {
				return err
			}
		}
	}

	onRetry := retry.OnRetry(func(n uint, err error) {
		if !errors.Is(ctx.Err(), context.Canceled) {
			zap.L().Error("run process", zap.Error(err), zap.Uint("attempts", n))
		}
	})

	return retry.Do(retryableFunc, retry.Context(ctx), retry.DelayType(retry.FixedDelay), retry.Delay(time.Second), retry.Attempts(30), onRetry)
}

func (s *server) process(ctx context.Context) error {
	// Load latest finalized block number from RPC.
	latestFinalizedBlock, err := s.ethereumClient.BlockByNumber(ctx, big.NewInt(rpc.FinalizedBlockNumber.Int64()))
	if err != nil {
		zap.L().Error("get latest finalized block from rpc", zap.Any("server", s.Name()), zap.Error(err))

		return err
	}

	blockEnd := latestFinalizedBlock.NumberU64()

	// If the block number matches the previous one, wait until a new block is minted.
	if blockEnd <= s.blockNumber {
		blockConfirmationTime := 10 * time.Second
		zap.L().Info(
			"waiting for a new block to be minted",
			zap.Uint64("block.number.local", s.blockNumber),
			zap.Uint64("block.number.latest", blockEnd),
			zap.Duration("block.confirmationTime", blockConfirmationTime),
		)

		timer := time.NewTimer(blockConfirmationTime)
		<-timer.C

		return nil
	}

	// Fetch logs from the previous block number to the latest block number.
	logs, err := s.fetchLogs(ctx, s.blockNumber, blockEnd)

	if err != nil {
		zap.L().Error("fetch logs", zap.Any("server", s.Name()), zap.Error(err))

		return err
	}

	// Presence of logs indicates the start of a new epoch.
	if len(logs) > 0 {
		if err = s.processLogs(ctx, logs); err != nil {
			zap.L().Error("process logs", zap.Any("server", s.Name()), zap.Error(err))

			return err
		}

		// Wait for a new epoch
		newEpochWaitTime := 17 * time.Hour

		zap.L().Info(
			"waiting for a new epoch",
			zap.Duration("newEpochWaitTime", newEpochWaitTime),
		)

		timer := time.NewTimer(newEpochWaitTime)
		<-timer.C
	}

	s.blockNumber = blockEnd

	return nil
}

func (s *server) fetchLogs(ctx context.Context, blockStart, blockEnd uint64) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			s.settlementContractAddress,
		},
		FromBlock: new(big.Int).SetUint64(blockStart),
		ToBlock:   new(big.Int).SetUint64(blockEnd),
		Topics: [][]common.Hash{
			{
				l2.EventHashStakingV1RewardDistributed,
			},
		},
	}

	logs, err := s.ethereumClient.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("filter logs: %w", err)
	}

	sort.SliceStable(logs, func(i, j int) bool {
		return logs[i].BlockNumber > logs[j].BlockNumber
	})

	return logs, nil
}

// processLogs processes the logs to maintain the epoch data.
func (s *server) processLogs(ctx context.Context, logs []types.Log) error {
	// Retrieve the current epoch from the settlement contract.
	currentEpoch, err := s.settlementContract.CurrentEpoch(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("get current epoch: %w", err)
	}

	// Parse the RewardDistributed event to get the epoch.
	event, err := s.stakingContract.ParseRewardDistributed(logs[0])
	if err != nil {
		return fmt.Errorf("parse RewardDistributed event: %w", err)
	}

	if currentEpoch.Cmp(event.Epoch) == 0 {
		if err = s.simpleEnforcer.MaintainEpochData(ctx, event.Epoch.Int64()); err != nil {
			return err
		}

		zap.L().Info("maintain new epoch data completed", zap.Int64("epoch", event.Epoch.Int64()))
	}

	return nil
}

func New(redis *redis.Client, ethereumClient *ethclient.Client, blockNumber uint64, simpleEnforcer *enforcer.SimpleEnforcer, stakingContract *stakingv2.Staking, settlementContract *l2.Settlement, settlementContractAddress common.Address) service.Server {
	return &server{
		cronJob:                   cronjob.New(redis, Name, 1*time.Minute),
		blockNumber:               blockNumber,
		settlementContract:        settlementContract,
		stakingContract:           stakingContract,
		simpleEnforcer:            simpleEnforcer,
		ethereumClient:            ethereumClient,
		settlementContractAddress: settlementContractAddress,
	}
}

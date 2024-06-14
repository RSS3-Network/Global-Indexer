package epochfresher

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/contract/l2"
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
	stakingContract           *l2.Staking
	settlementContract        *l2.Settlement
	ethereumClient            *ethclient.Client
	settlementContractAddress common.Address
}

func (s *server) Name() string {
	return Name
}

func (s *server) Run(ctx context.Context) error {
	for {
		blockEnd, err := s.ethereumClient.BlockNumber(ctx)
		if err != nil {
			return fmt.Errorf("get block number: %w", err)
		}

		if blockEnd <= s.blockNumber {
			blockConfirmationTime := time.Second
			zap.L().Info(
				"waiting for a new block to be minted",
				zap.Uint64("block.number.local", s.blockNumber),
				zap.Uint64("block.number.latest", blockEnd),
				zap.Duration("block.confirmationTime", blockConfirmationTime),
			)

			timer := time.NewTimer(blockConfirmationTime)

			<-timer.C

			continue
		}

		logs, err := s.fetchLogs(ctx, s.blockNumber, blockEnd)

		if err != nil {
			return fmt.Errorf("fetch logs: %w", err)
		}

		if len(logs) > 0 {
			currentEpoch, err := s.settlementContract.CurrentEpoch(&bind.CallOpts{})

			if err != nil {
				return fmt.Errorf("get current epoch: %w", err)
			}

			event, err := s.stakingContract.ParseRewardDistributed(logs[0])
			if err != nil {
				return fmt.Errorf("parse RewardDistributed event: %w", err)
			}

			if currentEpoch.Cmp(event.Epoch) == 0 {
				if err = s.simpleEnforcer.MaintainNewEpochData(ctx, event.Epoch.Int64()); err != nil {
					return err
				}

				zap.L().Info("maintain new epoch data completed", zap.Int64("epoch", event.Epoch.Int64()))
			}
		}

		s.blockNumber = blockEnd
	}
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
				l2.EventHashStakingRewardDistributed,
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

func New(redis *redis.Client, ethereumClient *ethclient.Client, blockNumber uint64, simpleEnforcer *enforcer.SimpleEnforcer, stakingContract *l2.Staking, settlementContract *l2.Settlement, settlementContractAddress common.Address) service.Server {
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

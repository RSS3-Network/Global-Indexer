package l2

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/shopspring/decimal"
)

func (s *server) indexStakingLog(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	switch eventHash := log.Topics[0]; eventHash {
	case l2.EventHashStakingStaked:
		return s.indexStakingStakedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	case l2.EventHashStakingUnstakeClaimed:
		return s.indexStakingUnstakeClaimedLog(ctx, header, transaction, receipt, log, databaseTransaction)
	default:
		// Discard all unsupported events.
		// l2.EventHashStakingDeposited
		// l2.EventHashStakingWithdrawRequested
		// l2.EventHashStakingWithdrawalClaimed
		return nil
	}
}

func (s *server) indexStakingStakedLog(ctx context.Context, _ *types.Header, _ *types.Transaction, _ *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	event, err := s.contractStaking.ParseStaked(*log)
	if err != nil {
		return fmt.Errorf("parse Staked event: %w", err)
	}

	stakeStaker, err := databaseTransaction.FindStakeStaker(ctx, event.User, event.NodeAddr)
	if err != nil {
		return fmt.Errorf("find stake staker: %w", err)
	}

	stakeStaker.Value = stakeStaker.Value.Add(decimal.NewFromBigInt(event.Amount, 0))

	return databaseTransaction.SaveStakeStaker(ctx, stakeStaker)
}

func (s *server) indexStakingUnstakeClaimedLog(ctx context.Context, _ *types.Header, _ *types.Transaction, _ *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	event, err := s.contractStaking.ParseUnstakeClaimed(*log)
	if err != nil {
		return fmt.Errorf("parse UnstakeClaimed event: %w", err)
	}

	stakeStaker, err := databaseTransaction.FindStakeStaker(ctx, event.User, event.NodeAddr)
	if err != nil {
		return fmt.Errorf("find stake staker: %w", err)
	}

	stakeStaker.Value = stakeStaker.Value.Sub(decimal.NewFromBigInt(event.UnstakeAmount, 0))

	return databaseTransaction.SaveStakeStaker(ctx, stakeStaker)
}

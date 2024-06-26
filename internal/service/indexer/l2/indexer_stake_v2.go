package l2

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rss3-network/global-indexer/internal/database"
)

func (s *server) indexStakingV2Log(_ context.Context, _ *types.Header, _ *types.Transaction, _ *types.Receipt, _ *types.Log, _ database.Client) error {
	panic("implement me")
}

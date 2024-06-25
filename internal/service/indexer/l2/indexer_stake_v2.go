package l2

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rss3-network/global-indexer/internal/database"
)

func (s *server) indexStakingV2Log(ctx context.Context, header *types.Header, transaction *types.Transaction, receipt *types.Receipt, log *types.Log, databaseTransaction database.Client) error {
	panic("implement me")
}

package dialer

import (
	"context"
	"fmt"

	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/database/dialer/cockroachdb"
)

func Dial(ctx context.Context, config *config.Database) (database.Client, error) {
	switch config.Driver {
	case database.DriverCockroachDB:
		return cockroachdb.Dial(ctx, config.URI)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", config.Driver)
	}
}

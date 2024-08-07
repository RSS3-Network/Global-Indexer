package provider

import (
	"context"
	"fmt"

	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/database/dialer"
)

func ProvideDatabaseClient(configFile *config.File) (database.Client, error) {
	databaseClient, err := dialer.Dial(context.TODO(), configFile.Database)
	if err != nil {
		return nil, fmt.Errorf("dial to database: %w", err)
	}

	if err := databaseClient.Migrate(context.TODO()); err != nil {
		return nil, fmt.Errorf("mrigate database: %w", err)
	}

	return databaseClient, nil
}

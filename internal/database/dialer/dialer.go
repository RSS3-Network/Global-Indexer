package dialer

import (
	"context"
	"fmt"
)

func Dial(ctx context.Context, config *database.Config) (database.Client, error) {
	switch config.Driver {
	case database.DriverCockroachDB:
		return cockroachdb.Dial(ctx, config.URI, *config.Partition)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", config.Driver)
	}
}

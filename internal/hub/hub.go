package hub

import (
	"context"
	"github.com/naturalselectionlabs/global-indexer/internal/database"
)

type Hub struct {
	databaseClient database.Client
}

func NewHub(_ context.Context, databaseClient database.Client) *Hub {
	return &Hub{
		databaseClient: databaseClient,
	}
}

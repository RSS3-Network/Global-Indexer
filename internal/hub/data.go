package hub

import (
	"context"
	"github.com/naturalselectionlabs/global-indexer/schema"
)

func (h *Hub) registerNode(ctx context.Context, request *RegisterNodeRequest) error {
	node := &schema.Node{
		Address:  request.Address,
		Endpoint: request.Endpoint,
		Stream:   request.Stream,
		Config:   request.Config,
	}

	// Query node from chain
	// TODO

	// Save node to database
	return h.databaseClient.SaveNode(ctx, node)
}

package enforcer

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rss3-network/node/schema/worker/decentralized"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
)

// getNodeWorkerStatus retrieves the worker status for the node.
func (e *SimpleEnforcer) getNodeWorkerStatus(ctx context.Context, endpoint, accessToken string) (*WorkersStatusResponse, error) {
	fullURL := endpoint + "/workers_status"

	body, err := e.httpClient.FetchWithMethod(ctx, http.MethodGet, fullURL, accessToken, nil)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	response := &WorkersStatusResponse{}

	if err = json.Unmarshal(data, response); err != nil {
		return nil, err
	}

	// Set the platform for the Farcaster network.
	for i, w := range response.Data.Decentralized {
		if w.Network == network.Farcaster {
			response.Data.Decentralized[i].Platform = decentralized.PlatformFarcaster
			response.Data.Decentralized[i].Tags = []tag.Tag{tag.Social}
		}
	}

	return response, nil
}

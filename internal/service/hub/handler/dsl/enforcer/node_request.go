package enforcer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/rss3-network/node/v2/schema/worker/decentralized"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
)

// getRSSHubNodeStatus retrieves the RSSHub node status.
func (e *SimpleEnforcer) getRSSHubNodeStatus(ctx context.Context, endpoint, accessToken string) (bool, error) {
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return false, fmt.Errorf("invalid RSS endpoint: %w", err)
	}

	baseURL.Path = path.Join(baseURL.Path, "healthz")
	if accessToken != "" {
		query := baseURL.Query()
		query.Set("key", accessToken)
		baseURL.RawQuery = query.Encode()
	}

	body, _, err := e.httpClient.FetchWithMethod(ctx, http.MethodGet, baseURL.String(), "", nil)
	if err != nil {
		return false, fmt.Errorf("failed to fetch RSS healthz: %w", err)
	}
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}

	return strings.Contains(strings.ToLower(strings.TrimSpace(string(data))), "ok"), nil
}

// getNodeWorkerStatus retrieves the worker status for the node.
func (e *SimpleEnforcer) getNodeWorkerStatus(ctx context.Context, versionStr, endpoint, accessToken string) (*WorkersStatusResponse, error) {
	curVersion, _ := version.NewVersion(versionStr)

	var prefix string
	if minVersion, _ := version.NewVersion("1.1.2"); curVersion.GreaterThanOrEqual(minVersion) {
		prefix = "operators/"
	}

	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}

	fullURL := endpoint + prefix + "workers_status"

	body, _, err := e.httpClient.FetchWithMethod(ctx, http.MethodGet, fullURL, accessToken, nil)
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

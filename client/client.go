package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/naturalselectionlabs/global-indexer/internal/hub"
)

type Client struct {
	endpoint   *url.URL
	httpClient *http.Client
}

type RegisterNodeRequest hub.RegisterNodeRequest

func (c *Client) RegisterNode(ctx context.Context, request RegisterNodeRequest) error {
	var response string

	if err := c.sendRequest(ctx, "/nodes/register", nil, request, &response); err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	return nil
}

type HeartbeatRequest hub.NodeHeartbeatRequest

func (c *Client) NodeHeartbeat(ctx context.Context, request HeartbeatRequest) error {
	var response string

	if err := c.sendRequest(ctx, "/nodes/heartbeat", nil, request, &response); err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	return nil
}

func (c *Client) sendRequest(ctx context.Context, path string, values url.Values, body any, result any) error {
	internalURL := *c.endpoint
	internalURL.Path = path
	internalURL.RawQuery = values.Encode()

	requestBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, internalURL.String(), bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func NewClient(endpoint string) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		endpoint:   u,
		httpClient: http.DefaultClient,
	}, nil
}

package node

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/form/v4"
)

var _ Builder = (*builder)(nil)

type Builder interface {
	GetRSSHubPath(param, query string, nodes []Cache) (map[common.Address]string, error)
	GetActivityByIDPath(query ActivityRequest, nodes []Cache) (map[common.Address]string, error)
	GetAccountActivitiesPath(query AccountActivitiesRequest, nodes []Cache) (map[common.Address]string, error)
}

type builder struct {
	encoder *form.Encoder
}

func (c *builder) GetRSSHubPath(param, query string, nodes []Cache) (map[common.Address]string, error) {
	endpointMap, err := c.buildPath(fmt.Sprintf("/rss/%s?%s", param, query), nil, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

func (c *builder) GetActivityByIDPath(query ActivityRequest, nodes []Cache) (map[common.Address]string, error) {
	endpointMap, err := c.buildPath(fmt.Sprintf("/decentralized/tx/%s", query.ID), query, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

func (c *builder) GetAccountActivitiesPath(query AccountActivitiesRequest, nodes []Cache) (map[common.Address]string, error) {
	endpointMap, err := c.buildPath(fmt.Sprintf("/decentralized/%s", query.Account), query, nodes)
	if err != nil {
		return nil, fmt.Errorf("build path: %w", err)
	}

	return endpointMap, nil
}

func (c *builder) buildPath(path string, query any, nodes []Cache) (map[common.Address]string, error) {
	if query != nil {
		values, err := c.encoder.Encode(query)

		if err != nil {
			return nil, fmt.Errorf("build params %w", err)
		}

		path = fmt.Sprintf("%s?%s", path, values.Encode())
	}

	urls := make(map[common.Address]string, len(nodes))

	for _, node := range nodes {
		fullURL, err := url.JoinPath(node.Endpoint, path)
		if err != nil {
			return nil, fmt.Errorf("failed to join path for node %s: %w", node.Address, err)
		}

		decodedURL, err := url.QueryUnescape(fullURL)
		if err != nil {
			return nil, fmt.Errorf("failed to unescape url for node %s: %w", node.Address, err)
		}

		urls[common.HexToAddress(node.Address)] = decodedURL
	}

	return urls, nil
}

func NewPathBuilder() Builder {
	return &builder{
		encoder: form.NewEncoder(),
	}
}

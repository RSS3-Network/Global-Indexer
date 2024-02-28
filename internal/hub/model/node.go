package model

import (
	"fmt"
	"net/url"

	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

func NewNode(node *schema.Node, baseURL url.URL) *schema.Node {
	if node.Avatar != nil {
		node.Avatar.Image = baseURL.JoinPath(fmt.Sprintf("/nodes/%s/avatar.svg", node.Address)).String()
	}

	return node
}

func NewNodes(nodes []*schema.Node, baseURL url.URL) []*schema.Node {
	nodeModels := make([]*schema.Node, 0, len(nodes))
	for _, node := range nodes {
		nodeModels = append(nodeModels, NewNode(node, baseURL))
	}

	return nodeModels
}

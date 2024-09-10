package nta

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
)

type NodeRequest struct {
	Address common.Address `param:"node_address" validate:"required"`
}

type BatchNodeRequest struct {
	Cursor        *string          `query:"cursor"`
	Limit         int              `query:"limit" validate:"min=1,max=50" default:"10"`
	NodeAddresses []common.Address `query:"node_addresses"`
}

type NodeResponseData *schema.Node

type NodesResponseData []*schema.Node

func NewNode(node *schema.Node, baseURL url.URL) NodeResponseData {
	if node.Avatar != nil {
		node.Avatar.Image = baseURL.JoinPath(fmt.Sprintf("/nta/nodes/%s/avatar.svg", node.Address)).String()
	}

	if node.HideTaxRate {
		node.TaxRateBasisPoints = nil
	}

	return node
}

func NewNodes(nodes []*schema.Node, baseURL url.URL) NodesResponseData {
	nodeModels := make([]*schema.Node, 0, len(nodes))
	for _, node := range nodes {
		nodeModels = append(nodeModels, NewNode(node, baseURL))
	}

	return nodeModels
}

package nta

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
)

type NodeRequest struct {
	Address common.Address `param:"id" validate:"required"`
}

type BatchNodeRequest struct {
	Cursor      *string          `query:"cursor"`
	Limit       int              `query:"limit" validate:"min=1,max=50" default:"10"`
	NodeAddress []common.Address `query:"nodeAddress"`
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

// NodeStatusTransitionError represents an error when attempting to transition a Node to an invalid status.
type NodeStatusTransitionError struct {
	From, To schema.NodeStatus
}

// Error returns a string representation of the NodeStatusTransitionError.
// TODO: move to a more appropriate location.
func (err *NodeStatusTransitionError) Error() string {
	return fmt.Sprintf("invalid status transition from %s to %s", err.From, err.To)
}

// See https://www.figma.com/file/2PCGRBkIRuQ7VmttXyT6gB/Epoch-workflow?type=whiteboard&node-id=0-1&t=uiVv3wIktG5NAHCz-0
// for the state machine diagram.
var transitions = map[schema.NodeStatus][]schema.NodeStatus{
	schema.NodeStatusRegistered: {schema.NodeStatusRegistered, schema.NodeStatusOnline, schema.NodeStatusExited},
	schema.NodeStatusOnline:     {schema.NodeStatusOnline, schema.NodeStatusExiting, schema.NodeStatusExited, schema.NodeStatusSlashed, schema.NodeStatusOffline},
	schema.NodeStatusExiting:    {schema.NodeStatusExiting, schema.NodeStatusExited},
	schema.NodeStatusSlashed:    {schema.NodeStatusSlashed, schema.NodeStatusOnline, schema.NodeStatusOffline},
	schema.NodeStatusOffline:    {schema.NodeStatusOffline, schema.NodeStatusOnline, schema.NodeStatusExited},
	schema.NodeStatusExited:     {schema.NodeStatusExited, schema.NodeStatusRegistered},
}

func isValidTransition(from, to schema.NodeStatus) bool {
	for _, validTo := range transitions[from] {
		if to == validTo {
			return true
		}
	}

	return false
}

// UpdateNodeStatus updates the status of a given node if the transition is valid.
// It returns a NodeStatusTransitionError when the transition is invalid.
func UpdateNodeStatus(node *schema.Node, newStatus schema.NodeStatus) error {
	if isValidTransition(node.Status, newStatus) {
		node.Status = newStatus
		return nil
	}

	return &NodeStatusTransitionError{node.Status, newStatus}
}

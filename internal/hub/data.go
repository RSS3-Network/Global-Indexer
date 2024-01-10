package hub

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/global-indexer/common/ethereum"
	"github.com/naturalselectionlabs/global-indexer/schema"
	"github.com/samber/lo"
)

func (h *Hub) getNode(ctx context.Context, address common.Address) (*schema.Node, error) {
	node, err := h.databaseClient.FindNode(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("get node %s: %w", address, err)
	}

	nodeInfo, err := h.stakingContract.GetNode(&bind.CallOpts{}, address)
	if err != nil {
		return nil, fmt.Errorf("get node from chain: %w", err)
	}

	node.Name = nodeInfo.Name
	node.Description = nodeInfo.Description
	node.TaxFraction = nodeInfo.TaxFraction

	return node, nil
}

func (h *Hub) getNodes(ctx context.Context, addresses []common.Address, cursor *string) ([]*schema.Node, error) {
	nodes, err := h.databaseClient.FindNodes(ctx, addresses, cursor)
	if err != nil {
		return nil, fmt.Errorf("get nodes: %w", err)
	}

	lo.ForEach(nodes, func(node *schema.Node, _ int) {
		nodeInfo, err := h.stakingContract.GetNode(&bind.CallOpts{}, node.Address)
		if err != nil {
			return
		}

		node.Name = nodeInfo.Name
		node.Description = nodeInfo.Description
		node.TaxFraction = nodeInfo.TaxFraction
	})

	return nodes, nil
}

func (h *Hub) registerNode(ctx context.Context, request *RegisterNodeRequest) error {
	node := &schema.Node{
		Address:  request.Address,
		Endpoint: request.Endpoint,
		Stream:   request.Stream,
		Config:   request.Config,
	}

	// Check node from chain.
	nodeInfo, err := h.stakingContract.GetNode(&bind.CallOpts{}, request.Address)
	if err != nil {
		return fmt.Errorf("get node from chain: %w", err)
	}

	if nodeInfo.Account == ethereum.AddressGenesis {
		return fmt.Errorf("node: %s has not been registered on the chain", request.Address.String())
	}

	node.IsPublicGood = nodeInfo.PublicGood

	// Save node to database.
	return h.databaseClient.SaveNode(ctx, node)
}

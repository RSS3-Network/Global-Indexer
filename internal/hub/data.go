package hub

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/ethereum"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

var message = "I, %s, am signing this message for registering my intention to operate an RSS3 Node."

func (h *Hub) getNode(ctx context.Context, address common.Address) (*schema.Node, error) {
	node, err := h.databaseClient.FindNode(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("get node %s: %w", address, err)
	}

	if node == nil {
		return nil, nil
	}

	nodeInfo, err := h.stakingContract.GetNode(&bind.CallOpts{}, address)
	if err != nil {
		return nil, fmt.Errorf("get node from chain: %w", err)
	}

	node.Name = nodeInfo.Name
	node.Description = nodeInfo.Description
	node.TaxRateBasisPoints = nodeInfo.TaxRateBasisPoints
	node.OperationPoolTokens = nodeInfo.OperationPoolTokens.String()
	node.StakingPoolTokens = nodeInfo.StakingPoolTokens.String()
	node.TotalShares = nodeInfo.TotalShares.String()
	node.SlashedTokens = nodeInfo.SlashedTokens.String()

	return node, nil
}

func (h *Hub) getNodes(ctx context.Context, request *BatchNodeRequest) ([]*schema.Node, error) {
	nodes, err := h.databaseClient.FindNodes(ctx, request.NodeAddress, request.Cursor, request.Limit)
	if err != nil {
		return nil, fmt.Errorf("get nodes: %w", err)
	}

	addresses := lo.Map(nodes, func(node *schema.Node, _ int) common.Address {
		return node.Address
	})

	nodeInfo, err := h.stakingContract.GetNodes(&bind.CallOpts{}, addresses)
	if err != nil {
		return nil, fmt.Errorf("get nodes from chain: %w", err)
	}

	nodeInfoMap := lo.SliceToMap(nodeInfo, func(node l2.DataTypesNode) (common.Address, l2.DataTypesNode) {
		return node.Account, node
	})

	for _, node := range nodes {
		if nodeInfo, exists := nodeInfoMap[node.Address]; exists {
			node.Name = nodeInfo.Name
			node.Description = nodeInfo.Description
			node.TaxRateBasisPoints = nodeInfo.TaxRateBasisPoints
			node.OperationPoolTokens = nodeInfo.OperationPoolTokens.String()
			node.StakingPoolTokens = nodeInfo.StakingPoolTokens.String()
			node.TotalShares = nodeInfo.TotalShares.String()
			node.SlashedTokens = nodeInfo.SlashedTokens.String()
		}
	}

	return nodes, nil
}

func (h *Hub) register(ctx context.Context, request *RegisterNodeRequest) error {
	// Check signature.
	if err := h.checkSignature(ctx, request.Address, hexutil.MustDecode(request.Signature)); err != nil {
		return err
	}

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
		return fmt.Errorf("node: %s has not been registered on the chain", strings.ToLower(request.Address.String()))
	}

	if strings.Compare(nodeInfo.OperationPoolTokens.String(), decimal.NewFromInt(10000).Mul(decimal.NewFromInt(1e18)).String()) < 0 {
		return fmt.Errorf("insufficient operation pool tokens")
	}

	node.IsPublicGood = nodeInfo.PublicGood
	node.LastHeartbeatTimestamp = time.Now().Unix()
	node.Status = schema.StatusOnline

	// Save node to database.
	return h.databaseClient.SaveNode(ctx, node)
}

func (h *Hub) heartbeat(ctx context.Context, request *NodeHeartbeatRequest) error {
	// Check signature.
	if err := h.checkSignature(ctx, request.Address, hexutil.MustDecode(request.Signature)); err != nil {
		return fmt.Errorf("check signature: %w", err)
	}

	// Check node from database.
	node, err := h.databaseClient.FindNode(ctx, request.Address)
	if err != nil {
		return fmt.Errorf("get node %s from database: %w", request.Address, err)
	}

	if node == nil {
		return fmt.Errorf("node %s not found", request.Address)
	}

	node.LastHeartbeatTimestamp = request.Timestamp
	node.Status = schema.StatusOnline

	// Save node to database.
	return h.databaseClient.SaveNode(ctx, node)
}

func (h *Hub) checkSignature(_ context.Context, address common.Address, signature []byte) error {
	message := fmt.Sprintf(message, strings.ToLower(address.Hex()))
	data := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(data)).Bytes()

	if signature[crypto.RecoveryIDOffset] == 27 || signature[crypto.RecoveryIDOffset] == 28 {
		signature[crypto.RecoveryIDOffset] -= 27
	}

	pubKey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return fmt.Errorf("failed to parse signature: %w", err)
	}

	result := crypto.PubkeyToAddress(*pubKey)

	if address != result {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

package signer

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	client *rpc.Client
}

func NewSignerClient(endpoint string) (*Client, error) {
	var httpClient *http.Client

	rpcClient, err := rpc.DialOptions(context.Background(), endpoint, rpc.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	signer := &Client{client: rpcClient}
	// Check if reachable
	res, err := signer.pingVersion()
	if err != nil {
		return nil, err
	}

	if res != "ok" {
		return nil, fmt.Errorf("signer service unreachable")
	}

	return signer, nil
}

func (s *Client) pingVersion() (string, error) {
	var v string

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	if err := s.client.CallContext(ctx, &v, "health_status"); err != nil {
		return "", err
	}

	return v, nil
}

func (s *Client) SignTransaction(ctx context.Context, chainID *big.Int, from common.Address, tx *types.Transaction) (*types.Transaction, error) {
	args := NewTransactionArgsFromTransaction(chainID, from, tx)

	var result hexutil.Bytes
	if err := s.client.CallContext(ctx, &result, "eth_signTransaction", args); err != nil {
		return nil, fmt.Errorf("eth_signTransaction failed: %w", err)
	}

	signed := &types.Transaction{}
	if err := signed.UnmarshalBinary(result); err != nil {
		return nil, err
	}

	return signed, nil
}

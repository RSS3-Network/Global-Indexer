package crypto

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	gisigner "github.com/naturalselectionlabs/rss3-global-indexer/common/signer"
)

func PrivateKeySignerFn(key *ecdsa.PrivateKey, chainID *big.Int) bind.SignerFn {
	from := crypto.PubkeyToAddress(key.PublicKey)
	signer := types.LatestSignerForChainID(chainID)

	return func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		if address != from {
			return nil, bind.ErrNotAuthorized
		}

		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
		if err != nil {
			return nil, err
		}

		return tx.WithSignature(signer, signature)
	}
}

type SignerFn func(context.Context, common.Address, *types.Transaction) (*types.Transaction, error)

type SignerFactory func(chainID *big.Int) SignerFn

func NewSignerFactory(privateKey, endpoint, address string) (SignerFactory, common.Address, error) {
	var (
		signer      SignerFactory
		fromAddress common.Address
	)

	if endpoint != "" && address != "" {
		signerClient, err := gisigner.NewSignerClient(endpoint)
		if err != nil {
			return nil, common.Address{}, fmt.Errorf("failed to create the signer client: %w", err)
		}

		fromAddress = common.HexToAddress(address)
		signer = func(chainID *big.Int) SignerFn {
			return func(ctx context.Context, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
				if !bytes.Equal(address[:], fromAddress[:]) {
					return nil, fmt.Errorf("attempting to sign for %s, expected %s: ", address, address)
				}

				return signerClient.SignTransaction(ctx, chainID, address, tx)
			}
		}
	} else {
		var (
			privKey *ecdsa.PrivateKey

			err error
		)

		if privateKey == "" {
			return nil, common.Address{}, fmt.Errorf("at least specify a private key")
		}

		privKey, err = crypto.HexToECDSA(strings.TrimPrefix(privateKey, "0x"))
		if err != nil {
			return nil, common.Address{}, fmt.Errorf("failed to parse the private key: %w", err)
		}

		fromAddress = crypto.PubkeyToAddress(privKey.PublicKey)
		signer = func(chainID *big.Int) SignerFn {
			s := PrivateKeySignerFn(privKey, chainID)
			return func(_ context.Context, addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
				return s(addr, tx)
			}
		}
	}

	return signer, fromAddress, nil
}

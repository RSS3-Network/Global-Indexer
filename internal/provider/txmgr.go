package provider

import (
	"fmt"
	"math/big"
	"time"

	gicrypto "github.com/rss3-network/global-indexer/common/crypto"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/internal/client/ethereum"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/config/flag"
	"github.com/spf13/viper"
)

func ProvideTxManager(config *config.File, ethereumMultiChainClient *ethereum.MultiChainClient) (*txmgr.SimpleTxManager, error) {
	signerFactory, from, err := gicrypto.NewSignerFactory(config.Settler.PrivateKey, config.Settler.SignerEndpoint, config.Settler.WalletAddress)
	if err != nil {
		return nil, fmt.Errorf("create signer: %w", err)
	}

	defaultTxConfig := txmgr.Config{
		ResubmissionTimeout:       20 * time.Second,
		FeeLimitMultiplier:        5,
		TxSendTimeout:             5 * time.Minute,
		TxNotInMempoolTimeout:     1 * time.Hour,
		NetworkTimeout:            5 * time.Minute,
		ReceiptQueryInterval:      500 * time.Millisecond,
		NumConfirmations:          5,
		SafeAbortNonceTooLowCount: 3,
	}

	chainID := new(big.Int).SetUint64(viper.GetUint64(flag.KeyChainIDL2))

	ethereumClient, err := ethereumMultiChainClient.Get(chainID.Uint64())
	if err != nil {
		return nil, fmt.Errorf("load l2 ethereum client: %w", err)
	}

	txManager, err := txmgr.NewSimpleTxManager(defaultTxConfig, chainID, nil, ethereumClient, from, signerFactory(chainID))
	if err != nil {
		return nil, fmt.Errorf("create tx manager %w", err)
	}

	return txManager, nil
}

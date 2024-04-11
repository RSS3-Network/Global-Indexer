package flag

import "github.com/rss3-network/global-indexer/internal/client/ethereum"

const (
	KeyConfig = "config"
	KeyServer = "server"

	KeyChainIDL1 = "chain-id.l1"
	KeyChainIDL2 = "chain-id.l2"
)

const (
	ValueChainIDL1 = ethereum.ChainIDEthereum
	ValueChainIDL2 = ethereum.ChainIDVSL
)

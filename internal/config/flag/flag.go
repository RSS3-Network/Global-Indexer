package flag

import (
	"github.com/rss3-network/global-indexer/contract/l1"
	"github.com/rss3-network/global-indexer/contract/l2"
)

const (
	KeyConfig = "config"
	KeyServer = "server"

	KeyChainIDL1 = "chain-id.l1"
	KeyChainIDL2 = "chain-id.l2"
)

const (
	ValueChainIDL1 = l1.ChainIDMainnet
	ValueChainIDL2 = l2.ChainIDMainnet
)

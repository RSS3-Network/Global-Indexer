package database

import (
	"context"
	"database/sql"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

var (
	ErrorRowNotFound = errors.New("row not found")
)

type Client interface {
	Session
	Transaction

	FindCheckpoint(ctx context.Context, chainID uint64) (*schema.Checkpoint, error)
	SaveCheckpoint(ctx context.Context, checkpoint *schema.Checkpoint) error

	FindNode(ctx context.Context, nodeAddress common.Address) (*schema.Node, error)
	FindNodes(ctx context.Context, nodeAddresses []common.Address, cursor *string, limit int) ([]*schema.Node, error)
	SaveNode(ctx context.Context, node *schema.Node) error
	UpdateNodesStatus(ctx context.Context, lastHeartbeatTimestamp int64) error

	FindBridgeTransaction(ctx context.Context, id common.Hash) (*schema.BridgeTransaction, error)
	FindBridgeTransactions(ctx context.Context) ([]*schema.BridgeTransaction, error)
	FindBridgeTransactionsByAddress(ctx context.Context, address common.Address) ([]*schema.BridgeTransaction, error)
	FindBridgeEventsByID(ctx context.Context, id common.Hash) (*schema.BridgeEvent, error)
	FindBridgeEventsByIDs(ctx context.Context, ids []common.Hash) ([]*schema.BridgeEvent, error)
	FindNodeStat(ctx context.Context) (*schema.Stat, error)
	FindNodeStats(ctx context.Context, nodeAddresses []common.Address) ([]*schema.Stat, error)
	FindNodeStatsByType(ctx context.Context, fullNode, rssNode *bool, limit int) ([]*schema.Stat, error)
	SaveNodeStat(ctx context.Context, stat *schema.Stat) error
	SaveNodeStats(ctx context.Context, stats []*schema.Stat) error

	FindNodeIndexers(ctx context.Context, nodeAddresses []common.Address, networks, workers []string) ([]*schema.Indexer, error)
	SaveNodeIndexers(ctx context.Context, indexers []*schema.Indexer) error
	DeleteNodeIndexers(ctx context.Context, nodeAddress common.Address) error

	SaveBridgeTransaction(ctx context.Context, bridgeTransaction *schema.BridgeTransaction) error
	SaveBridgeEvent(ctx context.Context, bridgeEvent *schema.BridgeEvent) error

	FindStakeTransaction(ctx context.Context, id common.Hash) (*schema.StakeTransaction, error)
	FindStakeTransactions(ctx context.Context) ([]*schema.StakeTransaction, error)
	FindStakeTransactionsByUser(ctx context.Context, address common.Address) ([]*schema.StakeTransaction, error)
	FindStakeTransactionsByNode(ctx context.Context, address common.Address) ([]*schema.StakeTransaction, error)
	FindStakeEventsByID(ctx context.Context, id common.Hash) ([]*schema.StakeEvent, error)
	FindStakeEventsByIDs(ctx context.Context, ids []common.Hash) ([]*schema.StakeEvent, error)
	FindStakeChipsByOwner(ctx context.Context, owner common.Address) ([]*schema.StakeChip, error)
	FindStakeChipsByNode(ctx context.Context, node common.Address) ([]*schema.StakeChip, error)
	SaveStakeTransaction(ctx context.Context, stakeTransaction *schema.StakeTransaction) error
	SaveStakeEvent(ctx context.Context, stakeEvent *schema.StakeEvent) error
	SaveStakeChips(ctx context.Context, stakeChips ...*schema.StakeChip) error
	UpdateStakeChipsOwner(ctx context.Context, owner common.Address, stakeChips ...*big.Int) error
}

type Session interface {
	Migrate(ctx context.Context) error
	WithTransaction(ctx context.Context, transactionFunction func(ctx context.Context, client Client) error, transactionOptions ...*sql.TxOptions) error
	Begin(ctx context.Context, transactionOptions ...*sql.TxOptions) (Client, error)
}

type Transaction interface {
	Rollback() error
	Commit() error
}

var _ goose.Logger = (*SugaredLogger)(nil)

type SugaredLogger struct {
	Logger *zap.SugaredLogger
}

func (s SugaredLogger) Fatalf(format string, v ...interface{}) {
	s.Logger.Fatalf(format, v...)
}

func (s SugaredLogger) Printf(format string, v ...interface{}) {
	s.Logger.Infof(format, v...)
}

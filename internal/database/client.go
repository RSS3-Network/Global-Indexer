package database

import (
	"context"
	"database/sql"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrorRowNotFound = errors.New("row not found")
)

type Client interface {
	Session
	Transaction

	RollbackBlock(ctx context.Context, chainID, blockNUmber uint64) error

	FindCheckpoint(ctx context.Context, chainID uint64) (*schema.Checkpoint, error)
	SaveCheckpoint(ctx context.Context, checkpoint *schema.Checkpoint) error

	FindNode(ctx context.Context, nodeAddress common.Address) (*schema.Node, error)
	FindNodes(ctx context.Context, nodeAddresses []common.Address, status *schema.Status, cursor *string, limit int) ([]*schema.Node, error)
	FindNodeAvatar(ctx context.Context, nodeAddress common.Address) (*l2.ChipsTokenMetadata, error)
	SaveNode(ctx context.Context, node *schema.Node) error
	UpdateNodesStatus(ctx context.Context, lastHeartbeatTimestamp int64) error

	FindNodeStat(ctx context.Context, nodeAddress common.Address) (*schema.Stat, error)
	FindNodeStats(ctx context.Context, query *schema.StatQuery) ([]*schema.Stat, error)
	FindNodeIndexers(ctx context.Context, nodeAddresses []common.Address, networks, workers []string) ([]*schema.Indexer, error)
	FindNodeSnapshots(ctx context.Context) ([]*schema.NodeSnapshot, error)
	SaveNodeStat(ctx context.Context, stat *schema.Stat) error
	SaveNodeStats(ctx context.Context, stats []*schema.Stat) error
	SaveNodeIndexers(ctx context.Context, indexers []*schema.Indexer) error
	SaveNodeSnapshot(ctx context.Context, nodeSnapshot *schema.NodeSnapshot) error
	DeleteNodeIndexers(ctx context.Context, nodeAddress common.Address) error

	FindBridgeTransaction(ctx context.Context, query schema.BridgeTransactionQuery) (*schema.BridgeTransaction, error)
	FindBridgeTransactions(ctx context.Context, query schema.BridgeTransactionsQuery) ([]*schema.BridgeTransaction, error)
	FindBridgeEvents(ctx context.Context, query schema.BridgeEventsQuery) ([]*schema.BridgeEvent, error)
	SaveBridgeTransaction(ctx context.Context, bridgeTransaction *schema.BridgeTransaction) error
	SaveBridgeEvent(ctx context.Context, bridgeEvent *schema.BridgeEvent) error

	FindStakeTransaction(ctx context.Context, query schema.StakeTransactionQuery) (*schema.StakeTransaction, error)
	FindStakeTransactions(ctx context.Context, query schema.StakeTransactionsQuery) ([]*schema.StakeTransaction, error)
	FindStakeEvents(ctx context.Context, query schema.StakeEventsQuery) ([]*schema.StakeEvent, error)
	FindStakeChip(ctx context.Context, query schema.StakeChipQuery) (*schema.StakeChip, error)
	FindStakeChips(ctx context.Context, query schema.StakeChipsQuery) ([]*schema.StakeChip, error)
	FindStakeStakings(ctx context.Context, query schema.StakeStakingsQuery) ([]*schema.StakeStaking, error)
	FindStakeSnapshots(ctx context.Context) ([]*schema.StakeSnapshot, error)
	SaveStakeTransaction(ctx context.Context, stakeTransaction *schema.StakeTransaction) error
	SaveStakeEvent(ctx context.Context, stakeEvent *schema.StakeEvent) error
	SaveStakeChips(ctx context.Context, stakeChips ...*schema.StakeChip) error
	SaveStakeSnapshot(ctx context.Context, stakeSnapshot *schema.StakeSnapshot) error
	UpdateStakeChipsOwner(ctx context.Context, owner common.Address, stakeChips ...*big.Int) error

	SaveEpoch(ctx context.Context, epoch *schema.Epoch) error
	FindEpochs(ctx context.Context, limit int, cursor *string) ([]*schema.Epoch, error)
	FindEpochTransactions(ctx context.Context, id uint64, itemsLimit int, cursor *string) ([]*schema.Epoch, error)
	FindEpochTransaction(ctx context.Context, transactionHash common.Hash, itemsLimit int, cursor *string) (*schema.Epoch, error)
	FindEpochNodeRewards(ctx context.Context, nodeAddress common.Address, limit int, cursor *string) ([]*schema.Epoch, error)

	SaveEpochTrigger(ctx context.Context, epochTrigger *schema.EpochTrigger) error
	FindLatestEpochTrigger(ctx context.Context) (*schema.EpochTrigger, error)

	SaveBillingRecordDeposited(ctx context.Context, billingRecord *schema.BillingRecordDeposited) error
	SaveBillingRecordWithdrawal(ctx context.Context, billingRecord *schema.BillingRecordWithdrawal) error
	SaveBillingRecordCollected(ctx context.Context, billingRecord *schema.BillingRecordCollected) error

	PrepareBillingCollectTokens(ctx context.Context, nowTime time.Time) (*map[common.Address]schema.BillingCollectDataPerAddress, error)
	PrepareBillingWithdrawTokens(ctx context.Context) (*map[common.Address]float64, error)
	UpdateBillingRuLimit(ctx context.Context, succeededUsersWithRu map[common.Address]int64) error

	GatewayDeposit(ctx context.Context, address common.Address, ruIncrease int64) (bool, error)

	Raw() *gorm.DB
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

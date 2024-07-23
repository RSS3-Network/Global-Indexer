package database

import (
	"context"
	"database/sql"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pressly/goose/v3"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	ErrorRowNotFound = errors.New("row not found")
)

type Client interface {
	Session
	Transaction

	RollbackBlock(ctx context.Context, chainID, blockNumber uint64) error

	FindCheckpoint(ctx context.Context, chainID uint64) (*schema.Checkpoint, error)
	SaveCheckpoint(ctx context.Context, checkpoint *schema.Checkpoint) error

	FindNode(ctx context.Context, nodeAddress common.Address) (*schema.Node, error)
	FindNodes(ctx context.Context, query schema.FindNodesQuery) ([]*schema.Node, error)
	FindNodeAvatar(ctx context.Context, nodeAddress common.Address) (*l2.ChipsTokenMetadata, error)
	SaveNode(ctx context.Context, node *schema.Node) error
	UpdateNodesStatusOffline(ctx context.Context, lastHeartbeatTimestamp int64) error
	UpdateNodesHideTaxRate(ctx context.Context, nodeAddress common.Address, hideTaxRate bool) error
	UpdateNodesScore(ctx context.Context, nodes []*schema.Node) error
	UpdateNodePublicGood(ctx context.Context, nodeAddress common.Address, isPublicGood bool) error

	BatchUpdateNodes(ctx context.Context, data []*schema.BatchUpdateNode) error
	SaveNodeEvent(ctx context.Context, nodeEvent *schema.NodeEvent) error
	FindNodeEvents(ctx context.Context, nodeEventsQuery *schema.NodeEventsQuery) ([]*schema.NodeEvent, error)
	DeleteNodeEventsByBlockNumber(ctx context.Context, blockNumber uint64) error
	UpdateNodeEventsFinalizedByBlockNumber(ctx context.Context, blockNumber uint64) error

	FindNodeStat(ctx context.Context, nodeAddress common.Address) (*schema.Stat, error)
	FindNodeStats(ctx context.Context, query *schema.StatQuery) ([]*schema.Stat, error)
	SaveNodeStat(ctx context.Context, stat *schema.Stat) error
	SaveNodeStats(ctx context.Context, stats []*schema.Stat) error
	FindNodeWorkers(ctx context.Context, query *schema.WorkerQuery) ([]*schema.Worker, error)
	SaveNodeWorkers(ctx context.Context, workers []*schema.Worker) error
	UpdateNodeWorkerActive(ctx context.Context) error
	SaveNodeInvalidResponses(ctx context.Context, nodeInvalidResponses []*schema.NodeInvalidResponse) error

	FindNodeCountSnapshots(ctx context.Context) ([]*schema.NodeSnapshot, error)
	SaveNodeCountSnapshot(ctx context.Context, nodeSnapshot *schema.NodeSnapshot) error
	FindStakerCountSnapshots(ctx context.Context) ([]*schema.StakerCountSnapshot, error)
	SaveStakerCountSnapshot(ctx context.Context, stakeSnapshot *schema.StakerCountSnapshot) error
	FindStakerProfitSnapshots(ctx context.Context, query schema.StakerProfitSnapshotsQuery) ([]*schema.StakerProfitSnapshot, error)
	SaveStakerProfitSnapshots(ctx context.Context, stakerProfitSnapshots []*schema.StakerProfitSnapshot) error
	FindStakerCountRecentEpochs(ctx context.Context, recentEpochs int) (map[common.Address]*schema.StakeRecentCount, error)
	FindOperatorProfitSnapshots(ctx context.Context, query schema.OperatorProfitSnapshotsQuery) ([]*schema.OperatorProfitSnapshot, error)
	SaveOperatorProfitSnapshots(ctx context.Context, operatorProfitSnapshots []*schema.OperatorProfitSnapshot) error
	SaveNodeAPYSnapshots(ctx context.Context, nodeAPYSnapshots []*schema.NodeAPYSnapshot) error
	FindEpochAPYSnapshots(ctx context.Context, query schema.EpochAPYSnapshotQuery) ([]*schema.EpochAPYSnapshot, error)
	SaveEpochAPYSnapshot(ctx context.Context, epochAPYSnapshots *schema.EpochAPYSnapshot) error
	FindEpochAPYSnapshotsAverage(ctx context.Context) (decimal.Decimal, error)

	FindBridgeTransaction(ctx context.Context, query schema.BridgeTransactionQuery) (*schema.BridgeTransaction, error)
	FindBridgeTransactions(ctx context.Context, query schema.BridgeTransactionsQuery) ([]*schema.BridgeTransaction, error)
	FindBridgeEvents(ctx context.Context, query schema.BridgeEventsQuery) ([]*schema.BridgeEvent, error)
	UpdateBridgeTransactionsFinalizedByBlockNumber(ctx context.Context, chainID, blockNumber uint64) error
	UpdateBridgeEventsFinalizedByBlockNumber(ctx context.Context, chainID, blockNumber uint64) error
	SaveBridgeTransaction(ctx context.Context, bridgeTransaction *schema.BridgeTransaction) error
	SaveBridgeEvent(ctx context.Context, bridgeEvent *schema.BridgeEvent) error
	DeleteBridgeTransactionsByBlockNumber(ctx context.Context, chainID, blockNumber uint64) error
	DeleteBridgeEventsByBlockNumber(ctx context.Context, chainID, blockNumber uint64) error

	FindStakeTransaction(ctx context.Context, query schema.StakeTransactionQuery) (*schema.StakeTransaction, error)
	FindStakeTransactions(ctx context.Context, query schema.StakeTransactionsQuery) ([]*schema.StakeTransaction, error)
	FindStakeEvents(ctx context.Context, query schema.StakeEventsQuery) ([]*schema.StakeEvent, error)
	FindStakeChip(ctx context.Context, query schema.StakeChipQuery) (*schema.StakeChip, error)
	FindStakeChips(ctx context.Context, query schema.StakeChipsQuery) ([]*schema.StakeChip, error)
	UpdateStakeTransactionsFinalizedByBlockNumber(ctx context.Context, blockNumber uint64) error
	UpdateStakeEventsFinalizedByBlockNumber(ctx context.Context, blockNumber uint64) error
	UpdateStakeChipsFinalizedByBlockNumber(ctx context.Context, blockNumber uint64) error
	DeleteStakeChipsByBlockNumber(ctx context.Context, blockNumber uint64) error
	FindStakeStakings(ctx context.Context, query schema.StakeStakingsQuery) ([]*schema.StakeStaking, error)
	SaveStakeTransaction(ctx context.Context, stakeTransaction *schema.StakeTransaction) error
	SaveStakeEvent(ctx context.Context, stakeEvent *schema.StakeEvent) error
	DeleteStakeTransactionsByBlockNumber(ctx context.Context, blockNumber uint64) error
	DeleteStakeEventsByBlockNumber(ctx context.Context, blockNumber uint64) error
	SaveStakeChips(ctx context.Context, stakeChips ...*schema.StakeChip) error
	UpdateStakeChipsOwner(ctx context.Context, owner common.Address, stakeChips ...*big.Int) error

	SaveEpoch(ctx context.Context, epoch *schema.Epoch) error
	FindEpochs(ctx context.Context, limit int, cursor *string) ([]*schema.Epoch, error)
	FindEpochTransactions(ctx context.Context, id uint64, itemsLimit int, cursor *string) ([]*schema.Epoch, error)
	FindEpochTransaction(ctx context.Context, transactionHash common.Hash, itemsLimit int, cursor *string) (*schema.Epoch, error)
	FindEpochNodeRewards(ctx context.Context, nodeAddress common.Address, limit int, cursor *string) ([]*schema.Epoch, error)

	SaveEpochTrigger(ctx context.Context, epochTrigger *schema.EpochTrigger) error
	FindLatestEpochTrigger(ctx context.Context) (*schema.EpochTrigger, error)

	FindAverageTaxSubmissions(ctx context.Context, query schema.AverageTaxRateSubmissionQuery) ([]*schema.AverageTaxRateSubmission, error)
	SaveAverageTaxSubmission(ctx context.Context, averageTaxSubmission *schema.AverageTaxRateSubmission) error
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

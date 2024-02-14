package cockroachdb

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/pressly/goose/v3"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"moul.io/zapgorm2"
)

var _ database.Client = (*client)(nil)

//go:embed migration/*.sql
var migrationFS embed.FS

type client struct {
	database *gorm.DB
}

func (c *client) Migrate(ctx context.Context) error {
	goose.SetBaseFS(migrationFS)
	goose.SetTableName("versions")
	goose.SetLogger(&database.SugaredLogger{Logger: zap.L().Sugar()})

	if err := goose.SetDialect(new(postgres.Dialector).Name()); err != nil {
		return fmt.Errorf("set migration dialect: %w", err)
	}

	connector, err := c.database.DB()
	if err != nil {
		return fmt.Errorf("get database connector: %w", err)
	}

	return goose.UpContext(ctx, connector, "migration")
}

func (c *client) WithTransaction(ctx context.Context, transactionFunction func(ctx context.Context, client database.Client) error, transactionOptions ...*sql.TxOptions) error {
	transaction, err := c.Begin(ctx, transactionOptions...)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	if err := transactionFunction(ctx, transaction); err != nil {
		_ = transaction.Rollback()

		return fmt.Errorf("execute transaction: %w", err)
	}

	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (c *client) Begin(ctx context.Context, transactionOptions ...*sql.TxOptions) (database.Client, error) {
	transaction := c.database.WithContext(ctx).Begin(transactionOptions...)
	if err := transaction.Error; err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	return &client{database: transaction}, nil
}

func (c *client) Rollback() error {
	return c.database.Rollback().Error
}

func (c *client) Commit() error {
	return c.database.Commit().Error
}

func (c *client) RollbackBlock(ctx context.Context, chainID, blockNUmber uint64) error {
	databaseClient := c.database.WithContext(ctx)

	// Delete the bridge data.
	if err := databaseClient.
		Where(`"chain_id" = ? AND "block_number" >= ?`, chainID, blockNUmber).
		Delete(&table.BridgeTransaction{}).
		Error; err != nil {
		return fmt.Errorf("delete bridge transactions: %w", err)
	}

	if err := databaseClient.
		Where(`"chain_id" = ? AND "block_number" >= ?`, chainID, blockNUmber).
		Delete(&table.BridgeEvent{}).
		Error; err != nil {
		return fmt.Errorf("delete bridge events: %w", err)
	}

	// Delete the stake data.
	if err := databaseClient.
		Where(`"block_number" >= ?`, blockNUmber).
		Error; err != nil {
		return fmt.Errorf("delete bridge transactions: %w", err)
	}

	if err := databaseClient.
		Where(`"block_number" >= ?`, blockNUmber).
		Delete(&table.StakeEvent{}).
		Error; err != nil {
		return fmt.Errorf("delete bridge events: %w", err)
	}

	return nil
}

func (c *client) FindNode(ctx context.Context, nodeAddress common.Address) (*schema.Node, error) {
	var node table.Node

	if err := c.database.WithContext(ctx).First(&node, "address = ?", nodeAddress).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		return nil, nil
	}

	return node.Export()
}

func (c *client) FindNodes(ctx context.Context, nodeAddresses []common.Address, cursor *string, limit int) ([]*schema.Node, error) {
	databaseStatement := c.database.WithContext(ctx)

	if cursor != nil {
		var nodeCursor *table.Node

		if err := c.database.WithContext(ctx).First(&nodeCursor, "address = ?", common.HexToAddress(lo.FromPtr(cursor))).Error; err != nil {
			return nil, fmt.Errorf("get node cursor: %w", err)
		}

		databaseStatement = databaseStatement.Where("created_at < ?", nodeCursor.CreatedAt)
	}

	if len(nodeAddresses) > 0 {
		databaseStatement = databaseStatement.Where("address IN ?", nodeAddresses)
	}

	var nodes table.Nodes

	if err := databaseStatement.Limit(limit).Order("created_at DESC").Find(&nodes).Error; err != nil {
		return nil, fmt.Errorf("find nodes: %w", err)
	}

	return nodes.Export()
}

func (c *client) SaveNode(ctx context.Context, data *schema.Node) error {
	var nodes table.Node

	if err := nodes.Import(data); err != nil {
		return err
	}

	// Save node.
	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "address",
			},
		},
		UpdateAll: true,
	}

	return c.database.WithContext(ctx).Clauses(onConflict).Create(&nodes).Error
}

func (c *client) UpdateNodesStatus(ctx context.Context, lastHeartbeatTimestamp int64) error {
	return c.WithTransaction(ctx, func(ctx context.Context, client database.Client) error {
		for {
			result := c.database.WithContext(ctx).Model(&table.Node{}).
				Where("last_heartbeat_timestamp < ? and status = ?", time.Unix(lastHeartbeatTimestamp, 0), schema.StatusOnline).
				Update("status", schema.StatusOffline).Limit(1000)
			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return nil
			}
		}
	})
}

func (c *client) FindNodeStat(ctx context.Context, nodeAddress common.Address) (*schema.Stat, error) {
	var stat table.Stat

	if err := c.database.WithContext(ctx).First(&stat, "address = ?", nodeAddress).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		return nil, nil
	}

	return stat.Export()
}

func (c *client) FindNodeStats(ctx context.Context, query *schema.StatQuery) ([]*schema.Stat, error) {
	var stats table.Stats

	databaseStatement, err := c.buildNodeStatQuery(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("build find node stats: %w", err)
	}

	if err := databaseStatement.Find(&stats).Error; err != nil {
		return nil, fmt.Errorf("find nodes: %w", err)
	}

	return stats.Export()
}

func (c *client) buildNodeStatQuery(ctx context.Context, query *schema.StatQuery) (*gorm.DB, error) {
	databaseStatement := c.database.WithContext(ctx)

	if query.Cursor != nil {
		var statCursor *table.Stat

		if err := databaseStatement.First(&statCursor, "address = ?", common.HexToAddress(lo.FromPtr(query.Cursor))).Error; err != nil {
			return nil, fmt.Errorf("get node cursor: %w", err)
		}

		databaseStatement = databaseStatement.Where(clause.Gt{
			Column: "created_at",
			Value:  statCursor.CreatedAt,
		})
	}

	if query.Address != nil {
		databaseStatement = databaseStatement.Where(clause.Eq{
			Column: "address",
			Value:  query.Address,
		})
	}

	if len(query.AddressList) > 0 {
		databaseStatement = databaseStatement.Where(clause.IN{
			Column: "address",
			Values: lo.ToAnySlice(query.AddressList),
		})
	}

	if query.IsFullNode != nil {
		databaseStatement = databaseStatement.Where(clause.Eq{
			Column: "is_full_node",
			Value:  query.IsFullNode,
		})
	}

	if query.IsRssNode != nil {
		databaseStatement = databaseStatement.Where(clause.Eq{
			Column: "is_rss_node",
			Value:  query.IsRssNode,
		})
	}

	if query.Limit != nil {
		databaseStatement = databaseStatement.Limit(*query.Limit)
	}

	if query.ValidRequest != nil {
		databaseStatement = databaseStatement.Where(clause.Lt{
			Column: "epoch_invalid_request_count",
			Value:  query.ValidRequest,
		})
	}

	orderByPointsClause := clause.OrderByColumn{
		Column: clause.Column{
			Name: "points",
		},
	}

	orderByCreatedAtClause := clause.OrderByColumn{
		Column: clause.Column{
			Name: "created_at",
		},
	}

	if query.PointsOrder != nil && strings.EqualFold(*query.PointsOrder, "DESC") {
		orderByPointsClause.Desc = true

		databaseStatement = databaseStatement.Order(orderByPointsClause)
	} else {
		databaseStatement = databaseStatement.Order(orderByCreatedAtClause)
	}

	return databaseStatement, nil
}

func (c *client) SaveNodeStat(ctx context.Context, stat *schema.Stat) error {
	var stats table.Stat

	if err := stats.Import(stat); err != nil {
		return err
	}

	// Save node stat.
	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "address",
			},
		},
		UpdateAll: true,
	}

	return c.database.WithContext(ctx).Clauses(onConflict).Create(&stats).Error
}

func (c *client) SaveNodeStats(ctx context.Context, stats []*schema.Stat) error {
	var tStats table.Stats

	if err := tStats.Import(stats); err != nil {
		return err
	}

	// Save node indexers.
	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "address",
			},
		},
		UpdateAll: true,
	}

	return c.database.WithContext(ctx).Clauses(onConflict).CreateInBatches(tStats, math.MaxUint8).Error
}

func (c *client) DeleteNodeIndexers(ctx context.Context, nodeAddress common.Address) error {
	return c.database.WithContext(ctx).Where("address = ?", nodeAddress).Delete(&table.Indexer{}).Error
}

func (c *client) FindNodeIndexers(ctx context.Context, nodeAddresses []common.Address, networks, workers []string) ([]*schema.Indexer, error) {
	var indexers table.Indexers

	databaseStatement := c.database.WithContext(ctx)

	if len(nodeAddresses) > 0 {
		databaseStatement = databaseStatement.Where("address IN ?", nodeAddresses)
	}

	if len(networks) > 0 {
		databaseStatement = databaseStatement.Where("network IN ?", networks)
	}

	if len(workers) > 0 {
		databaseStatement = databaseStatement.Where("worker IN ?", workers)
	}

	if err := databaseStatement.Find(&indexers).Error; err != nil {
		return nil, fmt.Errorf("find nodes: %w", err)
	}

	return indexers.Export()
}

func (c *client) SaveNodeIndexers(ctx context.Context, indexers []*schema.Indexer) error {
	var tIndexers table.Indexers

	if err := tIndexers.Import(indexers); err != nil {
		return err
	}

	return c.database.WithContext(ctx).CreateInBatches(tIndexers, math.MaxUint8).Error
}

// Dial dials a database.
func Dial(_ context.Context, dataSourceName string) (database.Client, error) {
	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault()

	config := gorm.Config{
		Logger: logger,
	}

	databaseClient, err := gorm.Open(postgres.Open(dataSourceName), &config)
	if err != nil {
		return nil, fmt.Errorf("dial database: %w", err)
	}

	return &client{
		database: databaseClient,
	}, nil
}

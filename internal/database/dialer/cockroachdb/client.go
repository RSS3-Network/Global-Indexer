package cockroachdb

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"math"

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

	return c.database.WithContext(ctx).Clauses(onConflict).Save(&nodes).Error
}

func (c *client) FindNodeStat(_ context.Context) (*schema.Stat, error) {
	return nil, nil
}

func (c *client) FindNodeStats(ctx context.Context, nodeAddresses []common.Address) ([]*schema.Stat, error) {
	databaseStatement := c.database.WithContext(ctx)

	if len(nodeAddresses) > 0 {
		databaseStatement = databaseStatement.Where("address IN ?", nodeAddresses)
	}

	var stats table.Stats

	if err := databaseStatement.Limit(3).Order("points DESC").Find(&stats).Error; err != nil {
		return nil, fmt.Errorf("find nodes: %w", err)
	}

	return stats.Export()
}

func (c *client) FindNodeStatsByType(ctx context.Context, fullNode, rssNode *bool, limit int) ([]*schema.Stat, error) {
	databaseStatement := c.database.WithContext(ctx)

	if fullNode != nil {
		databaseStatement = databaseStatement.Where("is_full_node = ?", *fullNode)
	}

	if rssNode != nil {
		databaseStatement = databaseStatement.Where("is_rss_node = ?", *rssNode)
	}

	var stats table.Stats

	if err := databaseStatement.Limit(limit).Order("points DESC").Find(&stats).Error; err != nil {
		return nil, fmt.Errorf("find nodes: %w", err)
	}

	return stats.Export()
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

	return c.database.WithContext(ctx).Clauses(onConflict).Save(&stats).Error
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
	return c.database.WithContext(ctx).Delete(&table.Indexer{}, nodeAddress).Error
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

	// Save node indexers.
	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "address",
			},
		},
		UpdateAll: true,
	}

	return c.database.WithContext(ctx).Clauses(onConflict).CreateInBatches(tIndexers, math.MaxUint8).Error
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

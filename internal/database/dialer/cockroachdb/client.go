package cockroachdb

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/global-indexer/internal/database"
	"github.com/naturalselectionlabs/global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/global-indexer/schema"
	"github.com/pressly/goose/v3"
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
		return nil, err
	}

	return node.Export()
}

func (c *client) FindNodes(ctx context.Context, nodeAddresses []common.Address) ([]*schema.Node, error) {
	var nodes table.Nodes

	if err := c.database.WithContext(ctx).Find(&nodes, "address IN ?", nodeAddresses).Error; err != nil {
		return nil, err
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

func (c *client) FindNodeStats(ctx context.Context) ([]*schema.Stat, error) {
	var stats table.Stats

	if err := c.database.WithContext(ctx).Find(&stats).Error; err != nil {
		return nil, err
	}

	return stats.Export()
}

func (c *client) SaveNodeStat(_ context.Context, _ *schema.Stat) error {
	return nil
}

func (c *client) SaveNodeStats(_ context.Context, _ []*schema.Stat) error {
	return nil
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

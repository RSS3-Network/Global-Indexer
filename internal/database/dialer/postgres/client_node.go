package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/database/dialer/postgres/table"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) FindNode(ctx context.Context, nodeAddress common.Address) (*schema.Node, error) {
	var node table.Node

	if err := c.database.WithContext(ctx).First(&node, "address = ?", nodeAddress).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, err
	}

	return node.Export()
}

func (c *client) FindNodes(ctx context.Context, query schema.FindNodesQuery) ([]*schema.Node, error) {
	databaseStatement := c.database.WithContext(ctx)

	if query.Cursor != nil {
		var nodeCursor *table.Node

		if err := c.database.WithContext(ctx).First(&nodeCursor, "address = ?", common.HexToAddress(lo.FromPtr(query.Cursor))).Error; err != nil {
			return nil, fmt.Errorf("get Node cursor: %w", err)
		}

		if query.OrderByScore {
			databaseStatement = databaseStatement.Where("score < ? OR (score = ? AND created_at < ?)", nodeCursor.Score, nodeCursor.Score, nodeCursor.CreatedAt)
		} else {
			databaseStatement = databaseStatement.Where("created_at < ?", nodeCursor.CreatedAt)
		}
	}

	if query.Type != nil {
		databaseStatement = databaseStatement.Where("type = ?", query.Type.String())
	}

	if query.Status != nil {
		databaseStatement = databaseStatement.Where("status = ?", query.Status.String())
	}

	if len(query.NodeAddresses) > 0 {
		databaseStatement = databaseStatement.Where("address IN ?", query.NodeAddresses)
	}

	if query.Limit != nil {
		databaseStatement = databaseStatement.Limit(*query.Limit)
	}

	if query.OrderByScore {
		databaseStatement = databaseStatement.Order("score DESC, created_at DESC")
	} else {
		databaseStatement = databaseStatement.Order("created_at DESC")
	}

	var nodes table.Nodes

	if err := databaseStatement.Find(&nodes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, err
	}

	return nodes.Export()
}

func (c *client) FindNodeAvatar(ctx context.Context, nodeAddress common.Address) (*l2.ChipsTokenMetadata, error) {
	var node table.Node

	if err := c.database.WithContext(ctx).Model(&table.Node{}).Where("address = ?", nodeAddress).First(&node).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, err
	}

	var avatar l2.ChipsTokenMetadata
	if err := json.Unmarshal(node.Avatar, &avatar); len(node.Avatar) > 0 && err != nil {
		return nil, fmt.Errorf("unmarshal node avatar: %w", err)
	}

	return &avatar, nil
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

func (c *client) SaveNodeCountSnapshot(ctx context.Context, nodeSnapshot *schema.NodeSnapshot) error {
	databaseClient := c.database.WithContext(ctx)

	if err := databaseClient.
		Table((*table.Node).TableName(nil)).
		Count(&nodeSnapshot.Count).
		Error; err != nil {
		return fmt.Errorf("query count: %w", err)
	}

	var value table.NodeSnapshot
	if err := value.Import(*nodeSnapshot); err != nil {
		return fmt.Errorf("import node snapshot: %w", err)
	}

	return databaseClient.
		Table((*table.NodeSnapshot).TableName(nil)).
		Create(nodeSnapshot).
		Error
}

func (c *client) UpdateNodesStatusOffline(ctx context.Context, lastHeartbeatTimestamp int64) error {
	return c.WithTransaction(ctx, func(ctx context.Context, _ database.Client) error {
		for {
			result := c.database.WithContext(ctx).Model(&table.Node{}).
				Where("last_heartbeat_timestamp < ? and status = ?", time.Unix(lastHeartbeatTimestamp, 0), schema.NodeStatusOnline).
				Update("status", schema.NodeStatusOffline).Limit(1000)
			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return nil
			}
		}
	})
}

func (c *client) UpdateNodesHideTaxRate(ctx context.Context, nodeAddress common.Address, hideTaxRate bool) error {
	return c.database.
		WithContext(ctx).
		Model((*table.Node)(nil)).
		Where("address = ?", nodeAddress).
		Update("hideTaxRate", hideTaxRate).
		Error
}

func (c *client) UpdateNodesScore(ctx context.Context, nodes []*schema.Node) error {
	var tNodes table.Nodes

	if err := tNodes.Import(nodes); err != nil {
		return err
	}

	// Update node scores.
	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "address",
			},
		},
		DoUpdates: clause.AssignmentColumns([]string{"score"}),
	}

	return c.database.WithContext(ctx).Clauses(onConflict).CreateInBatches(tNodes, math.MaxUint8).Error
}

func (c *client) BatchUpdateNodes(ctx context.Context, data []*schema.BatchUpdateNode) error {
	rawSQL := "UPDATE node_info SET apy = CASE address"
	values := make([]interface{}, 0)

	for _, value := range data {
		rawSQL += " WHEN ? THEN ?"

		values = append(values, value.Address, value.Apy)
	}

	addresses := make([]common.Address, len(data))
	for i, value := range data {
		addresses[i] = value.Address
	}

	rawSQL += " END WHERE address IN (?)"

	values = append(values, addresses)

	return c.database.WithContext(ctx).Exec(rawSQL, values...).Error
}

func (c *client) UpdateNodePublicGood(ctx context.Context, nodeAddress common.Address, isPublicGood bool) error {
	return c.database.
		WithContext(ctx).
		Model((*table.Node)(nil)).
		Where("address = ?", nodeAddress).
		Update("is_public_good", isPublicGood).
		Error
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
		return nil, fmt.Errorf("build Find node stats: %w", err)
	}

	if err := databaseStatement.Find(&stats).Error; err != nil {
		return nil, fmt.Errorf("find Nodes: %w", err)
	}

	return stats.Export()
}

func (c *client) buildNodeStatQuery(ctx context.Context, query *schema.StatQuery) (*gorm.DB, error) {
	databaseStatement := c.database.WithContext(ctx)

	if query.Cursor != nil {
		var statCursor *table.Stat

		if err := databaseStatement.First(&statCursor, "address = ?", common.HexToAddress(lo.FromPtr(query.Cursor))).Error; err != nil {
			return nil, fmt.Errorf("get Node cursor: %w", err)
		}

		if query.PointsOrder != nil && strings.EqualFold(*query.PointsOrder, "DESC") {
			databaseStatement = databaseStatement.Where("points < ? OR (points = ? AND created_at < ?)", statCursor.Points, statCursor.Points, statCursor.CreatedAt)
		} else {
			databaseStatement = databaseStatement.Where("created_at < ?", statCursor.CreatedAt)
		}
	}

	if query.Address != nil {
		databaseStatement = databaseStatement.Where(clause.Eq{
			Column: "address",
			Value:  query.Address,
		})
	}

	if len(query.Addresses) > 0 {
		databaseStatement = databaseStatement.Where(clause.IN{
			Column: "address",
			Values: lo.ToAnySlice(query.Addresses),
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

	if query.PointsOrder != nil && strings.EqualFold(*query.PointsOrder, "DESC") {
		databaseStatement = databaseStatement.Order("points DESC, created_at DESC")
	} else {
		databaseStatement = databaseStatement.Order("created_at DESC")
	}

	return databaseStatement, nil
}

func (c *client) SaveNodeInvalidResponses(ctx context.Context, nodeInvalidResponse []*schema.NodeInvalidResponse) error {
	var tNodeInvalidResponses table.NodeInvalidResponses

	tNodeInvalidResponses.Import(nodeInvalidResponse)

	return c.database.WithContext(ctx).CreateInBatches(tNodeInvalidResponses, math.MaxUint8).Error
}

func (c *client) FindNodeCountSnapshots(ctx context.Context) ([]*schema.NodeSnapshot, error) {
	databaseClient := c.database.WithContext(ctx)

	var nodeSnapshots []*table.NodeSnapshot

	if err := databaseClient.
		Order(`"date" DESC`).
		Limit(100). // TODO Replace this constant with a query parameter.
		Find(&nodeSnapshots).Error; err != nil {
		return nil, err
	}

	values := make([]*schema.NodeSnapshot, 0, len(nodeSnapshots))

	for _, nodeSnapshot := range nodeSnapshots {
		value, err := nodeSnapshot.Export()
		if err != nil {
			return nil, fmt.Errorf("export node snapshot: %w", err)
		}

		values = append(values, value)
	}

	return values, nil
}

func (c *client) SaveNodeStat(ctx context.Context, stat *schema.Stat) error {
	var stats table.Stat

	if err := stats.Import(stat); err != nil {
		return err
	}

	// Save Node stat.
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

	// Save Node indexers.
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

func (c *client) UpdateNodeWorkerActive(ctx context.Context) error {
	return c.database.WithContext(ctx).Model(&table.Worker{}).Where("is_active = ?", true).Update("is_active", false).Error
}

func (c *client) FindNodeWorkers(ctx context.Context, query *schema.WorkerQuery) ([]*schema.Worker, error) {
	var workers table.Workers

	databaseStatement := c.database.WithContext(ctx)

	if query.IsActive != nil {
		databaseStatement = databaseStatement.Where("is_active = ?", query.IsActive)
	}

	if query.EpochID > 0 {
		databaseStatement = databaseStatement.Where("epoch_id = ?", query.EpochID)
	}

	if len(query.NodeAddresses) > 0 {
		databaseStatement = databaseStatement.Where("address IN ?", query.NodeAddresses)
	}

	if len(query.Networks) > 0 {
		databaseStatement = databaseStatement.Where("network IN ?", query.Networks)
	}

	if len(query.Names) > 0 {
		databaseStatement = databaseStatement.Where("name IN ?", query.Names)
	}

	if err := databaseStatement.Find(&workers).Error; err != nil {
		return nil, fmt.Errorf("find node worker : %w", err)
	}

	return workers.Export(), nil
}

func (c *client) SaveNodeWorkers(ctx context.Context, workers []*schema.Worker) error {
	var tWorkers table.Workers

	tWorkers.Import(workers)

	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "epoch_id",
			},
			{
				Name: "address",
			},
			{
				Name: "network",
			},
			{
				Name: "name",
			},
		},
		UpdateAll: true,
	}

	return c.database.WithContext(ctx).Clauses(onConflict).CreateInBatches(tWorkers, math.MaxUint8).Error
}

func (c *client) SaveNodeEvent(ctx context.Context, nodeEvent *schema.NodeEvent) error {
	var event table.NodeEvent

	if err := event.Import(*nodeEvent); err != nil {
		return fmt.Errorf("import node event: %w", err)
	}

	// Save Node stat.
	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "transaction_hash",
			},
			{
				Name: "transaction_index",
			},
			{
				Name: "log_index",
			},
		},
		UpdateAll: true,
	}

	return c.database.WithContext(ctx).Clauses(onConflict).Create(&event).Error
}

func (c *client) FindNodeEvents(ctx context.Context, query *schema.NodeEventsQuery) ([]*schema.NodeEvent, error) {
	databaseStatement := c.database.WithContext(ctx)

	if query.Cursor != nil {
		key := strings.Split(*query.Cursor, ":")
		if len(key) != 3 {
			return nil, fmt.Errorf("invalid cursor: %s", *query.Cursor)
		}

		var nodeEvent *table.NodeEvent

		if err := c.database.WithContext(ctx).Where("transaction_hash = ?", key[0]).
			Where("transaction_index = ?", key[1]).
			Where("log_index = ?", key[2]).
			First(&nodeEvent).Error; err != nil {
			return nil, fmt.Errorf("get Node cursor: %w", err)
		}

		databaseStatement = databaseStatement.Where("block_number < ?", nodeEvent.BlockNumber).
			Or("block_number = ? AND transaction_index < ?", nodeEvent.BlockNumber, nodeEvent.TransactionIndex).
			Or("block_number = ? AND transaction_index < ? AND log_index < ?", nodeEvent.BlockNumber, nodeEvent.TransactionIndex, nodeEvent.LogIndex)
	}

	if query.NodeAddress != nil {
		databaseStatement = databaseStatement.Where("address_from = ?", query.NodeAddress)
	}

	if query.Finalized != nil {
		databaseStatement = databaseStatement.Where("finalized = ?", query.Finalized)
	}

	if query.Type != nil {
		databaseStatement = databaseStatement.Where("type = ?", query.Type)
	}

	if query.Limit != nil {
		databaseStatement = databaseStatement.Limit(*query.Limit)
	}

	var events table.NodeEvents

	if err := databaseStatement.Order("block_number DESC, transaction_index DESC, log_index DESC").Find(&events).Error; err != nil {
		return nil, err
	}

	return events.Export()
}

func (c *client) FindOperatorProfitSnapshots(ctx context.Context, query schema.OperatorProfitSnapshotsQuery) ([]*schema.OperatorProfitSnapshot, error) {
	databaseClient := c.database.WithContext(ctx).Table((*table.OperatorProfitSnapshot).TableName(nil))

	if query.Operator != nil {
		databaseClient = databaseClient.Where("operator = ?", *query.Operator)
	}

	if query.Cursor != nil {
		databaseClient = databaseClient.Where("id < ?", query.Cursor)
	}

	if query.BeforeDate != nil {
		databaseClient = databaseClient.Where("date <= ?", query.BeforeDate)
	}

	if query.AfterDate != nil {
		databaseClient = databaseClient.Where("date >= ?", query.AfterDate)
	}

	if query.Limit != nil {
		databaseClient = databaseClient.Limit(*query.Limit)
	}

	var snapshots table.OperatorProfitSnapshots

	if len(query.Dates) > 0 {
		var (
			queries []string
			values  []interface{}
		)

		for _, date := range query.Dates {
			queries = append(queries, `(SELECT * FROM "node"."operator_profit_snapshots" WHERE "date" >= ? and "operator" = ? ORDER BY "date" LIMIT 1)`)
			values = append(values, date, query.Operator)
		}

		// Combine all queries with UNION ALL
		fullQuery := strings.Join(queries, " UNION ALL ")

		// Execute the combined query
		if err := databaseClient.Raw(fullQuery, values...).Scan(&snapshots).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, database.ErrorRowNotFound
			}

			return nil, fmt.Errorf("find rows: %w", err)
		}
	} else {
		if err := databaseClient.Order("epoch_id DESC, id DESC").Find(&snapshots).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, database.ErrorRowNotFound
			}

			return nil, fmt.Errorf("find rows: %w", err)
		}
	}

	return snapshots.Export()
}

func (c *client) SaveOperatorProfitSnapshots(ctx context.Context, snapshots []*schema.OperatorProfitSnapshot) error {
	var value table.OperatorProfitSnapshots

	if err := value.Import(snapshots); err != nil {
		return fmt.Errorf("import operator profit snapshots: %w", err)
	}

	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "operator",
			},
			{
				Name: "epoch_id",
			},
		},
		UpdateAll: true,
	}

	return c.database.WithContext(ctx).Clauses(onConflict).CreateInBatches(value, math.MaxUint8).Error
}

func (c *client) SaveNodeAPYSnapshots(ctx context.Context, nodeAPYSnapshots []*schema.NodeAPYSnapshot) error {
	var value table.NodeAPYSnapshots

	if err := value.Import(nodeAPYSnapshots); err != nil {
		return fmt.Errorf("import node APY snapshots: %w", err)
	}

	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "node_address",
			},
			{
				Name: "epoch_id",
			},
		},
		UpdateAll: true,
	}

	return c.database.WithContext(ctx).Clauses(onConflict).CreateInBatches(value, math.MaxUint8).Error
}

func (c *client) DeleteNodeEventsByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return c.database.
		WithContext(ctx).
		Delete(new(table.NodeEvent), `"block_number" = ? AND NOT "finalized"`, blockNumber).
		Error
}

func (c *client) UpdateNodeEventsFinalizedByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return c.database.
		WithContext(ctx).
		Table((*table.NodeEvent).TableName(nil)).
		Where(`"block_number" < ? AND NOT "finalized"`, blockNumber).
		Update("finalized", true).
		Error
}

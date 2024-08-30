package cockroachdb

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/common/ethereum"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/sourcegraph/conc/pool"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) FindStakeTransaction(ctx context.Context, query schema.StakeTransactionQuery) (*schema.StakeTransaction, error) {
	var row table.StakeTransaction

	databaseClient := c.database.WithContext(ctx)

	if query.ID != nil {
		databaseClient = databaseClient.Where(`"id" = ?`, query.ID.String())
	}

	if query.User != nil {
		databaseClient = databaseClient.Where(`"user" = ?`, query.User.String())
	}

	if query.Node != nil {
		databaseClient = databaseClient.Where(`"node" = ?`, query.Node.String())
	}

	if query.Address != nil {
		databaseClient = databaseClient.Where(`"user" = ? OR "node" = ?`, query.Address.String())
	}

	if query.Type != nil {
		databaseClient = databaseClient.Where(`"type" = ?`, query.Type)
	}

	if err := databaseClient.Order(`"block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC`).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, fmt.Errorf("find stake transaction: %w", err)
	}

	return row.Export()
}

func (c *client) FindStakeTransactions(ctx context.Context, query schema.StakeTransactionsQuery) ([]*schema.StakeTransaction, error) {
	databaseClient := c.database.WithContext(ctx)

	if query.Cursor != nil {
		var cursor table.StakeTransaction
		if err := databaseClient.Where(`"id" = ?`, query.Cursor.String()).First(&cursor).Error; err != nil {
			return nil, fmt.Errorf("query cursor: %w", err)
		}

		databaseClient = databaseClient.Where(
			`("block_number" < ?) OR ("block_number" = ? AND "transaction_index" < ?)`,
			cursor.BlockNumber,
			cursor.BlockNumber, cursor.TransactionIndex,
		)
	}

	if query.IDs != nil {
		databaseClient = databaseClient.Where(`"id" = ?`, lo.Map(query.IDs, func(id common.Hash, _ int) string {
			return id.String()
		}))
	}

	if query.User != nil {
		databaseClient = databaseClient.Where(`"user" = ?`, query.User.String())
	}

	if query.Node != nil {
		databaseClient = databaseClient.Where(`"node" = ?`, query.Node.String())
	}

	if query.Address != nil {
		databaseClient = databaseClient.Where(`"user" = ? OR "node" = ?`, query.Address.String())
	}

	if query.Type != nil {
		databaseClient = databaseClient.Where(`"type" = ?`, query.Type)
	}

	if query.BlockTimestamp != nil {
		databaseClient = databaseClient.Where(`"block_timestamp" >= ?`, query.BlockTimestamp)
	}

	if query.Pending != nil && *query.Pending {
		subQuery := c.database.WithContext(ctx).
			Select("TRUE").
			Table((*table.StakeEvent).TableName(nil)).
			Where(`"transactions"."id" = "events"."id" AND "events"."type" IN ('withdraw_claimed', 'unstake_claimed')`)

		databaseClient = databaseClient.
			Where(`"type" IN (?, ?)`, schema.StakeTransactionTypeUnstake, schema.StakeTransactionTypeWithdraw).
			Not(`EXISTS (?)`, subQuery)
	}

	if query.Order != "" {
		databaseClient = databaseClient.Order(query.Order)
	} else {
		databaseClient = databaseClient.Order(`"block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC`)
	}

	if query.Limit != 0 {
		databaseClient = databaseClient.Limit(query.Limit)
	}

	var rows []table.StakeTransaction

	if err := databaseClient.Find(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, fmt.Errorf("find stake transactions: %w", err)
	}

	results := make([]*schema.StakeTransaction, 0, len(rows))

	for _, row := range rows {
		result, err := row.Export()
		if err != nil {
			return nil, fmt.Errorf("export stake transaction: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (c *client) FindStakeEvents(ctx context.Context, query schema.StakeEventsQuery) ([]*schema.StakeEvent, error) {
	databaseClient := c.database.WithContext(ctx)

	if len(query.IDs) > 0 {
		databaseClient = databaseClient.Where(`"id" IN ?`, lo.Map(query.IDs, func(id common.Hash, _ int) string {
			return id.String()
		}))
	}

	var rows []table.StakeEvent
	if err := databaseClient.Order(`"block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC`).Find(&rows).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, fmt.Errorf("find stake events: %w", err)
	}

	results := make([]*schema.StakeEvent, 0, len(rows))

	for _, row := range rows {
		result, err := row.Export()
		if err != nil {
			return nil, fmt.Errorf("export stake event: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (c *client) FindStakeChips(ctx context.Context, query schema.StakeChipsQuery) ([]*schema.StakeChip, error) {
	databaseClient := c.database.WithContext(ctx).Table((*table.StakeChip).TableName(nil))

	if query.BlockNumber != nil {
		databaseClient = databaseClient.Where(`"block_number" <= ?`, query.BlockNumber)
	}

	if query.Cursor != nil {
		databaseClient = databaseClient.Where(`"id" > ?`, query.Cursor.String())
	}

	if len(query.IDs) > 0 {
		databaseClient = databaseClient.Where(`"id" IN ?`, lo.Map(query.IDs, func(id *big.Int, _ int) uint64 { return id.Uint64() }))
	}

	if query.Node != nil {
		databaseClient = databaseClient.Where(`"node" = ?`, query.Node.String())
	}

	if query.Owner != nil {
		databaseClient = databaseClient.Where(`"owner" = ?`, query.Owner.String())
	}

	if query.Limit != nil {
		databaseClient = databaseClient.Limit(*query.Limit)
	}

	databaseClient = databaseClient.Order("id ASC")

	if query.DistinctOwner {
		subQuery := databaseClient

		databaseClient = c.database.WithContext(ctx).Table((*table.StakeChip).TableName(nil)).
			Select("DISTINCT ON (owner) *").Order("owner, id DESC").Where("id IN (?)", subQuery.Select("id"))
	}

	var rows []*table.StakeChip
	if err := databaseClient.Find(&rows).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find rows: %w", err)
	}

	results := make([]*schema.StakeChip, 0, len(rows))

	for _, row := range rows {
		result, err := row.Export()
		if err != nil {
			return nil, fmt.Errorf("export row: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (c *client) FindStakerCount(ctx context.Context, query schema.StakeChipsQuery) (int64, error) {
	databaseClient := c.database.WithContext(ctx).Table((*table.StakeChip).TableName(nil)).
		Distinct(`"owner"`).
		Where(`"owner" != ?`, ethereum.AddressGenesis.String())

	if query.BlockNumber != nil {
		databaseClient = databaseClient.Where(`"block_number" <= ?`, query.BlockNumber)
	}

	if query.Cursor != nil {
		databaseClient = databaseClient.Where(`"id" > ?`, query.Cursor.String())
	}

	if len(query.IDs) > 0 {
		databaseClient = databaseClient.Where(`"id" IN ?`, lo.Map(query.IDs, func(id *big.Int, _ int) uint64 { return id.Uint64() }))
	}

	if query.Node != nil {
		databaseClient = databaseClient.Where(`"node" = ?`, query.Node.String())
	}

	if query.Owner != nil {
		databaseClient = databaseClient.Where(`"owner" = ?`, query.Owner.String())
	}

	if query.Limit != nil {
		databaseClient = databaseClient.Limit(*query.Limit)
	}

	var count int64

	if err := databaseClient.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (c *client) FindStakeChip(ctx context.Context, query schema.StakeChipQuery) (*schema.StakeChip, error) {
	databaseClient := c.database.WithContext(ctx)

	if query.ID != nil {
		databaseClient = databaseClient.Where(`"id" = ?`, query.ID.String())
	}

	var row table.StakeChip
	if err := databaseClient.First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		return nil, fmt.Errorf("find stake chip: %w", err)
	}

	result, err := row.Export()
	if err != nil {
		return nil, fmt.Errorf("export row: %w", err)
	}

	return result, nil
}

func (c *client) DeleteStakeChipsByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return c.database.
		WithContext(ctx).
		Delete(new(table.StakeChip), `"block_number" = ? AND NOT "finalized"`, blockNumber).
		Error
}

func (c *client) FindStakeStakings(ctx context.Context, query schema.StakeStakingsQuery) ([]*schema.StakeStaking, error) {
	databaseClient := c.database.WithContext(ctx)

	if query.Cursor != nil {
		cursor, err := base64.StdEncoding.DecodeString(*query.Cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid curosr: %w", err)
		}

		splits := strings.Split(string(cursor), "-")
		if length := len(splits); length != 3 {
			return nil, fmt.Errorf("invalid curosr length: %d", length)
		}

		databaseClient = databaseClient.Where(
			`"value" < @value OR ("value" = @value AND "staker" > @staker) OR ("value" = @value AND "staker" = @staker AND "node" > @node)`,
			sql.Named("value", splits[0]),
			sql.Named("staker", splits[1]),
			sql.Named("node", splits[2]),
		)
	}

	if query.Staker != nil {
		databaseClient = databaseClient.Where(`"staker" = ?`, query.Staker.String())
	}

	if query.Node != nil {
		databaseClient = databaseClient.Where(`"node" = ?`, query.Node.String())
	}

	var stakeStakings []*table.StakeStaking
	if err := databaseClient.
		Limit(query.Limit).
		Order(`"value" DESC, "staker", "node"`).
		Find(&stakeStakings).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	resultsPool := pool.NewWithResults[*schema.StakeStaking]().WithContext(ctx).WithFirstError().WithCancelOnError()

	for _, stakeStaking := range stakeStakings {
		stakeStaking := stakeStaking

		resultsPool.Go(func(ctx context.Context) (*schema.StakeStaking, error) {
			databaseClient := c.database.WithContext(ctx)

			var stakeChips []*table.StakeChip
			if err := databaseClient.
				Where(`"owner" = ? AND "node" = ?`, stakeStaking.Staker, stakeStaking.Node).
				Order(`"id" DESC`).
				Limit(5).
				Find(&stakeChips).Error; err != nil {
				return nil, err
			}

			stakeStaking, err := stakeStaking.Export()
			if err != nil {
				return nil, fmt.Errorf("export stake staking: %w", err)
			}

			stakeStaking.Chips.Showcase = lo.Map(stakeChips, func(stakeChip *table.StakeChip, _ int) *schema.StakeChip {
				return lo.Must(stakeChip.Export())
			})

			return stakeStaking, nil
		})
	}

	results, err := resultsPool.Wait()
	if err != nil {
		return nil, err
	}

	slices.SortStableFunc(results, func(left, right *schema.StakeStaking) int {
		if n := right.Value.Cmp(left.Value); n != 0 { // DESC
			return n
		}

		if n := strings.Compare(left.Staker.String(), right.Staker.String()); n != 0 { // ASC
			return n
		}

		return strings.Compare(left.Node.String(), right.Node.String()) // ASC
	})

	return results, nil
}

func (c *client) FindStakeStaker(ctx context.Context, address common.Address) (*schema.StakeStaker, error) {
	databaseTransaction := c.database.WithContext(ctx).Begin(&sql.TxOptions{ReadOnly: true})
	defer databaseTransaction.Rollback()

	var totalStakedTokens decimal.Decimal

	/*
		SELECT
		    coalesce(sum(
		        CASE
		            WHEN transactions.type = 'stake' AND events.type = 'staked' THEN value
		            WHEN transactions.type = 'unstake' AND events.type = 'claimed' THEN -value
		            ELSE 0
		        END
		        ), 0) AS total_staked_tokens
		FROM stake.transactions
		         LEFT JOIN stake.events ON transactions.id = events.id
		WHERE transactions.user = $1 AND transactions.finalized;
	*/

	if err := databaseTransaction.Debug().
		Select(`
			coalesce(sum(
				CASE
					WHEN transactions.type = ? AND events.type = ? THEN value
					WHEN transactions.type = ? AND events.type = ? THEN -value
					ELSE 0
				END
			), 0) AS total_staked_tokens`,
			schema.StakeTransactionTypeStake, schema.StakeEventTypeStakeStaked,
			schema.StakeTransactionTypeUnstake, schema.StakeEventTypeUnstakeClaimed,
		).
		Table((*table.StakeTransaction).TableName(nil)).
		Joins("LEFT JOIN stake.events ON transactions.id = events.id").
		Where(`"transactions"."user" = ? AND transactions.finalized`, address.String()).
		Scan(&totalStakedTokens).
		Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	/*
		SELECT
			count(node) AS staked_nodes,
			sum(count) 	AS owned_chips,
			sum(value) 	AS stake_tokens
		FROM stake.stakings
		WHERE staker = $1;
	*/

	type StakeStakingAggregate struct {
		StakedNodes  uint64
		OwnedChips   uint64
		StakedTokens decimal.Decimal
	}

	var aggregate StakeStakingAggregate

	if err := databaseTransaction.
		Select("count(node) AS staked_nodes, sum(count) AS owned_chips, sum(value) AS staked_tokens").
		Table((*table.StakeStaking).TableName(nil)).
		Where(`"staker" = ?`, address.String()).
		Scan(&aggregate).
		Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	_ = databaseTransaction.Commit().Error

	stakeStaker := schema.StakeStaker{
		Address:             address,
		TotalStakedNodes:    aggregate.StakedNodes,
		TotalChips:          aggregate.OwnedChips,
		TotalStakedTokens:   totalStakedTokens,
		CurrentStakedTokens: aggregate.StakedTokens,
	}

	return &stakeStaker, nil
}

func (c *client) FindStakerCountSnapshots(ctx context.Context) ([]*schema.StakerCountSnapshot, error) {
	databaseClient := c.database.WithContext(ctx)

	var stakeSnapshots []*table.StakerCountSnapshot

	if err := databaseClient.
		Order(`"date" DESC`).
		Limit(100). // FIXME: Replace this constant with a query parameter.
		Find(&stakeSnapshots).Error; err != nil {
		return nil, err
	}

	values := make([]*schema.StakerCountSnapshot, 0, len(stakeSnapshots))

	for _, stakeSnapshot := range stakeSnapshots {
		value, err := stakeSnapshot.Export()
		if err != nil {
			return nil, fmt.Errorf("export staker count snapshots: %w", err)
		}

		values = append(values, value)
	}

	return values, nil
}

func (c *client) FindStakerCountRecentEpochs(ctx context.Context, recentEpochs int) (map[common.Address]*schema.StakeRecentCount, error) {
	// Get the block number of the recent epoch.
	subQuery := c.database.
		WithContext(ctx).
		Table((*table.Epoch).TableName(nil)).
		Select(`"block_number"`).
		Order(`"id" DESC`).
		Offset(recentEpochs).
		Limit(1)

	// Gets the count of unique stakers for each node.
	databaseClient := c.database.
		WithContext(ctx).
		Table((*table.StakeTransaction).TableName(nil)).
		Select(`"node", count(DISTINCT "user"),sum("value") as "stake"`).
		Where(`"block_number" >= coalesce((?), 0) AND "type" = 'stake'`, subQuery).
		Group(`"node"`)

	// Define a row struct to store the result.
	type row struct {
		Node  string          `gorm:"column:node"`
		Count uint64          `gorm:"column:count"`
		Stake decimal.Decimal `gorm:"column:stake"`
	}

	// SELECT "node", count(DISTINCT "user"), sum("value") as "stake"
	// FROM "stake"."transactions"
	// WHERE "block_number" >= coalesce((SELECT "block_number" FROM "epoch" ORDER BY "id" DESC LIMIT 1 OFFSET @recentEpochs), 0)
	//   AND "type" = 'stake'
	// GROUP BY "node"

	var rows []row
	if err := databaseClient.Find(&rows).Error; err != nil {
		return nil, err
	}

	// Converts the rows into a map of Node address to their staker counts.
	result := lo.SliceToMap(rows, func(row row) (common.Address, *schema.StakeRecentCount) {
		return common.HexToAddress(row.Node), &schema.StakeRecentCount{
			StakerCount: row.Count,
			StakeValue:  row.Stake,
		}
	})

	return result, nil
}

func (c *client) SaveStakeTransaction(ctx context.Context, stakeTransaction *schema.StakeTransaction) error {
	var value table.StakeTransaction
	if err := value.Import(*stakeTransaction); err != nil {
		return fmt.Errorf("import stake transaction: %w", err)
	}

	clauses := []clause.Expression{
		clause.OnConflict{
			Columns: []clause.Column{
				{
					Name: "id",
				},
				{
					Name: "type",
				},
			},
			UpdateAll: true,
		},
	}

	return c.database.WithContext(ctx).Clauses(clauses...).Create(&value).Error
}

func (c *client) SaveStakeEvent(ctx context.Context, stakeEvent *schema.StakeEvent) error {
	var value table.StakeEvent
	if err := value.Import(*stakeEvent); err != nil {
		return fmt.Errorf("import stake event: %w", err)
	}

	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "transaction_hash",
			},
			{
				Name: "log_index",
			},
			{
				Name: "id",
			},
		},
		UpdateAll: true,
	}

	return c.database.WithContext(ctx).Clauses(onConflict).Create(&value).Error
}

func (c *client) SaveStakeChips(ctx context.Context, stakeChips ...*schema.StakeChip) error {
	values := make([]*table.StakeChip, 0, len(stakeChips))

	clauses := []clause.Expression{
		clause.OnConflict{
			UpdateAll: true,
			Columns: []clause.Column{
				{
					Name: "id",
				},
			},
		},
	}

	for _, stakeChip := range stakeChips {
		var value table.StakeChip

		if err := value.Import(*stakeChip); err != nil {
			return fmt.Errorf("import stake chip: %w", err)
		}

		values = append(values, &value)
	}

	return c.database.WithContext(ctx).Clauses(clauses...).Create(&values).Error
}

func (c *client) UpdateStakeChipsOwner(ctx context.Context, owner common.Address, stakeChipIDs ...*big.Int) error {
	ids := lo.Map(stakeChipIDs, func(stakeChipID *big.Int, _ int) decimal.Decimal {
		return decimal.NewFromBigInt(stakeChipID, 0)
	})

	return c.database.WithContext(ctx).Model((*table.StakeChip)(nil)).Where(`"id" IN ?`, ids).UpdateColumn("owner", owner.String()).Error
}

func (c *client) SaveStakerCountSnapshot(ctx context.Context, stakeSnapshot *schema.StakerCountSnapshot) error {
	databaseClient := c.database.WithContext(ctx)

	if err := databaseClient.
		Table((*table.StakeChip).TableName(nil)).
		Distinct(`"owner"`).
		Where(`"owner" != ?`, ethereum.AddressGenesis.String()).
		Count(&stakeSnapshot.Count).
		Error; err != nil {
		return fmt.Errorf("query count: %w", err)
	}

	var value table.StakerCountSnapshot
	if err := value.Import(*stakeSnapshot); err != nil {
		return fmt.Errorf("import stakers_count snapshot: %w", err)
	}

	return databaseClient.
		Table((*table.StakerCountSnapshot).TableName(nil)).
		Create(stakeSnapshot).
		Error
}

func (c *client) FindStakerProfitSnapshots(ctx context.Context, query schema.StakerProfitSnapshotsQuery) ([]*schema.StakerProfitSnapshot, error) {
	databaseClient := c.database.WithContext(ctx).Table((*table.StakerProfitSnapshot).TableName(nil))

	if query.Cursor != nil {
		databaseClient = databaseClient.Where(`"id" < ?`, query.Cursor)
	}

	if query.OwnerAddress != nil {
		databaseClient = databaseClient.Where(`"owner_address" = ?`, query.OwnerAddress)
	}

	if query.EpochID != nil {
		databaseClient = databaseClient.Where(`"epoch_id" = ?`, query.EpochID)
	}

	if query.BeforeDate != nil {
		databaseClient = databaseClient.Where(`"date" <= ?`, query.BeforeDate)
	}

	if query.AfterDate != nil {
		databaseClient = databaseClient.Where(`"date" >= ?`, query.AfterDate)
	}

	if query.EpochIDs != nil {
		databaseClient = databaseClient.Where(`"epoch_id" IN ?`, query.EpochIDs)
	}

	if query.Limit != nil {
		databaseClient = databaseClient.Limit(*query.Limit)
	}

	var rows []*table.StakerProfitSnapshot

	if len(query.Dates) > 0 {
		var (
			queries []string
			values  []interface{}
		)

		for _, date := range query.Dates {
			queries = append(queries, `(SELECT * FROM "stake"."profit_snapshots" WHERE "date" >= ? AND "owner_address" = ? ORDER BY "date" LIMIT 1)`)
			values = append(values, date, query.OwnerAddress)
		}

		// Combine all queries with UNION ALL
		fullQuery := strings.Join(queries, " UNION ALL ")

		// Execute the combined query
		if err := databaseClient.Raw(fullQuery, values...).Scan(&rows).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, database.ErrorRowNotFound
			}

			return nil, fmt.Errorf("find rows: %w", err)
		}
	} else {
		if err := databaseClient.Order("epoch_id DESC, id DESC").Find(&rows).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, database.ErrorRowNotFound
			}

			return nil, fmt.Errorf("find rows: %w", err)
		}
	}

	results := make([]*schema.StakerProfitSnapshot, 0, len(rows))

	for _, row := range rows {
		result, err := row.Export()
		if err != nil {
			return nil, fmt.Errorf("export row: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}

func (c *client) SaveStakerProfitSnapshots(ctx context.Context, snapshots []*schema.StakerProfitSnapshot) error {
	var value table.StakerProfitSnapshots

	if err := value.Import(snapshots); err != nil {
		return fmt.Errorf("import staker profit snapshots: %w", err)
	}

	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "owner_address",
			},
			{
				Name: "epoch_id",
			},
		},
		UpdateAll: true,
	}

	return c.database.WithContext(ctx).Clauses(onConflict).Create(&value).Error
}

func (c *client) DeleteStakeTransactionsByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return c.database.
		WithContext(ctx).
		Delete(new(table.StakeTransaction), `"block_number" = ? AND NOT "finalized"`, blockNumber).
		Error
}

func (c *client) DeleteStakeEventsByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return c.database.
		WithContext(ctx).
		Delete(new(table.StakeEvent), `"block_number" = ? AND NOT "finalized"`, blockNumber).
		Error
}

func (c *client) UpdateStakeTransactionsFinalizedByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return c.database.
		WithContext(ctx).
		Table((*table.StakeTransaction).TableName(nil)).
		Where(`"block_number" < ? AND NOT "finalized"`, blockNumber).
		Update("finalized", true).
		Error
}

func (c *client) UpdateStakeEventsFinalizedByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return c.database.
		WithContext(ctx).
		Table((*table.StakeEvent).TableName(nil)).
		Where(`"block_number" < ? AND NOT "finalized"`, blockNumber).
		Update("finalized", true).
		Error
}

func (c *client) UpdateStakeChipsFinalizedByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return c.database.
		WithContext(ctx).
		Table((*table.StakeChip).TableName(nil)).
		Where(`"block_number" < ? AND NOT "finalized"`, blockNumber).
		Update("finalized", true).
		Error
}

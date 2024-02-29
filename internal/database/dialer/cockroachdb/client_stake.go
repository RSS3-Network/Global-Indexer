package cockroachdb

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/ethereum"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
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

	if query.Pending != nil {
		subQuery := c.database.WithContext(ctx).
			Select("TRUE").
			Table((*table.StakeEvent).TableName(nil)).
			Where(`"transactions"."id" = "events"."id" AND "events"."type" = 'claimed'`)

		databaseClient = databaseClient.
			Where(`"type" IN (?, ?)`, schema.StakeTransactionTypeUnstake, schema.StakeTransactionTypeWithdraw).
			Not(`EXISTS (?)`, subQuery)
	}

	var rows []table.StakeTransaction

	if err := databaseClient.Order(`"block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC`).Limit(query.Limit).Find(&rows).Error; err != nil {
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
		databaseClient.Where(`"id" IN ?`, lo.Map(query.IDs, func(id common.Hash, _ int) string {
			return id.String()
		}))
	}

	var rows []table.StakeEvent
	if err := c.database.WithContext(ctx).Order(`"block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC`).Find(&rows).Error; err != nil {
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
	databaseClient := c.database.WithContext(ctx)

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

	var rows []*table.StakeChip
	if err := databaseClient.Order(`"id" ASC`).Find(&rows).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

func (c *client) FindStakeStakings(ctx context.Context, query schema.StakeStakingsQuery) ([]*schema.StakeStaking, error) {
	databaseClient := c.database.
		WithContext(ctx).
		Table((*table.StakeChip).TableName(nil))

	databaseClient = databaseClient.Where(`"owner" != ?`, ethereum.AddressGenesis.String())

	if query.Cursor != nil {
		cursor, err := base64.StdEncoding.DecodeString(*query.Cursor)
		if err != nil {
			return nil, err
		}

		splits := strings.Split(string(cursor), "-")
		if len(splits) != 2 {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}

		databaseClient = databaseClient.Where(
			`"owner" > ? OR ("owner" = ? AND "node" > ?)`,
			splits[0], splits[0], splits[1],
		)
	}

	if query.Staker != nil {
		databaseClient = databaseClient.Where(`"owner" = ?`, query.Staker.String())
	}

	if query.Node != nil {
		databaseClient = databaseClient.Where(`"node" = ?`, query.Node.String())
	}

	type StakeStaking struct {
		Owner string          `gorm:"column:owner"`
		Node  string          `gorm:"column:node"`
		Value decimal.Decimal `gorm:"column:value"`
		Count uint64          `gorm:"column:count"`
	}

	var stakeStakings []*StakeStaking
	if err := databaseClient.
		Select(`"owner", "node", count(*) AS "count", sum("value") AS "value"`).
		Group(`"owner", "node"`).
		Order(`"count" DESC, "owner", "node"`).
		Limit(query.Limit).
		Find(&stakeStakings).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	results := make([]*schema.StakeStaking, 0, len(stakeStakings))

	for _, stakeStaking := range stakeStakings {
		databaseClient := c.database.WithContext(ctx)

		var stakeChips []*table.StakeChip
		if err := databaseClient.
			Where(
				`"owner" = ? AND "node" = ? AND "owner" != ?`, stakeStaking.Owner, stakeStaking.Node, ethereum.AddressGenesis.String(),
			).
			Order(`"id"`).
			Limit(5).
			Find(&stakeChips).Error; err != nil {
			return nil, err
		}

		results = append(results, &schema.StakeStaking{
			Staker: common.HexToAddress(stakeStaking.Owner),
			Node:   common.HexToAddress(stakeStaking.Node),
			Value:  stakeStaking.Value,
			Chips: schema.StakeStakingChips{
				Total: stakeStaking.Count,
				Showcase: lo.Map(stakeChips, func(stakeChip *table.StakeChip, _ int) *schema.StakeChip {
					result, _ := stakeChip.Export()

					return result
				}),
			},
		})
	}

	return results, nil
}

func (c *client) FindStakeSnapshots(ctx context.Context) ([]*schema.StakeSnapshot, error) {
	databaseClient := c.database.WithContext(ctx)

	var stakeSnapshots []*table.StakeSnapshot

	if err := databaseClient.
		Order(`"date" DESC`).
		Limit(100). // TODO Replace this constant with a query parameter.
		Find(&stakeSnapshots).Error; err != nil {
		return nil, err
	}

	values := make([]*schema.StakeSnapshot, 0, len(stakeSnapshots))

	for _, stakeSnapshot := range stakeSnapshots {
		value, err := stakeSnapshot.Export()
		if err != nil {
			return nil, fmt.Errorf("export stake snapshot: %w", err)
		}

		values = append(values, value)
	}

	return values, nil
}

func (c *client) SaveStakeTransaction(ctx context.Context, stakeTransaction *schema.StakeTransaction) error {
	var value table.StakeTransaction
	if err := value.Import(*stakeTransaction); err != nil {
		return fmt.Errorf("import stake transaction: %w", err)
	}

	return c.database.WithContext(ctx).Create(&value).Error
}

func (c *client) SaveStakeEvent(ctx context.Context, stakeEvent *schema.StakeEvent) error {
	var value table.StakeEvent
	if err := value.Import(*stakeEvent); err != nil {
		return fmt.Errorf("import stake event: %w", err)
	}

	return c.database.WithContext(ctx).Create(&value).Error
}

func (c *client) SaveStakeChips(ctx context.Context, stakeChips ...*schema.StakeChip) error {
	values := make([]*table.StakeChip, 0, len(stakeChips))

	for _, stakeChip := range stakeChips {
		var value table.StakeChip

		if err := value.Import(*stakeChip); err != nil {
			return fmt.Errorf("import stake chip: %w", err)
		}

		values = append(values, &value)
	}

	return c.database.WithContext(ctx).Create(&values).Error
}

func (c *client) UpdateStakeChipsOwner(ctx context.Context, owner common.Address, stakeChipIDs ...*big.Int) error {
	ids := lo.Map(stakeChipIDs, func(stakeChipID *big.Int, _ int) decimal.Decimal {
		return decimal.NewFromBigInt(stakeChipID, 0)
	})

	return c.database.WithContext(ctx).Model((*table.StakeChip)(nil)).Where(`"id" IN ?`, ids).UpdateColumn("owner", owner.String()).Error
}

func (c *client) SaveStakeSnapshot(ctx context.Context, stakeSnapshot *schema.StakeSnapshot) error {
	databaseClient := c.database.WithContext(ctx)

	if err := databaseClient.
		Table((*table.StakeChip).TableName(nil)).
		Distinct(`"owner"`).
		Where(`"owner" != ?`, ethereum.AddressGenesis.String()).
		Count(&stakeSnapshot.Count).
		Error; err != nil {
		return fmt.Errorf("query count: %w", err)
	}

	var value table.StakeSnapshot
	if err := value.Import(*stakeSnapshot); err != nil {
		return fmt.Errorf("import stake snapshot: %w", err)
	}

	return databaseClient.
		Table((*table.StakeSnapshot).TableName(nil)).
		Create(stakeSnapshot).
		Error
}

package cockroachdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) SaveEpoch(ctx context.Context, epoch *schema.Epoch) error {
	// Save epoch.
	var data table.Epoch
	if err := data.Import(epoch); err != nil {
		zap.L().Error("import epoch", zap.Error(err), zap.Any("epoch", epoch))

		return err
	}

	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "transaction_hash",
			},
		},
		UpdateAll: true,
	}

	if err := c.database.WithContext(ctx).Clauses(onConflict).Create(&data).Error; err != nil {
		zap.L().Error("insert epoch", zap.Error(err), zap.Any("epoch", epoch))

		return err
	}

	// Save epoch items.
	var items table.EpochItems
	if err := items.Import(epoch.RewardedNodes); err != nil {
		zap.L().Error("import epoch items", zap.Error(err), zap.Any("epoch", epoch))

		return err
	}

	onConflict = clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "transaction_hash",
			},
			{
				Name: "index",
			},
		},
		UpdateAll: true,
	}

	if err := c.database.WithContext(ctx).Clauses(onConflict).CreateInBatches(&items, 500).Error; err != nil {
		zap.L().Error("insert epoch items", zap.Error(err), zap.Any("epoch", epoch))

		return err
	}

	return nil
}

func (c *client) FindEpochs(ctx context.Context, query *schema.FindEpochsQuery) ([]*schema.Epoch, error) {
	var data table.Epochs

	subQuery := c.database.WithContext(ctx).Model(&table.Epoch{})

	if query.Distinct != nil && *query.Distinct {
		subQuery = subQuery.Select("DISTINCT id")
	} else {
		subQuery = subQuery.Select("id")
	}

	if query.EpochID != nil {
		subQuery = subQuery.Where("id = ?", *query.EpochID)
	}

	if query.BlockNumber != nil {
		subQuery = subQuery.Where("block_number = ?", *query.BlockNumber)
	}

	if query.Finalized != nil {
		subQuery = subQuery.Where("finalized = ?", *query.Finalized)
	}

	if query.Cursor != nil {
		subQuery = subQuery.Where("id < ?", query.Cursor)
	}

	if query.Limit != nil {
		subQuery = subQuery.Limit(*query.Limit)
	}

	subQuery = subQuery.Order("id DESC")

	if err := c.database.WithContext(ctx).Model(&table.Epoch{}).Where("id IN (?)", subQuery).Order("id DESC, block_number DESC").Find(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		zap.L().Error("find epochs", zap.Error(err))

		return nil, err
	}

	return data.Export(nil)
}

func (c *client) FindEpochTransactions(ctx context.Context, id uint64, itemsLimit int, cursor *string) ([]*schema.Epoch, error) {
	// Find epoch transactions by id.
	databaseStatement := c.database.WithContext(ctx).Model(&table.Epoch{})

	if cursor != nil {
		var transaction *table.Epoch

		if err := c.database.WithContext(ctx).First(&transaction, "transaction_hash = ?", cursor).Error; err != nil {
			return nil, fmt.Errorf("find epoch cursor: %w", err)
		}

		databaseStatement = databaseStatement.Where("block_number < ?", transaction.BlockNumber)
	}

	var data table.Epochs

	if err := databaseStatement.Where("id = ?", id).Order("block_number DESC").Find(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		zap.L().Error("find epoch", zap.Error(err), zap.Uint64("id", id))

		return nil, err
	}

	// Find epoch items by transaction_hash.
	hashes := lo.Map(data, func(x *table.Epoch, _ int) string {
		return x.TransactionHash
	})

	var items table.EpochItems

	databaseStatement = c.database.WithContext(ctx).Model(&table.NodeRewardRecord{}).Where("transaction_hash IN (?)", hashes).Where("index <= ?", itemsLimit)

	if err := databaseStatement.Order("index ASC").Limit(itemsLimit).Find(&items).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("find epoch items", zap.Error(err), zap.Any("hashes", hashes))

			return data.Export(nil)
		}

		zap.L().Error("find epoch items", zap.Error(err), zap.Uint64("id", id))

		return nil, err
	}

	epochItems, err := items.Export()
	if err != nil {
		zap.L().Error("export epoch items", zap.Error(err), zap.Uint64("id", id))

		return nil, err
	}

	return data.Export(epochItems)
}

func (c *client) FindEpochTransaction(ctx context.Context, transactionHash common.Hash, itemsLimit int, cursor *string) (*schema.Epoch, error) {
	var data table.Epoch

	if err := c.database.WithContext(ctx).Model(&table.Epoch{}).Where("transaction_hash = ?", transactionHash.String()).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		zap.L().Error("find epoch", zap.Error(err), zap.Any("transactionHash", transactionHash))

		return nil, err
	}

	// Find epoch items by transaction_hash.
	var items table.EpochItems

	databaseStatement := c.database.WithContext(ctx).Model(&table.NodeRewardRecord{}).Where("transaction_hash = ?", transactionHash.String())

	if cursor != nil {
		databaseStatement = databaseStatement.Where("index > ?", cursor)
	}

	if err := databaseStatement.Order("index ASC").Limit(itemsLimit).Find(&items).Error; err != nil {
		zap.L().Error("find epoch items", zap.Error(err), zap.Any("transaction_hash", transactionHash))

		return nil, err
	}

	epochItems, err := items.Export()
	if err != nil {
		zap.L().Error("export epoch items", zap.Error(err), zap.Any("transaction_hash", transactionHash))

		return nil, err
	}

	return data.Export(epochItems)
}

func (c *client) FindEpochNodeRewards(ctx context.Context, nodeAddress common.Address, limit int, cursor *string) ([]*schema.Epoch, error) {
	// Find epoch items by nodeAddress.
	var items table.EpochItems

	databaseStatement := c.database.WithContext(ctx).Model(&table.NodeRewardRecord{}).Where("node_address = ?", nodeAddress.String())

	if cursor != nil {
		databaseStatement = databaseStatement.Where("epoch_id < ?", cursor)
	}

	if err := databaseStatement.Limit(limit).Order("epoch_id DESC, index ASC").Find(&items).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		zap.L().Error("find epoch items", zap.Error(err), zap.String("nodeAddress", nodeAddress.String()))

		return nil, err
	}

	epochIDs := make([]uint64, 0, len(items))
	itemsMap := make(map[uint64][]*schema.RewardedNode, len(items))

	for _, item := range items {
		data, err := item.Export()
		if err != nil {
			zap.L().Error("export epoch item", zap.Error(err), zap.String("nodeAddress", nodeAddress.String()), zap.Any("item", item))

			return nil, err
		}

		if _, ok := itemsMap[item.EpochID]; !ok {
			itemsMap[item.EpochID] = make([]*schema.RewardedNode, 0, 1)
		}

		itemsMap[item.EpochID] = append(itemsMap[item.EpochID], data)
		epochIDs = append(epochIDs, item.EpochID)
	}

	// Find epochs by epochIDs.
	var epochs table.Epochs

	if err := c.database.WithContext(ctx).Model(&table.Epoch{}).Where("id IN ?", epochIDs).Order("id DESC").Find(&epochs).Error; err != nil {
		zap.L().Error("find epochs", zap.Error(err), zap.Any("epochIDs", epochIDs))

		return nil, err
	}

	result := make([]*schema.Epoch, 0, len(epochs))

	for _, epoch := range epochs {
		data, err := epoch.Export(itemsMap[epoch.ID])
		if err != nil {
			zap.L().Error("export epoch", zap.Error(err), zap.Any("epoch", epoch))

			return nil, err
		}

		result = append(result, data)
	}

	return result, nil
}

func (c *client) UpdateEpochsFinalizedByBlockNumber(ctx context.Context, blockNumber uint64) error {
	return c.database.
		WithContext(ctx).
		Table((*table.Epoch).TableName(nil)).
		Where(`"block_number" < ? AND NOT "finalized"`, blockNumber).
		Update("finalized", true).
		Error
}

func (c *client) DeleteEpochsByBlockNumber(ctx context.Context, blockNumber uint64) error {
	epoch, err := c.FindEpochs(ctx, &schema.FindEpochsQuery{
		BlockNumber: lo.ToPtr(blockNumber),
	})
	if err != nil {
		zap.L().Error("find epochs by block number", zap.Error(err), zap.Uint64("blockNumber", blockNumber))

		return err
	}

	if len(epoch) == 0 {
		return nil
	}

	if err = c.database.WithContext(ctx).Where(`block_number = ? AND NOT "finalized"`, blockNumber).Delete(&table.Epoch{}).Error; err != nil {
		zap.L().Error("delete epochs by block number", zap.Error(err), zap.Uint64("blockNumber", blockNumber))

		return err
	}

	transactionHashes := lo.Map(epoch, func(x *schema.Epoch, _ int) string {
		return x.TransactionHash.String()
	})

	if err = c.database.WithContext(ctx).Where("transaction_hash IN (?)", transactionHashes).Delete(&table.NodeRewardRecord{}).Error; err != nil {
		zap.L().Error("delete epoch items by block number", zap.Error(err), zap.Uint64("blockNumber", blockNumber))

		return err
	}

	return nil
}

func (c *client) SaveEpochTrigger(ctx context.Context, epochTrigger *schema.EpochTrigger) error {
	// Save epoch trigger.
	var data table.EpochTrigger
	if err := data.Import(epochTrigger); err != nil {
		zap.L().Error("import epoch trigger", zap.Error(err), zap.Any("epochTrigger", epochTrigger))

		return err
	}

	if err := c.database.WithContext(ctx).Create(&data).Error; err != nil {
		zap.L().Error("insert epoch trigger", zap.Error(err), zap.Any("epochTrigger", epochTrigger))

		return err
	}

	return nil
}

func (c *client) FindLatestEpochTrigger(ctx context.Context) (*schema.EpochTrigger, error) {
	var data table.EpochTrigger

	if err := c.database.WithContext(ctx).Order("created_at DESC").First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		zap.L().Error("find latest epoch trigger", zap.Error(err))

		return nil, err
	}

	return data.Export()
}

func (c *client) FindEpochTriggers(ctx context.Context, epochID uint64) ([]*schema.EpochTrigger, error) {
	var data table.EpochTriggers

	if err := c.database.WithContext(ctx).Model(&table.EpochTrigger{}).Where("epoch_id = ?", epochID).Order("created_at ASC").Find(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		zap.L().Error("find epoch triggers", zap.Error(err), zap.Uint64("epochID", epochID))

		return nil, err
	}

	return data.Export()
}

func (c *client) FindEpochAPYSnapshots(ctx context.Context, query schema.EpochAPYSnapshotQuery) ([]*schema.EpochAPYSnapshot, error) {
	var data table.EpochAPYSnapshots

	databaseStatement := c.database.WithContext(ctx).Model(&table.EpochAPYSnapshot{})

	if query.EpochID != nil {
		databaseStatement = databaseStatement.Where("epoch_id = ?", *query.EpochID)
	}

	if query.Limit != nil {
		databaseStatement = databaseStatement.Limit(*query.Limit)
	}

	if err := databaseStatement.Order("epoch_id DESC").Find(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		zap.L().Error("find epoch apy snapshots", zap.Error(err), zap.Any("query", query))

		return nil, err
	}

	return data.Export()
}

func (c *client) SaveEpochAPYSnapshot(ctx context.Context, epochAPYSnapshot *schema.EpochAPYSnapshot) error {
	// Save epoch APY snapshot.
	var data table.EpochAPYSnapshot
	if err := data.Import(epochAPYSnapshot); err != nil {
		zap.L().Error("import epoch APY snapshot", zap.Error(err), zap.Any("epochAPYSnapshot", epochAPYSnapshot))

		return err
	}

	if err := c.database.WithContext(ctx).Create(&data).Error; err != nil {
		zap.L().Error("insert epoch APY snapshot", zap.Error(err), zap.Any("epochAPYSnapshot", epochAPYSnapshot))

		return err
	}

	return nil
}

func (c *client) FindEpochAPYSnapshotsAverage(ctx context.Context) (decimal.Decimal, error) {
	var avgAPY decimal.Decimal

	if err := c.database.WithContext(ctx).Model(&table.EpochAPYSnapshot{}).
		Select("AVG(apy) as avg_apy").Row().Scan(&avgAPY); err != nil {
		zap.L().Error("retrieve and calculate average APY", zap.Error(err))

		return decimal.Zero, fmt.Errorf("retrieve and calculate average APY: %w", err)
	}

	return avgAPY, nil
}

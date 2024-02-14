package cockroachdb

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
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
	if err := items.Import(epoch.RewardItems); err != nil {
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

func (c *client) FindEpochs(ctx context.Context, limit int, cursor *string) ([]*schema.Epoch, error) {
	var data table.Epochs

	databaseStatement := c.database.WithContext(ctx).Model(&table.Epoch{})

	if cursor != nil {
		databaseStatement = databaseStatement.Where("id < ?", cursor)
	}

	if err := databaseStatement.Order("id DESC").Limit(limit).Find(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		zap.L().Error("find epochs", zap.Error(err))

		return nil, err
	}

	return data.Export()
}

func (c *client) FindEpoch(ctx context.Context, id uint64, itemsLimit int, cursor *string) (*schema.Epoch, error) {
	var data table.Epoch

	if err := c.database.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, database.ErrorRowNotFound
		}

		zap.L().Error("find epoch", zap.Error(err), zap.Uint64("id", id))

		return nil, err
	}

	var items table.EpochItems

	databaseStatement := c.database.WithContext(ctx).Model(&table.EpochItem{}).Where("epoch_id = ?", id)

	if cursor != nil {
		databaseStatement = databaseStatement.Where("index > ?", cursor)
	}

	if err := databaseStatement.Order("index ASC").Limit(itemsLimit).Find(&items).Error; err != nil {
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

func (c *client) FindEpochNodeRewards(ctx context.Context, nodeAddress common.Address, limit int, cursor *string) ([]*schema.Epoch, error) {
	// Find epoch items by nodeAddress.
	var items table.EpochItems

	databaseStatement := c.database.WithContext(ctx).Model(&table.EpochItem{}).Where("node_address = ?", nodeAddress.String())

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
	itemsMap := make(map[uint64][]*schema.EpochItem, len(items))

	for _, item := range items {
		data, err := item.Export()
		if err != nil {
			zap.L().Error("export epoch item", zap.Error(err), zap.String("nodeAddress", nodeAddress.String()), zap.Any("item", item))

			return nil, err
		}

		if _, ok := itemsMap[item.EpochID]; !ok {
			itemsMap[item.EpochID] = make([]*schema.EpochItem, 0, 1)
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

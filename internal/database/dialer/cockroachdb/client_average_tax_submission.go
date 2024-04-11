package cockroachdb

import (
	"context"

	"github.com/rss3-network/global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/rss3-network/global-indexer/schema"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

// SaveAverageTaxSubmission Save records of average tax submissions
func (c *client) SaveAverageTaxSubmission(ctx context.Context, submission *schema.AverageTaxRateSubmission) error {
	var data table.AverageTaxRateSubmission
	if err := data.Import(submission); err != nil {
		zap.L().Error("import average tax submission", zap.Error(err), zap.Any("submission", submission))

		return err
	}

	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "epoch_id",
			},
		},
		UpdateAll: true,
	}

	if err := c.database.WithContext(ctx).Clauses(onConflict).Create(&data).Error; err != nil {
		zap.L().Error("insert average tax submission", zap.Error(err), zap.Any("submission", submission))

		return err
	}

	return nil
}

// FindAverageTaxSubmissions Find records of average tax submissions
func (c *client) FindAverageTaxSubmissions(ctx context.Context, query schema.AverageTaxRateSubmissionQuery) ([]*schema.AverageTaxRateSubmission, error) {
	databaseStatement := c.database.WithContext(ctx).Table((*table.AverageTaxRateSubmission).TableName(nil))

	if query.EpochID != nil {
		databaseStatement = databaseStatement.Where("epoch_id = ?", *query.EpochID)
	}

	if query.Limit != nil {
		databaseStatement = databaseStatement.Limit(*query.Limit)
	}

	var submissions table.AverageTaxSubmissions

	if err := databaseStatement.Order("epoch_id DESC").Find(&submissions).Error; err != nil {
		zap.L().Error("find average tax submissions", zap.Error(err), zap.Any("query", query))

		return nil, err
	}

	return submissions.Export()
}

package distributor

import (
	"context"

	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"go.uber.org/zap"
)

// processRSSHubResults processes responses for RSSHub requests.
func (d *Distributor) processRSSHubResponses(responses []*model.DataResponse) {
	if err := d.simpleEnforcer.VerifyResponses(context.Background(), responses); err != nil {
		zap.L().Error("fail to verify rss hub responses", zap.Any("responses", len(responses)))
	} else {
		zap.L().Info("complete rss hub responses verify", zap.Any("responses", len(responses)))
	}
}

// processActivityResults processes responses for Activity requests.
func (d *Distributor) processActivityResponses(responses []*model.DataResponse) {
	if err := d.simpleEnforcer.VerifyResponses(context.Background(), responses); err != nil {
		zap.L().Error("fail to verify activity id responses ", zap.Any("responses", len(responses)))
	} else {
		zap.L().Info("complete activity id responses verify", zap.Any("responses", len(responses)))
	}
}

// processActivitiesResponses processes responses for Activities requests.
func (d *Distributor) processActivitiesResponses(responses []*model.DataResponse) {
	ctx := context.Background()

	if err := d.simpleEnforcer.VerifyResponses(ctx, responses); err != nil {
		zap.L().Error("fail to verify activity responses", zap.Any("responses", len(responses)))

		return
	}

	zap.L().Info("complete activity responses verify", zap.Any("responses", len(responses)))

	d.simpleEnforcer.VerifyPartialResponses(ctx, responses)
}

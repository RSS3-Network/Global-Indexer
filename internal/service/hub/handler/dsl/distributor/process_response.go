package distributor

import (
	"context"

	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"go.uber.org/zap"
)

// processRSSHubResults processes the RSS Hub responses.
func (d *Distributor) processRSSHubResponses(responses []*model.DataResponse) {
	if err := d.simpleEnforcer.VerifyResponses(context.Background(), responses); err != nil {
		zap.L().Error("fail to verify rss hub responses", zap.Any("responses", len(responses)))
	} else {
		zap.L().Info("complete rss hub responses verify", zap.Any("responses", len(responses)))
	}
}

// processActivityResults processes activity data retrieval responses.
func (d *Distributor) processActivityResponses(responses []*model.DataResponse) {
	if err := d.simpleEnforcer.VerifyResponses(context.Background(), responses); err != nil {
		zap.L().Error("fail to verify activity id responses ", zap.Any("responses", len(responses)))
	} else {
		zap.L().Info("complete activity id responses verify", zap.Any("responses", len(responses)))
	}
}

// processActivitiesResults processes account activities data retrieval responses.
func (d *Distributor) processActivitiesResponses(responses []*model.DataResponse) {
	ctx := context.Background()

	if err := d.simpleEnforcer.VerifyResponses(ctx, responses); err != nil {
		zap.L().Error("fail to verify activity responses", zap.Any("responses", len(responses)))

		return
	}

	zap.L().Info("complete activity responses verify", zap.Any("responses", len(responses)))

	d.simpleEnforcer.VerifyPartialResponses(ctx, responses)
}

package distributor

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// processRSSHubResults processes responses for RSSHub requests.
func (d *Distributor) processRSSHubResponses(responses []*model.DataResponse) {
	if err := d.simpleEnforcer.VerifyResponses(context.Background(), responses); err != nil {
		zap.L().Error("fail to verify rss hub responses", zap.Any("responses", len(responses)))
	} else {
		_ = d.processNodeFailureResponse(context.Background(), responses)

		zap.L().Info("complete rss hub responses verify", zap.Any("responses", len(responses)))
	}
}

// processActivityResults processes responses for Activity requests.
func (d *Distributor) processActivityResponses(responses []*model.DataResponse) {
	if err := d.simpleEnforcer.VerifyResponses(context.Background(), responses); err != nil {
		zap.L().Error("fail to verify activity id responses ", zap.Any("responses", len(responses)))
	} else {
		_ = d.processNodeFailureResponse(context.Background(), responses)

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

	epochID := d.processNodeFailureResponse(context.Background(), responses)

	zap.L().Info("complete activity responses verify", zap.Any("responses", len(responses)))

	d.simpleEnforcer.VerifyPartialResponses(ctx, epochID, responses)
}

// processNodeFailureResponse finds the valid response data and saves the failure responses.
func (d *Distributor) processNodeFailureResponse(ctx context.Context, responses []*model.DataResponse) uint64 {
	epochID, err := d.getLatestEpochID(ctx)
	if err != nil {
		zap.L().Error("get latest epoch event from database", zap.Error(err))
		return 0
	}

	validatorNode, validatorRequest, validatorResponse := getValidResponseData(responses)

	d.saveFailureResponses(ctx, epochID, validatorNode, validatorRequest, validatorResponse, responses)

	return epochID
}

// getLatestEpochID returns the recent epoch ID.
func (d *Distributor) getLatestEpochID(ctx context.Context) (uint64, error) {
	epochEvent, err := d.databaseClient.FindEpochs(ctx, 1, nil)
	if err != nil && !errors.Is(err, database.ErrorRowNotFound) {
		return 0, err
	}

	if len(epochEvent) > 0 {
		return epochEvent[0].ID, nil
	}

	return 0, nil
}

// getValidResponseData returns the valid response data which valid points are greater than 0.
func getValidResponseData(responses []*model.DataResponse) (common.Address, string, json.RawMessage) {
	for _, response := range responses {
		if response.ValidPoint > 0 {
			return response.Address, response.Endpoint, response.Data
		}
	}

	return common.Address{}, "", nil
}

// saveFailureResponses saves the responses which invalid points are greater than 0 and status is challengeable.
func (d *Distributor) saveFailureResponses(ctx context.Context, epochID uint64, validatorNode common.Address, validatorRequest string, validatorResponse json.RawMessage, responses []*model.DataResponse) {
	for _, response := range responses {
		if response.InvalidPoint > 0 {
			err := d.databaseClient.SaveNodeFailureResponse(ctx, &schema.NodeFailureResponse{
				EpochID:           epochID,
				Status:            schema.NodeFailureResponseStatusChallengeable,
				ValidatorNode:     validatorNode,
				ValidatorRequest:  validatorRequest,
				ValidatorResponse: validatorResponse,
				VerifiedNode:      response.Address,
				VerifiedRequest:   response.Endpoint,
				VerifiedResponse:  lo.Ternary(response.Err != nil, json.RawMessage(response.Err.Error()), response.Data),
			})
			if err != nil {
				zap.L().Error("save node failure response", zap.Error(err))
			}
		}
	}
}

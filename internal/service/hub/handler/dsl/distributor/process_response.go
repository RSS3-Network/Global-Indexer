package distributor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

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
		_ = d.processNodeInvalidResponse(context.Background(), responses)

		zap.L().Info("complete rss hub responses verify", zap.Any("responses", len(responses)))
	}
}

// processActivityResults processes responses for Activity requests.
func (d *Distributor) processActivityResponses(responses []*model.DataResponse) {
	if err := d.simpleEnforcer.VerifyResponses(context.Background(), responses); err != nil {
		zap.L().Error("fail to verify activity id responses ", zap.Any("responses", len(responses)))
	} else {
		_ = d.processNodeInvalidResponse(context.Background(), responses)

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

	epochID := d.processNodeInvalidResponse(context.Background(), responses)

	if epochID == 0 {
		return
	}

	zap.L().Info("complete activity responses verify", zap.Any("responses", len(responses)))

	d.simpleEnforcer.VerifyPartialResponses(ctx, epochID, responses)
}

// processNodeInvalidResponse finds the valid response data and saves the invalid responses.
func (d *Distributor) processNodeInvalidResponse(ctx context.Context, responses []*model.DataResponse) uint64 {
	validatorNodes, request, validatorResponse, err := getValidResponseData(responses)
	if err != nil {
		zap.L().Error("get valid response data", zap.Error(err))
		return 0
	}

	// If all responses are valid, return 0.
	if len(validatorNodes) == len(responses) {
		return 0
	}

	epochID, err := d.getLatestEpochID(ctx)
	if err != nil {
		zap.L().Error("get latest epoch event from database", zap.Error(err))
		return 0
	}

	d.saveInvalidResponses(ctx, epochID, validatorNodes, request, validatorResponse, responses)

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
func getValidResponseData(responses []*model.DataResponse) ([]common.Address, string, json.RawMessage, error) {
	var (
		validatorNodes []common.Address
		data           json.RawMessage
	)

	for _, response := range responses {
		if response.ValidPoint > 0 {
			validatorNodes = append(validatorNodes, response.Address)
			data = response.Data
		}
	}

	request, err := extractPathAndParams(responses[0].Endpoint)
	if err != nil {
		return nil, "", nil, err
	}

	return validatorNodes, request, data, nil
}

// extractPathAndParams extracts the path and params from the endpoint.
func extractPathAndParams(endpoint string) (string, error) {
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", err
	}

	return strings.TrimPrefix(endpoint, parsedURL.Scheme+"://"+parsedURL.Host), nil
}

// saveInvalidResponses saves the responses which invalid points are greater than 0 and status is challengeable.
func (d *Distributor) saveInvalidResponses(ctx context.Context, epochID uint64, validatorNodes []common.Address, request string, validatorResponse json.RawMessage, responses []*model.DataResponse) {
	for _, response := range responses {
		if response.InvalidPoint > 0 {
			err := d.databaseClient.SaveNodeInvalidResponse(ctx, &schema.NodeInvalidResponse{
				EpochID:           epochID,
				InvalidType:       lo.Ternary(response.Err != nil, schema.NodeInvalidResponseTypeError, schema.NodeInvalidResponseTypeInconsistent),
				ValidatorNodes:    validatorNodes,
				Request:           request,
				ValidatorResponse: validatorResponse,
				Node:              response.Address,
				InvalidResponse:   lo.Ternary(response.Err != nil, json.RawMessage(response.Err.Error()), response.Data),
			})

			if err != nil {
				zap.L().Error("save node invalid response", zap.Error(err))
			}
		}
	}
}

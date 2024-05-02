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
		validatorNodes = make([]common.Address, 0)
		data           = json.RawMessage("{}")
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
		zap.L().Error("parsing URL", zap.Error(err))
		return "", err
	}

	return strings.TrimPrefix(endpoint, parsedURL.Scheme+"://"+parsedURL.Host), nil
}

// saveInvalidResponses saves the responses which invalid points are greater than 0.
func (d *Distributor) saveInvalidResponses(ctx context.Context, epochID uint64, validatorNodes []common.Address, request string, validatorResponse json.RawMessage, responses []*model.DataResponse) {
	var (
		nodeInvalidResponses = make([]*schema.NodeInvalidResponse, 0, len(responses))
		err                  error
	)

	for _, response := range responses {
		if response.InvalidPoint == 0 {
			continue
		}

		typeValue := schema.NodeInvalidResponseTypeInconsistent
		responseValue := response.Data

		if response.Err != nil {
			typeValue = schema.NodeInvalidResponseTypeError
			responseValue, err = json.Marshal(fmt.Sprintf(`{"error_message": "%s"}`, response.Err))

			if err != nil {
				zap.L().Error("json marshaling", zap.Error(err))

				responseValue = json.RawMessage(`{"error_message": "error response"}`)
			}
		}

		nodeInvalidResponse := &schema.NodeInvalidResponse{
			EpochID:           epochID,
			Type:              typeValue,
			ValidatorNodes:    validatorNodes,
			Request:           request,
			ValidatorResponse: validatorResponse,
			Node:              response.Address,
			Response:          responseValue,
		}

		nodeInvalidResponses = append(nodeInvalidResponses, nodeInvalidResponse)
	}

	if err := d.databaseClient.SaveNodeInvalidResponses(ctx, nodeInvalidResponses); err != nil {
		zap.L().Error("save node invalid response", zap.Error(err))
	}
}

package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/form/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/httpx"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"go.uber.org/zap"
)

type Router interface {
	// BuildPath builds the path for the request and returns a map of node addresses to their respective paths
	BuildPath(path string, query any, nodes []model.Cache) (map[common.Address]string, error)
	// DistributeRequest sends the request to the nodes and processes the results
	DistributeRequest(ctx context.Context, nodeMap map[common.Address]string, processResults func([]model.DataResponse)) (model.DataResponse, error)
}

type SimpleRouter struct {
	httpClient httpx.Client
}

func (r *SimpleRouter) BuildPath(path string, query any, nodes []model.Cache) (map[common.Address]string, error) {
	if query != nil {
		values, err := form.NewEncoder().Encode(query)

		if err != nil {
			return nil, fmt.Errorf("build params %w", err)
		}

		path = fmt.Sprintf("%s?%s", path, values.Encode())
	}

	urls := make(map[common.Address]string, len(nodes))

	for _, node := range nodes {
		fullURL, err := url.JoinPath(node.Endpoint, path)
		if err != nil {
			return nil, fmt.Errorf("failed to join path for node %s: %w", node.Address, err)
		}

		decodedURL, err := url.QueryUnescape(fullURL)
		if err != nil {
			return nil, fmt.Errorf("failed to unescape url for node %s: %w", node.Address, err)
		}

		urls[common.HexToAddress(node.Address)] = decodedURL
	}

	return urls, nil
}

func (r *SimpleRouter) DistributeRequest(ctx context.Context, nodeMap map[common.Address]string, processResults func([]*model.DataResponse)) (model.DataResponse, error) {
	// firstResult is a channel that will be used to send the first result
	var firstResult = make(chan model.DataResponse, 1)

	// Distribute the request to the nodes
	r.distribute(ctx, nodeMap, processResults, firstResult)

	select {
	case result := <-firstResult:
		close(firstResult)
		return result, nil
	case <-ctx.Done():
		return model.DataResponse{Err: fmt.Errorf("failed to retrieve node data, please retry")}, nil
	}
}

func (r *SimpleRouter) distribute(ctx context.Context, nodeMap map[common.Address]string, processResults func([]*model.DataResponse), firstResult chan<- model.DataResponse) {
	var (
		waitGroup sync.WaitGroup
		mu        sync.Mutex

		// results contains all the returned results
		results []*model.DataResponse
		// resultSent is used to ensure that the first result is sent only once
		resultSent bool
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for address, endpoint := range nodeMap {
		waitGroup.Add(1)

		go func(address common.Address, endpoint string) {
			defer waitGroup.Done()

			result := &model.DataResponse{Address: address}
			// Fetch the data from the node.
			body, err := r.httpClient.Fetch(ctx, endpoint)

			if err != nil {
				zap.L().Error("failed to fetch request", zap.String("node", address.String()), zap.Error(err))

				result.Err = err
			} else {
				// Read the response body.
				data, readErr := io.ReadAll(body)

				if readErr != nil {
					zap.L().Error("failed to read response body", zap.String("node", address.String()), zap.Error(readErr))

					result.Err = readErr
				} else {
					activity := &model.ActivityResponse{}
					activities := &model.ActivitiesResponse{}

					// Check if the node's data is valid.
					if !validateData(data, activity) && !validateData(data, activities) {
						zap.L().Error("failed to parse response", zap.String("node", address.String()))

						result.Err = fmt.Errorf("invalid data")
					} else {
						// If the data is non-null, set the result as valid.
						if activity.Data != nil || activities.Data != nil {
							result.Valid = true
						}

						result.Data = data
					}
				}
			}

			sendResult(&mu, &results, result, &resultSent, firstResult, len(nodeMap))
		}(address, endpoint)
	}

	waitGroup.Wait()
	processResults(results)
}

func sendResult(mu *sync.Mutex, results *[]*model.DataResponse, result *model.DataResponse, resultSent *bool, firstResult chan<- model.DataResponse, nodeMapLen int) {
	mu.Lock()
	defer mu.Unlock()

	*results = append(*results, result)

	if !*resultSent {
		// If the result is valid (no error and contains data), send it as the first valid result.
		if result.Err == nil && result.Valid {
			firstResult <- *result

			*resultSent = true

			return
		}

		// If all the results have been received
		if len(*results) == nodeMapLen {
			for _, res := range *results {
				if res.Err != nil {
					continue
				}

				firstResult <- *res

				return
			}

			firstResult <- *(*results)[0]
		}
	}
}

func validateData(data []byte, target any) bool {
	var errRes model.ErrResponse
	if err := json.Unmarshal(data, &errRes); err == nil && errRes.ErrorCode != "" {
		return false
	}

	var notFound model.NotFoundResponse
	if err := json.Unmarshal(data, &notFound); err == nil && notFound.Message != "" {
		return false
	}

	if err := json.Unmarshal(data, target); err == nil {
		return true
	}

	return false
}

func NewSimpleRouter(httpClient httpx.Client) (*SimpleRouter, error) {
	return &SimpleRouter{
		httpClient: httpClient,
	}, nil
}

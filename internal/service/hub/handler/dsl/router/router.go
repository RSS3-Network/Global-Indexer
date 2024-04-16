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
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"go.uber.org/zap"
)

type Router interface {
	// BuildPath builds the path for the request and returns a map of node addresses to their respective paths
	BuildPath(path string, query any, nodes []model.NodeEndpointCache) (map[common.Address]string, error)
	// DistributeRequest sends the request to the nodes and processes the results
	DistributeRequest(ctx context.Context, nodeMap map[common.Address]string, processResults func([]model.DataResponse)) (model.DataResponse, error)
}

type SimpleRouter struct {
	httpClient httputil.Client
}

func (r *SimpleRouter) BuildPath(path string, query any, nodes []model.NodeEndpointCache) (map[common.Address]string, error) {
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

func (r *SimpleRouter) DistributeRequest(ctx context.Context, nodeMap map[common.Address]string, processResponses func([]*model.DataResponse)) (model.DataResponse, error) {
	// firstResponse is a channel that will be used to send the first response
	var firstResponse = make(chan model.DataResponse, 1)

	// Distribute the request to the nodes
	r.distribute(ctx, nodeMap, processResponses, firstResponse)

	select {
	case response := <-firstResponse:
		close(firstResponse)
		return response, nil
	case <-ctx.Done():
		return model.DataResponse{Err: fmt.Errorf("failed to retrieve node data, please retry")}, nil
	}
}

// distribute sends the request to the nodes and processes the responses
func (r *SimpleRouter) distribute(ctx context.Context, nodeMap map[common.Address]string, processResponses func([]*model.DataResponse), firstResponse chan<- model.DataResponse) {
	var (
		waitGroup sync.WaitGroup
		mu        sync.Mutex

		// responses contains all the returned responses
		responses []*model.DataResponse
		// responseSent is used to ensure that the first response is sent only once
		responseSent bool
	)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, httputil.DefaultTimeout)

	defer cancel()

	for address, endpoint := range nodeMap {
		waitGroup.Add(1)

		go func(address common.Address, endpoint string) {
			defer waitGroup.Done()

			response := &model.DataResponse{Address: address}
			// Fetch the data from the node.
			body, err := r.httpClient.Fetch(ctx, endpoint)

			if err != nil {
				zap.L().Error("failed to fetch request", zap.String("node", address.String()), zap.Error(err))

				response.Err = err
			} else {
				// Read the response body.
				data, readErr := io.ReadAll(body)

				if readErr != nil {
					zap.L().Error("failed to read response body", zap.String("node", address.String()), zap.Error(readErr))

					response.Err = readErr
				} else {
					activity := &model.ActivityResponse{}
					activities := &model.ActivitiesResponse{}

					// Check if the node's data is valid.
					if !validateData(data, activity) && !validateData(data, activities) {
						zap.L().Error("failed to parse response", zap.String("node", address.String()))

						response.Err = fmt.Errorf("invalid data")
					} else {
						// If the data is non-null, set the result as valid.
						if activity.Data != nil || activities.Data != nil {
							response.Valid = true
						}

						response.Data = data
					}
				}
			}

			sendResponse(&mu, &responses, response, &responseSent, firstResponse, len(nodeMap))
		}(address, endpoint)
	}

	waitGroup.Wait()
	// Process the responses to calculate the actual request of each node
	processResponses(responses)
}

// sendResponse sends the first valid response to the firstResponse channel
// If all the responses are invalid, the first response will be the first response received
func sendResponse(mu *sync.Mutex, responses *[]*model.DataResponse, response *model.DataResponse, responseSent *bool, firstResponse chan<- model.DataResponse, nodesRequested int) {
	mu.Lock()
	defer mu.Unlock()

	*responses = append(*responses, response)

	if !*responseSent {
		// If the response is valid (no error and contains data), send it as the first valid response.
		if response.Err == nil && response.Valid {
			firstResponse <- *response

			*responseSent = true

			return
		}

		// If all the results have been received
		if len(*responses) == nodesRequested {
			for _, res := range *responses {
				if res.Err != nil {
					continue
				}

				firstResponse <- *res

				return
			}

			firstResponse <- *(*responses)[0]
		}
	}
}

func validateData(data []byte, target any) bool {
	var errRes model.ErrResponse
	if err := json.Unmarshal(data, &errRes); err == nil && errRes.ErrorCode != "" {
		return false
	}

	if err := json.Unmarshal(data, target); err == nil {
		return true
	}

	return false
}

func NewSimpleRouter(httpClient httputil.Client) *SimpleRouter {
	return &SimpleRouter{
		httpClient: httpClient,
	}
}

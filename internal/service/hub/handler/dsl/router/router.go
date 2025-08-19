package router

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"go.uber.org/zap"
)

type Router interface {
	// BuildPath builds the path for the request and returns a map of Node addresses to their respective paths
	BuildPath(path string, query any, nodes []model.NodeEndpointCache) (map[common.Address]string, error)
	// DistributeRequest sends the request to the Nodes and processes the results
	DistributeRequest(ctx context.Context, nodeMap map[common.Address]string, processResults func([]model.DataResponse)) (model.DataResponse, error)
}

type SimpleRouter struct {
	httpClient httputil.Client
}

func (r *SimpleRouter) BuildPath(method, path string, query url.Values, nodes []*model.NodeEndpointCache, body []byte, isRssNode bool) (map[common.Address]model.RequestMeta, error) {
	if method == http.MethodGet && query != nil {
		path = fmt.Sprintf("%s?%s", path, query.Encode())
	}

	urls := make(map[common.Address]model.RequestMeta, len(nodes))

	for _, node := range nodes {
		fullURL := buildFullURL(node.Endpoint, path)

		if method != http.MethodPost {
			decodedURL, err := url.QueryUnescape(fullURL)
			if err != nil {
				return nil, fmt.Errorf("failed to unescape url for node %s: %w", node.Address, err)
			}

			fullURL = decodedURL
		}

		urls[common.HexToAddress(node.Address)] = model.RequestMeta{
			Method:      method,
			Endpoint:    fullURL,
			AccessToken: node.AccessToken,
			Body:        body,
			IsRssNode:   isRssNode,
		}
	}

	return urls, nil
}

func buildFullURL(endpoint, urlPath string) string {
	// ensure the endpoint ends with a "/"
	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}

	// ensure the urlPath does not start with a "/"
	urlPath = strings.TrimPrefix(urlPath, "/")

	// join the endpoint and urlPath
	return endpoint + urlPath
}

func (r *SimpleRouter) DistributeRequest(ctx context.Context, nodeMap map[common.Address]model.RequestMeta, processResponses func([]*model.DataResponse)) (model.DataResponse, error) {
	// firstResponse is a channel that will be used to send the first response
	var firstResponse = make(chan model.DataResponse, 1)

	// Distribute the request to the Nodes
	r.distribute(ctx, nodeMap, processResponses, firstResponse)

	select {
	case response := <-firstResponse:
		close(firstResponse)
		return response, nil
	case <-ctx.Done():
		return model.DataResponse{Err: fmt.Errorf("failed to retrieve node data, please retry")}, nil
	}
}

// distribute sends the request to the Nodes and processes the responses
func (r *SimpleRouter) distribute(ctx context.Context, nodeMap map[common.Address]model.RequestMeta, processResponses func([]*model.DataResponse), firstResponse chan<- model.DataResponse) {
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

	for address, requestMeta := range nodeMap {
		waitGroup.Add(1)

		go func(address common.Address, requestMeta model.RequestMeta) {
			defer waitGroup.Done()

			response := &model.DataResponse{Address: address, Endpoint: requestMeta.Endpoint}
			// Fetch the data from the Node.
			body, headers, err := r.httpClient.FetchWithMethod(ctx, requestMeta.Method, requestMeta.Endpoint, requestMeta.AccessToken, bytes.NewReader(requestMeta.Body))

			if err != nil {
				zap.L().Error("failed to fetch request", zap.String("node", address.String()), zap.Error(err))

				response.Err = err
			} else {
				// Read the response body.
				data, readErr := io.ReadAll(body)

				zap.L().Info("fetch request", zap.String("node", address.String()), zap.String("endpoint", requestMeta.Endpoint), zap.String("method", requestMeta.Method))

				if readErr != nil {
					zap.L().Error("failed to read response body", zap.String("node", address.String()), zap.Error(readErr))

					response.Err = readErr
				} else {
					if requestMeta.IsRssNode {
						response.Data = data
						response.IsRssNode = requestMeta.IsRssNode
						response.Etag = headers.Get("Etag")
					} else {
						var v interface{}
						err = json.Unmarshal(data, &v)

						if err != nil {
							zap.L().Error("failed to unmarshal response body", zap.String("node", address.String()), zap.Error(err))

							response.Err = fmt.Errorf("invalid data")
						}

						if _, ok := v.([]interface{}); ok {
							zap.L().Info("response is an array", zap.String("node", address.String()))

							response.Data = data
						} else {
							activity := &model.ActivityResponse{}
							activities := &model.ActivitiesResponse{}

							// Check if the Node's data is valid.
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
				}
			}

			sendResponse(&mu, &responses, response, &responseSent, firstResponse, len(nodeMap))
		}(address, requestMeta)
	}

	waitGroup.Wait()

	canceled := false
	// If the context is canceled manually, do not process the responses
	for _, response := range responses {
		if errors.Is(response.Err, httputil.ErrorManuallyCanceled) {
			canceled = true
			break
		}
	}

	if !canceled {
		zap.L().Info("begin to process responses", zap.Any("responses", len(responses)))
		// Process the responses to calculate the actual request of each node
		go processResponses(responses)
	}
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

package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/form/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/httpx"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"go.uber.org/zap"
)

type Router interface {
	BuildPath(path string, query any, nodes []model.Cache) (map[common.Address]string, error)
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

func (r *SimpleRouter) DistributeRequest(_ context.Context, nodeMap map[common.Address]string, processResults func([]model.DataResponse)) (model.DataResponse, error) {
	var (
		waitGroup   sync.WaitGroup
		firstResult = make(chan model.DataResponse, 1)
		results     []model.DataResponse
		mu          sync.Mutex
	)

	for address, endpoint := range nodeMap {
		waitGroup.Add(1)

		go func(address common.Address, endpoint string) {
			defer waitGroup.Done()

			body, err := r.httpClient.Fetch(context.Background(), endpoint)
			if err != nil {
				zap.L().Error("fetch request error", zap.Any("node", address.String()), zap.Error(err))

				mu.Lock()
				results = append(results, model.DataResponse{Address: address, Err: err})

				if len(results) == len(nodeMap) {
					firstResult <- model.DataResponse{Address: address, Data: []byte(model.MessageNodeDataFailed)}
				}

				mu.Unlock()

				return
			}

			data, err := io.ReadAll(body)
			if err != nil {
				return
			}

			flagActivities, _ := r.validateActivities(data)
			flagActivity, _ := r.validateActivity(data)

			if !flagActivities && !flagActivity {
				zap.L().Error("response parse error", zap.Any("node", address.String()))

				mu.Lock()
				results = append(results, model.DataResponse{Address: address, Err: fmt.Errorf("invalid data")})

				if len(results) == len(nodeMap) {
					firstResult <- model.DataResponse{Address: address, Data: data}
				}
				mu.Unlock()

				return
			}

			mu.Lock()
			results = append(results, model.DataResponse{Address: address, Data: data, First: true})
			mu.Unlock()

			select {
			case firstResult <- model.DataResponse{Address: address, Data: data}:
			default:
			}
		}(address, endpoint)
	}

	go func() {
		waitGroup.Wait()
		close(firstResult)
		processResults(results)
	}()

	select {
	case result := <-firstResult:
		return result, nil
	case <-time.After(time.Second * 3):
		return model.DataResponse{Data: []byte(model.MessageNodeDataFailed)}, fmt.Errorf("timeout waiting for results")
	}
}

func (r *SimpleRouter) validateActivities(data []byte) (bool, *model.ActivitiesResponse) {
	var (
		res      model.ActivitiesResponse
		errRes   model.ErrResponse
		notFound model.NotFoundResponse
	)

	if err := json.Unmarshal(data, &errRes); err != nil {
		return false, nil
	}

	if errRes.ErrorCode != "" {
		return false, nil
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return false, nil
	}

	if err := json.Unmarshal(data, &notFound); err != nil {
		return false, nil
	}

	if notFound.Message != "" {
		return false, nil
	}

	return true, &res
}

func (r *SimpleRouter) validateActivity(data []byte) (bool, *model.ActivityResponse) {
	var (
		res      model.ActivityResponse
		errRes   model.ErrResponse
		notFound model.NotFoundResponse
	)

	if err := json.Unmarshal(data, &errRes); err != nil {
		return false, nil
	}

	if errRes.ErrorCode != "" {
		return false, nil
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return false, nil
	}

	if err := json.Unmarshal(data, &notFound); err != nil {
		return false, nil
	}

	if notFound.Message != "" {
		return false, nil
	}

	return true, &res
}

func NewSimpleRouter(httpClient httpx.Client) (*SimpleRouter, error) {
	return &SimpleRouter{
		httpClient: httpClient,
	}, nil
}

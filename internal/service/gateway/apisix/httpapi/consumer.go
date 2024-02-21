package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Documents: https://apisix.apache.org/zh/docs/apisix/admin-api/#consumer

const ConsumerAPIBase = "/apisix/admin/consumers"

type KeyAuthPlugin struct {
	Key string `json:"key"`
	//Header string `json:"header"` // Configure in route plugin
}

type ConsumerPropsInput struct {
	Username string `json:"username"`
	GroupID  string `json:"group_id"`
	Plugins  struct {
		// Enable key-auth plugin by default
		KeyAuth KeyAuthPlugin `json:"key-auth"`
	} `json:"plugins"`
	Description *string           `json:"desc,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type ConsumerProps struct {
	ConsumerPropsInput

	CreateTime int64 `json:"create_time"` // Timestamp
	UpdateTime int64 `json:"update_time"` // Timestamp
}

type ConsumerResponse struct {
	Key   string        `json:"key"`
	Value ConsumerProps `json:"value"`

	CreatedIndex  *int `json:"createdIndex,omitempty"`
	ModifiedIndex *int `json:"modifiedIndex,omitempty"`
}

func (s *Service) consumerUsername(keyID uint64) string {
	return fmt.Sprintf("key_%d", keyID)
}

func (s *Service) RecoverKeyIDFromConsumerUsername(username string) (int, error) {
	return strconv.Atoi(strings.Replace(username, "key_", "", 1))
}

func (s *Service) CheckConsumer(ctx context.Context, keyID uint64) (*ConsumerResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s%s/%s", s.Config.APISixAdminEndpoint, ConsumerAPIBase, s.consumerUsername(keyID)),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", s.Config.APISixAdminKey)

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var errBody APISixErr
		err = json.NewDecoder(res.Body).Decode(&errBody)

		if err != nil {
			// Even failed to decode error body
			return nil, err
		}

		// else API wrongly called
		return nil, errors.New(errBody.ErrorMsg)
	}

	var cProps ConsumerResponse

	err = json.NewDecoder(res.Body).Decode(&cProps)

	if err != nil {
		return nil, err
	}

	return &cProps, nil
}

func (s *Service) NewConsumer(ctx context.Context, keyID uint64, key string, userAddress string) error {
	// Check consumer group
	_, err := s.CheckConsumerGroup(ctx, userAddress)

	if err != nil {
		if errors.Is(err, ErrNoSuchConsumerGroup) {
			// Create consumer group
			err = s.NewConsumerGroup(ctx, userAddress)
			if err != nil {
				return err
			}
		}

		return err
	}

	desc := fmt.Sprintf("Consumer %d for user (address): %s", keyID, userAddress)
	cProps := ConsumerPropsInput{
		Username:    s.consumerUsername(keyID),
		GroupID:     userAddress,
		Description: &desc,
		Labels:      nil,
	}
	cProps.Plugins.KeyAuth.Key = key
	//cProps.Plugins.KeyAuth.Header = Config.APIGatewayKeyHeader

	reqBytes, err := json.Marshal(&cProps)

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT",
		fmt.Sprintf("%s%s", s.Config.APISixAdminEndpoint, ConsumerAPIBase),
		bytes.NewReader(reqBytes),
	)

	if err != nil {
		return err
	}

	req.Header.Set("X-API-KEY", s.Config.APISixAdminKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		var errBody APISixErr
		err = json.NewDecoder(res.Body).Decode(&errBody)

		if err != nil {
			// Even failed to decode error body
			return err
		}

		// else API wrongly called

		return errors.New(errBody.ErrorMsg)
	}

	return nil
}

func (s *Service) DeleteConsumer(ctx context.Context, keyID uint64) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE",
		fmt.Sprintf("%s%s/%s", s.Config.APISixAdminEndpoint, ConsumerAPIBase, s.consumerUsername(keyID)),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("X-API-KEY", s.Config.APISixAdminKey)

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		var errBody APISixErr
		err = json.NewDecoder(res.Body).Decode(&errBody)

		if err != nil {
			// Even failed to decode error body
			return err
		}
		// else API wrongly called
		return errors.New(errBody.ErrorMsg)
	}

	return nil
}

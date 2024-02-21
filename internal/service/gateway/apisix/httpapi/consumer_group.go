package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Documents: https://apisix.apache.org/zh/docs/apisix/admin-api/#consumer-group

const ConsumerGroupAPIBase = "/apisix/admin/consumer_groups"

type ConsumerGroupPropsInput struct {
	Plugins     map[string]any    `json:"plugins"`
	Description *string           `json:"desc,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type ConsumerGroupProps struct {
	ConsumerGroupPropsInput

	ID         *string `json:"id,omitempty"`
	CreateTime int64   `json:"create_time"` // TimeStamp
	UpdateTime int64   `json:"update_time"` // TimeStamp
}

type ConsumerGroupResponse struct {
	Key   string             `json:"key"`
	Value ConsumerGroupProps `json:"value"`

	CreatedIndex  *int `json:"createdIndex,omitempty"`
	ModifiedIndex *int `json:"modifiedIndex,omitempty"`
}

var ErrNoSuchConsumerGroup = errors.New("no such consumer group")

func (s *Service) CheckConsumerGroup(ctx context.Context, userAddress string) (*ConsumerGroupResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s%s/%s", s.Config.APISixAdminEndpoint, ConsumerGroupAPIBase, userAddress),
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
		if res.StatusCode == http.StatusNotFound {
			// No such consumer group
			return nil, ErrNoSuchConsumerGroup
		}
		// else
		var errBody APISixErr
		err = json.NewDecoder(res.Body).Decode(&errBody)

		if err != nil {
			// Even failed to decode error body
			return nil, err
		}

		// else API wrongly called

		return nil, errors.New(errBody.ErrorMsg)
	}

	var cgProps ConsumerGroupResponse

	err = json.NewDecoder(res.Body).Decode(&cgProps)

	if err != nil {
		return nil, err
	}

	return &cgProps, nil
}

func (s *Service) NewConsumerGroup(ctx context.Context, userAddress string) error {
	desc := fmt.Sprintf("Consumer group for user (address): %s", userAddress)
	reqBytes, err := json.Marshal(&ConsumerGroupPropsInput{
		Plugins:     map[string]any{}, // TODO: add account level based traffic limit plugins, etc.
		Description: &desc,
		Labels:      nil, // TODO: add level tags, etc.
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT",
		fmt.Sprintf("%s%s/%s", s.Config.APISixAdminEndpoint, ConsumerGroupAPIBase, userAddress),
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

func (s *Service) DeleteConsumerGroup(ctx context.Context, userAddress string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE",
		fmt.Sprintf("%s%s/%s", s.Config.APISixAdminEndpoint, ConsumerGroupAPIBase, userAddress),
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

//type URIBlockerPluginProps struct {
//	BlockRules      []string `json:"block_rules"`
//	RejectedCode    int      `json:"rejected_code"`
//	RejectedMsg     string   `json:"rejected_msg"`
//	CaseInsensitive bool     `json:"case_insensitive"`
//}
//
//var uriBlockerPluginProps = URIBlockerPluginProps{
//	BlockRules:      []string{"/*"},
//	RejectedCode:    http.StatusPaymentRequired,
//	RejectedMsg:     "Time to top up!",
//	CaseInsensitive: false,
//}

type LimitCountPluginProps struct {
	Count                int    `json:"count"`
	TimeWindow           int    `json:"time_window"`
	RejectedCode         int    `json:"rejected_code"`
	RejectedMsg          string `json:"rejected_msg"`
	AllowDegradation     bool   `json:"allow_degradation"`
	KeyType              string `json:"key_type"`
	Policy               string `json:"policy"`
	ShowLimitQuotaHeader bool   `json:"show_limit_quota_header"`
}

var limitCountPluginProps = LimitCountPluginProps{
	Count:                10,
	TimeWindow:           60,
	RejectedCode:         http.StatusPaymentRequired,
	RejectedMsg:          "Time to top up!",
	AllowDegradation:     false,
	KeyType:              "var",
	Policy:               "local",
	ShowLimitQuotaHeader: true,
}

func (s *Service) PauseConsumerGroup(ctx context.Context, userAddress string) error {
	// Add a uri-blocker plugin
	reqBytes, err := json.Marshal(&map[string]any{
		//"uri-blocker": uriBlockerPluginProps,
		"limit-count": limitCountPluginProps,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH",
		fmt.Sprintf("%s%s/%s/plugins", s.Config.APISixAdminEndpoint, ConsumerGroupAPIBase, userAddress),
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

	if res.StatusCode != http.StatusOK {
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

func (s *Service) ResumeConsumerGroup(ctx context.Context, userAddress string) error {
	// Remove all plugins
	req, err := http.NewRequestWithContext(ctx, "PATCH",
		fmt.Sprintf("%s%s/%s/plugins", s.Config.APISixAdminEndpoint, ConsumerGroupAPIBase, userAddress),
		bytes.NewReader([]byte("{}")), // Empty
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

	if res.StatusCode != http.StatusOK {
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

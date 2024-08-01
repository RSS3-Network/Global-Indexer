package httputil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/samber/lo"
)

var (
	ErrorNoResults        = errors.New("no results")
	ErrorManuallyCanceled = errors.New("request was manually canceled")
	ErrorTimeout          = errors.New("request timed out")
)

const (
	DefaultTimeout  = 3 * time.Second
	DefaultAttempts = 3
)

type Client interface {
	FetchWithMethod(ctx context.Context, method, path, authorization string, body io.Reader) (io.ReadCloser, error)
}

var _ Client = (*httpClient)(nil)

type httpClient struct {
	httpClient *http.Client
	attempts   uint
}

func (h *httpClient) FetchWithMethod(ctx context.Context, method, path, authorization string, body io.Reader) (readCloser io.ReadCloser, err error) {
	var bodyBytes []byte
	// Read the body into a byte slice to be able to retry the request
	if body != nil {
		bodyBytes, _ = io.ReadAll(body)
	}

	retryableFunc := func() error {
		readCloser, err = h.fetchWithMethod(ctx, method, path, authorization, bytes.NewReader(bodyBytes))
		return err
	}

	retryIfFunc := func(err error) bool {
		nonRetryableErrors := []error{
			ErrorNoResults,
			ErrorManuallyCanceled,
			ErrorTimeout,
		}

		for _, nonRetryableErr := range nonRetryableErrors {
			if errors.Is(err, nonRetryableErr) {
				return false
			}
		}

		return true
	}

	if err = retry.Do(retryableFunc, retry.Attempts(h.attempts), retry.RetryIf(retryIfFunc)); err != nil {
		return nil, err
	}

	return readCloser, nil
}

func (h *httpClient) fetchWithMethod(ctx context.Context, method, path, authorization string, body io.Reader) (io.ReadCloser, error) {
	request, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	if method == http.MethodPost {
		request.Header.Set("Content-Type", "application/json")
	}

	if authorization != "" {
		request.Header.Set("Authorization", authorization)
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		if cause := context.Cause(ctx); errors.Is(cause, context.Canceled) {
			return nil, ErrorManuallyCanceled
		} else if errors.Is(cause, context.DeadlineExceeded) {
			return nil, ErrorTimeout
		}

		return nil, fmt.Errorf("send request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		defer lo.Try(response.Body.Close)

		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return response.Body, nil
}

func NewHTTPClient(options ...ClientOption) (Client, error) {
	instance := httpClient{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		attempts: DefaultAttempts,
	}

	for _, option := range options {
		if err := option(&instance); err != nil {
			return nil, fmt.Errorf("apply options: %w", err)
		}
	}

	return &instance, nil
}

type ClientOption func(*httpClient) error

func WithAttempts(attempts uint) ClientOption {
	return func(h *httpClient) error {
		h.attempts = attempts

		return nil
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(h *httpClient) error {
		h.httpClient.Timeout = timeout

		return nil
	}
}

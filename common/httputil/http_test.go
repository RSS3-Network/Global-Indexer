package httputil_test

import (
	"context"
	"net/http"
	"sync"
	"testing"

	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/stretchr/testify/require"
)

var (
	setupOnce  sync.Once
	httpClient httputil.Client
)

func setup(t *testing.T) {
	setupOnce.Do(func() {
		var err error

		httpClient, err = httputil.NewHTTPClient()
		require.NoError(t, err)
	})
}

func TestHTTPClient_FetchWithMethod(t *testing.T) {
	t.Parallel()

	setup(t)

	type arguments struct {
		url string
	}

	testcases := []struct {
		name      string
		arguments arguments
	}{
		{
			name: "Fetch Arweave",
			arguments: arguments{
				url: "https://arweave.net/aMAYipJXf9rVHnwRYnNF7eUCxBc1zfkaopBt5TJwLWw",
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			_, _, err := httpClient.FetchWithMethod(context.TODO(), http.MethodGet, testcase.arguments.url, "", nil)
			require.NoError(t, err)
		})
	}
}

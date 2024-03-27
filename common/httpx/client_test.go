package httpx_test

import (
	"context"
	"sync"
	"testing"

	"github.com/naturalselectionlabs/rss3-global-indexer/common/httpx"
	"github.com/stretchr/testify/require"
)

var (
	setupOnce  sync.Once
	httpClient httpx.Client
)

func setup(t *testing.T) {
	setupOnce.Do(func() {
		var err error

		httpClient, err = httpx.NewHTTPClient()
		require.NoError(t, err)
	})
}

func TestHTTPClient_Fetch(t *testing.T) {
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
		{
			name: "Fetch External Api",
			arguments: arguments{
				url: "https://data.lens.phaver.com/api/lens/posts/1fdcc7ce-91a7-4af7-8022-13132842a5ec",
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			_, err := httpClient.Fetch(context.TODO(), testcase.arguments.url)
			require.NoError(t, err)
		})
	}
}

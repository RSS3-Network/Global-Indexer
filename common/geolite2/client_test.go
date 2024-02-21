package geolite2_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/naturalselectionlabs/rss3-global-indexer/common/geolite2"
	"github.com/stretchr/testify/require"
)

func TestNodeLocal(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		endpoint string
	}{
		{
			name:     "ip",
			endpoint: "1.2.3.4",
		},
		{
			name:     "domain",
			endpoint: "gi.rss3.dev",
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			c, err := geolite2.NewClient("GeoLite2-City.mmdb")
			require.NoError(t, err)

			locals, err := c.LookupLocal(context.Background(), testcase.endpoint)
			require.NoError(t, err)

			data, _ := json.Marshal(locals)
			t.Log(testcase.endpoint, string(data))
		})
	}
}

package geolite2_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/rss3-network/global-indexer/common/geolite2"
	"github.com/rss3-network/global-indexer/internal/config"
)

func TestNodeLocal(t *testing.T) {
	t.Parallel()

	c := geolite2.NewClient(&config.GeoIP{
		//Account:    976142,
		LicenseKey: "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",
		File:       "./mmdb/GeoLite2-City.mmdb",
	})

	testcases := []struct {
		name     string
		endpoint string
	}{
		{
			name:     "ip",
			endpoint: "86.25.237.241",
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

			locals, _ := c.LookupNodeLocation(context.Background(), testcase.endpoint)
			//require.NoError(t, err)

			data, _ := json.Marshal(locals)
			t.Log(testcase.endpoint, string(data))
		})
	}
}

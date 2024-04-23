package provider

import (
	"github.com/rss3-network/global-indexer/common/httputil"
)

func ProvideHTTPClient() (httputil.Client, error) {
	return httputil.NewHTTPClient()
}

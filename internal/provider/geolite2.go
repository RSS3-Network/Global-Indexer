package provider

import (
	"github.com/rss3-network/global-indexer/common/geolite2"
	"github.com/rss3-network/global-indexer/internal/config"
)

func ProvideGeoIP2(configFile *config.File) (*geolite2.Client, error) {
	return geolite2.NewClient(configFile.GeoIP)
}

package provider

import (
	"github.com/naturalselectionlabs/rss3-global-indexer/common/geolite2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
)

func ProvideGeoIP2(configFile *config.File) (*geolite2.Client, error) {
	return geolite2.NewClient(configFile.GeoIP)
}

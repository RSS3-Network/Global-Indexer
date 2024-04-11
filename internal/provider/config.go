package provider

import (
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/config/flag"
	"github.com/spf13/viper"
)

func ProvideConfig() (*config.File, error) {
	return config.Setup(viper.GetString(flag.KeyConfig))
}

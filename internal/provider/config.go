package provider

import (
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config/flag"
	"github.com/spf13/viper"
)

func ProvideConfig() (*config.File, error) {
	return config.Setup(viper.GetString(flag.KeyConfig))
}

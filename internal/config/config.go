package config

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"gopkg.in/yaml.v3"
)

const (
	Environment            = "environment"
	EnvironmentDevelopment = "development"
	EnvironmentProduction  = "production"
)

type File struct {
	Environment string       `yaml:"environment" validate:"required" default:"development"`
	Database    *Database    `yaml:"database"`
	Redis       *Redis       `yaml:"redis"`
	RSS3Chain   *RSS3Chain   `yaml:"rss3_chain"`
	Epoch       *Epoch       `yaml:"epoch"`
	GeoIP       *GeoIP       `yaml:"geo_ip"`
	Gateway     *Gateway     `yaml:"gateway"`
	APISixAdmin *APISixAdmin `yaml:"apisix_admin" validate:"required"`
	Billing     *Billing     `yaml:"billing"`
}

type Database struct {
	Driver database.Driver `mapstructure:"driver" validate:"required" default:"cockroachdb"`
	URI    string          `mapstructure:"uri" validate:"required" default:"postgres://root@localhost:26257/defaultdb"`
}

type Redis struct {
	URI string `mapstructure:"uri" validate:"required" default:"redis://localhost:6379/0"`
}

type RSS3Chain struct {
	EndpointL1     string `yaml:"endpoint_l1" validate:"required"`
	EndpointL2     string `yaml:"endpoint_l2" validate:"required"`
	BlockThreadsL1 uint64 `yaml:"block_threads_l1" default:"1"`
}

type Epoch struct {
	WalletPrivateKey string `yaml:"wallet_private_key" validate:"required"`
	GasLimit         uint64 `yaml:"gas_limit" default:"2500000"`
}

type GeoIP struct {
	Account    int    `yaml:"account" validate:"required"`
	LicenseKey string `yaml:"license_key" validate:"required"`
	File       string `yaml:"file" validate:"required" default:"./common/geolite2/mmdb/GeoLite2-City.mmdb"`
}

type Gateway struct {
	API struct {
		Listen struct {
			Host     string `yaml:"host" default:"0.0.0.0"`
			Port     uint64 `yaml:"port" default:"5555"`
			PromPort uint64 `yaml:"prom_port" default:"9000"`
		} `yaml:"listen"`
		JWTKey     string `yaml:"jwt_key" validate:"required"`
		SIWEDomain string `yaml:"siwe_domain" validate:"required"`
	} `yaml:"api" validate:"required"`
	APISixKafka struct {
		Brokers string `yaml:"brokers" validate:"required"`
		Topic   string `yaml:"topic" validate:"required"`
	} `yaml:"apisix_kafka" validate:"required"`
}

type APISixAdmin struct {
	Endpoint string `yaml:"endpoint" validate:"required"`
	Key      string `yaml:"key" validate:"required"`
}

type Billing struct {
	CollectTokenTo    string `yaml:"collect_token_to" validate:"required"`
	SlackNotification struct {
		BotToken       string `yaml:"bot_token"`
		Channel        string `yaml:"channel"`
		BlockchainScan string `yaml:"blockchain_scan" validate:"required"`
	} `yaml:"slack_notification"`
}

func Setup(configFilePath string) (*File, error) {
	// Read config file.
	config, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	// Unmarshal config file.
	var configFile File
	if err := yaml.Unmarshal(config, &configFile); err != nil {
		return nil, fmt.Errorf("unmarshal config file: %w", err)
	}

	// Set default values.yaml.
	if err := defaults.Set(&configFile); err != nil {
		return nil, fmt.Errorf("set default values.yaml: %w", err)
	}

	// Validate config values.yaml.
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&configFile); err != nil {
		return nil, fmt.Errorf("validate config file: %w", err)
	}

	return &configFile, nil
}

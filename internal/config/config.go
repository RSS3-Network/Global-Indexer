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
	Environment string           `yaml:"environment" validate:"required" default:"development"`
	Database    *Database        `yaml:"database"`
	Redis       *Redis           `yaml:"redis"`
	RSS3Chain   *RSS3ChainConfig `yaml:"rss3_chain"`
	Gateway     *GatewayConfig   `yaml:"gateway"`
}

type Database struct {
	Driver database.Driver `mapstructure:"driver" validate:"required" default:"cockroachdb"`
	URI    string          `mapstructure:"uri" validate:"required" default:"postgres://root@localhost:26257/defaultdb"`
}

type Redis struct {
	URI string `mapstructure:"uri" validate:"required" default:"redis://localhost:6379/0"`
}

type RSS3ChainConfig struct {
	EndpointL1     string `yaml:"endpoint_l1" validate:"required"`
	EndpointL2     string `yaml:"endpoint_l2" validate:"required"`
	BlockThreadsL1 uint64 `yaml:"block_threads_l1" default:"1"`
}

type GatewayConfig struct {
	API struct {
		Listen struct {
			Host     string `yaml:"host"`
			Port     uint64 `yaml:"port"`
			PromPort uint64 `yaml:"prom_port"`
		} `yaml:"listen"`
		JWTKey     string `yaml:"jwt_key"`
		SIWEDomain string `yaml:"siwe_domain"`
	} `yaml:"api"`
	Billing struct {
		CollectTokenTo    string `yaml:"collect_token_to"`
		SlackNotification struct {
			BotToken       string `yaml:"bot_token"`
			Channel        string `yaml:"channel"`
			BlockchainScan string `yaml:"blockchain_scan"`
		} `yaml:"slack_notification"`
	} `yaml:"billing"`
	APISix struct {
		Admin struct {
			Endpoint string `yaml:"endpoint"`
			Key      string `yaml:"key"`
		} `yaml:"admin"`
		Kafka struct {
			Brokers string `yaml:"brokers"`
			Topic   string `yaml:"topic"`
		} `yaml:"kafka"`
	} `yaml:"apisix"`
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

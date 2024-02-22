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
	NameService *NameService `yaml:"name_service"`
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

type NameService struct {
	Ens  *EndpointConfig `yaml:"ens"`
	Csb  *EndpointConfig `yaml:"csb"`
	Lens *EndpointConfig `yaml:"lens"`
	Fc   *EndpointConfig `yaml:"fc"`
}

type EndpointConfig struct {
	Endpoint string `yaml:"endpoint" validate:"required"`
	APIkey   string `yaml:"api_key"`
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

package config

import (
	"fmt"
	"math"
	"os"
	"unsafe"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/handler/dsl/model"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	Environment            = "environment"
	EnvironmentDevelopment = "development"
	EnvironmentProduction  = "production"
)

type File struct {
	Environment    string          `yaml:"environment" validate:"required" default:"development"`
	Database       *Database       `yaml:"database"`
	Redis          *Redis          `yaml:"redis"`
	RSS3Chain      *RSS3Chain      `yaml:"rss3_chain"`
	Settler        *Settler        `yaml:"settler"`
	Distributor    *Distributor    `yaml:"distributor"`
	SpecialRewards *SpecialRewards `yaml:"special_rewards"`
	GeoIP          *GeoIP          `yaml:"geo_ip"`
	RPC            *RPC            `yaml:"rpc"`
	Telemetry      *Telemetry      `json:"telemetry"`
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
	BlockThreadsL2 uint64 `yaml:"block_threads_l2" default:"1"`
}

type Settler struct {
	PrivateKey     string `yaml:"private_key"`
	WalletAddress  string `yaml:"wallet_address"`
	SignerEndpoint string `yaml:"signer_endpoint"`
	// EpochIntervalInHours
	EpochIntervalInHours int    `yaml:"epoch_interval_in_hours" default:"18"`
	GasLimit             uint64 `yaml:"gas_limit" default:"2500000"`
	// BatchSize is the number of Nodes to process in each batch.
	// This is to prevent the contract call from running out of gas.
	BatchSize int `yaml:"batch_size" default:"200"`
}

type Distributor struct {
	// The number of demotions that trigger a slashing.
	MaxDemotionCount int `yaml:"max_demotion_count"`
	// The required number of qualified Nodes
	QualifiedNodeCount int `yaml:"qualified_node_count" default:"3"`
	// The required number of verifications before a request is considered valid
	VerificationCount int `yaml:"verification_count" default:"3"`
}

type SpecialRewards struct {
	GiniCoefficient       float64 `yaml:"gini_coefficient" validate:"required"`
	StakerFactor          float64 `yaml:"staker_factor" validate:"required"`
	NodeThreshold         float64 `yaml:"node_threshold" validate:"required"`
	EpochLimit            int     `yaml:"epoch_limit" validate:"required"`
	Rewards               float64 `yaml:"rewards" validate:"required"`
	RewardsCeiling        float64 `yaml:"rewards_ceiling" validate:"required"`
	RewardsRatioActive    float64 `yaml:"rewards_ratio_active" validate:"required"`
	RewardsRatioOperation float64 `yaml:"rewards_ratio_operation" validate:"required"`
}

type GeoIP struct {
	Account    int    `yaml:"account"`
	LicenseKey string `yaml:"license_key"`
	File       string `yaml:"file" validate:"required" default:"./common/geolite2/mmdb/GeoLite2-City.mmdb"`
}

type RPC struct {
	RPCNetwork *RPCNetwork `yaml:"network"`
}

type RPCNetwork struct {
	Ethereum  *RPCEndpoint `yaml:"ethereum"`
	Crossbell *RPCEndpoint `yaml:"crossbell"`
	Polygon   *RPCEndpoint `yaml:"polygon"`
	Farcaster *RPCEndpoint `yaml:"farcaster"`
}

type RPCEndpoint struct {
	Endpoint string `yaml:"endpoint" validate:"required"`
	APIkey   string `yaml:"api_key"`
}

type Telemetry struct {
	Endpoint string `yaml:"endpoint" validate:"required"`
	Insecure bool   `yaml:"insecure"`
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

	// Initialize some common global variables.
	initGlobalVars(&configFile)

	return &configFile, nil
}

func initGlobalVars(file *File) {
	if file.Distributor.MaxDemotionCount == 0 {
		var i int

		if unsafe.Sizeof(i) == 4 {
			file.Distributor.MaxDemotionCount = math.MaxInt32
		} else {
			file.Distributor.MaxDemotionCount = math.MaxInt64
		}
	}

	model.DemotionCountBeforeSlashing = file.Distributor.MaxDemotionCount
	model.RequiredVerificationCount = file.Distributor.VerificationCount
	model.RequiredQualifiedNodeCount = file.Distributor.QualifiedNodeCount

	zap.L().Info("init constants", zap.Any("MaxDemotionCount", model.DemotionCountBeforeSlashing), zap.Any("VerificationCount", model.RequiredVerificationCount), zap.Any("QualifiedNodeCount", model.RequiredQualifiedNodeCount))
}

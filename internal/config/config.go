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
	Environment   string         `yaml:"environment" validate:"required" default:"development"`
	Database      *Database      `yaml:"database"`
	Redis         *Redis         `yaml:"redis"`
	RSS3Chain     *RSS3Chain     `yaml:"rss3_chain"`
	Settler       *Settler       `yaml:"settler"`
	Distributor   *Distributor   `yaml:"distributor"`
	Rewards       *Rewards       `yaml:"rewards"`
	ActiveScores  *ActiveScores  `yaml:"active_scores"`
	GeoIP         *GeoIP         `yaml:"geo_ip"`
	RPC           *RPC           `yaml:"rpc"`
	Telemetry     *Telemetry     `json:"telemetry"`
	TokenPriceAPI *TokenPriceAPI `yaml:"token_price_api"`
}

type Database struct {
	Driver database.Driver `mapstructure:"driver" validate:"required" default:"postgres"`
	URI    string          `mapstructure:"uri" validate:"required" default:"postgres://root@localhost:5432/postgres"`
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
	BatchSize            int `yaml:"batch_size" default:"200"`
	ProductionStartEpoch int `yaml:"production_start_epoch" default:"227"`
	GracePeriodEpochs    int `yaml:"grace_period_epochs" default:"28"`
}

type Distributor struct {
	// The number of demotions that triggers a slashing.
	MaxDemotionCount int `yaml:"max_demotion_count" default:"4"`
	// The number of nodes required to meet the criteria during distribution.
	QualifiedNodeCount int `yaml:"qualified_node_count" default:"3"`
	// The number of verification activities selected during the second verification.
	VerificationCount int `yaml:"verification_count" default:"3"`
	ToleranceSeconds  int `yaml:"tolerance_seconds" default:"1200"`
}

type Rewards struct {
	OperationRewards float64         `yaml:"operation_rewards" validate:"required"`
	OperationScore   *OperationScore `yaml:"operation_score" validate:"required"`
}

type OperationScore struct {
	Distribution *Distribution `yaml:"distribution" validate:"required"`
	Data         *Data         `yaml:"data" validate:"required"`
	Stability    *Stability    `yaml:"stability" validate:"required"`
}

type Distribution struct {
	Weight        float64 `yaml:"weight" validate:"required"`
	WeightInvalid float64 `yaml:"weight_invalid" validate:"required"`
}

type Data struct {
	Weight         float64 `yaml:"weight" validate:"required"`
	WeightNetwork  float64 `yaml:"weight_network" validate:"required"`
	WeightIndexer  float64 `yaml:"weight_indexer" validate:"required"`
	WeightActivity float64 `yaml:"weight_activity" validate:"required"`
}

type Stability struct {
	Weight        float64 `yaml:"weight" validate:"required"`
	WeightUptime  float64 `yaml:"weight_uptime" validate:"required"`
	WeightVersion float64 `yaml:"weight_version" validate:"required"`
}

type ActiveScores struct {
	GiniCoefficient float64 `yaml:"gini_coefficient" validate:"required"`
	StakerFactor    float64 `yaml:"staker_factor" validate:"required"`
	EpochLimit      int     `yaml:"epoch_limit" validate:"required"`
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

type TokenPriceAPI struct {
	Endpoint  string `yaml:"endpoint" validate:"required"`
	AuthToken string `yaml:"auth_token"`
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
	if file.Distributor.MaxDemotionCount == -1 {
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
	model.ToleranceSeconds = file.Distributor.ToleranceSeconds

	zap.L().Info("init constants", zap.Any("MaxDemotionCount", model.DemotionCountBeforeSlashing), zap.Any("VerificationCount", model.RequiredVerificationCount), zap.Any("QualifiedNodeCount", model.RequiredQualifiedNodeCount), zap.Any("ToleranceSeconds", model.ToleranceSeconds))
}

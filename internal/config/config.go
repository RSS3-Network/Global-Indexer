package config

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/naturalselectionlabs/global-indexer/internal/database"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	Environment            = "environment"
	EnvironmentDevelopment = "development"
	EnvironmentProduction  = "production"
)

type File struct {
	Environment string           `yaml:"environment" validate:"required" default:"development"`
	Database    *database.Config `yaml:"database"`
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
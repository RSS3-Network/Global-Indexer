package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config/flag"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var flags *pflag.FlagSet

var command = cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		flags = cmd.PersistentFlags()

		config, err := config.Setup(lo.Must(flags.GetString(flag.KeyConfig)))
		if err != nil {
			return fmt.Errorf("setup config file: %w", err)
		}

		// Dial and migrate database.
		databaseClient, err := dialer.Dial(cmd.Context(), config.Database)
		if err != nil {
			return fmt.Errorf("dial database: %w", err)
		}

		if err := databaseClient.Migrate(cmd.Context()); err != nil {
			return fmt.Errorf("migrate database: %w", err)
		}

		// Dial rss3 ethereum client.
		ethereumClient, err := ethclient.DialContext(cmd.Context(), config.RSS3Chain.Endpoint)
		if err != nil {
			return fmt.Errorf("dial rss3 ethereum client: %w", err)
		}

		hub, err := hub.NewServer(cmd.Context(), databaseClient, ethereumClient)
		if err != nil {
			return fmt.Errorf("new hub server: %w", err)
		}

		return hub.Run(cmd.Context())
	},
}

func initializeLogger() {
	if os.Getenv(config.Environment) == config.EnvironmentDevelopment {
		zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	} else {
		zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	}
}

func init() {
	initializeLogger()

	command.PersistentFlags().String(flag.KeyConfig, "./deploy/config.yaml", "config file path")
}

func main() {
	if err := command.ExecuteContext(context.Background()); err != nil {
		zap.L().Fatal("execute command", zap.Error(err))
	}
}

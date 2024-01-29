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
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/indexer"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler"
	"github.com/redis/go-redis/v9"
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
		ethereumClient, err := ethclient.DialContext(cmd.Context(), config.RSS3Chain.EndpointL2)
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

var indexCommand = &cobra.Command{
	Use: "index",
	RunE: func(cmd *cobra.Command, args []string) error {
		flags = cmd.PersistentFlags()

		config, err := config.Setup(lo.Must(flags.GetString(flag.KeyConfig)))
		if err != nil {
			return fmt.Errorf("setup config file: %w", err)
		}

		databaseClient, err := dialer.Dial(cmd.Context(), config.Database)
		if err != nil {
			return err
		}

		if err := databaseClient.Migrate(cmd.Context()); err != nil {
			return fmt.Errorf("migrate database: %w", err)
		}

		instance, err := indexer.New(databaseClient, *config.RSS3Chain)
		if err != nil {
			return err
		}

		return instance.Run(cmd.Context())
	},
}

var schedulerCommand = &cobra.Command{
	Use: "scheduler",
	RunE: func(cmd *cobra.Command, args []string) error {
		flags = cmd.PersistentFlags()

		config, err := config.Setup(lo.Must(flags.GetString(flag.KeyConfig)))
		if err != nil {
			return fmt.Errorf("setup config file: %w", err)
		}

		databaseClient, err := dialer.Dial(cmd.Context(), config.Database)
		if err != nil {
			return err
		}

		if err := databaseClient.Migrate(cmd.Context()); err != nil {
			return fmt.Errorf("migrate database: %w", err)
		}

		options, err := redis.ParseURL(config.Redis.URI)
		if err != nil {
			return fmt.Errorf("parse redis uri: %w", err)
		}

		redisClient := redis.NewClient(options)

		instance, err := scheduler.New(lo.Must(flags.GetString(flag.KeyServer)), databaseClient, redisClient)
		if err != nil {
			return err
		}

		return instance.Run(cmd.Context())
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

	command.AddCommand(indexCommand)
	command.AddCommand(schedulerCommand)
	command.PersistentFlags().String(flag.KeyConfig, "./deploy/config.yaml", "config file path")

	indexCommand.PersistentFlags().String(flag.KeyConfig, "./deploy/config.yaml", "config file path")

	schedulerCommand.PersistentFlags().String(flag.KeyConfig, "./deploy/config.yaml", "config file path")
	schedulerCommand.PersistentFlags().String(flag.KeyServer, "detector", "server name")
}

func main() {
	if err := command.ExecuteContext(context.Background()); err != nil {
		zap.L().Fatal("execute command", zap.Error(err))
	}
}

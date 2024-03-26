package main

import (
	"context"
	"fmt"
	"os"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config/flag"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/epoch"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/hub"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/indexer"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var command = cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		server := service.NewServer(
			hub.Module,
			fx.Provide(hub.NewServer),
		)

		if err := server.Start(cmd.Context()); err != nil {
			return fmt.Errorf("start server: %w", err)
		}

		server.Wait()

		return nil
	},
}

var indexCommand = &cobra.Command{
	Use: "index",
	RunE: func(cmd *cobra.Command, args []string) error {
		server := service.NewServer(
			indexer.Module,
			fx.Provide(indexer.NewServer),
		)

		if err := server.Start(cmd.Context()); err != nil {
			return fmt.Errorf("start server: %w", err)
		}

		server.Wait()

		return nil
	},
}

var schedulerCommand = &cobra.Command{
	Use: "scheduler",
	RunE: func(cmd *cobra.Command, args []string) error {
		server := service.NewServer(
			scheduler.Module,
			fx.Provide(scheduler.NewServer),
		)

		if err := server.Start(cmd.Context()); err != nil {
			return fmt.Errorf("start server: %w", err)
		}

		server.Wait()

		return nil
	},
}

var epochCommand = &cobra.Command{
	Use: "epoch",
	RunE: func(cmd *cobra.Command, args []string) error {
		server := service.NewServer(
			epoch.Module,
			fx.Provide(epoch.NewServer),
		)

		if err := server.Start(cmd.Context()); err != nil {
			return fmt.Errorf("start server: %w", err)
		}

		server.Wait()

		return nil
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
	command.AddCommand(epochCommand)

	command.PersistentFlags().String(flag.KeyConfig, "./deploy/config.yaml", "config file path")
	command.PersistentFlags().Uint64(flag.KeyChainIDL1, flag.ValueChainIDL1, "l1 chain id")
	command.PersistentFlags().Uint64(flag.KeyChainIDL2, flag.ValueChainIDL2, "l2 chain id")

	indexCommand.PersistentFlags().String(flag.KeyConfig, "./deploy/config.yaml", "config file path")
	schedulerCommand.PersistentFlags().String(flag.KeyConfig, "./deploy/config.yaml", "config file path")
	schedulerCommand.PersistentFlags().String(flag.KeyServer, "detector", "server name")
	epochCommand.PersistentFlags().String(flag.KeyConfig, "./deploy/config.yaml", "config file path")
}

func main() {
	if err := command.ExecuteContext(context.Background()); err != nil {
		zap.L().Fatal("execute command", zap.Error(err))
	}
}

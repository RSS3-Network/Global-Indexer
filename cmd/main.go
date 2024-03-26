package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/geolite2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/cache"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config/flag"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/constant"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/nameresolver"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/epoch"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/hub"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/indexer"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/scheduler"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
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

		// Dial redis.
		options, err := redis.ParseURL(config.Redis.URI)
		if err != nil {
			return fmt.Errorf("parse redis uri: %w", err)
		}

		redisClient := redis.NewClient(options)

		geoLite2, err := geolite2.NewClient(config.GeoIP)
		if err != nil {
			return fmt.Errorf("new geo lite2 client: %w", err)
		}

		nameService, err := nameresolver.NewNameResolver(cmd.Context(), config.RPC.RPCNetwork)
		if err != nil {
			return fmt.Errorf("init name resolver: %w", err)
		}

		hub, err := hub.NewServer(cmd.Context(), databaseClient, ethereumClient, redisClient, geoLite2, nameService)
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

		if config.Telemetry != nil {
			if err := setupOpenTelemetry("indexer", config.Telemetry); err != nil {
				return fmt.Errorf("setup opentelemetry tracer")
			}
		}

		databaseClient, err := dialer.Dial(cmd.Context(), config.Database)
		if err != nil {
			return err
		}

		options, err := redis.ParseURL(config.Redis.URI)
		if err != nil {
			return fmt.Errorf("parse redis uri: %w", err)
		}

		cacheClient := cache.New(redis.NewClient(options))

		if err := databaseClient.Migrate(cmd.Context()); err != nil {
			return fmt.Errorf("migrate database: %w", err)
		}

		instance, err := indexer.New(databaseClient, cacheClient, *config.RSS3Chain)
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

		// Dial rss3 ethereum client.
		ethereumClient, err := ethclient.DialContext(cmd.Context(), config.RSS3Chain.EndpointL2)
		if err != nil {
			return fmt.Errorf("dial rss3 ethereum client: %w", err)
		}

		instance, err := scheduler.New(lo.Must(flags.GetString(flag.KeyServer)), databaseClient, redisClient, ethereumClient)
		if err != nil {
			return err
		}

		return instance.Run(cmd.Context())
	},
}

var epochCommand = &cobra.Command{
	Use: "epoch",
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

		instance, err := epoch.New(cmd.Context(), databaseClient, redisClient, *config)
		if err != nil {
			return err
		}

		return instance.Run(cmd.Context())
	},
}

func setupOpenTelemetry(serviceName string, config *config.Telemetry) error {
	options := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(config.Endpoint),
	}

	if config.Insecure {
		options = append(options, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptrace.New(context.Background(), otlptracehttp.NewClient(options...))
	if err != nil {
		return err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(constant.BuildServiceName(serviceName)),
			semconv.ServiceVersionKey.String(constant.BuildServiceVersion()),
		)),
	)

	otel.SetTracerProvider(tracerProvider)

	return nil
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
	indexCommand.PersistentFlags().String(flag.KeyConfig, "./deploy/config.yaml", "config file path")
	schedulerCommand.PersistentFlags().String(flag.KeyConfig, "./deploy/config.yaml", "config file path")
	schedulerCommand.PersistentFlags().String(flag.KeyServer, "detector", "server name")
	epochCommand.PersistentFlags().String(flag.KeyConfig, "./deploy/config.yaml", "config file path")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err := command.ExecuteContext(ctx); err != nil {
		zap.L().Fatal("execute command", zap.Error(err))
	}
}

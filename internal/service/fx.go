package service

import (
	"context"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/constant"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/provider"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewServer(options ...fx.Option) *fx.App {
	return fx.New(
		fx.Options(options...),
		fx.Provide(provider.ProvideConfig),
		fx.Provide(provider.ProvideOpenTelemetryTracer),
		fx.Invoke(InjectLifecycle),
		fx.Invoke(InjectOpenTelemetry),
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ZapLogger{
				Logger: zap.L(),
			}
		}),
	)
}

func InjectLifecycle(lifecycle fx.Lifecycle, server Server) {
	constant.ServiceName = server.Name()

	hook := fx.Hook{
		OnStart: func(ctx context.Context) error {
			return server.Run(ctx)
		},
	}

	lifecycle.Append(hook)
}

func InjectOpenTelemetry(tracerProvider trace.TracerProvider) {
	otel.SetTracerProvider(tracerProvider)
}

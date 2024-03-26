package service

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewServer(options ...fx.Option) *fx.App {
	return fx.New(
		fx.Options(options...),
		fx.Invoke(InjectLifecycle),
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ZapLogger{
				Logger: zap.L(),
			}
		}),
	)
}

func InjectLifecycle(lifecycle fx.Lifecycle, server Server) {
	hook := fx.Hook{
		OnStart: func(ctx context.Context) error {
			return server.Run(ctx)
		},
	}

	lifecycle.Append(hook)
}

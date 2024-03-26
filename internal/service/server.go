package service

import "context"

type Server interface {
	Name() string
	Run(ctx context.Context) error
}

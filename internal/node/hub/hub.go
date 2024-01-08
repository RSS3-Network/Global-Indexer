package hub

import "context"

type Hub struct {
}

func NewHub(_ context.Context) *Hub {
	return &Hub{}
}

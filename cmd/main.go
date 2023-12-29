package main

import (
	"context"
	"github.com/naturalselectionlabs/global-indexer/internal/hub"
)

func main() {
	server, err := hub.NewServer(context.Background())
	if err != nil {
		panic(err)
	}

	if err := server.Run(context.Background()); err != nil {
		panic(err)
	}
}

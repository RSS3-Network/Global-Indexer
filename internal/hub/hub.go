package hub

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/global-indexer/common/ethereum/contract/staking"
	"github.com/naturalselectionlabs/global-indexer/internal/database"
)

type Hub struct {
	databaseClient  database.Client
	stakingContract *staking.Staking
}

func NewHub(_ context.Context, databaseClient database.Client) *Hub {
	return &Hub{
		databaseClient:  databaseClient,
		stakingContract: stakingContract,
	}, nil
}

package hub

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/ethereum/contract/staking"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
)

type Hub struct {
	databaseClient  database.Client
	stakingContract *staking.Staking
}

var _ echo.Validator = (*Validator)(nil)

var defaultValidator = &Validator{
	validate: validator.New(),
}

type Validator struct {
	validate *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func NewHub(_ context.Context, databaseClient database.Client, ethereumClient *ethclient.Client) (*Hub, error) {
	stakingContract, err := staking.NewStaking(staking.AddressStaking, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new staking contract: %w", err)
	}

	return &Hub{
		databaseClient:  databaseClient,
		stakingContract: stakingContract,
	}, nil
}

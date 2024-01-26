package hub

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/shedlock"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/job"
	"github.com/naturalselectionlabs/rss3-global-indexer/provider/node"
)

type Hub struct {
	databaseClient  database.Client
	stakingContract *l2.Staking
	pathBuilder     node.Builder
	httpClient      *http.Client
	employer        *shedlock.Employer
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
	stakingContract, err := l2.NewStaking(l2.AddressStakingProxy, ethereumClient)
	if err != nil {
		return nil, fmt.Errorf("new staking contract: %w", err)
	}

	employer := shedlock.New()

	if err = addCronJobs(databaseClient, stakingContract, employer); err != nil {
		return nil, fmt.Errorf("add cron jobs: %w", err)
	}

	return &Hub{
		databaseClient:  databaseClient,
		stakingContract: stakingContract,
		pathBuilder:     node.NewPathBuilder(),
		httpClient:      http.DefaultClient,
		employer:        employer,
	}, nil
}

func addCronJobs(databaseClient database.Client, stakingContract *l2.Staking, employer *shedlock.Employer) error {
	jobs := []job.Job{
		job.NewSortNodesJob(databaseClient, stakingContract),
	}

	for _, cronjob := range jobs {
		if err := employer.AddJob(cronjob.Name(), cronjob.Spec(), cronjob.Timeout(), job.NewCronJob(employer, cronjob)); err != nil {
			return err
		}
	}

	return nil
}

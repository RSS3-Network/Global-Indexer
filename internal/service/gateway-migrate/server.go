package gateway_migrate

import (
	"context"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
)

type Server struct {
	databaseClient database.Client
}

func (s *Server) Run(ctx context.Context) error {
	return s.databaseClient.Raw().WithContext(ctx).AutoMigrate(
		&table.GatewayAccount{},
		&table.GatewayKey{},
		&table.GatewayConsumptionLog{},
		&table.GatewayPendingWithdrawRequest{},

		&table.BillingRecordDeposited{},
		&table.BillingRecordWithdrawal{},
		&table.BillingRecordCollected{},
	)
}

func New(databaseClient database.Client) (service.Server, error) {
	instance := Server{
		databaseClient: databaseClient,
	}

	return &instance, nil
}

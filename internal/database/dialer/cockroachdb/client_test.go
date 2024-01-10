package cockroachdb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/global-indexer/internal/database"
	"github.com/naturalselectionlabs/global-indexer/internal/database/dialer"
	"github.com/naturalselectionlabs/global-indexer/schema"
	"github.com/naturalselectionlabs/rss3-node/config"
	"github.com/naturalselectionlabs/rss3-node/schema/filter"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/cockroachdb"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		driver      database.Driver
		nodeCreated *schema.Node
	}{
		{
			name:   "cockroach",
			driver: database.DriverCockroachDB,
			nodeCreated: &schema.Node{
				Address: common.HexToAddress("0xc98D64DA73a6616c42117b582e832812e7B8D57F"),
				Stream: &config.Stream{
					Enable: lo.ToPtr(true),
					Driver: "kafka",
					Topic:  "node.feeds",
					URI:    "localhost:9092",
				},
				Config: &config.Node{
					RSS: []*config.Module{
						{
							Network:  filter.NetworkRSS,
							Endpoint: "https://node.rss3.dev",
						},
					},
					Decentralized: []*config.Module{
						{
							Network:  filter.NetworkEthereum,
							Endpoint: "https://rpc.ankr.com/eth",
						},
					},
				},
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			var container *gnomock.Container
			var dataSourceName string
			var err error

			for {
				container, dataSourceName, err = createContainer(context.Background(), testcase.driver)
				if err == nil {
					break
				}
			}

			t.Cleanup(func() {
				require.NoError(t, gnomock.Stop(container))
			})

			// Dial the database.
			client, err := dialer.Dial(context.Background(), &database.Config{
				Driver: testcase.driver,
				URI:    dataSourceName,
			})

			require.NoError(t, err)
			require.NotNil(t, client)

			// Migrate the database.
			require.NoError(t, client.Migrate(context.Background()))

			// Save node.
			require.NoError(t, client.SaveNode(context.Background(), testcase.nodeCreated))

			// Find node.
			nodeFound, err := client.FindNode(context.Background(), testcase.nodeCreated.Address)
			require.NoError(t, err)
			require.NotEmpty(t, nodeFound.Address)

			// Find nodes.
			nodesFound, err := client.FindNodes(context.Background(), []common.Address{testcase.nodeCreated.Address}, nil)
			require.NoError(t, err)
			require.Equal(t, 1, len(nodesFound))

			// Update node.
			testcase.nodeCreated.Stream.URI = "localhost:9093"
			require.NoError(t, client.SaveNode(context.Background(), testcase.nodeCreated))

			// Find node.
			nodeFound, err = client.FindNode(context.Background(), testcase.nodeCreated.Address)
			require.NoError(t, err)
			require.Equal(t, testcase.nodeCreated.Stream.URI, nodeFound.Stream.URI)
		})
	}
}

func createContainer(ctx context.Context, driver database.Driver) (container *gnomock.Container, dataSourceName string, err error) {
	conf := database.Config{
		Driver: driver,
	}

	switch driver {
	case database.DriverCockroachDB:
		preset := cockroachdb.Preset(
			cockroachdb.WithDatabase("test"),
			cockroachdb.WithVersion("v23.1.8"),
		)

		// Use a health check function to wait for the database to be ready.
		healthcheckFunc := func(ctx context.Context, container *gnomock.Container) error {
			conf.URI = formatContainerURI(container)

			client, err := dialer.Dial(ctx, &conf)
			if err != nil {
				return err
			}

			transaction, err := client.Begin(ctx)
			if err != nil {
				return err
			}

			defer lo.Try(transaction.Rollback)

			return nil
		}

		container, err = gnomock.Start(preset, gnomock.WithContext(ctx), gnomock.WithHealthCheck(healthcheckFunc))
		if err != nil {
			return nil, "", err
		}

		return container, formatContainerURI(container), nil
	default:
		return nil, "", fmt.Errorf("unsupported driver: %s", driver)
	}
}

func formatContainerURI(container *gnomock.Container) string {
	return fmt.Sprintf(
		"postgres://root@%s:%d/%s?sslmode=disable",
		container.Host,
		container.DefaultPort(),
		"test",
	)
}

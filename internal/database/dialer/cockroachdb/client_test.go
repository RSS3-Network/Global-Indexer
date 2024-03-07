package cockroachdb_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
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
				ID:      big.NewInt(1),
				Address: common.HexToAddress("0xc98D64DA73a6616c42117b582e832812e7B8D57F"),
				Stream: json.RawMessage(`
				{
				   "Driver":"kafka",
				   "Enable":false,
				   "Topic":"rss3.node.feeds",
				   "URI":"localhost:9092"
				}`),
				Config: json.RawMessage(`
				{
				   "Decentralized":[
					  {
						 "Endpoint":"https://rpc.ankr.com/eth",
						 "IPFSGateways":null,
						 "Network":"ethereum",
						 "Parameters":{
							"block_number_start":null,
							"block_number_target":null
						 },
						 "Worker":"fallback"
					  }
				   ],
				   "Federated":null,
				   "RSS":[
					  {
						 "Endpoint":"https://rsshub.app/",
						 "IPFSGateways":null,
						 "Network":"rss",
						 "Parameters":{
							"authentication":{
							   "access_code":null,
							   "access_key":null,
							   "password":null,
							   "username":null
							}
						 },
						 "Worker":"unknown"
					  }
				   ]
				}`),
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
			client, err := dialer.Dial(context.Background(), &config.Database{
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
			nodesFound, err := client.FindNodes(context.Background(), []common.Address{testcase.nodeCreated.Address}, nil, nil, 10)
			require.NoError(t, err)
			require.Equal(t, 1, len(nodesFound))

			// Update node.
			testcase.nodeCreated.Stream = json.RawMessage(`{}`)
			require.NoError(t, client.SaveNode(context.Background(), testcase.nodeCreated))

			// Find node.
			nodeFound, err = client.FindNode(context.Background(), testcase.nodeCreated.Address)
			require.NoError(t, err)
			require.Equal(t, testcase.nodeCreated.Stream, nodeFound.Stream)
		})
	}
}

func createContainer(ctx context.Context, driver database.Driver) (container *gnomock.Container, dataSourceName string, err error) {
	conf := config.Database{
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

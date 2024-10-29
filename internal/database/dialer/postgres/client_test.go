package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/adrianbrad/psqldocker"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/database/dialer"
	"github.com/rss3-network/global-indexer/schema"
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
			name:   "postgres",
			driver: database.DriverPostgres,
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

			var (
				container      *psqldocker.Container
				dataSourceName string
				err            error
			)

			for {
				container, dataSourceName, err = createContainer(context.Background(), testcase.driver)
				if err == nil {
					break
				}
			}

			t.Cleanup(func() {
				require.NoError(t, container.Close())
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
			nodesFound, err := client.FindNodes(context.Background(), schema.FindNodesQuery{
				NodeAddresses: []common.Address{testcase.nodeCreated.Address},
			})
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

func createContainer(_ context.Context, driver database.Driver) (container *psqldocker.Container, dataSourceName string, err error) {
	switch driver {
	case database.DriverPostgres:
		c, err := psqldocker.NewContainer(
			"user",
			"password",
			"test",
		)
		if err != nil {
			return nil, "", fmt.Errorf("create psql container: %w", err)
		}

		return c, formatContainerURI(c), nil
	default:
		return nil, "", fmt.Errorf("unsupported driver: %s", driver)
	}
}

func formatContainerURI(container *psqldocker.Container) string {
	return fmt.Sprintf(
		"postgres://user:password@%s:%s/%s?sslmode=disable",
		"127.0.0.1",
		container.Port(),
		"test",
	)
}

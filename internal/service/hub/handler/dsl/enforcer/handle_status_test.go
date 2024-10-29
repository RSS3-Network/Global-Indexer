package enforcer

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-version"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	workerStatusNodeIndexing  = `{"data":{"decentralized":[{"network":"farcaster","worker":"core","tags":null,"platform":"Unknown","status":"Indexing","remote_state":1718215438006,"indexed_state":1718215435040}],"rss":null,"federated":null}}`
	workerStatusNodeUnhealthy = `{"data":{"decentralized":[{"network":"farcaster","worker":"core","tags":null,"platform":"Unknown","status":"Unhealthy","remote_state":0,"indexed_state":0}],"rss":null,"federated":null}}`
	workerStatusNodeOnline    = `{"data":{"decentralized":[{"network":"farcaster","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":0,"indexed_state":0}],"rss":null,"federated":null}}`
	workerStatusNodeRssOnline = `{"data":{"decentralized":null,"rss":{"worker_id":"","worker":null,"network":"rss","tags":["rss"],"platform":"Unknown","status":"Ready","remote_state":0,"indexed_state":0,"index_count":0},"federated":null}}`
)

func TestDetermineStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		workerStatus    string
		initialStatus   schema.NodeStatus
		minVersion      string
		version         string
		expectedStatus  schema.NodeStatus
		expectedErrPath string
	}{
		{
			name:            "WorkerStatusUnavailable",
			workerStatus:    "",
			initialStatus:   schema.NodeStatusOffline,
			minVersion:      "1.0.0",
			version:         "1.0.0",
			expectedStatus:  schema.NodeStatusOffline,
			expectedErrPath: "workers_status",
		},
		{
			name:            "WorkerStatusUnavailable",
			workerStatus:    "",
			initialStatus:   schema.NodeStatusInitializing,
			minVersion:      "1.0.0",
			version:         "1.0.0",
			expectedStatus:  schema.NodeStatusInitializing,
			expectedErrPath: "",
		},
		{
			name:            "RegisteredToOutdated",
			workerStatus:    workerStatusNodeIndexing,
			initialStatus:   schema.NodeStatusRegistered,
			minVersion:      "1.1.0",
			version:         "1.0.0",
			expectedStatus:  schema.NodeStatusOutdated,
			expectedErrPath: "",
		},
		{
			name:            "RegisteredStaysRegistered",
			workerStatus:    workerStatusNodeUnhealthy,
			initialStatus:   schema.NodeStatusRegistered,
			minVersion:      "1.0.0",
			version:         "1.0.0",
			expectedStatus:  schema.NodeStatusRegistered,
			expectedErrPath: "",
		},
		{
			name:            "RegisteredToInitializing",
			workerStatus:    workerStatusNodeIndexing,
			initialStatus:   schema.NodeStatusRegistered,
			minVersion:      "1.0.0",
			version:         "1.0.0",
			expectedStatus:  schema.NodeStatusInitializing,
			expectedErrPath: "",
		},
		{
			name:            "OutdatedStaysOutdated",
			workerStatus:    workerStatusNodeIndexing,
			initialStatus:   schema.NodeStatusOutdated,
			minVersion:      "1.1.0",
			version:         "1.0.0",
			expectedStatus:  schema.NodeStatusOutdated,
			expectedErrPath: "",
		},
		{
			name:            "OutdatedToRegistered",
			workerStatus:    workerStatusNodeUnhealthy,
			initialStatus:   schema.NodeStatusOutdated,
			minVersion:      "1.0.0",
			version:         "1.0.0",
			expectedStatus:  schema.NodeStatusRegistered,
			expectedErrPath: "",
		},
		{
			name:            "OutdatedToInitializing",
			workerStatus:    workerStatusNodeIndexing,
			initialStatus:   schema.NodeStatusOutdated,
			minVersion:      "1.0.0",
			version:         "1.0.0",
			expectedStatus:  schema.NodeStatusInitializing,
			expectedErrPath: "",
		},
		{
			name:            "InitializingToOutdated",
			workerStatus:    workerStatusNodeIndexing,
			initialStatus:   schema.NodeStatusInitializing,
			minVersion:      "1.1.0",
			version:         "1.0.0",
			expectedStatus:  schema.NodeStatusOutdated,
			expectedErrPath: "",
		},
		{
			name:            "InitializingToRegistered",
			workerStatus:    workerStatusNodeUnhealthy,
			initialStatus:   schema.NodeStatusInitializing,
			minVersion:      "1.0.0",
			version:         "1.0.0",
			expectedStatus:  schema.NodeStatusRegistered,
			expectedErrPath: "",
		},
		{
			name:            "InitializingStaysInitializing",
			workerStatus:    workerStatusNodeIndexing,
			initialStatus:   schema.NodeStatusInitializing,
			minVersion:      "1.0.0",
			version:         "1.0.0",
			expectedStatus:  schema.NodeStatusInitializing,
			expectedErrPath: "",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := setupMockClient(tt.workerStatus)
			enforcer := &SimpleEnforcer{httpClient: mockClient}
			node := &schema.Node{
				Address:     common.Address{},
				Endpoint:    "http://localhost:8080",
				AccessToken: "token",
				Status:      tt.initialStatus,
				Version:     tt.version,
			}

			minVersion, _ := version.NewVersion(tt.minVersion)
			status, errPath := enforcer.determineStatus(ctx, node, minVersion)

			assert.Equal(t, tt.expectedStatus, status)
			assert.Equal(t, tt.expectedErrPath, errPath)
		})
	}
}

func setupMockClient(workerStatus string) *MockHTTPClient {
	mockClient := new(MockHTTPClient)
	if workerStatus == "" {
		mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").
			Return(io.NopCloser(bytes.NewReader([]byte(""))), errors.New("workers_status"))
	} else {
		mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").
			Return(io.NopCloser(bytes.NewReader([]byte(workerStatus))), nil)
	}

	return mockClient
}

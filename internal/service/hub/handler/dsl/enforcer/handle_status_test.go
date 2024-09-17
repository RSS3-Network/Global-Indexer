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
	nodeInfo                  = `{"data":{"operator":"0x5fdfd813ad20a90ba0972dd300ac9071c296b851","version":{"tag":"v1.0.0","commit":"8b36c72"}}}`
	nodeOutdated              = `{"data":{"operator":"0x5fdfd813ad20a90ba0972dd300ac9071c296b851","version":{"tag":"v0.9.0","commit":"8b36c72"}}}`
	workerStatusNodeIndexing  = `{"data":{"decentralized":[{"network":"farcaster","worker":"core","tags":null,"platform":"Unknown","status":"Indexing","remote_state":1718215438006,"indexed_state":1718215435040}],"rss":null,"federated":null}}`
	workerStatusNodeUnhealthy = `{"data":{"decentralized":[{"network":"farcaster","worker":"core","tags":null,"platform":"Unknown","status":"Unhealthy","remote_state":0,"indexed_state":0}],"rss":null,"federated":null}}`
	workerStatusNodeOnline    = `{"data":{"decentralized":[{"network":"farcaster","worker":"core","tags":null,"platform":"Unknown","status":"Ready","remote_state":0,"indexed_state":0}],"rss":null,"federated":null}}`
)

func TestDetermineStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		workerStatus    string
		nodeInfo        string
		initialStatus   schema.NodeStatus
		minVersion      string
		expectedStatus  schema.NodeStatus
		expectedErrPath string
	}{
		{
			name:            "NodeInfoUnavailable",
			workerStatus:    workerStatusNodeIndexing,
			nodeInfo:        "",
			initialStatus:   schema.NodeStatusInitializing,
			minVersion:      "1.0.0",
			expectedStatus:  schema.NodeStatusOffline,
			expectedErrPath: "info",
		},
		{
			name:            "WorkerStatusUnavailable",
			workerStatus:    "",
			nodeInfo:        nodeInfo,
			initialStatus:   schema.NodeStatusInitializing,
			minVersion:      "1.0.0",
			expectedStatus:  schema.NodeStatusOffline,
			expectedErrPath: "workers_status",
		},
		{
			name:            "RegisteredToOutdated",
			workerStatus:    workerStatusNodeIndexing,
			nodeInfo:        nodeInfo,
			initialStatus:   schema.NodeStatusRegistered,
			minVersion:      "1.1.0",
			expectedStatus:  schema.NodeStatusOutdated,
			expectedErrPath: "",
		},
		{
			name:            "RegisteredStaysRegistered",
			workerStatus:    workerStatusNodeUnhealthy,
			nodeInfo:        nodeInfo,
			initialStatus:   schema.NodeStatusRegistered,
			minVersion:      "1.0.0",
			expectedStatus:  schema.NodeStatusRegistered,
			expectedErrPath: "",
		},
		{
			name:            "RegisteredToInitializing",
			workerStatus:    workerStatusNodeIndexing,
			nodeInfo:        nodeInfo,
			initialStatus:   schema.NodeStatusRegistered,
			minVersion:      "1.0.0",
			expectedStatus:  schema.NodeStatusInitializing,
			expectedErrPath: "",
		},
		{
			name:            "OutdatedStaysOutdated",
			workerStatus:    workerStatusNodeIndexing,
			nodeInfo:        nodeInfo,
			initialStatus:   schema.NodeStatusOutdated,
			minVersion:      "1.1.0",
			expectedStatus:  schema.NodeStatusOutdated,
			expectedErrPath: "",
		},
		{
			name:            "OutdatedToRegistered",
			workerStatus:    workerStatusNodeUnhealthy,
			nodeInfo:        nodeInfo,
			initialStatus:   schema.NodeStatusOutdated,
			minVersion:      "1.0.0",
			expectedStatus:  schema.NodeStatusRegistered,
			expectedErrPath: "",
		},
		{
			name:            "OutdatedToInitializing",
			workerStatus:    workerStatusNodeIndexing,
			nodeInfo:        nodeInfo,
			initialStatus:   schema.NodeStatusOutdated,
			minVersion:      "1.0.0",
			expectedStatus:  schema.NodeStatusInitializing,
			expectedErrPath: "",
		},
		{
			name:            "InitializingToOutdated",
			workerStatus:    workerStatusNodeIndexing,
			nodeInfo:        nodeInfo,
			initialStatus:   schema.NodeStatusInitializing,
			minVersion:      "1.1.0",
			expectedStatus:  schema.NodeStatusOutdated,
			expectedErrPath: "",
		},
		{
			name:            "InitializingToRegistered",
			workerStatus:    workerStatusNodeUnhealthy,
			nodeInfo:        nodeInfo,
			initialStatus:   schema.NodeStatusInitializing,
			minVersion:      "1.0.0",
			expectedStatus:  schema.NodeStatusRegistered,
			expectedErrPath: "",
		},
		{
			name:            "InitializingStaysInitializing",
			workerStatus:    workerStatusNodeIndexing,
			nodeInfo:        nodeInfo,
			initialStatus:   schema.NodeStatusInitializing,
			minVersion:      "1.0.0",
			expectedStatus:  schema.NodeStatusInitializing,
			expectedErrPath: "",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockClient := setupMockClient(tt.workerStatus, tt.nodeInfo)
			enforcer := &SimpleEnforcer{httpClient: mockClient}
			stat := &schema.Stat{
				Address:     common.Address{},
				Endpoint:    "http://localhost:8080",
				AccessToken: "token",
				Epoch:       1,
				Status:      tt.initialStatus,
			}

			minVersion, _ := version.NewVersion(tt.minVersion)
			status, errPath := enforcer.determineStatus(ctx, stat, minVersion)

			assert.Equal(t, tt.expectedStatus, status)
			assert.Equal(t, tt.expectedErrPath, errPath)
		})
	}
}

func setupMockClient(workerStatus, nodeInfo string) *MockHTTPClient {
	mockClient := new(MockHTTPClient)
	if workerStatus == "" {
		mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").
			Return(io.NopCloser(bytes.NewReader([]byte(""))), errors.New("workers_status"))
	} else {
		mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").
			Return(io.NopCloser(bytes.NewReader([]byte(workerStatus))), nil)
	}

	if nodeInfo == "" {
		mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/info").
			Return(io.NopCloser(bytes.NewReader([]byte(""))), errors.New("info"))
	} else {
		mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/info").
			Return(io.NopCloser(bytes.NewReader([]byte(nodeInfo))), nil)
	}

	return mockClient
}

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

func Test_NodeInfoUnavailableReturnsOffline(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockClient := new(MockHTTPClient)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeIndexing))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/info").Return(io.NopCloser(bytes.NewReader([]byte(``))), errors.New("info"))

	enforcer := &SimpleEnforcer{
		httpClient: mockClient,
	}
	stat := &schema.Stat{
		Address:     common.Address{0},
		Endpoint:    "http://localhost:8080",
		AccessToken: "token",
		Epoch:       1,
		Status:      schema.NodeStatusInitializing,
	}

	minVersion, _ := version.NewVersion("1.0.0")

	status, errPath := enforcer.determineStatus(ctx, stat, minVersion)

	assert.Equal(t, schema.NodeStatusOffline, status)
	assert.Equal(t, "info", errPath)
}

func Test_NodeWorkerStatusUnavailableReturnsOffline(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockClient := new(MockHTTPClient)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(``))), errors.New("workers_status"))
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/info").Return(io.NopCloser(bytes.NewReader([]byte(nodeInfo))), nil)

	enforcer := &SimpleEnforcer{
		httpClient: mockClient,
	}
	stat := &schema.Stat{
		Address:     common.Address{0},
		Endpoint:    "http://localhost:8080",
		AccessToken: "token",
		Epoch:       1,
		Status:      schema.NodeStatusInitializing,
	}

	minVersion, _ := version.NewVersion("1.0.0")

	status, errPath := enforcer.determineStatus(ctx, stat, minVersion)

	assert.Equal(t, schema.NodeStatusOffline, status)
	assert.Equal(t, "workers_status", errPath)
}

func Test_NodeStatusRegisteredReturnsOutdated(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockClient := new(MockHTTPClient)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeIndexing))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/info").Return(io.NopCloser(bytes.NewReader([]byte(nodeInfo))), nil)

	enforcer := &SimpleEnforcer{
		httpClient: mockClient,
	}
	stat := &schema.Stat{
		Address:     common.Address{0},
		Endpoint:    "http://localhost:8080",
		AccessToken: "token",
		Epoch:       1,
		Status:      schema.NodeStatusRegistered,
	}

	minVersion, _ := version.NewVersion("1.1.0")

	status, errPath := enforcer.determineStatus(ctx, stat, minVersion)

	assert.Equal(t, schema.NodeStatusOutdated, status)
	assert.Equal(t, "", errPath)
}

func Test_NodeStatusRegisteredReturnsRegistered(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockClient := new(MockHTTPClient)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeUnhealthy))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/info").Return(io.NopCloser(bytes.NewReader([]byte(nodeInfo))), nil)

	enforcer := &SimpleEnforcer{
		httpClient: mockClient,
	}
	stat := &schema.Stat{
		Address:     common.Address{0},
		Endpoint:    "http://localhost:8080",
		AccessToken: "token",
		Epoch:       1,
		Status:      schema.NodeStatusRegistered,
	}

	minVersion, _ := version.NewVersion("1.0.0")

	status, errPath := enforcer.determineStatus(ctx, stat, minVersion)

	assert.Equal(t, schema.NodeStatusRegistered, status)
	assert.Equal(t, "", errPath)
}

func Test_NodeStatusRegisteredReturnsInitializing(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockClient := new(MockHTTPClient)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeIndexing))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/info").Return(io.NopCloser(bytes.NewReader([]byte(nodeInfo))), nil)

	enforcer := &SimpleEnforcer{
		httpClient: mockClient,
	}
	stat := &schema.Stat{
		Address:     common.Address{0},
		Endpoint:    "http://localhost:8080",
		AccessToken: "token",
		Epoch:       1,
		Status:      schema.NodeStatusRegistered,
	}

	minVersion, _ := version.NewVersion("1.0.0")

	status, errPath := enforcer.determineStatus(ctx, stat, minVersion)

	assert.Equal(t, schema.NodeStatusInitializing, status)
	assert.Equal(t, "", errPath)
}

func Test_NodeStatusOutdatedReturnsOutdated(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockClient := new(MockHTTPClient)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeIndexing))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/info").Return(io.NopCloser(bytes.NewReader([]byte(nodeInfo))), nil)

	enforcer := &SimpleEnforcer{
		httpClient: mockClient,
	}
	stat := &schema.Stat{
		Address:     common.Address{0},
		Endpoint:    "http://localhost:8080",
		AccessToken: "token",
		Epoch:       1,
		Status:      schema.NodeStatusOutdated,
	}

	minVersion, _ := version.NewVersion("1.1.0")

	status, errPath := enforcer.determineStatus(ctx, stat, minVersion)

	assert.Equal(t, schema.NodeStatusOutdated, status)
	assert.Equal(t, "", errPath)
}

func Test_NodeStatusOutdatedReturnsRegistered(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockClient := new(MockHTTPClient)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeUnhealthy))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/info").Return(io.NopCloser(bytes.NewReader([]byte(nodeInfo))), nil)

	enforcer := &SimpleEnforcer{
		httpClient: mockClient,
	}
	stat := &schema.Stat{
		Address:     common.Address{0},
		Endpoint:    "http://localhost:8080",
		AccessToken: "token",
		Epoch:       1,
		Status:      schema.NodeStatusOutdated,
	}

	minVersion, _ := version.NewVersion("1.0.0")

	status, errPath := enforcer.determineStatus(ctx, stat, minVersion)

	assert.Equal(t, schema.NodeStatusRegistered, status)
	assert.Equal(t, "", errPath)
}

func Test_NodeStatusOutdatedReturnsInitializing(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockClient := new(MockHTTPClient)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeIndexing))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/info").Return(io.NopCloser(bytes.NewReader([]byte(nodeInfo))), nil)

	enforcer := &SimpleEnforcer{
		httpClient: mockClient,
	}
	stat := &schema.Stat{
		Address:     common.Address{0},
		Endpoint:    "http://localhost:8080",
		AccessToken: "token",
		Epoch:       1,
		Status:      schema.NodeStatusOutdated,
	}

	minVersion, _ := version.NewVersion("1.0.0")

	status, errPath := enforcer.determineStatus(ctx, stat, minVersion)

	assert.Equal(t, schema.NodeStatusInitializing, status)
	assert.Equal(t, "", errPath)
}

func Test_NodeStatusInitializingReturnsOutdated(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockClient := new(MockHTTPClient)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeIndexing))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/info").Return(io.NopCloser(bytes.NewReader([]byte(nodeInfo))), nil)

	enforcer := &SimpleEnforcer{
		httpClient: mockClient,
	}
	stat := &schema.Stat{
		Address:     common.Address{0},
		Endpoint:    "http://localhost:8080",
		AccessToken: "token",
		Epoch:       1,
		Status:      schema.NodeStatusInitializing,
	}

	minVersion, _ := version.NewVersion("1.1.0")

	status, errPath := enforcer.determineStatus(ctx, stat, minVersion)

	assert.Equal(t, schema.NodeStatusOutdated, status)
	assert.Equal(t, "", errPath)
}

func Test_NodeStatusInitializingReturnsRegistered(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockClient := new(MockHTTPClient)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeUnhealthy))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/info").Return(io.NopCloser(bytes.NewReader([]byte(nodeInfo))), nil)

	enforcer := &SimpleEnforcer{
		httpClient: mockClient,
	}
	stat := &schema.Stat{
		Address:     common.Address{0},
		Endpoint:    "http://localhost:8080",
		AccessToken: "token",
		Epoch:       1,
		Status:      schema.NodeStatusInitializing,
	}

	minVersion, _ := version.NewVersion("1.0.0")

	status, errPath := enforcer.determineStatus(ctx, stat, minVersion)

	assert.Equal(t, schema.NodeStatusRegistered, status)
	assert.Equal(t, "", errPath)
}

func Test_NodeStatusInitializingReturnsInitializing(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mockClient := new(MockHTTPClient)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/workers_status").Return(io.NopCloser(bytes.NewReader([]byte(workerStatusNodeIndexing))), nil)
	mockClient.On("FetchWithMethod", mock.Anything, "http://localhost:8080/info").Return(io.NopCloser(bytes.NewReader([]byte(nodeInfo))), nil)

	enforcer := &SimpleEnforcer{
		httpClient: mockClient,
	}
	stat := &schema.Stat{
		Address:     common.Address{0},
		Endpoint:    "http://localhost:8080",
		AccessToken: "token",
		Epoch:       1,
		Status:      schema.NodeStatusInitializing,
	}

	minVersion, _ := version.NewVersion("1.0.0")

	status, errPath := enforcer.determineStatus(ctx, stat, minVersion)

	assert.Equal(t, schema.NodeStatusInitializing, status)
	assert.Equal(t, "", errPath)
}

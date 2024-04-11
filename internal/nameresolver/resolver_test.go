package nameresolver_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/nameresolver"
)

func Test_Resolve(t *testing.T) {
	t.Parallel()

	resolverConfig := &config.RPC{
		RPCNetwork: &config.RPCNetwork{
			Ethereum: &config.RPCEndpoint{
				Endpoint: "https://rpc.ankr.com/eth",
			},
			Crossbell: &config.RPCEndpoint{
				Endpoint: "https://rpc.crossbell.io",
			},
			Polygon: &config.RPCEndpoint{
				Endpoint: "https://rpc.ankr.com/polygon",
			},
			//Farcaster: &config.RPCEndpoint{
			//	Endpoint: "https://nemes.farcaster.xyz:2281",
			//},
		},
	}

	nr, _ := nameresolver.NewNameResolver(context.Background(), resolverConfig.RPCNetwork)

	type arguments struct {
		ns string
	}

	tests := []struct {
		name   string
		input  arguments
		output string
		err    error
	}{
		{
			name:   "unregister eth",
			input:  arguments{"qwerfdsazxcv.eth"},
			output: "",
			err:    fmt.Errorf("%s", nameresolver.ErrUnregisterName),
		},
		{
			name:   "unregister csb",
			input:  arguments{"qwerfdsazxcv.csb"},
			output: "",
			err:    fmt.Errorf("%s", nameresolver.ErrUnregisterName),
		},
		{
			name:   "unregister lens",
			input:  arguments{"qwerfdsazxcv.lens"},
			output: "",
			err:    fmt.Errorf("%s", nameresolver.ErrUnregisterName),
		},
		//{
		//	name:   "unregister farcaster",
		//	input:  arguments{"qwerfdsazxcv.fc"},
		//	output: "",
		//	err:    fmt.Errorf("%s", nameresolver.ErrUnregisterName),
		//},
		{
			name:   "unsupport name service .xxx",
			input:  arguments{"qwerfdsazxcv.xxx"},
			output: "",
			err:    fmt.Errorf("%s:%s", nameresolver.ErrUnSupportName, "qwerfdsazxcv.xxx"),
		},
		{
			name:   "resolve ens",
			input:  arguments{"vitalik.eth"},
			output: "0xd8da6bf26964af9d7eed9e03e53415d37aa96045",
			err:    nil,
		},
		{
			name:   "resolve csb",
			input:  arguments{"brucexx.csb"},
			output: "0x23c46e912b34C09c4bCC97F4eD7cDd762cee408A",
			err:    nil,
		},
		{
			name:   "resolve lens",
			input:  arguments{"diygod.lens"},
			output: "0xc8b960d09c0078c18dcbe7eb9ab9d816bcca8944",
			err:    nil,
		},
		//{
		//	name:   "resolve fc",
		//	input:  arguments{"brucexc.fc"},
		//	output: "0xe5d6216f0085a7f6b9b692e06cf5856e6fa41b55",
		//	err:    nil,
		//},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := nr.Resolve(context.Background(), tt.input.ns)
			if tt.err == nil {
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}
				if !strings.EqualFold(tt.output, output) {
					t.Errorf("Failure: %v => %v (expected %v)\n", tt.input, output, tt.output)
				}
			} else {
				if err == nil {
					t.Fatalf("missing expected error")
				}

				if !strings.Contains(err.Error(), tt.err.Error()) {
					t.Fatalf("Failure: unexpected error value %v, (expected %v)\n", err, tt.err)
				}
			}
		})
	}
}

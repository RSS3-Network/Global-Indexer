package hub

import (
	"errors"
	"testing"

	"github.com/naturalselectionlabs/rss3-global-indexer/internal/hub/model"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
)

// TestUpdateNodeStatus tests the UpdateNodeStatus function.
func TestUpdateNodeStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		initial    schema.NodeStatus
		newStatus  schema.NodeStatus
		shouldFail bool
	}{
		// valid transitions
		// NodeStatusRegistered
		{"Valid_Registered_Online", schema.NodeStatusRegistered, schema.NodeStatusOnline, false},
		{"Valid_Registered_Exited", schema.NodeStatusRegistered, schema.NodeStatusExited, false},

		// NodeStatusOnline
		{"Valid_Online_Exiting", schema.NodeStatusOnline, schema.NodeStatusExiting, false},
		{"Valid_Online_Exited", schema.NodeStatusOnline, schema.NodeStatusExited, false},
		{"Valid_Online_Slashed", schema.NodeStatusOnline, schema.NodeStatusSlashed, false},
		{"Valid_Online_Offline", schema.NodeStatusOnline, schema.NodeStatusOffline, false},

		// NodeStatusExiting
		{"Valid_Exiting_Exited", schema.NodeStatusExiting, schema.NodeStatusExited, false},

		// NodeStatusSlashed
		{"Valid_Slashed_Online", schema.NodeStatusSlashed, schema.NodeStatusOnline, false},
		{"Valid_Slashed_Offline", schema.NodeStatusSlashed, schema.NodeStatusOffline, false},

		// NodeStatusOffline
		{"Valid_Offline_Online", schema.NodeStatusOffline, schema.NodeStatusOnline, false},
		{"Valid_Offline_Exited", schema.NodeStatusOffline, schema.NodeStatusExited, false},

		// NodeStatusExited
		{"Valid_Exited_Registered", schema.NodeStatusExited, schema.NodeStatusRegistered, false},

		// invalid transitions
		// NodeStatusRegistered
		{"Invalid_Registered_Exiting", schema.NodeStatusRegistered, schema.NodeStatusExiting, true},
		{"Invalid_Registered_Slashed", schema.NodeStatusRegistered, schema.NodeStatusSlashed, true},
		{"Invalid_Registered_Offline", schema.NodeStatusRegistered, schema.NodeStatusOffline, true},

		// NodeStatusOnline
		{"Invalid_Online_Registered", schema.NodeStatusOnline, schema.NodeStatusRegistered, true},

		// NodeStatusExiting
		{"Invalid_Exiting_Online", schema.NodeStatusExiting, schema.NodeStatusOnline, true},
		{"Invalid_Exiting_Slashed", schema.NodeStatusExiting, schema.NodeStatusSlashed, true},
		{"Invalid_Exiting_Offline", schema.NodeStatusExiting, schema.NodeStatusOffline, true},
		{"Invalid_Exiting_Registered", schema.NodeStatusExiting, schema.NodeStatusRegistered, true},

		// NodeStatusSlashed
		{"Invalid_Slashed_Exiting", schema.NodeStatusSlashed, schema.NodeStatusExiting, true},
		{"Invalid_Slashed_Exited", schema.NodeStatusSlashed, schema.NodeStatusExited, true},
		{"Invalid_Slashed_Registered", schema.NodeStatusSlashed, schema.NodeStatusRegistered, true},

		// NodeStatusOffline
		{"Invalid_Offline_Exiting", schema.NodeStatusOffline, schema.NodeStatusExiting, true},
		{"Invalid_Offline_Slashed", schema.NodeStatusOffline, schema.NodeStatusSlashed, true},
		{"Invalid_Offline_Registered", schema.NodeStatusOffline, schema.NodeStatusRegistered, true},

		// NodeStatusExited
		{"Invalid_Exited_Online", schema.NodeStatusExited, schema.NodeStatusOnline, true},
		{"Invalid_Exited_Exiting", schema.NodeStatusExited, schema.NodeStatusExiting, true},
		{"Invalid_Exited_Slashed", schema.NodeStatusExited, schema.NodeStatusSlashed, true},
		{"Invalid_Exited_Offline", schema.NodeStatusExited, schema.NodeStatusOffline, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			node := &schema.Node{Status: tt.initial}
			err := model.UpdateNodeStatus(node, tt.newStatus)

			if tt.shouldFail {
				if err == nil {
					t.Errorf("expected an error, got none")
				} else {
					var nodeStatusTransitionError *model.NodeStatusTransitionError
					if !errors.As(err, &nodeStatusTransitionError) {
						t.Errorf("expected a schema.NodeStatusTransitionError, got %T", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				} else if node.Status != tt.newStatus {
					t.Errorf("expected status _ be %v, got %v", tt.newStatus, node.Status)
				}
			}
		})
	}
}

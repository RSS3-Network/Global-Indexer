package model

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/node/schema/worker"
	"github.com/rss3-network/protocol-go/schema"
	"github.com/rss3-network/protocol-go/schema/metadata"
	"github.com/rss3-network/protocol-go/schema/network"
	"github.com/rss3-network/protocol-go/schema/tag"
)

var (
	RssNodeCacheKey  = "nodes:rss"
	FullNodeCacheKey = "nodes:full"

	WorkerToNetworksMapKey  = "map:worker_to_networks"
	NetworkToWorkersMapKey  = "map:network_to_workers"
	PlatformToWorkersMapKey = "map:platform_to_workers"
	TagToWorkersMapKey      = "map:tag_to_workers"

	SubscribeNodeCacheKey = "epoch"

	// RequiredQualifiedNodeCount the required number of qualified Nodes
	RequiredQualifiedNodeCount = 3
	// RequiredVerificationCount the required number of verifications before a request is considered valid
	RequiredVerificationCount = 3
	// DemotionCountBeforeSlashing the number of demotions that trigger a slashing
	DemotionCountBeforeSlashing = 4

	// MutablePlatformMap is a map of mutable platforms which should be excluded from the data comparison.
	MutablePlatformMap = map[string]struct{}{
		worker.PlatformFarcaster.String(): {},
	}

	WorkerToNetworksMap  = make(map[worker.Worker][]string, len(worker.WorkerValues()))
	NetworkToWorkersMap  = make(map[network.Network][]string, len(network.NetworkValues()))
	PlatformToWorkersMap = make(map[worker.Platform][]string, len(worker.PlatformValues()))
	TagToWorkersMap      = make(map[tag.Tag][]string, len(tag.TagValues()))
)

// NodeEndpointCache stores the elements in the heap.
type NodeEndpointCache struct {
	Address      string `json:"address"`
	Endpoint     string `json:"endpoint"`
	Score        float64
	InvalidCount int64
}

// DataResponse represents the response returned by a Node.
// It is also used to store the verification result.
type DataResponse struct {
	Address  common.Address
	Endpoint string
	Data     []byte
	// A valid response must be non-null and non-error
	Valid bool
	Err   error
	// ValidPoint is the points given to the response
	ValidPoint int
	// InvalidPoint is the points given to the response when it is invalid
	InvalidPoint int
}

type ErrResponse struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

// ActivityResponse represents a single Activity in a response being returned to the requester.
type ActivityResponse struct {
	Data *Activity `json:"data"`
}

// ActivitiesResponse represents a list of Activity in a response being returned to the requester.
type ActivitiesResponse struct {
	Data []*Activity `json:"data"`
	Meta *MetaCursor `json:"meta,omitempty"`
}

type MetaCursor struct {
	Cursor string `json:"cursor"`
}

// Activity represents an activity.
type Activity struct {
	ID       string    `json:"id"`
	Owner    string    `json:"owner,omitempty"`
	Network  string    `json:"network"`
	Index    uint      `json:"index"`
	From     string    `json:"from"`
	To       string    `json:"to"`
	Tag      string    `json:"tag"`
	Type     string    `json:"type"`
	Platform string    `json:"platform,omitempty"`
	Actions  []*Action `json:"actions"`
}

// Action represents an action within an Activity.
type Action struct {
	Tag         string            `json:"tag"`
	Type        string            `json:"type"`
	Platform    string            `json:"platform,omitempty"`
	From        string            `json:"from"`
	To          string            `json:"to"`
	Metadata    metadata.Metadata `json:"metadata"`
	RelatedURLs []string          `json:"related_urls,omitempty"`
}

type Actions []*Action

var _ json.Unmarshaler = (*Action)(nil)

func (a *Action) UnmarshalJSON(bytes []byte) error {
	type ActionAlias Action

	type action struct {
		ActionAlias

		MetadataX json.RawMessage `json:"metadata"`
	}

	var temp action

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("unmarshal action: %w", err)
	}

	tagX, err := tag.TagString(temp.Tag)
	if err != nil {
		return fmt.Errorf("invalid action tag: %w", err)
	}

	typeX, err := schema.ParseTypeFromString(tagX, temp.Type)
	if err != nil {
		return fmt.Errorf("invalid action type: %w", err)
	}

	temp.Metadata, err = metadata.Unmarshal(typeX, temp.MetadataX)
	if err != nil {
		return fmt.Errorf("invalid action metadata: %w", err)
	}

	*a = Action(temp.ActionAlias)

	return nil
}

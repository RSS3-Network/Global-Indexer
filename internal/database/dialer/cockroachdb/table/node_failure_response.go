package table

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
)

type NodeFailureResponse struct {
	EpochID           uint64                           `gorm:"column:epoch_id;primaryKey"`
	Status            schema.NodeFailureResponseStatus `gorm:"column:status"`
	ValidatorNode     common.Address                   `gorm:"column:validator_node"`
	ValidatorRequest  string                           `gorm:"column:validator_request"`
	ValidatorResponse json.RawMessage                  `gorm:"column:validator_response"`
	VerifiedNode      common.Address                   `gorm:"column:verified_node"`
	VerifiedRequest   string                           `gorm:"column:verified_request"`
	VerifiedResponse  json.RawMessage                  `gorm:"column:verified_response"`
	CreatedAt         time.Time                        `gorm:"column:created_at"`
	UpdatedAt         time.Time                        `gorm:"column:updated_at"`
}

func (*NodeFailureResponse) TableName() string {
	return "node_response_failure"
}

func (n *NodeFailureResponse) Import(nodeResponseFailure *schema.NodeFailureResponse) {
	n.EpochID = nodeResponseFailure.EpochID
	n.Status = nodeResponseFailure.Status
	n.ValidatorNode = nodeResponseFailure.ValidatorNode
	n.ValidatorRequest = nodeResponseFailure.ValidatorRequest
	n.ValidatorResponse = nodeResponseFailure.ValidatorResponse
	n.VerifiedNode = nodeResponseFailure.VerifiedNode
	n.VerifiedRequest = nodeResponseFailure.VerifiedRequest
	n.VerifiedResponse = nodeResponseFailure.VerifiedResponse
}

func (n *NodeFailureResponse) Export() (*schema.NodeFailureResponse, error) {
	validatorRequest, err := extractPathAndParams(n.ValidatorRequest)
	if err != nil {
		return nil, err
	}

	verifiedRequest, err := extractPathAndParams(n.VerifiedRequest)
	if err != nil {
		return nil, err
	}

	return &schema.NodeFailureResponse{
		EpochID:           n.EpochID,
		Status:            n.Status,
		ValidatorNode:     n.ValidatorNode,
		ValidatorRequest:  validatorRequest,
		ValidatorResponse: n.ValidatorResponse,
		VerifiedNode:      n.VerifiedNode,
		VerifiedRequest:   verifiedRequest,
		VerifiedResponse:  n.VerifiedResponse,
		CreatedAt:         n.CreatedAt.Unix(),
	}, nil
}

func extractPathAndParams(endpoint string) (string, error) {
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return "", err
	}

	return strings.TrimPrefix(endpoint, parsedURL.Scheme+"://"+parsedURL.Host), nil
}

package nta

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

type RegisterNodeRequest struct {
	Address     common.Address  `json:"address" validate:"required"`
	Signature   string          `json:"signature" validate:"required"`
	Endpoint    string          `json:"endpoint" validate:"required"`
	Stream      json.RawMessage `json:"stream,omitempty"`
	Config      json.RawMessage `json:"config,omitempty"`
	Type        string          `json:"type" default:"normal"`
	AccessToken string          `json:"access_token" validate:"required_if=Type normal"`
}

type NodeHeartbeatRequest struct {
	Address   common.Address `json:"address" validate:"required"`
	Signature string         `json:"signature" validate:"required"`
	Endpoint  string         `json:"endpoint" validate:"required"`
	Timestamp int64          `json:"timestamp" validate:"required"`
}

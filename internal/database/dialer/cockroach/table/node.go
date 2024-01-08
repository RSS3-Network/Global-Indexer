package table

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
)

type Node struct {
	Address      common.Address  `gorm:"column:address;primaryKey"`
	Endpoint     string          `gorm:"column:endpoint"`
	IsPublicGood bool            `gorm:"column:is_public_good"`
	StreamURI    string          `gorm:"column:stream_uri"`
	Config       json.RawMessage `gorm:"column:config;type:jsonb"`
}

func (Node) TableName() string {
	return "node"
}

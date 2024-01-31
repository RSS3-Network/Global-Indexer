package model

import "github.com/rss3-network/protocol-go/schema/filter"

type NodeConfig struct {
	RSS           []*Module `json:"rss"`
	Federated     []*Module `json:"federated"`
	Decentralized []*Module `json:"decentralized"`
}

type Module struct {
	Network  filter.Network `json:"network"`
	Endpoint string         `json:"endpoint"`
	Worker   filter.Name    `json:"worker"`
}

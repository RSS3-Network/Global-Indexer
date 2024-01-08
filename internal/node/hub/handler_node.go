package hub

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/global-indexer/schema"
	"net/http"
)

func (h *Hub) GetNodesHandler(c echo.Context) error {
	request := BatchNodeRequest{}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	if request.NodeAddress != nil {
		var filteredData []*schema.Node
		for _, nodeAddress := range request.NodeAddress {
			for _, node := range data {
				if node.Address == nodeAddress {
					filteredData = append(filteredData, node)
				}
			}
		}
		data = filteredData
	}

	return c.JSON(http.StatusOK, BatchNodeResponse{
		Data:   data,
		Page:   0,
		Offset: 20,
		Total:  uint(len(data)),
	})
}

func (h *Hub) GetNodeHandler(c echo.Context) error {
	request := NodeRequest{}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	var node *schema.Node
	for _, n := range data {
		if n.Address == request.Address {
			node = n
			break
		}
	}

	if node == nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("node not found: %v", request.Address))
	}

	return c.JSON(http.StatusOK, NodeResponse{
		Data: node,
	})
}

type NodeRequest struct {
	Address common.Address `path:"id"`
}

type NodeResponse struct {
	Data *schema.Node `json:"data"`
}

type BatchNodeRequest struct {
	Page        uint             `query:"page"`
	Offset      uint             `query:"offset"`
	NodeAddress []common.Address `query:"node_address"`
}

type BatchNodeResponse struct {
	Data   []*schema.Node `json:"data"`
	Page   uint           `json:"page"`
	Offset uint           `json:"offset"`
	Total  uint           `json:"total"`
}

type Node struct {
	Address      common.Address `json:"address"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Endpoint     string         `json:"endpoint"`
	TaxFraction  uint64         `json:"taxFraction"`
	IsPublicGood bool           `json:"isPublicGood"`
	StreamURI    string         `json:"streamURI"`
}

var data = []*schema.Node{
	{
		Address:      common.HexToAddress("0x0"),
		Name:         "rss3-node",
		Description:  "rss3",
		Endpoint:     "https://node.rss3.dev/",
		TaxFraction:  10,
		IsPublicGood: false,
		StreamURI:    "https://node.rss3.dev/kafka",
	},
	{
		Address:      common.HexToAddress("0x1"),
		Name:         "google-node",
		Description:  "google",
		Endpoint:     "https://node.google.com/",
		TaxFraction:  10,
		IsPublicGood: true,
		StreamURI:    "https://node.google.com/kafka",
	},
}

package hub

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/global-indexer/schema"
	"github.com/naturalselectionlabs/rss3-node/config"
	"github.com/samber/lo"
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
		Cursor: data[len(data)-1].Address.String(),
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

func (h *Hub) RegisterNodeHandler(c echo.Context) error {
	var request RegisterNodeRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	if err := h.registerNode(c.Request().Context(), &request); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("register node failed: %v", err))
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("node registered: %v", request.Address))
}

type RegisterNodeRequest struct {
	Address      common.Address `json:"address"`
	Endpoint     string         `json:"endpoint"`
	IsPublicGood bool           `json:"isPublicGood"`
	Stream       *config.Stream `json:"stream"`
	Config       *config.Node   `json:"config"`
}

type NodeRequest struct {
	Address common.Address `path:"id"`
}

type NodeResponse struct {
	Data *schema.Node `json:"data"`
}

type BatchNodeRequest struct {
	Cursor      string           `query:"cursor"`
	NodeAddress []common.Address `query:"nodeAddress"`
}

type BatchNodeResponse struct {
	Data   []*schema.Node `json:"data"`
	Cursor string         `json:"cursor"`
}

var data = []*schema.Node{
	{
		Address:      common.HexToAddress("0x0"),
		Name:         "rss3-node",
		Description:  "rss3",
		Endpoint:     "https://node.rss3.dev",
		TaxFraction:  10,
		IsPublicGood: false,
		Stream: &config.Stream{
			Enable: lo.ToPtr(true),
			Driver: "kafka",
			Topic:  "rss3.node.feeds",
			URI:    "https://node.rss3.dev:9092",
		},
	},
	{
		Address:      common.HexToAddress("0x1"),
		Name:         "google-node",
		Description:  "google",
		Endpoint:     "https://node.google.com/",
		TaxFraction:  10,
		IsPublicGood: true,
		Stream: &config.Stream{
			Enable: lo.ToPtr(true),
			Driver: "kafka",
			Topic:  "rss3.node.feeds",
			URI:    "https://node.google.com:9092",
		},
	},
}

package hub

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Hub) GetNodesHandler(c echo.Context) error {
	request := NodeRequest{}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	data := []*Node{
		{
			Address:      "0x0",
			Name:         "rss3-node",
			Description:  "rss3",
			Endpoint:     "https://node.rss3.dev/",
			TaxFraction:  20,
			IsPublicGood: false,
			StreamURI:    "https://node.rss3.dev/kafka",
		},
		{
			Address:      "0x1",
			Name:         "google-node",
			Description:  "google",
			Endpoint:     "https://node.google.com/",
			TaxFraction:  10,
			IsPublicGood: true,
			StreamURI:    "https://node.google.com/kafka",
		},
	}

	if request.NodeAddress != nil {
		var filteredData []*Node

		for _, nodeAddress := range request.NodeAddress {
			for _, node := range data {
				if node.Address == nodeAddress {
					filteredData = append(filteredData, node)
				}
			}
		}

		data = filteredData
	}

	return c.JSON(http.StatusOK, NodeResponse{
		Data:   data,
		Page:   0,
		Offset: 20,
		Total:  uint(len(data)),
	})
}

type NodeRequest struct {
	Page        uint     `query:"page"`
	Offset      uint     `query:"offset"`
	NodeAddress []string `query:"address"`
}

type NodeResponse struct {
	Data   []*Node `json:"data"`
	Page   uint    `json:"page"`
	Offset uint    `json:"offset"`
	Total  uint    `json:"total"`
}

type Node struct {
	Address      string `json:"address"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Endpoint     string `json:"endpoint"`
	TaxFraction  uint64 `json:"taxFraction"`
	IsPublicGood bool   `json:"isPublicGood"`
	StreamURI    string `json:"streamURI"`
}

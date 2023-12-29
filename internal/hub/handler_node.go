package hub

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Hub) GetNodesHandler(c echo.Context) error {
	request := NodeRequest{}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	data := []*Node{
		{
			NodeAddress:     "0x0",
			NodeName:        "rss3-node",
			NodeDescription: "rss3",
			NodeEndpoint:    "https://node.rss3.dev/",
			NodeTax:         "1",
			NodeClaimedTax:  "1",
			IsPublicWelfare: false,
			StreamURI:       "https://node.rss3.dev/kafka",
			OperatorWebsite: "https://rss3.io/",
		},
		{
			NodeAddress:     "0x1",
			NodeName:        "google-node",
			NodeDescription: "google",
			NodeEndpoint:    "https://node.google.com/",
			NodeTax:         "1",
			NodeClaimedTax:  "1",
			IsPublicWelfare: true,
			StreamURI:       "https://node.google.com/kafka",
			OperatorWebsite: "https://google.com/",
		},
	}

	if request.NodeAddress != nil {
		var filteredData []*Node
		for _, nodeAddress := range request.NodeAddress {
			for _, node := range data {
				if node.NodeAddress == nodeAddress {
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
	NodeAddress []string `query:"node_address"`
}

type NodeResponse struct {
	Data   []*Node `json:"data"`
	Page   uint    `json:"page"`
	Offset uint    `json:"offset"`
	Total  uint    `json:"total"`
}

type Node struct {
	NodeAddress     string `json:"node_address"`
	NodeName        string `json:"node_name"`
	NodeDescription string `json:"node_description"`
	NodeEndpoint    string `json:"node_endpoint"`
	NodeTax         string `json:"node_tax"`
	NodeClaimedTax  string `json:"node_claimed_tax"`
	IsPublicWelfare bool   `json:"is_public_welfare"`
	StreamURI       string `json:"stream_uri"`
	OperatorWebsite string `json:"operator_website"`
}

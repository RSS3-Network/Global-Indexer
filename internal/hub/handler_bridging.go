package hub

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Hub) GetBridgingHandler(c echo.Context) error {
	request := BridgingRequest{}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	data := []*Bridging{
		{
			Sender:   "0xA",
			Receiver: "0xA",
			Block: &Block{
				Number:    1,
				Hash:      "0x01",
				Timestamp: uint64(time.Now().Unix()),
			},
			Transaction: &Transaction{
				Hash:   "0x02",
				Index:  70,
				Nonce:  147973,
				Status: "success",
			},
			Action: "deposit",
			Event: &BridgingEvent{
				Token: &Token{
					Address: &Address{
						Layer1: "0x1",
						Layer2: "0x2",
					},
				},
				Value:   "1000000000000000000",
				Decimal: 18,
			},
		},
		{
			Sender:   "0xB",
			Receiver: "0xB",
			Block: &Block{
				Number:    2,
				Hash:      "0x03",
				Timestamp: uint64(time.Now().Unix()),
			},
			Transaction: &Transaction{
				Hash:   "0x04",
				Index:  80,
				Nonce:  147977,
				Status: "failure",
			},
			Action: "stake",
			Event: &BridgingEvent{
				Token: &Token{
					Address: &Address{
						Layer1: "0x1",
						Layer2: "0x2",
					},
				},
				Value:   "1000000000000000000",
				Decimal: 18,
			},
		},
	}

	filteredData := make([]*Bridging, 0)

	for _, bridging := range data {
		if request.Address != "" && request.Address != bridging.Receiver && request.Address != bridging.Sender {
			continue
		}

		filteredData = append(filteredData, bridging)
	}

	return c.JSON(http.StatusOK, BridgingResponse{
		Data:   filteredData,
		Cursor: "",
	})
}

type BridgingRequest struct {
	Address string `query:"address"`
	Cursor  string `query:"cursor"`
}

type BridgingResponse struct {
	Data   []*Bridging `json:"data"`
	Cursor string      `json:"cursor"`
}

type Bridging struct {
	Sender      string         `json:"sender"`
	Receiver    string         `json:"receiver"`
	Block       *Block         `json:"block"`
	Transaction *Transaction   `json:"transaction"`
	Action      string         `json:"action"`
	Event       *BridgingEvent `json:"event"`
}

type BridgingEvent struct {
	Token   *Token `json:"token"`
	Value   string `json:"value"`
	Decimal uint   `json:"decimal"`
}

type Token struct {
	Address *Address `json:"address"`
}

type Address struct {
	Layer1 string `json:"layer1"`
	Layer2 string `json:"layer2"`
}

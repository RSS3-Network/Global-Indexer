package hub

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Hub) GetBridgingHandler(c echo.Context) error {
	request := BridgingRequest{}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	data := []*Bridging{
		{
			Sender:   "0x428AB2BA90Eba0a4Be7aF34C9Ac451ab061AC010",
			Receiver: "0x428AB2BA90Eba0a4Be7aF34C9Ac451ab061AC010",
			Block: &Block{
				Number:    18884279,
				Hash:      "0x2a43aaef44872d575f17ab0482e72b8da0083017d4d7ede5790399b10838e68d",
				Timestamp: 1703768603,
			},
			Transaction: &Transaction{
				Hash:   "0x427cc05cd02ef3e031cb89387bafd39fd91a440c81eceac9b14d40fdba4fd9cc",
				Index:  70,
				Nonce:  147973,
				Status: "success",
			},
			Action: "deposit",
			Event: &BridgingEvent{
				Token: &Token{
					Address: &Address{
						Layer1: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
						Layer2: "0x7F5c764cBc14f9669B88837ca1490cCa17c31607",
					},
					Value:   "132657840479",
					Decimal: 18,
				},
			},
		},
		{
			Sender:   "0x428AB2BA90Eba0a4Be7aF34C9Ac451ab061AC010",
			Receiver: "0x428AB2BA90Eba0a4Be7aF34C9Ac451ab061AC010",
			Block: &Block{
				Number:    18884279 + 1,
				Hash:      "0xd86e83801b9930975ede7d4d7103800ad8b27e2840ba71f575d27844feaa34a2",
				Timestamp: 1703768603 + 1,
			},
			Transaction: &Transaction{
				Hash:   "0xcc9df4abdf04d41b9caece18c044a19df93dfab78398bc130e3fe20dc50cc724",
				Index:  70,
				Nonce:  147973 + 1,
				Status: "success",
			},
			Action: "withdraw",
			Event: &BridgingEvent{
				Token: &Token{
					Address: &Address{
						Layer1: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
						Layer2: "0x7F5c764cBc14f9669B88837ca1490cCa17c31607",
					},
					Value:   "132657840479",
					Decimal: 18,
				},
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
	Token *Token `json:"token"`
	Data  string `json:"data"`
}

type Token struct {
	Address *Address `json:"address"`
	Value   string   `json:"value"`
	Decimal uint     `json:"decimal"`
}

type Address struct {
	Layer1 string `json:"layer1"`
	Layer2 string `json:"layer2"`
}

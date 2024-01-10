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
		// withdrawn
		{
			Sender:   "0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944",
			Receiver: "0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944",
			Token: Token{
				Address: Address{
					Layer1: "0x568F64582A377ea52d0067c4E430B9aE22A60473",
					Layer2: "0x4200000000000000000000000000000000000042",
				},
				Value:   "1000000000000000000", // 10 ^ 18
				Decimal: 18,
			},
			Data: "",
			Event: BridgingEvent{
				Deposit: &BridgingDepositStages{
					// initialize_deposit // https://sepolia.etherscan.io/tx/0x89d9438c50de4c0bb3d16bc1cbc9f9aade2ea61383764650badbc24aa4f6ba87
					Initialize: &BridgingStage{
						Block: Block{
							Hash:      "0x4b9edb02bbe2d11cdddf043410f6c0a12d26ce3bbb7f8c234600eebaff8b6819",
							Number:    5054034,
							Timestamp: 1704814068,
						},
						Transaction: Transaction{
							Hash:   "0x89d9438c50de4c0bb3d16bc1cbc9f9aade2ea61383764650badbc24aa4f6ba87",
							Index:  51,
							Nonce:  46,
							Status: "success",
						},
					},
					// finalize_deposit // https://scan.testnet.rss3.dev/tx/0x3e452b5f060bee5b9a858c1d43672cd47b58ed022ae9e6a96abd11514e70a48b
					Finalized: &BridgingStage{
						Block: Block{
							Hash:      "0xfcaa940d479facff6207e8eadfdb00788baca6306744800791f6a8e3195905b3",
							Number:    41694,
							Timestamp: 1704782364,
						},
						Transaction: Transaction{
							Hash:   "0x3e452b5f060bee5b9a858c1d43672cd47b58ed022ae9e6a96abd11514e70a48b",
							Index:  1,
							Nonce:  4,
							Status: "success",
						},
					},
				},
			},
		},
		{
			Sender:   "0x3B6D02A24Df681FFdf621D35D70ABa7adaAc07c1",
			Receiver: "0x3B6D02A24Df681FFdf621D35D70ABa7adaAc07c1",
			Token: Token{
				Address: Address{
					Layer1: "0xc575bd904d16a433624db98d01d5abd5c92d0f38",
					Layer2: "0x4200000000000000000000000000000000000010",
				},
				Value:   "1000000000000000000000", // 1000 * 10 ^ 18
				Decimal: 18,
			},
			Event: BridgingEvent{
				Withdraw: &BridgingWithdrawStages{
					// initialize_withdrawn // https://scan.testnet.rss3.dev/tx/0x788fd349c19c7ad3ae46234f5f35998e1f68b9dd9eb7a7a6b535286f8307ac1c
					Initialized: &BridgingStage{
						Block: Block{
							Number:    80918,
							Hash:      "0xbf2a3ea7e4b51746de9a36666b1c90278152472e3c3de9f25cf22ea57a671937",
							Timestamp: 1704860812,
						},
						Transaction: Transaction{
							Hash:   "0x788fd349c19c7ad3ae46234f5f35998e1f68b9dd9eb7a7a6b535286f8307ac1c",
							Index:  1,
							Nonce:  10,
							Status: "success",
						},
					},
					// prove_withdrawn // https://sepolia.etherscan.io/tx/0x0decd9659b3acd99e0c199a86eb7115f2fad0f800cb0cbbc76ae503d8523e85c
					Proved: &BridgingStage{
						Block: Block{
							Hash:      "0x958c59e7a92d5d68c975295d67a5dd249d4b9d6e352b8fe8d7b1f79a4a8e6cf8",
							Number:    5057720,
							Timestamp: 0,
						},
						Transaction: Transaction{
							Hash:   "0x0decd9659b3acd99e0c199a86eb7115f2fad0f800cb0cbbc76ae503d8523e85c",
							Index:  23,
							Nonce:  11,
							Status: "success",
						},
					},
					// finalize_withdrawn // https://sepolia.etherscan.io/tx/0x956bf19b609576725caf61d3056b482735b88a59afea772f666c4829f75017da
					Finalized: &BridgingStage{
						Block: Block{
							Hash:      "0x48a57e6aab93abb6248980c1a6ab9f75db2855dc1df02bd4dbc33b3f80fe87ec",
							Number:    5058517,
							Timestamp: 1704872532,
						},
						Transaction: Transaction{
							Hash:   "0x956bf19b609576725caf61d3056b482735b88a59afea772f666c4829f75017da",
							Index:  52,
							Nonce:  12,
							Status: "success",
						},
					},
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
	Sender   string        `json:"sender"`
	Receiver string        `json:"receiver"`
	Token    Token         `json:"token"`
	Data     string        `json:"data"`
	Event    BridgingEvent `json:"event"`
}

type BridgingEvent struct {
	Deposit  *BridgingDepositStages  `json:"deposit,omitempty"`
	Withdraw *BridgingWithdrawStages `json:"withdraw,omitempty"`
}

type BridgingStage struct {
	Block       Block       `json:"block"`
	Transaction Transaction `json:"transaction"`
}

type BridgingDepositStages struct {
	Initialize *BridgingStage `json:"initialized,omitempty"`
	Finalized  *BridgingStage `json:"finalized,omitempty"`
}

type BridgingWithdrawStages struct {
	Initialized *BridgingStage `json:"initialized,omitempty"`
	Proved      *BridgingStage `json:"proved,omitempty"`
	Finalized   *BridgingStage `json:"finalized,omitempty"`
}

type Token struct {
	Address Address `json:"address"`
	Value   string  `json:"value"`
	Decimal uint    `json:"decimal"`
}

type Address struct {
	Layer1 string `json:"layer1"`
	Layer2 string `json:"layer2"`
}

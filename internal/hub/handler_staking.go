package hub

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Hub) GetStakingHandler(c echo.Context) error {
	request := StakingRequest{}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	data := []*Staking{
		{
			UserAddress: "0xaA3f6537De9Ded08ACdC8a75773C70921B14bc34",
			NodeAddress: "0x43cb41B12907C37757a8CdCA80deD9eD7356f3Aa",
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
			Action: "staked",
			Event: &StakingEvent{
				Staked: &Staked{
					Value:        "1000000000000000000",
					StartTokenID: big.NewInt(0),
					EndTokenID:   big.NewInt(1),
				},
			},
		},
		{
			UserAddress: "0xaA3f6537De9Ded08ACdC8a75773C70921B14bc34",
			NodeAddress: "0x43cb41B12907C37757a8CdCA80deD9eD7356f3Aa",
			Block: &Block{
				Number:    18884279,
				Hash:      "0xd86e83801b9930975ede7d4d7103800ad8b27e2840ba71f575d27844feaa34a2",
				Timestamp: 1703768603,
			},
			Transaction: &Transaction{
				Hash:   "0xcc9df4abdf04d41b9caece18c044a19df93dfab78398bc130e3fe20dc50cc724",
				Index:  70,
				Nonce:  147973,
				Status: "success",
			},
			Action: "requestUnstake",
			Event: &StakingEvent{
				RequestUnstake: &RequestUnstake{
					ID: 0,
					ChipsIDs: []uint64{
						0, 1,
					},
				},
			},
		},
		{
			UserAddress: "0xaA3f6537De9Ded08ACdC8a75773C70921B14bc34",
			NodeAddress: "0x43cb41B12907C37757a8CdCA80deD9eD7356f3Aa",
			Block: &Block{
				Number:    18884279,
				Hash:      "0xd86e83801b9930975ede7d4d7103800ad8b27e2840ba71f575d27844feaa34a2",
				Timestamp: 1703768603,
			},
			Transaction: &Transaction{
				Hash:   "0xcc9df4abdf04d41b9caece18c044a19df93dfab78398bc130e3fe20dc50cc724",
				Index:  70,
				Nonce:  147973,
				Status: "success",
			},
			Action: "unstakeClaimed",
			Event: &StakingEvent{
				UnstakeClaimed: &UnstakeClaimed{
					RequestID:   0,
					NodeAddress: "0x43cb41B12907C37757a8CdCA80deD9eD7356f3Aa",
					User:        "0xaA3f6537De9Ded08ACdC8a75773C70921B14bc34",
					Amount:      "1000000000000000000",
				},
			},
		},
	}

	filteredData := make([]*Staking, 0)

	for _, staking := range data {
		if request.UserAddress != "" && request.UserAddress != staking.UserAddress {
			continue
		}

		if request.NodeAddress != "" && request.NodeAddress != staking.NodeAddress {
			continue
		}

		filteredData = append(filteredData, staking)
	}

	return c.JSON(http.StatusOK, StakingResponse{
		Data:   filteredData,
		Cursor: "",
	})
}

type StakingRequest struct {
	Cursor      string `query:"cursor"`
	UserAddress string `query:"userAddress"`
	NodeAddress string `query:"nodeAddress"`
}

type StakingResponse struct {
	Data   []*Staking `json:"data"`
	Cursor string     `json:"cursor"`
}

type Staking struct {
	UserAddress string        `json:"userAddress"`
	NodeAddress string        `json:"nodeAddress"`
	Block       *Block        `json:"block"`
	Transaction *Transaction  `json:"transaction"`
	Action      string        `json:"action"`
	Event       *StakingEvent `json:"event"`
}

type Block struct {
	Hash      string `json:"hash"`
	Number    uint64 `json:"number"`
	Timestamp uint64 `json:"timestamp"`
}

type Transaction struct {
	Hash   string `json:"hash"`
	Index  int64  `json:"index"`
	Nonce  int    `json:"nonce"`
	Status string `json:"status"`
}

type StakingEvent struct {
	Staked         *Staked         `json:"staked,omitempty"`
	RequestUnstake *RequestUnstake `json:"requestUnstake,omitempty"`
	UnstakeClaimed *UnstakeClaimed `json:"unstakeClaimed,omitempty"`
}

type Staked struct {
	Value        string   `json:"value"`
	StartTokenID *big.Int `json:"startTokenID"`
	EndTokenID   *big.Int `json:"endTokenID"`
}

type RequestUnstake struct {
	ID       uint64   `json:"id"`
	ChipsIDs []uint64 `json:"chipsIDs"`
}

type UnstakeClaimed struct {
	RequestID   int    `json:"requestID"`
	NodeAddress string `json:"nodeAddress"`
	User        string `json:"user"`
	Amount      string `json:"amount"`
}

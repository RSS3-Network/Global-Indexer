package hub

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Hub) GetStakingHandler(c echo.Context) error {
	request := StakingRequest{}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request: %v", err))
	}

	data := []*Staking{
		{
			UserAddress: "0xA",
			NodeAddress: "0x1",
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
			Action: "stake",
			Event: &StakingEvent{
				Stake: &Stake{
					Value:        "1000000000000000000",
					StartTokenID: big.NewInt(1),
					EndTokenID:   big.NewInt(4),
				},
			},
		},
		{
			UserAddress: "0xB",
			NodeAddress: "0x2",
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
			Event: &StakingEvent{
				Stake: &Stake{
					Value:        "1000000000000000000",
					StartTokenID: big.NewInt(2),
					EndTokenID:   big.NewInt(5),
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
	UserAddress string `query:"user_address"`
	NodeAddress string `query:"node_address"`
}

type StakingResponse struct {
	Data   []*Staking `json:"data"`
	Cursor string     `json:"cursor"`
}

type Staking struct {
	UserAddress string        `json:"user_address"`
	NodeAddress string        `json:"node_address"`
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
	Stake *Stake `json:"stake,omitempty"`
}

type Stake struct {
	Value        string   `json:"value"`
	StartTokenID *big.Int `json:"start_token_id"`
	EndTokenID   *big.Int `json:"end_token_id"`
}

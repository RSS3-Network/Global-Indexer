package nta

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/schema"
)

type GetStakerRequest struct {
	Address common.Address `param:"address" validate:"required"`
}

type GetStakerResponse schema.StakeStaker

func (n *NTA) GetStaker(c echo.Context) error {
	var request GetStakerRequest

	if err := c.Bind(&request); err != nil {
		return fmt.Errorf("bind request: %w", err)
	}

	if err := c.Validate(&request); err != nil {
		return fmt.Errorf("validate request: %w", err)
	}

	staker, err := n.databaseClient.FindStakeStaker(c.Request().Context(), request.Address)
	if err != nil {
		return fmt.Errorf("fetch stake staker by address %s: %w", request.Address, err)
	}

	response := GetStakerResponse(*staker)

	return c.JSON(http.StatusOK, response)
}

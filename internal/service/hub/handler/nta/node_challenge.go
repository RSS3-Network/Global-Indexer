package nta

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
)

var (
	registerMessage    = "I, %s, am signing this message for registering my intention to operate an RSS3 Node."
	hideTaxRateMessage = "I, %s, am signing this message for registering my intention to hide the tax rate on Explorer for my RSS3 Node."
)

func (n *NTA) GetNodeChallenge(c echo.Context) error {
	var request nta.NodeChallengeRequest

	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, fmt.Errorf("bind request: %w", err))
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, fmt.Errorf("validate failed: %w", err))
	}

	var data nta.NodeChallengeResponseData

	switch request.Type {
	case "":
		data = nta.NodeChallengeResponseData(fmt.Sprintf(registerMessage, strings.ToLower(request.Address.String())))
	case "hideTaxRate":
		data = nta.NodeChallengeResponseData(fmt.Sprintf(hideTaxRateMessage, strings.ToLower(request.Address.String())))
	default:
		return errorx.BadRequestError(c, fmt.Errorf("invalid challenge type: %s", request.Type))
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: data,
	})
}

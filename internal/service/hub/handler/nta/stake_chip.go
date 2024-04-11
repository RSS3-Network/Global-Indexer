package nta

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/hub/model/errorx"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/hub/model/nta"
	"github.com/naturalselectionlabs/rss3-global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

func (n *NTA) GetStakeChips(c echo.Context) error {
	var request nta.GetStakeChipsRequest
	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.InternalError(c, err)
	}

	stakeChipsQuery := schema.StakeChipsQuery{
		Cursor: request.Cursor,
		IDs:    request.IDs,
		Node:   request.Node,
		Owner:  request.Owner,
		Limit:  &request.Limit,
	}

	stakeChips, err := n.databaseClient.FindStakeChips(c.Request().Context(), stakeChipsQuery)
	if err != nil {
		return fmt.Errorf("find stake chips: %w", err)
	}

	// Get current chip values
	nodeAddresses := lo.Map(stakeChips, func(stakeChip *schema.StakeChip, _ int) common.Address {
		return stakeChip.Node
	})

	node, err := n.databaseClient.FindNodes(c.Request().Context(), schema.FindNodesQuery{
		NodeAddresses: nodeAddresses,
	})
	if err != nil {
		return fmt.Errorf("find nodes: %w", err)
	}

	values := lo.SliceToMap(node, func(node *schema.Node) (common.Address, decimal.Decimal) {
		return node.Address, node.MinTokensToStake
	})

	for _, chip := range stakeChips {
		chip.LatestValue = values[chip.Node]
	}

	var response nta.Response
	response.Data = lo.Map(stakeChips, func(stakeChip *schema.StakeChip, _ int) *nta.StakeChip {
		return nta.NewStakeChip(stakeChip, n.baseURL(c))
	})

	if length := len(stakeChips); length > 0 && length == request.Limit {
		response.Cursor = stakeChips[length-1].ID.String()
	}

	return c.JSON(http.StatusOK, response)
}

func (n *NTA) GetStakeChip(c echo.Context) error {
	var request nta.GetStakeChipRequest
	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.InternalError(c, err)
	}

	stakeChipQuery := schema.StakeChipQuery{
		ID: request.ID,
	}

	stakeChip, err := n.databaseClient.FindStakeChip(c.Request().Context(), stakeChipQuery)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNoContent)
		}

		return err
	}

	node, err := n.databaseClient.FindNode(c.Request().Context(), stakeChip.Node)
	if err != nil {
		return fmt.Errorf("find node: %w", err)
	}

	stakeChip.LatestValue = node.MinTokensToStake

	var response nta.Response
	response.Data = nta.NewStakeChip(stakeChip, n.baseURL(c))

	return c.JSON(http.StatusOK, response)
}

func (n *NTA) GetStakeChipImage(c echo.Context) error {
	var request nta.GetStakeChipsImageRequest
	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidateFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.InternalError(c, err)
	}

	stakeChipQuery := schema.StakeChipQuery{
		ID: request.ID,
	}

	chip, err := n.databaseClient.FindStakeChip(c.Request().Context(), stakeChipQuery)
	if err != nil {
		return fmt.Errorf("find stake chip: %w", err)
	}

	var metadata l2.ChipsTokenMetadata
	if err := json.Unmarshal(chip.Metadata, &metadata); err != nil {
		return fmt.Errorf("invalid metadata: %w", err)
	}

	data, found := strings.CutPrefix(metadata.Image, "data:image/svg+xml;base64,")
	if !found {
		return fmt.Errorf("invalid image")
	}

	content, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return fmt.Errorf("invalid data: %w", err)
	}

	return c.Blob(http.StatusOK, "image/svg+xml", content)
}

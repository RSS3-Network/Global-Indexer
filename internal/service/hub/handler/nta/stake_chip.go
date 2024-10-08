package nta

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/creasty/defaults"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/rss3-network/global-indexer/schema"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (n *NTA) GetStakeChips(c echo.Context) error {
	var request nta.GetStakeChipsRequest
	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.InternalError(c)
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
		zap.L().Error("find stake chips from database", zap.Error(err))

		return errorx.InternalError(c)
	}

	chipIDs := make([]*big.Int, len(stakeChips))

	for i, chip := range stakeChips {
		chipIDs[i] = chip.ID
	}

	chipsInfo, err := n.stakingContract.StakingV2GetChipsInfo(c.Request().Context(), nil, chipIDs)
	if err != nil {
		zap.L().Error("get chips info by multicall", zap.Error(err))

		return errorx.InternalError(c)
	}

	for i, chipInfo := range chipsInfo {
		stakeChips[i].LatestValue = decimal.NewFromBigInt(chipInfo.Tokens, 0)
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
		return errorx.ValidationFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.InternalError(c)
	}

	stakeChipQuery := schema.StakeChipQuery{
		ID: request.ChipID,
	}

	stakeChip, err := n.databaseClient.FindStakeChip(c.Request().Context(), stakeChipQuery)
	if err != nil {
		if errors.Is(err, database.ErrorRowNotFound) {
			return c.NoContent(http.StatusNoContent)
		}

		return errorx.InternalError(c)
	}

	chipInfo, err := n.stakingContract.GetChipInfo(&bind.CallOpts{Context: c.Request().Context()}, stakeChip.ID)
	if err != nil {
		zap.L().Error("get chip info from rpc", zap.Error(err), zap.String("chipID", stakeChip.ID.String()))

		return fmt.Errorf("get chip info: %w", err)
	}

	stakeChip.LatestValue = decimal.NewFromBigInt(chipInfo.Tokens, 0)

	var response nta.Response
	response.Data = nta.NewStakeChip(stakeChip, n.baseURL(c))

	return c.JSON(http.StatusOK, response)
}

// TODO: add redis cache
func (n *NTA) GetStakeChipImage(c echo.Context) error {
	var request nta.GetStakeChipsImageRequest
	if err := c.Bind(&request); err != nil {
		return errorx.BadParamsError(c, err)
	}

	if err := c.Validate(&request); err != nil {
		return errorx.ValidationFailedError(c, err)
	}

	if err := defaults.Set(&request); err != nil {
		return errorx.InternalError(c)
	}

	stakeChipQuery := schema.StakeChipQuery{
		ID: request.ChipID,
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

package handlers

import (
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database/dialer/cockroachdb/table"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/utils"
	"github.com/samber/lo"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
)

func (a *App) GetDepositHistory(ctx echo.Context, params oapi.GetDepositHistoryParams) error {
	user := ctx.Get("user").(*table.GatewayAccount)

	// Parse date
	limit, page := parseLimitPage(params.Limit, params.Page)

	// Query from database
	query := a.databaseClient.WithContext(ctx.Request().Context()).Model(&table.BillingRecordDeposited{}).Where("user = ?", user.Address)

	var totalCount int64
	err := query.Count(&totalCount).Error
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	var records []table.BillingRecordDeposited
	err = query.Order("block_timestamp DESC").Offset(limit * (page - 1)).Limit(limit).Find(&records).Error
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	var list []oapi.DepositRecord
	for _, record := range records {
		blockTs := record.BlockTimestamp.UnixMilli()
		amount, _ := record.Amount.Float64()
		list = append(list, oapi.DepositRecord{
			TxHash:         lo.ToPtr(record.TxHash.Hex()),
			BlockTimestamp: &blockTs,
			Index:          lo.ToPtr(int(record.Index)),
			Amount:         lo.ToPtr(float32(amount)),
		})
	}

	count := uint(len(list))
	pageCurrent := int64(page)
	pageMax := int64(math.Ceil(float64(totalCount) / float64(limit)))
	return ctx.JSON(http.StatusOK, oapi.DepositHistoryResponse{
		Count:       &count,
		PageCurrent: &pageCurrent,
		PageMax:     &pageMax,
		List:        &list,
	})
}

func (*App) GetWithdrawalHistory(ctx echo.Context, params oapi.GetWithdrawalHistoryParams) error {
	user := ctx.Get("user").(*model.Account)

	// Parse date
	limit, page := parseLimitPage(params.Limit, params.Page)

	// Query from database
	query := app.EntClient.WithdrawalRecord.Query().Where(
		withdrawalrecord.User(user.Address),
	)

	totalCount, err := query.Count(ctx.Request().Context())
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	records, err := query.Order(
		entschema.Desc(withdrawalrecord.FieldBlockTimestamp),
	).Offset(limit * (page - 1)).Limit(limit).All(ctx.Request().Context())
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	var list []oapi.WithdrawalRecord
	for _, record := range records {
		blockTs := record.BlockTimestamp.UnixMilli()
		list = append(list, oapi.WithdrawalRecord{
			TxHash:         &record.TxHash,
			BlockTimestamp: &blockTs,
			Index:          to.Uint_IntPtr(record.Index),
			User:           &record.User,
			Amount:         to.Float64_Float32Ptr(record.Amount),
			Fee:            to.Float64_Float32Ptr(record.Fee),
		})
	}

	count := uint(len(list))
	pageCurrent := int64(page)
	pageMax := int64(math.Ceil(float64(totalCount) / float64(limit)))
	return ctx.JSON(http.StatusOK, oapi.WithdrawalHistoryResponse{
		Count:       &count,
		PageCurrent: &pageCurrent,
		PageMax:     &pageMax,
		List:        &list,
	})
}

func (*App) GetCollectionHistory(ctx echo.Context, params oapi.GetCollectionHistoryParams) error {
	user := ctx.Get("user").(*model.Account)

	limit, page := parseLimitPage(params.Limit, params.Page)

	// Query from database
	query := app.EntClient.CollectionRecord.Query().Where(
		collectionrecord.User(user.Address),
	)

	totalCount, err := query.Count(ctx.Request().Context())
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	records, err := query.Order(
		entschema.Desc(collectionrecord.FieldBlockTimestamp),
	).Offset(limit * (page - 1)).Limit(limit).All(ctx.Request().Context())
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	var list []oapi.CollectionRecord
	for _, record := range records {
		blockTs := record.BlockTimestamp.UnixMilli()
		list = append(list, oapi.CollectionRecord{
			TxHash:         &record.TxHash,
			BlockTimestamp: &blockTs,
			Index:          to.Uint_IntPtr(record.Index),
			Amount:         to.Float64_Float32Ptr(record.Amount),
		})
	}

	count := uint(len(list))
	pageCurrent := int64(page)
	pageMax := int64(math.Ceil(float64(totalCount) / float64(limit)))
	return ctx.JSON(http.StatusOK, oapi.CollectionHistoryResponse{
		Count:       &count,
		PageCurrent: &pageCurrent,
		PageMax:     &pageMax,
		List:        &list,
	})
}

func (app *App) GetConsumptionHistoryByKey(ctx echo.Context, key int, params oapi.GetConsumptionHistoryByKeyParams) error {
	k, err := getKey(ctx, key)
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusUnauthorized)
	}

	since, until := parseDates(params.Since, params.Until)
	logs, err := k.QueryConsumptionLogs().
		Where(
			consumptionlog.ConsumptionDateGTE(since),
			consumptionlog.ConsumptionDateLTE(until),
		).
		Order(entschema.Desc(
			consumptionlog.FieldConsumptionDate,
		)).
		All(ctx.Request().Context())
	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	sinceInt64 := since.UnixMilli()
	untilInt64 := until.UnixMilli()
	resp := &oapi.ConsumptionLogResponse{
		Since:   &sinceInt64,
		Until:   &untilInt64,
		History: &[]oapi.ConsumptionLogByKey{},
	}
	if params.Merge != nil && *params.Merge {
		var (
			apiCalls int64 = 0
			ruUsed   int64 = 0
		)
		for _, log := range logs {
			apiCalls += log.APICalls
			ruUsed += log.RuUsed
		}
		*resp.History = append(*resp.History, oapi.ConsumptionLogByKey{
			KeyName:  &k.Name,
			ApiCalls: &apiCalls,
			RuUsed:   &ruUsed,
		})
	} else {
		for _, log := range logs {
			consumptionDate := log.ConsumptionDate.UnixMilli()
			*resp.History = append(*resp.History, oapi.ConsumptionLogByKey{
				KeyName:         &k.Name,
				ConsumptionDate: &consumptionDate,
				ApiCalls:        &log.APICalls,
				RuUsed:          &log.RuUsed,
			})
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (app *App) GetConsumptionHistoryByAccount(ctx echo.Context, params oapi.GetConsumptionHistoryByAccountParams) error {
	user := ctx.Get("user").(*model.Account)

	since, until := parseDates(params.Since, params.Until)

	results, err := user.GetUsageByDate(ctx.Request().Context(), since, until)

	if err != nil {
		return utils.SendJSONError(ctx, http.StatusInternalServerError)
	}

	sinceInt64 := since.UnixMilli()
	untilInt64 := until.UnixMilli()
	resp := &oapi.ConsumptionLogResponse{
		Since:   &sinceInt64,
		Until:   &untilInt64,
		History: &[]oapi.ConsumptionLogByKey{},
	}
	if results != nil {
		if params.Merge != nil && *params.Merge {
			var (
				apiCalls int64 = 0
				ruUsed   int64 = 0
			)
			for _, log := range *results {
				apiCalls += *log.ApiCalls
				ruUsed += *log.RuUsed
			}
			*resp.History = append(*resp.History, oapi.ConsumptionLogByKey{
				ApiCalls: &apiCalls,
				RuUsed:   &ruUsed,
			})
		} else {
			*resp.History = append(*resp.History, *results...)
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

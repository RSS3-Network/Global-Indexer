package nta

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type GetTokensSupplyResponse struct {
	TotalSupply decimal.Decimal `json:"total_supply"`
}

func (n *NTA) GetTokenSupply(c echo.Context) error {
	totalSupply, err := n.getTokensTotalSupply()
	if err != nil {
		zap.L().Error("get token total supply", zap.Error(err))
		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: GetTokensSupplyResponse{
			TotalSupply: totalSupply,
		},
	})
}

func (n *NTA) getTokensTotalSupply() (decimal.Decimal, error) {
	totalSupply, err := n.contractGovernanceToken.TotalSupply(nil)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get total supply: %w", err)
	}

	return decimal.NewFromBigInt(totalSupply, -18), nil
}

type GetTvlResponse struct {
	Tvl decimal.Decimal `json:"tvl"`
}

type TokenPrice struct {
	Rss3 struct {
		Usd float64 `json:"usd"`
	} `json:"rss3"`
	Weth struct {
		Usd float64 `json:"usd"`
	} `json:"weth"`
	Usdt struct {
		Usd float64 `json:"usd"`
	} `json:"usdt"`
	Usdc struct {
		Usd float64 `json:"usd"`
	} `json:"usdc"`
}

func (n *NTA) GetTvl(c echo.Context) error {
	tokenPriceMap, err := n.getTokenPrices(c.Request().Context())
	if err != nil {
		zap.L().Error("get token price", zap.Error(err))
		return errorx.InternalError(c)
	}

	tvl := decimal.Zero

	for name, token := range n.erc20TokenMap {
		balance, err := token.BalanceOf(nil, n.addressL1StandardBridgeProxy)
		if err != nil {
			zap.L().Error("get token balance", zap.String("token", name), zap.Error(err))
			return errorx.InternalError(c)
		}

		tokenValue := calculateTokenValue(name, balance, tokenPriceMap)
		tvl = tvl.Add(tokenValue)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: GetTvlResponse{Tvl: tvl},
	})
}

func (n *NTA) getTokenPrices(ctx context.Context) (map[string]string, error) {
	var tokenPriceMap map[string]string
	if err := n.cacheClient.Get(ctx, "token:price:map", &tokenPriceMap); err == nil {
		return tokenPriceMap, nil
	}

	body, err := n.httpClient.FetchWithMethod(ctx, http.MethodGet, n.configFile.TokenPriceAPI.Endpoint, n.configFile.TokenPriceAPI.AuthToken, nil)
	if err != nil {
		return nil, fmt.Errorf("get token price from coingecko: %w", err)
	}

	var tokenPrice TokenPrice
	if err = json.NewDecoder(body).Decode(&tokenPrice); err != nil {
		return nil, fmt.Errorf("parse token price from response body: %w", err)
	}

	tokenPriceMap = map[string]string{
		"rss3": strconv.FormatFloat(tokenPrice.Rss3.Usd, 'f', -1, 64),
		"weth": strconv.FormatFloat(tokenPrice.Weth.Usd, 'f', -1, 64),
	}

	if err = n.cacheClient.Set(ctx, "token:price:map", tokenPriceMap, 30*60*time.Second); err != nil {
		zap.L().Warn("set token price map to cache", zap.Error(err))
	}

	return tokenPriceMap, nil
}

func calculateTokenValue(name string, balance *big.Int, priceMap map[string]string) decimal.Decimal {
	switch name {
	case "rss3", "weth":
		price, _ := decimal.NewFromString(priceMap[name])
		return decimal.NewFromBigInt(balance, -18).Mul(price)
	case "usdt", "usdc":
		return decimal.NewFromBigInt(balance, -6)
	default:
		return decimal.Zero
	}
}

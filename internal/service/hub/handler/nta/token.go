package nta

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/rss3-network/global-indexer/contract/l1"
	"github.com/rss3-network/global-indexer/contract/l2"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/errorx"
	"github.com/rss3-network/global-indexer/internal/service/hub/model/nta"
	"github.com/shopspring/decimal"
	"github.com/sourcegraph/conc/pool"
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

	return decimal.NewFromBigInt(totalSupply, 0), nil
}

type GetTvlResponse struct {
	Tvl decimal.Decimal `json:"tvl"`
}

func (n *NTA) GetTvl(c echo.Context) error {
	ctx := c.Request().Context()
	tokenPriceMap, err := n.getTokenPrices(ctx)

	if err != nil {
		zap.L().Error("get token price", zap.Error(err))
		return errorx.InternalError(c)
	}

	var (
		tvl decimal.Decimal
		mu  sync.Mutex
	)

	p := pool.New().WithContext(ctx).WithCancelOnError().WithMaxGoroutines(10)

	for address, bind := range n.erc20TokenMap {
		address, bind := address, bind

		p.Go(func(_ context.Context) error {
			var (
				balance *big.Int
				err     error
			)

			if address == l2.ContractMap[n.chainL2ID].AddressPowerToken {
				balance, err = bind.TotalSupply(nil)
			} else {
				balance, err = bind.BalanceOf(nil, l1.ContractMap[n.chainL1ID].AddressL1StandardBridgeProxy)
			}

			if err != nil {
				zap.L().Error("get token balance", zap.String("token", address.String()), zap.Error(err))
				return err
			}

			tokenValue := n.calculateTokenValue(address, balance, tokenPriceMap)

			mu.Lock()
			tvl = tvl.Add(tokenValue)
			mu.Unlock()

			return nil
		})
	}

	if err := p.Wait(); err != nil {
		zap.L().Error("get tvl", zap.Error(err))
		return errorx.InternalError(c)
	}

	return c.JSON(http.StatusOK, nta.Response{
		Data: GetTvlResponse{Tvl: tvl},
	})
}

type TokenPrice struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			TokenPrices map[string]string `json:"token_prices"`
		} `json:"attributes"`
	} `json:"data"`
}

const (
	tokenPriceKey = "token:price:map"

	// Get list of supported networks from https://api.geckoterminal.com/api/v2/networks
	ethereumNetwork = "eth"
	rss3Network     = "rss3-vsl-mainnet"
)

func (n *NTA) getTokenPrices(ctx context.Context) (map[string]string, error) {
	tokenPriceMap := make(map[string]string, 3)
	if err := n.cacheClient.Get(ctx, tokenPriceKey, &tokenPriceMap); err == nil && len(tokenPriceMap) == 3 {
		return tokenPriceMap, nil
	}

	networks := []string{ethereumNetwork, rss3Network}
	endpoint := n.configFile.TokenPriceAPI.Endpoint

	for _, network := range networks {
		var addressList []string
		if network == ethereumNetwork {
			addressList = []string{
				l1.ContractMap[n.chainL1ID].AddressWETHToken.String(),
				l1.ContractMap[n.chainL1ID].AddressGovernanceTokenProxy.String(),
			}
		} else {
			addressList = []string{
				l2.ContractMap[n.chainL2ID].AddressPowerToken.String(),
			}
		}

		url := fmt.Sprintf("%s/simple/networks/%s/token_price/%s", endpoint, network, strings.Join(addressList, ","))

		body, err := n.httpClient.FetchWithMethod(ctx, http.MethodGet, url, n.configFile.TokenPriceAPI.AuthToken, nil)
		if err != nil {
			return nil, fmt.Errorf("get token price: %w", err)
		}

		var tokenPrice TokenPrice
		if err = json.NewDecoder(body).Decode(&tokenPrice); err != nil {
			return nil, fmt.Errorf("parse token price from response body: %w", err)
		}

		for address, price := range tokenPrice.Data.Attributes.TokenPrices {
			tokenPriceMap[common.HexToAddress(address).String()] = price
		}
	}

	if err := n.cacheClient.Set(ctx, tokenPriceKey, tokenPriceMap, 30*60*time.Second); err != nil {
		zap.L().Warn("set token price map to cache", zap.Error(err))
	}

	return tokenPriceMap, nil
}

func (n *NTA) calculateTokenValue(address common.Address, balance *big.Int, priceMap map[string]string) decimal.Decimal {
	switch address {
	case l1.ContractMap[n.chainL1ID].AddressGovernanceTokenProxy, l1.ContractMap[n.chainL1ID].AddressWETHToken, l2.ContractMap[n.chainL2ID].AddressPowerToken:
		price, _ := decimal.NewFromString(priceMap[address.String()])
		return decimal.NewFromBigInt(balance, -18).Mul(price)
	case l1.ContractMap[n.chainL1ID].AddressUSDTToken, l1.ContractMap[n.chainL1ID].AddressUSDCToken:
		return decimal.NewFromBigInt(balance, -6)
	default:
		return decimal.Zero
	}
}

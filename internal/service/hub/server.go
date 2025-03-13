package hub

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/geolite2"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/common/txmgr"
	"github.com/rss3-network/global-indexer/docs"
	"github.com/rss3-network/global-indexer/internal/client/ethereum"
	"github.com/rss3-network/global-indexer/internal/config"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/nameresolver"
	"github.com/rss3-network/global-indexer/internal/service"
	"go.uber.org/zap"
)

const Name = "hub"

const (
	DefaultHost = "0.0.0.0"
	DefaultPort = "80"
)

type Server struct {
	httpServer *echo.Echo
	hub        *Hub
}

func (s *Server) Name() string {
	return Name
}

func (s *Server) Run(_ context.Context) error {
	address := net.JoinHostPort(DefaultHost, DefaultPort)

	return s.httpServer.Start(address)
}

func NewServer(databaseClient database.Client, redisClient *redis.Client, geoLite2 *geolite2.Client, ethereumMultiChainClient *ethereum.MultiChainClient, nameService *nameresolver.NameResolver, httpClient httputil.Client, txManager *txmgr.SimpleTxManager, config *config.File) (service.Server, error) {
	hub, err := NewHub(context.Background(), databaseClient, redisClient, ethereumMultiChainClient, geoLite2, nameService, httpClient, txManager, config)
	if err != nil {
		return nil, fmt.Errorf("new hub: %w", err)
	}

	instance := Server{
		httpServer: echo.New(),
		hub:        hub,
	}

	{
		// setup prometheus metrics
		instance.httpServer.Use(echoprometheus.NewMiddleware(Name))
		instance.httpServer.GET("/metrics", echoprometheus.NewHandler())
	}

	instance.httpServer.HideBanner = true
	instance.httpServer.HidePort = true
	instance.httpServer.Validator = defaultValidator
	instance.httpServer.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	{
		docsFile, err := docs.Generate()
		if err != nil {
			zap.L().Error("generate docs error", zap.Error(err))
		}

		instance.httpServer.GET("/docs/openapi.json", func(c echo.Context) error {
			return c.Blob(http.StatusOK, "application/json", docsFile)
		})
	}

	nodes := instance.httpServer.Group("/nodes")
	{
		// Node registration and heartbeat
		nodes.POST("/register", instance.hub.nta.RegisterNode)
		nodes.POST("/heartbeat", instance.hub.nta.NodeHeartbeat)
	}

	// nta is short for Network Transparency API
	nta := instance.httpServer.Group("/nta")
	{
		bridge := nta.Group("/bridgings")
		{
			bridge.GET("/transactions", instance.hub.nta.GetBridgeTransactions)
			bridge.GET("/transactions/:transaction_hash", instance.hub.nta.GetBridgeTransaction)
		}

		chips := nta.Group("/chips")
		{
			chips.GET("", instance.hub.nta.GetStakeChips)
			chips.GET("/:chip_id", instance.hub.nta.GetStakeChip)
			chips.GET("/:chip_id/image.svg", instance.hub.nta.GetStakeChipImage)
		}

		epochs := nta.Group("/epochs")
		{
			epochs.GET("", instance.hub.nta.GetEpochs)
			epochs.GET("/:epoch_id", instance.hub.nta.GetEpoch)
			epochs.GET("/:node_address/rewards", instance.hub.nta.GetEpochNodeRewards)
			epochs.GET("/distributions/:transaction_hash", instance.hub.nta.GetEpochDistribution)
			epochs.GET("/apy", instance.hub.nta.GetEpochsAPY)
		}

		networks := nta.Group("/networks")
		{
			networks.GET("/config", instance.hub.nta.GetNetworkConfig)
			networks.GET("/assets", instance.hub.nta.GetAssets)
		}

		nodes := nta.Group("/nodes")
		{
			nodes.GET("", instance.hub.nta.GetNodes)
			nodes.GET("/:node_address", instance.hub.nta.GetNode)
			nodes.GET("/:node_address/avatar.svg", instance.hub.nta.GetNodeAvatar)
			nodes.GET("/:node_address/challenge", instance.hub.nta.GetNodeChallenge)
			nodes.GET("/:node_address/events", instance.hub.nta.GetNodeEvents)
			nodes.GET("/:node_address/operation/profit", instance.hub.nta.GetNodeOperationProfit)

			nodes.POST("/:node_address/hide_tax_rate", instance.hub.nta.PostNodeHideTaxRate)
		}

		snapshots := nta.Group("/snapshots")
		{
			snapshots.GET("/nodes/count", instance.hub.nta.GetNodeCountSnapshots)
			snapshots.GET("/nodes/operation/profit", instance.hub.nta.GetNodeOperationProfitSnapshots)
			snapshots.GET("/stakers/count", instance.hub.nta.GetStakerCountSnapshots)
			snapshots.GET("/stakers/profit", instance.hub.nta.GetStakerProfitSnapshots)
			snapshots.GET("/epochs/apy", instance.hub.nta.GetEpochsAPYSnapshots)
		}

		stake := nta.Group("/stakings")
		{
			stake.GET("/:staker_address/profit", instance.hub.nta.GetStakerProfit)
			// FIXME: GetStakeStakings needs to be refactored
			// see https://github.com/RSS3-Network/Global-Indexer/issues/233
			stake.GET("/stakings", instance.hub.nta.GetStakeStakings)
			stake.GET("/:staker_address/stat", instance.hub.nta.GetStakingStat)
			stake.GET("/transactions", instance.hub.nta.GetStakeTransactions)
			stake.GET("/transactions/:transaction_hash", instance.hub.nta.GetStakeTransaction)
		}

		token := nta.Group("/token")
		{
			token.GET("/supply", instance.hub.nta.GetTokenSupply)
			token.GET("/tvl", instance.hub.nta.GetTvl)
		}

		dsl := nta.Group("/dsl")
		{
			dsl.GET("/total_requests", instance.hub.nta.GetDslTotalRequests)
		}
	}

	dsl := instance.httpServer.Group("")
	{
		rss := dsl.Group("/rss")
		{
			rss.GET("/*", instance.hub.dsl.GetRSSHub)
		}

		agentdata := dsl.Group("/agentdata")
		{
			agentdata.GET("/*", instance.hub.dsl.GetAI)
		}

		decentralized := dsl.Group("/decentralized")
		{
			decentralized.GET("/tx/:id", instance.hub.dsl.GetDecentralizedActivity)
			decentralized.GET("/:account", instance.hub.dsl.GetDecentralizedAccountActivities)
			decentralized.GET("/network/:network", instance.hub.dsl.GetDecentralizedNetworkActivities)
			decentralized.GET("/platform/:platform", instance.hub.dsl.GetDecentralizedPlatformActivities)
			decentralized.POST("/accounts", instance.hub.dsl.BatchGetDecentralizedAccountsActivities)
		}

		federated := dsl.Group("/federated")
		{
			federated.GET("/tx/:id", instance.hub.dsl.GetFederatedActivity)
			federated.GET("/:account", instance.hub.dsl.GetFederatedAccountActivities)
			federated.GET("/network/:network", instance.hub.dsl.GetFederatedNetworkActivities)
			federated.GET("/platform/:platform", instance.hub.dsl.GetFederatedPlatformActivities)
			federated.POST("/accounts", instance.hub.dsl.BatchGetFederatedAccountsActivities)
		}
	}

	return &instance, nil
}

package hub

import (
	"context"
	"fmt"
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rss3-network/global-indexer/common/geolite2"
	"github.com/rss3-network/global-indexer/common/httputil"
	"github.com/rss3-network/global-indexer/docs"
	"github.com/rss3-network/global-indexer/internal/client/ethereum"
	"github.com/rss3-network/global-indexer/internal/database"
	"github.com/rss3-network/global-indexer/internal/nameresolver"
	"github.com/rss3-network/global-indexer/internal/service"
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

func NewServer(databaseClient database.Client, redisClient *redis.Client, geoLite2 *geolite2.Client, ethereumMultiChainClient *ethereum.MultiChainClient, nameService *nameresolver.NameResolver, httpClient httputil.Client) (service.Server, error) {
	hub, err := NewHub(context.Background(), databaseClient, redisClient, ethereumMultiChainClient, geoLite2, nameService, httpClient)
	if err != nil {
		return nil, fmt.Errorf("new hub: %w", err)
	}

	instance := Server{
		httpServer: echo.New(),
		hub:        hub,
	}

	instance.httpServer.HideBanner = true
	instance.httpServer.HidePort = true
	instance.httpServer.Validator = defaultValidator
	instance.httpServer.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	{
		instance.httpServer.FileFS("/docs/openapi.yaml", "openapi.yaml", docs.EmbedFS)
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
		bridge := nta.Group("/bridge")
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
			networks.GET("", instance.hub.nta.GetNetworks)
			networks.GET("/:network_name/list_workers", instance.hub.nta.GetWorkersByNetwork)
			networks.GET("/:network_name/workers/:worker_name", instance.hub.nta.GetWorkerDetail)
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

		{
			snapshots.GET("/nodes/count", instance.hub.nta.GetNodeCountSnapshots)
			snapshots.GET("/nodes/operation/profit", instance.hub.nta.GetNodeOperationProfitSnapshots)
			snapshots.GET("/stakers/count", instance.hub.nta.GetStakerCountSnapshots)
			snapshots.GET("/stakers/profit", instance.hub.nta.GetStakerProfitSnapshots)
			snapshots.GET("/epochs/apy", instance.hub.nta.GetEpochsAPYSnapshots)

			snapshots.POST("/nodes/min_tokens_to_stake", instance.hub.nta.BatchGetNodeMinTokensToStakeSnapshots)
		}

		stake := nta.Group("/stake")
		{
			stake := nta.Group("/stake")
			stake.GET("/:staker_address/profit", instance.hub.nta.GetStakeOwnerProfit)
			stake.GET("/stakings", instance.hub.nta.GetStakeStakings)
			stake.GET("/transactions", instance.hub.nta.GetStakeTransactions)
			stake.GET("/transactions/:transaction_hash", instance.hub.nta.GetStakeTransaction)
		}
	}

	dsl := instance.httpServer.Group("")
	{
		dsl.GET("/rss/*", instance.hub.dsl.GetRSSHub)
		dsl.GET("/decentralized/tx/:id", instance.hub.dsl.GetActivity)
		dsl.GET("/decentralized/:account", instance.hub.dsl.GetAccountActivities)
	}

	return &instance, nil
}

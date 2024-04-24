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

	// register router
	instance.httpServer.POST("/nodes/register", instance.hub.nta.RegisterNode)
	instance.httpServer.POST("/nodes/heartbeat", instance.hub.nta.NodeHeartbeat)

	nta := instance.httpServer.Group("/nta")
	{
		nta.GET("/networks", instance.hub.nta.GetNetworks)
		nta.GET("/networks/:network/workers", instance.hub.nta.GetWorkersByNetwork)

		nta.GET("/nodes", instance.hub.nta.GetNodes)
		nta.GET("/nodes/:id", instance.hub.nta.GetNode)
		nta.GET("/nodes/:id/events", instance.hub.nta.GetNodeEvents)
		nta.GET("/nodes/:id/challenge", instance.hub.nta.GetNodeChallenge)
		nta.POST("/nodes/:id/hideTaxRate", instance.hub.nta.PostNodeHideTaxRate)
		nta.GET("/nodes/:id/avatar.svg", instance.hub.nta.GetNodeAvatar)
		nta.GET("/operation/:operator/profits", instance.hub.nta.GetOperatorProfit)

		nta.GET("/bridge/transactions", instance.hub.nta.GetBridgeTransactions)
		nta.GET("/bridge/transactions/:id", instance.hub.nta.GetBridgeTransaction)

		nta.GET("/stake/transactions", instance.hub.nta.GetStakeTransactions)
		nta.GET("/stake/transactions/:id", instance.hub.nta.GetStakeTransaction)
		nta.GET("/stake/stakings", instance.hub.nta.GetStakeStakings)
		nta.GET("/stake/:owner/profits", instance.hub.nta.GetStakeOwnerProfit)

		nta.GET("/chips", instance.hub.nta.GetStakeChips)
		nta.GET("/chips/:id", instance.hub.nta.GetStakeChip)
		nta.GET("/chips/:id/image.svg", instance.hub.nta.GetStakeChipImage)

		nta.GET("/epochs", instance.hub.nta.GetEpochs)
		nta.GET("/epochs/:id", instance.hub.nta.GetEpoch)
		nta.GET("/epochs/distributions/:transaction", instance.hub.nta.GetEpochDistribution)
		nta.GET("/epochs/:node/rewards", instance.hub.nta.GetEpochNodeRewards)

		nta.GET("/snapshots/nodes/count", instance.hub.nta.GetNodeCountSnapshots)
		nta.POST("/snapshots/nodes/minTokensToStake", instance.hub.nta.BatchGetNodeMinTokensToStakeSnapshots)
		nta.GET("/snapshots/stakers/count", instance.hub.nta.GetStakerCountSnapshots)
		nta.GET("/snapshots/stakers/profits", instance.hub.nta.GetStakerProfitsSnapshots)
		nta.GET("/snapshots/operators/profits", instance.hub.nta.GetOperatorProfitsSnapshots)
	}

	dsl := instance.httpServer.Group("")
	{
		dsl.GET("/rss/*", instance.hub.dsl.GetRSSHub)
		dsl.GET("/decentralized/tx/:id", instance.hub.dsl.GetActivity)
		dsl.GET("/decentralized/:account", instance.hub.dsl.GetAccountActivities)
	}

	return &instance, nil
}

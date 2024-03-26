package hub

import (
	"context"
	"fmt"
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/geolite2"
	"github.com/naturalselectionlabs/rss3-global-indexer/docs"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/client/ethereum"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/nameresolver"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	"github.com/redis/go-redis/v9"
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

func NewServer(databaseClient database.Client, redisClient *redis.Client, geoLite2 *geolite2.Client, ethereumMultiChainClient *ethereum.MultiChainClient, nameService *nameresolver.NameResolver) (service.Server, error) {
	hub, err := NewHub(databaseClient, redisClient, ethereumMultiChainClient, geoLite2, nameService)
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
	instance.httpServer.GET("/nodes", instance.hub.GetNodesHandler)
	instance.httpServer.GET("/nodes/:id", instance.hub.GetNodeHandler)
	instance.httpServer.GET("/nodes/:id/events", instance.hub.GetNodeEventsHandler)
	instance.httpServer.GET("/nodes/:id/challenge", instance.hub.GetNodeChallengeHandler)
	instance.httpServer.POST("/nodes/:id/hideTaxRate", instance.hub.PostNodeHideTaxRateHandler)
	instance.httpServer.GET("/nodes/:id/avatar.svg", instance.hub.GetNodeAvatarHandler)
	instance.httpServer.POST("/nodes/register", instance.hub.RegisterNodeHandler)
	instance.httpServer.POST("/nodes/heartbeat", instance.hub.NodeHeartbeatHandler)

	instance.httpServer.GET("/bridge/transactions", instance.hub.GetBridgeTransactions)
	instance.httpServer.GET("/bridge/transactions/:id", instance.hub.GetBridgeTransaction)

	instance.httpServer.GET("/stake/transactions", instance.hub.GetStakeTransactions)
	instance.httpServer.GET("/stake/transactions/:id", instance.hub.GetStakeTransaction)
	instance.httpServer.GET("/stake/stakings", instance.hub.GetStakeStakings)
	instance.httpServer.GET("/stake/:owner/profits", instance.hub.GetStakeOwnerProfit)

	instance.httpServer.GET("/chips", instance.hub.GetStakeChips)
	instance.httpServer.GET("/chips/:id", instance.hub.GetStakeChip)
	instance.httpServer.GET("/chips/:id/image.svg", instance.hub.GetStakeChipImage)

	instance.httpServer.GET("/epochs", instance.hub.GetEpochsHandler)
	instance.httpServer.GET("/epochs/:id", instance.hub.GetEpochHandler)
	instance.httpServer.GET("/epochs/distributions/:transaction", instance.hub.GetEpochDistributionHandler)
	instance.httpServer.GET("/epochs/:node/rewards", instance.hub.GetEpochNodeRewardsHandler)

	instance.httpServer.GET("/snapshots/nodes/count", instance.hub.GetNodeCountSnapshots)
	instance.httpServer.POST("/snapshots/nodes/minTokensToStake", instance.hub.BatchGetNodeMinTokensToStakeSnapshots)
	instance.httpServer.GET("/snapshots/stakers/count", instance.hub.GetStakerCountSnapshots)
	instance.httpServer.GET("/snapshots/stakers/profits", instance.hub.GetStakerProfitsSnapshots)

	instance.httpServer.GET("/rss/*", instance.hub.GetRSSHubHandler)
	instance.httpServer.GET("/decentralized/tx/:id", instance.hub.GetActivityHandler)
	instance.httpServer.GET("/decentralized/:account", instance.hub.GetAccountActivitiesHandler)

	return &instance, nil
}

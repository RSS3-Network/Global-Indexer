package hub

import (
	"context"
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/naturalselectionlabs/rss3-global-indexer/common/geolite2"
	"github.com/naturalselectionlabs/rss3-global-indexer/docs"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/nameresolver"
	"github.com/redis/go-redis/v9"
)

const (
	DefaultHost = "0.0.0.0"
	DefaultPort = "80"
)

type Server struct {
	httpServer *echo.Echo
	hub        *Hub
}

func (s *Server) Run(_ context.Context) error {
	address := net.JoinHostPort(DefaultHost, DefaultPort)

	return s.httpServer.Start(address)
}

func NewServer(ctx context.Context, databaseClient database.Client, ethereumClient *ethclient.Client, redisClient *redis.Client, geoLite2 *geolite2.Client, nameService *nameresolver.NameResolver) (*Server, error) {
	hub, err := NewHub(ctx, databaseClient, ethereumClient, redisClient, geoLite2, nameService)
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
	instance.httpServer.GET("/nodes", instance.hub.GetNodes)
	instance.httpServer.GET("/nodes/:id", instance.hub.GetNode)
	instance.httpServer.GET("/nodes/:id/events", instance.hub.GetNodeEvents)
	instance.httpServer.GET("/nodes/:id/challenge", instance.hub.GetNodeChallenge)
	instance.httpServer.POST("/nodes/:id/hideTaxRate", instance.hub.PostNodeHideTaxRate)
	instance.httpServer.GET("/nodes/:id/avatar.svg", instance.hub.GetNodeAvatar)
	instance.httpServer.POST("/nodes/register", instance.hub.RegisterNode)
	instance.httpServer.POST("/nodes/heartbeat", instance.hub.NodeHeartbeat)
	instance.httpServer.GET("/operation/:operator/profits", instance.hub.GetOperatorProfit)

	instance.httpServer.GET("/bridge/transactions", instance.hub.GetBridgeTransactions)
	instance.httpServer.GET("/bridge/transactions/:id", instance.hub.GetBridgeTransaction)

	instance.httpServer.GET("/stake/transactions", instance.hub.GetStakeTransactions)
	instance.httpServer.GET("/stake/transactions/:id", instance.hub.GetStakeTransaction)
	instance.httpServer.GET("/stake/stakings", instance.hub.GetStakeStakings)
	instance.httpServer.GET("/stake/:owner/profits", instance.hub.GetStakeOwnerProfit)

	instance.httpServer.GET("/chips", instance.hub.GetStakeChips)
	instance.httpServer.GET("/chips/:id", instance.hub.GetStakeChip)
	instance.httpServer.GET("/chips/:id/image.svg", instance.hub.GetStakeChipImage)

	instance.httpServer.GET("/epochs", instance.hub.GetEpochs)
	instance.httpServer.GET("/epochs/:id", instance.hub.GetEpoch)
	instance.httpServer.GET("/epochs/distributions/:transaction", instance.hub.GetEpochDistribution)
	instance.httpServer.GET("/epochs/:node/rewards", instance.hub.GetEpochNodeRewards)

	instance.httpServer.GET("/snapshots/nodes/count", instance.hub.GetNodeCountSnapshots)
	instance.httpServer.POST("/snapshots/nodes/minTokensToStake", instance.hub.BatchGetNodeMinTokensToStakeSnapshots)
	instance.httpServer.GET("/snapshots/stakers/count", instance.hub.GetStakerCountSnapshots)
	instance.httpServer.GET("/snapshots/stakers/profits", instance.hub.GetStakerProfitsSnapshots)
	instance.httpServer.GET("/snapshots/operators/profits", instance.hub.GetOperatorProfitsSnapshots)

	instance.httpServer.GET("/rss/*", instance.hub.GetRSSHub)
	instance.httpServer.GET("/decentralized/tx/:id", instance.hub.GetActivity)
	instance.httpServer.GET("/decentralized/:account", instance.hub.GetAccountActivities)

	return &instance, nil
}

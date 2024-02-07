package hub

import (
	"context"
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
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

func NewServer(ctx context.Context, databaseClient database.Client, ethereumClient *ethclient.Client, redisClient *redis.Client) (*Server, error) {
	hub, err := NewHub(ctx, databaseClient, ethereumClient, redisClient)
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

	// register router
	instance.httpServer.GET("/nodes", instance.hub.GetNodesHandler)
	instance.httpServer.GET("/nodes/:id", instance.hub.GetNodeHandler)
	instance.httpServer.GET("/nodes/:id/challenge", instance.hub.GetNodeChallengeHandler)
	instance.httpServer.POST("/nodes/register", instance.hub.RegisterNodeHandler)
	instance.httpServer.POST("/nodes/heartbeat", instance.hub.NodeHeartbeatHandler)

	instance.httpServer.GET("/bridge/transactions", instance.hub.GetBridgeTransactions)
	instance.httpServer.GET("/bridge/transactions/:id", instance.hub.GetBridgeTransaction)

	instance.httpServer.GET("/stake/transactions", instance.hub.GetStakeTransactions)
	instance.httpServer.GET("/stake/transactions/:id", instance.hub.GetStakeTransaction)

	instance.httpServer.GET("/epochs", instance.hub.GetEpochsHandler)
	instance.httpServer.GET("/epochs/:id", instance.hub.GetEpochHandler)
	instance.httpServer.GET("/epochs/:node/rewards", instance.hub.GetEpochNodeRewardsHandler)

	instance.httpServer.GET("/nodes/:node/chips", instance.hub.GetStakeNodeChips)
	instance.httpServer.GET("/wallets/:wallet/chips", instance.hub.GetStakeWalletChips)

	instance.httpServer.GET("/snapshot/nodes", instance.hub.GetNodeSnapshots)
	instance.httpServer.GET("/snapshot/stakers", instance.hub.GetStakeSnapshots)

	instance.httpServer.GET("/rss/*", instance.hub.GetRSSHubHandler)
	instance.httpServer.GET("/decentralized/tx/:id", instance.hub.GetActivityHandler)
	instance.httpServer.GET("/decentralized/:account", instance.hub.GetAccountActivitiesHandler)

	return &instance, nil
}

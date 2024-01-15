package hub

import (
	"context"
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/naturalselectionlabs/global-indexer/internal/database"
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

func NewServer(ctx context.Context, databaseClient database.Client, ethereumClient *ethclient.Client) (*Server, error) {
	hub, err := NewHub(ctx, databaseClient, ethereumClient)
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
	instance.httpServer.POST("/nodes/register", instance.hub.RegisterNodeHandler)
	instance.httpServer.POST("/nodes/heartbeat", instance.hub.NodeHeartbeatHandler)

	instance.httpServer.GET("/staking", instance.hub.GetStakingHandler)
	instance.httpServer.GET("/bridging", instance.hub.GetBridgingHandler)

	instance.httpServer.GET(PathGetRSSHub, instance.hub.GetRSSHubHandler)
	instance.httpServer.GET(PathGetDecentralizedTx, instance.hub.GetActivityHandler)
	instance.httpServer.GET(PathGetDecentralizedActivities, instance.hub.GetAccountActivitiesHandler)

	return &instance, nil
}

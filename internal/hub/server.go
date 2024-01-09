package hub

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/naturalselectionlabs/global-indexer/internal/database"
	"net"
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

func NewServer(ctx context.Context, databaseClient database.Client) (*Server, error) {
	instance := Server{
		httpServer: echo.New(),
		hub:        NewHub(ctx, databaseClient),
	}

	instance.httpServer.HideBanner = true
	instance.httpServer.HidePort = true
	instance.httpServer.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	// register router
	instance.httpServer.GET("/nodes", instance.hub.GetNodesHandler)
	instance.httpServer.GET("/nodes/:id", instance.hub.GetNodeHandler)
	instance.httpServer.POST("/nodes/register", instance.hub.RegisterNodeHandler)
	instance.httpServer.GET("/staking", instance.hub.GetStakingHandler)
	instance.httpServer.GET("/bridging", instance.hub.GetBridgingHandler)

	return &instance, nil
}

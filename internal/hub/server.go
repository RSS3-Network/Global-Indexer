package hub

import (
	"context"
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func NewServer(ctx context.Context) (*Server, error) {
	instance := Server{
		httpServer: echo.New(),
		hub:        NewHub(ctx),
	}

	instance.httpServer.HideBanner = true
	instance.httpServer.HidePort = true
	instance.httpServer.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	// register router
	instance.httpServer.GET("/nodes", instance.hub.GetNodesHandler)
	instance.httpServer.GET("/staking", instance.hub.GetStakingHandler)
	instance.httpServer.GET("/bridging", instance.hub.GetBridgingHandler)

	return &instance, nil
}

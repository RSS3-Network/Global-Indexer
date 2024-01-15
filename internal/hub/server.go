package hub

import (
	"context"
	"net"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/naturalselectionlabs/global-indexer/internal/database"
	"github.com/naturalselectionlabs/global-indexer/provider/node"
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

func NewServer(ctx context.Context, databaseClient database.Client, pathBuilder node.Builder) (*Server, error) {
	instance := Server{
		httpServer: echo.New(),
		hub:        NewHub(ctx, databaseClient, pathBuilder),
	}

	instance.httpServer.HideBanner = true
	instance.httpServer.HidePort = true
	instance.httpServer.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	instance.httpServer.Validator = &Validator{
		validate: validator.New(),
	}

	// register router
	instance.httpServer.GET(PathGetNodes, instance.hub.GetNodesHandler)
	instance.httpServer.GET(PathGetNode, instance.hub.GetNodeHandler)
	instance.httpServer.POST(PathNodesRegister, instance.hub.RegisterNodeHandler)
	instance.httpServer.GET(PathStaking, instance.hub.GetStakingHandler)
	instance.httpServer.GET(PathBridging, instance.hub.GetBridgingHandler)

	instance.httpServer.GET(PathGetRSSHub, instance.hub.GetRSSHubHandler)
	instance.httpServer.GET(PathGetDecentralizedTx, instance.hub.GetActivityHandler)
	instance.httpServer.GET(PathGetDecentralizedActivities, instance.hub.GetAccountActivitiesHandler)

	return &instance, nil
}

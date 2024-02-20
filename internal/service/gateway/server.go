package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
	apisixHTTPAPI "github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/apisix/httpapi"
	apisixKafkaLog "github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/apisix/kafkalog"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/gen/oapi"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/handlers"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/jwt"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/middlewares"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/siwe"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/swagger"
	"github.com/redis/go-redis/v9"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	config         config.GatewayConfig
	redis          *redis.Client
	databaseClient database.Client
}

func (s *Server) Run(ctx context.Context) error {
	errorPool := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	// Run echo server.
	errorPool.Go(func(ctx context.Context) error {

		// Initialize APISIX configurations
		apisixAPIService, err := apisixHTTPAPI.New(
			s.config.APISix.Admin.Endpoint,
			s.config.APISix.Admin.Key,
		)
		if err != nil {
			return err
		}

		// Prepare JWT
		jwtClient, err := jwt.New(s.config.API.JWTKey)
		if err != nil {
			return err
		}

		// Prepare SIWE
		siweClient, err := siwe.New(s.config.API.SIWEDomain, s.redis)
		if err != nil {
			return err
		}

		// Prepare echo
		e := echo.New()
		echoHandler, err := handlers.NewApp(
			apisixAPIService,
			s.redis,
			s.databaseClient.Raw(),
			jwtClient,
			siweClient,
		)
		if err != nil {
			return err
		}

		// Configure middlewares
		s.configureMiddlewares(e, echoHandler, jwtClient)

		// Connect to kafka for access logs
		kafkaService, err := apisixKafkaLog.New(
			strings.Split(s.config.APISix.Kafka.Brokers, ","),
			s.config.APISix.Kafka.Topic,
		)
		if err != nil {
			// Failed to Initialize kafka consumer
			return err
		}

		err = kafkaService.Start(echoHandler.ProcessAccessLog)
		if err != nil {
			// Failed to start kafka consumer
			return err
		}

		// Start echo API server
		return e.Start(fmt.Sprintf("%s:%d", s.config.API.Listen.Host, s.config.API.Listen.Port))
	})

	errorChan := make(chan error)
	go func() { errorChan <- errorPool.Wait() }()

	select {
	case err := <-errorChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func New(databaseClient database.Client, redis *redis.Client, config config.GatewayConfig) (service.Server, error) {
	instance := Server{
		config:         config,
		redis:          redis,
		databaseClient: databaseClient,
	}

	return &instance, nil
}

func (s *Server) configureMiddlewares(e *echo.Echo, app *handlers.App, jwtClient *jwt.JWT) {
	oapi.RegisterHandlers(e, app)

	// Add api docs
	if os.Getenv(config.Environment) == config.EnvironmentDevelopment {
		swg, err := oapi.GetSwagger()
		if err != nil {
			// Log but ignore
			zap.L().Error("get swagger doc", zap.Error(err))
		}
		swgJSON, err := swg.MarshalJSON()
		if err != nil {
			// Log but ignore
			zap.L().Error("marshal swagger doc", zap.Error(err))
		}
		e.Pre(swagger.SwaggerDoc("/", swgJSON))
	}

	// Check user authentication
	e.Use(middlewares.UserAuthenticationMiddleware(s.databaseClient.Raw(), jwtClient))

	e.HTTPErrorHandler = customHTTPErrorHandler
}

func customHTTPErrorHandler(err error, c echo.Context) {
	// ignore user cancelled error
	switch {
	case errors.Is(err, context.Canceled):
		_ = c.NoContent(0)
	case errors.Is(err, gorm.ErrRecordNotFound):
		_ = JSONResponseMsg(c, err.Error(), http.StatusNotFound)
	case errors.Is(err, gorm.ErrInvalidField):
		_ = JSONResponseMsg(c, err.Error(), http.StatusBadRequest)
	case errors.Is(err, errors.New(http.StatusText(http.StatusUnauthorized))) && err.Error() == http.StatusText(http.StatusUnauthorized):
		_ = JSONResponseMsg(c, "Your credentials have expired.", http.StatusUnauthorized)
	case strings.Contains(err.Error(), "Path was not found"):
		_ = JSONResponseMsg(c, err.Error(), http.StatusNotFound)
	}

	c.Echo().DefaultHTTPErrorHandler(err, c)
}

func JSONResponseMsg(ctx echo.Context, msg string, code int) error {
	return ctx.JSON(code, map[string]interface{}{
		"msg":    msg,
		"errors": struct{}{},
	})
}

package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	apisixHTTPAPI "github.com/naturalselectionlabs/rss3-global-indexer/internal/apisix/httpapi"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/config"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/database"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service"
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
)

type Server struct {
	config               config.Gateway
	redis                *redis.Client
	databaseClient       database.Client
	apisixHTTPAPIService *apisixHTTPAPI.Service
}

func (s *Server) Run(ctx context.Context) error {
	errorPool := pool.New().WithContext(ctx).WithCancelOnError().WithFirstError()

	// Run echo server.
	errorPool.Go(func(ctx context.Context) error {
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
			s.apisixHTTPAPIService,
			s.redis,
			s.databaseClient.Raw(),
			jwtClient,
			siweClient,
		)

		if err != nil {
			return err
		}

		// Configure middlewares
		configureMiddlewares(e, echoHandler, jwtClient, s.databaseClient.Raw(), s.apisixHTTPAPIService)

		// Connect to kafka for access logs
		kafkaService, err := apisixKafkaLog.New(
			strings.Split(s.config.APISixKafka.Brokers, ","),
			s.config.APISixKafka.Topic,
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

func New(databaseClient database.Client, redis *redis.Client, apisixHTTPAPIService *apisixHTTPAPI.Service, config config.Gateway) (service.Server, error) {
	instance := Server{
		config:               config,
		redis:                redis,
		databaseClient:       databaseClient,
		apisixHTTPAPIService: apisixHTTPAPIService,
	}

	return &instance, nil
}

func configureMiddlewares(e *echo.Echo, app *handlers.App, jwtClient *jwt.JWT, databaseClient *gorm.DB, apiSixAPIService *apisixHTTPAPI.Service) {
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
	e.Use(middlewares.UserAuthenticationMiddleware(databaseClient, apiSixAPIService, jwtClient))

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

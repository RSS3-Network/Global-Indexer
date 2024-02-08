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
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/middlewares"
	"github.com/redis/go-redis/v9"
	"github.com/sourcegraph/conc/pool"
	"gorm.io/gorm"
	"log"
	"net/http"
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

		// Prepare echo
		e := echo.New()
		echoHandler := handlers.NewApp(
			s.config.API.SIWEDomain,
			apisixAPIService,
			s.redis,
			s.databaseClient.Raw(),
		)

		// Configure middlewares
		configureMiddlewares(e, echoHandler)

		// Connect to kafka for access logs
		kafkaService, err := apisixKafkaLog.New(
			strings.Split(s.config.APISix.Kafka.Brokers, ","),
			s.config.APISix.Kafka.Topic,
		)
		if err != nil {
			// Failed to Initialize kafka consumer
			log.Panic(err)
		}

		err = kafkaService.Start(handlers.ProcessAccessLog)
		if err != nil {
			// Failed to start kafka consumer
			log.Panic(err)
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

func configureMiddlewares(e *echo.Echo, impls *handlers.App) {
	oapi.RegisterHandlers(e, impls)

	// Check user authentication
	e.Use(middlewares.UserAuthenticationMiddleware)

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

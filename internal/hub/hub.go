package hub

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/naturalselectionlabs/global-indexer/internal/database"
	"github.com/naturalselectionlabs/global-indexer/provider/node"
)

type Hub struct {
	databaseClient database.Client
	pathBuilder    node.Builder
	httpClient     *http.Client
}

var _ echo.Validator = &Validator{}

type Validator struct {
	validate *validator.Validate
	once     sync.Once
}

func (v *Validator) Validate(i interface{}) error {
	if v.validate == nil {
		v.once.Do(func() {
			v.validate = validator.New()
		})
	}

	return v.validate.Struct(i)
}

func NewHub(_ context.Context, databaseClient database.Client, pathBuilder node.Builder) *Hub {
	return &Hub{
		databaseClient: databaseClient,
		pathBuilder:    pathBuilder,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

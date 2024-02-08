package httpapi

import (
	"fmt"
)

type HTTPAPIService struct {
	Config *cfg
}

func New(adminEndpoint string, adminKey string) (*HTTPAPIService, error) {

	if adminEndpoint == "" {
		return nil, fmt.Errorf("missing admin endpoint")
	}

	if adminKey == "" {
		return nil, fmt.Errorf("missing admin key")
	}

	return &HTTPAPIService{
		Config: &cfg{
			APISixAdminEndpoint: adminEndpoint,
			APISixAdminKey:      adminKey,
		},
	}, nil
}

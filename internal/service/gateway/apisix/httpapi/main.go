package httpapi

import (
	"fmt"
)

type Service struct {
	Config *cfg
}

func New(adminEndpoint string, adminKey string) (*Service, error) {
	if adminEndpoint == "" {
		return nil, fmt.Errorf("missing admin endpoint")
	}

	if adminKey == "" {
		return nil, fmt.Errorf("missing admin key")
	}

	return &Service{
		Config: &cfg{
			APISixAdminEndpoint: adminEndpoint,
			APISixAdminKey:      adminKey,
		},
	}, nil
}

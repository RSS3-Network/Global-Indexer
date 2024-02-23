package httpapi

import (
	"fmt"
)

type Client struct {
	Config *cfg
}

func New(adminEndpoint string, adminKey string) (*Client, error) {
	if adminEndpoint == "" {
		return nil, fmt.Errorf("missing admin endpoint")
	}

	if adminKey == "" {
		return nil, fmt.Errorf("missing admin key")
	}

	return &Client{
		Config: &cfg{
			APISixAdminEndpoint: adminEndpoint,
			APISixAdminKey:      adminKey,
		},
	}, nil
}

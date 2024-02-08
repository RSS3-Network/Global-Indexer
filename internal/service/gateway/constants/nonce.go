package constants

import "time"

const (
	NONCE_LIFE       = 5 * time.Minute
	NONCE_KEY_PREFIX = "apigateway:nonce"
)

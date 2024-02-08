package siwe

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/naturalselectionlabs/api-gateway/app"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/constants"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/utils"
	"github.com/naturalselectionlabs/api-gateway/app/oapi/variables"
	"github.com/spruceid/siwe-go"
)

func ValidateSIWESignature(ctx context.Context, rawMessage, signature string) (string, int, error) {
	// Parse a SIWE Message
	message, err := siwe.ParseMessage(rawMessage)
	if err != nil {
		return "", 0, err
	}

	// Verify nonce
	nonce := message.GetNonce()
	if err = ConsumeNonce(ctx, nonce); err != nil {
		return "", 0, err
	}

	// Verifying and Authenticating a SIWE Message
	_, err = message.Verify(signature, &variables.SIWEDomain, nil, nil)
	if err != nil {
		return "", 0, err
	}

	return message.GetAddress().Hex(), message.GetChainID(), nil
}

func buildNonceKey(nonce string) string {
	return fmt.Sprintf("%s:%s", constants.NONCE_KEY_PREFIX, nonce)
}

func GetNonce(ctx context.Context) (string, error) {
	// Generate nonce
	nonce := utils.RandString(16)

	// Save into redis

	if err := app.RedisExt.Client(ctx).Set(
		buildNonceKey(nonce),
		"",
		constants.NONCE_LIFE,
	).Err(); err != nil {
		return "", err
	}

	return nonce, nil
}

func ConsumeNonce(ctx context.Context, nonce string) error {

	// Check if nonce exist
	_, err := app.RedisExt.Client(ctx).Get(buildNonceKey(nonce)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// No such key
			return errors.New("no such nonce")
		} else {
			return err
		}
	}

	app.RedisExt.Client(ctx).Del(buildNonceKey(nonce))

	return nil
}

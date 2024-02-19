package siwe

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/constants"
	"github.com/naturalselectionlabs/rss3-global-indexer/internal/service/gateway/utils"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/spruceid/siwe-go"
)

type SIWE struct {
	domain      string
	redisClient *redis.Client
}

func New(domain string, redisClient *redis.Client) (*SIWE, error) {
	return &SIWE{
		domain:      domain,
		redisClient: redisClient,
	}, nil
}

func (s *SIWE) Domain() string {
	return s.domain
}

func (s *SIWE) ValidateSIWESignature(ctx context.Context, rawMessage, signature string) (*common.Address, int, error) {
	// Parse a SIWE Message
	message, err := siwe.ParseMessage(rawMessage)
	if err != nil {
		return nil, 0, err
	}

	// Verify nonce
	nonce := message.GetNonce()
	if err = s.ConsumeNonce(ctx, nonce); err != nil {
		return nil, 0, err
	}

	// Verifying and Authenticating a SIWE Message
	_, err = message.Verify(signature, &s.domain, nil, nil)
	if err != nil {
		return nil, 0, err
	}

	return lo.ToPtr(message.GetAddress()), message.GetChainID(), nil
}

func (s *SIWE) buildNonceKey(nonce string) string {
	return fmt.Sprintf("%s:%s", constants.NONCE_KEY_PREFIX, nonce)
}

func (s *SIWE) GetNonce(ctx context.Context) (string, error) {
	// Generate nonce
	nonce := utils.RandString(16)

	// Save into redis

	if err := s.redisClient.Set(
		ctx,
		s.buildNonceKey(nonce),
		"",
		constants.NONCE_LIFE,
	).Err(); err != nil {
		return "", err
	}

	return nonce, nil
}

func (s *SIWE) ConsumeNonce(ctx context.Context, nonce string) error {

	// Check if nonce exist
	_, err := s.redisClient.Get(ctx, s.buildNonceKey(nonce)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// No such key
			return errors.New("no such nonce")
		} else {
			return err
		}
	}

	s.redisClient.Del(ctx, s.buildNonceKey(nonce))

	return nil
}

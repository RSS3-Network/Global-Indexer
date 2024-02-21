// AUTOGENERATED FILE (not)
// Everything here has its meaning, don't let golang-ci ruin them

package httpapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Consumer(t *testing.T) {

	ctx := context.Background()

	// Init configs
	s, _ := New("http://localhost:9180", "edd1c9f034335f136f87ad84b625c8f1")

	userAddr := "0x9E85cb8f7606dA552cdEb42130c315F3fed26625"

	// Create a new consumer group
	err := s.NewConsumerGroup(ctx, userAddr)
	assert.NoError(t, err)

	// Check consumer group
	cgInfo, err := s.CheckConsumerGroup(ctx, userAddr)
	assert.NoError(t, err)
	assert.Equal(t, *cgInfo.Value.ID, userAddr)

	// TODO: Further test before token created

	// Create 2 new consumers
	key1 := uuid.New().String()
	keyID1 := uint64(114514)

	key2 := uuid.New().String()
	keyID2 := uint64(1919810)

	err = s.NewConsumer(ctx, keyID1, key1, userAddr)
	assert.NoError(t, err)

	err = s.NewConsumer(ctx, keyID2, key2, userAddr)
	assert.NoError(t, err)

	// TODO: Further test after token created

	// Get consumer info
	cInfo, err := s.CheckConsumer(ctx, keyID1)
	assert.NoError(t, err)
	assert.Equal(t, cInfo.Value.GroupID, userAddr)
	assert.Equal(t, cInfo.Value.Username, s.consumerUsername(keyID1))
	cInfo, err = s.CheckConsumer(ctx, keyID2)
	assert.NoError(t, err)
	assert.Equal(t, cInfo.Value.GroupID, userAddr)
	assert.Equal(t, cInfo.Value.Username, s.consumerUsername(keyID2))

	// Pause consumer group
	err = s.PauseConsumerGroup(ctx, userAddr)
	assert.NoError(t, err)

	// TODO: Further test after pause

	// Resume paused consumer group
	err = s.ResumeConsumerGroup(ctx, userAddr)
	assert.NoError(t, err)

	// TODO: Further test after resume

	// Delete one consumer
	err = s.DeleteConsumer(ctx, keyID1)
	assert.NoError(t, err)

	// TODO: Further test after delete

	// Delete consumer group
	err = s.DeleteConsumerGroup(ctx, userAddr)
	assert.Error(t, err, fmt.Sprintf("can not delete this consumer group, consumer [%d] is still using it now", keyID2))

	// Delete another consumer
	err = s.DeleteConsumer(ctx, keyID2)
	assert.NoError(t, err)

	// TODO: Further test after delete

	// Delete consumer group
	err = s.DeleteConsumerGroup(ctx, userAddr)
	assert.NoError(t, err)

	cgInfo, err = s.CheckConsumerGroup(ctx, userAddr)
	assert.Nil(t, cgInfo)
	assert.Error(t, err, "Key not found")

	// Create consumer with non-exist group should create group automatically

	key3 := uuid.New().String()
	keyID3 := uint64(25565)

	err = s.NewConsumer(ctx, keyID3, key3, userAddr)
	assert.NoError(t, err)

	cgInfo, err = s.CheckConsumerGroup(ctx, userAddr)
	assert.NoError(t, err)
	assert.Equal(t, *cgInfo.Value.ID, userAddr)

	err = s.DeleteConsumer(ctx, keyID3)
	assert.NoError(t, err)

	err = s.DeleteConsumerGroup(ctx, userAddr)
	assert.NoError(t, err)
}

func Test_UsernameChangeHelper(t *testing.T) {

	// Init configs
	s, _ := New("http://localhost:9180", "edd1c9f034335f136f87ad84b625c8f1")

	keyID := uint64(114514)
	res, err := s.RecoverKeyIDFromConsumerUsername(s.consumerUsername(keyID))
	assert.Equal(t, keyID, res)
	assert.NoError(t, err)
}

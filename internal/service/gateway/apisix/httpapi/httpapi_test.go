package httpapi

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Consumer(t *testing.T) {

	// Init configs
	s, _ := New("http://localhost:9180", "edd1c9f034335f136f87ad84b625c8f1")

	userAddr := "0x9E85cb8f7606dA552cdEb42130c315F3fed26625"

	// Create a new consumer group
	err := s.NewConsumerGroup(userAddr)
	assert.NoError(t, err)

	// Check consumer group
	cgInfo, err := s.CheckConsumerGroup(userAddr)
	assert.NoError(t, err)
	assert.Equal(t, *cgInfo.Value.ID, userAddr)

	// TODO: Further test before token created

	// Create 2 new consumers
	key1 := uuid.New().String()
	keyId1 := uint64(114514)

	key2 := uuid.New().String()
	keyId2 := uint64(1919810)

	err = s.NewConsumer(keyId1, key1, userAddr)
	assert.NoError(t, err)

	err = s.NewConsumer(keyId2, key2, userAddr)
	assert.NoError(t, err)

	// TODO: Further test after token created

	// Get consumer info
	cInfo, err := s.CheckConsumer(keyId1)
	assert.NoError(t, err)
	assert.Equal(t, cInfo.Value.GroupID, userAddr)
	assert.Equal(t, cInfo.Value.Username, s.consumerUsername(keyId1))
	cInfo, err = s.CheckConsumer(keyId2)
	assert.NoError(t, err)
	assert.Equal(t, cInfo.Value.GroupID, userAddr)
	assert.Equal(t, cInfo.Value.Username, s.consumerUsername(keyId2))

	// Pause consumer group
	err = s.PauseConsumerGroup(userAddr)
	assert.NoError(t, err)

	// TODO: Further test after pause

	// Resume paused consumer group
	err = s.ResumeConsumerGroup(userAddr)
	assert.NoError(t, err)

	// TODO: Further test after resume

	// Delete one consumer
	err = s.DeleteConsumer(keyId1)
	assert.NoError(t, err)

	// TODO: Further test after delete

	// Delete consumer group
	err = s.DeleteConsumerGroup(userAddr)
	assert.Error(t, err, fmt.Sprintf("can not delete this consumer group, consumer [%d] is still using it now", keyId2))

	// Delete another consumer
	err = s.DeleteConsumer(keyId2)
	assert.NoError(t, err)

	// TODO: Further test after delete

	// Delete consumer group
	err = s.DeleteConsumerGroup(userAddr)
	assert.NoError(t, err)

	cgInfo, err = s.CheckConsumerGroup(userAddr)
	assert.Nil(t, cgInfo)
	assert.Error(t, err, "Key not found")

	// Create consumer with non-exist group should create group automatically

	key3 := uuid.New().String()
	keyId3 := uint64(25565)

	err = s.NewConsumer(keyId3, key3, userAddr)
	assert.NoError(t, err)

	cgInfo, err = s.CheckConsumerGroup(userAddr)
	assert.NoError(t, err)
	assert.Equal(t, *cgInfo.Value.ID, userAddr)

	err = s.DeleteConsumer(keyId3)
	assert.NoError(t, err)

	err = s.DeleteConsumerGroup(userAddr)
	assert.NoError(t, err)

}

func Test_UsernameChangeHelper(t *testing.T) {
	// Init configs
	s, _ := New("http://localhost:9180", "edd1c9f034335f136f87ad84b625c8f1")

	keyId := uint64(114514)
	res, err := s.RecoverKeyIDFromConsumerUsername(s.consumerUsername(keyId))
	assert.Equal(t, keyId, res)
	assert.NoError(t, err)
}

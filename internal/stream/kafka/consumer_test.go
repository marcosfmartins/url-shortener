package kafka

import (
	"context"
	"encoding/json"
	"github.com/marcosfmartins/url-shortener/internal/entity"
	"github.com/marcosfmartins/url-shortener/pkg/id"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewConsumer(t *testing.T) {
	topic, err := id.GenerateID()
	assert.NoError(t, err)

	consumer := NewConsumer(kafkaURL, topic, topic)
	assert.NotNil(t, consumer)
	defer consumer.Close()
}

func TestConsumer_GetMessage(t *testing.T) {
	ctx := context.Background()
	topic, err := id.GenerateID()
	assert.NoError(t, err)

	consumer := NewConsumer(kafkaURL, topic, topic)
	assert.NotNil(t, consumer)
	defer consumer.Close()

	obj := NewPublisher(kafkaURL, topic)
	assert.NotNil(t, obj)
	defer func() {
		err := obj.Close()
		assert.NoError(t, err)
	}()

	expected := &entity.AccessEvent{ID: topic, Timestamp: time.Now()}
	obj.Publish(ctx, expected)
	time.Sleep(1 * time.Second)
	err = obj.Publish(ctx, expected)
	assert.NoError(t, err)

	msg, err := consumer.GetMessage(ctx)
	assert.NoError(t, err)

	result := &entity.AccessEvent{}
	err = json.Unmarshal(msg.Value(), result)
	assert.NoError(t, err)

	err = consumer.Commit(ctx, msg)
	assert.NoError(t, err)
}

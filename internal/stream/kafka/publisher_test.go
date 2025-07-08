package kafka

import (
	"context"
	"encoding/json"
	"github.com/marcosfmartins/url-shortener/internal/entity"
	"github.com/marcosfmartins/url-shortener/pkg/id"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"time"
)

var kafkaURL []string

func init() {
	url := os.Getenv("KAFKA_BROKER")
	kafkaURL = strings.Split(url, ",")
}

func TestNewPublisher(t *testing.T) {
	topic, err := id.GenerateID()
	assert.NoError(t, err)

	obj := NewPublisher(kafkaURL, topic)
	assert.NotNil(t, obj)
	defer func() {
		err := obj.Close()
		assert.NoError(t, err)
	}()
}

func TestPublisher_Publish(t *testing.T) {
	ctx := context.Background()
	topic, err := id.GenerateID()
	assert.NoError(t, err)

	obj := NewPublisher(kafkaURL, topic)
	assert.NotNil(t, obj)
	defer func() {
		err := obj.Close()
		assert.NoError(t, err)
	}()

	expected := &entity.AccessEvent{
		ID:        topic,
		Timestamp: time.Now(),
	}

	obj.Publish(ctx, expected) //To create topic
	time.Sleep(1 * time.Second)

	err = obj.Publish(ctx, expected)
	assert.NoError(t, err)

	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaURL,
		Topic:   topic,
		GroupID: topic,
	})
	defer consumer.Close()

	msg, err := consumer.ReadMessage(ctx)
	assert.NoError(t, err)

	result := &entity.AccessEvent{}
	err = json.Unmarshal(msg.Value, result)
	assert.NoError(t, err)

	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Timestamp.Format(time.RFC3339Nano), result.Timestamp.Format(time.RFC3339Nano))
}

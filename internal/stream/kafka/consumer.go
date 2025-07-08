package kafka

import (
	"context"
	"errors"
	"github.com/marcosfmartins/url_shortener/internal/entity"
	"github.com/segmentio/kafka-go"
)

const MaxBytes = 10e6 // 10MB

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(url []string, topic, consumerGroup string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  url,
			GroupID:  consumerGroup,
			Topic:    topic,
			MaxBytes: MaxBytes,
		}),
	}
}

func (c *Consumer) GetMessage(ctx context.Context) (entity.Message, error) {
	msg, err := c.reader.FetchMessage(ctx)
	return &Message{msg: msg}, err
}

func (c *Consumer) Commit(ctx context.Context, msg entity.Message) error {
	rawMSg, ok := msg.(*Message)
	if !ok {
		return errors.New("invalid message")
	}

	return c.reader.CommitMessages(ctx, rawMSg.msg)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

type Message struct {
	msg kafka.Message
}

func (m *Message) Value() []byte {
	return m.msg.Value
}

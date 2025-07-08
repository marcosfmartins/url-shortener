package kafka

import (
	"context"
	"github.com/marcosfmartins/url_shortener/internal/entity"
	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(url []string, topic string) *Publisher {
	return &Publisher{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(url...),
			Topic:                  topic,
			Async:                  true,
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *Publisher) Publish(ctx context.Context, msg entity.Event) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: msg.Bytes(),
	})
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}

package entity

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, msg Event) error
	Close() error
}

type Consumer interface {
	GetMessage(ctx context.Context) (Message, error)
	Commit(ctx context.Context, msg Message) error
	Close() error
}

type Message interface {
	Value() []byte
}

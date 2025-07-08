package entity

import (
	"context"
	"time"
)

type URLRepository interface {
	Insert(ctx context.Context, url *URL) error
	FindByID(ctx context.Context, id string) (*URL, error)
	DeleteByID(ctx context.Context, id string) error
	Increment(ctx context.Context, obj []URL) error
}

type CacheRepository interface {
	Get(ctx context.Context, id string) (string, error)
	Set(ctx context.Context, key, value string, TTL time.Duration) error
	Delete(ctx context.Context, id string) error
}

package redis

import (
	"context"
	"github.com/marcosfmartins/url_shortener/internal/entity"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type Cache struct {
	client *redis.Client
}

func NewCache(URL string) *Cache {
	return &Cache{
		client: redis.NewClient(&redis.Options{Addr: URL}),
	}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil && !strings.Contains(err.Error(), "nil") {
		return "", entity.NotFoundError.WithError(err)
	}

	return val, nil
}

func (c *Cache) Set(ctx context.Context, key, value string, TTL time.Duration) error {
	err := c.client.Set(ctx, key, value, TTL).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, id).Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}

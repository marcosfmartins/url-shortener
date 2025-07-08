package redis

import (
	"context"
	"github.com/marcosfmartins/url-shortener/pkg/id"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	URL string
	ttl = 1 * time.Second
	ctx = context.Background()
)

func init() {
	URL = os.Getenv("REDIS_URL")
}

func TestNewCache(t *testing.T) {
	cache := NewCache(URL)
	assert.NotNil(t, cache)
	defer cache.Close()
}

func TestCache_Set(t *testing.T) {
	cache := NewCache(URL)
	assert.NotNil(t, cache)

	key, _ := id.GenerateID()

	err := cache.Set(ctx, key, "1", ttl)
	assert.NoError(t, err)

	result := cache.client.Get(ctx, key)

	value, err := result.Result()
	assert.NoError(t, err)

	assert.Equal(t, "1", value)
}

func TestCache_Get(t *testing.T) {
	t.Run("sucess", func(t *testing.T) {
		cache := NewCache(URL)
		assert.NotNil(t, cache)

		key, _ := id.GenerateID()

		err := cache.client.Set(ctx, key, "1", ttl).Err()
		assert.NoError(t, err)

		value, err := cache.Get(ctx, key)
		assert.NoError(t, err)

		assert.Equal(t, "1", value)
	})

	t.Run("not exist", func(t *testing.T) {
		cache := NewCache(URL)
		assert.NotNil(t, cache)

		key, _ := id.GenerateID()

		value, err := cache.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, "", value)
	})
}

func TestCache_Delete(t *testing.T) {
	cache := NewCache(URL)
	assert.NotNil(t, cache)

	key, _ := id.GenerateID()

	err := cache.client.Set(ctx, key, "1", ttl).Err()
	assert.NoError(t, err)

	err = cache.Delete(ctx, key)
	assert.NoError(t, err)
}

package service

import (
	"bytes"
	"context"
	"github.com/marcosfmartins/url-shortener/internal/entity"
	"github.com/marcosfmartins/url-shortener/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	ctx = context.Background()
	log = logger.NewZerologAdapter()
)

func TestNewURLService(t *testing.T) {
	obj := NewURLService(nil, nil, nil, nil)
	assert.NotNil(t, obj)
}

func TestURL_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		urlRepoMock := &entity.URLRepositoryMock{
			InsertFn: func(ctx context.Context, url *entity.URL) error {
				assert.NotNil(t, url.ID)
				assert.Len(t, url.ID, 10)
				assert.Equal(t, "https://example.com", url.Original)
				assert.NotEmpty(t, url.CreatedAt)
				return nil
			},
		}

		obj := NewURLService(log, urlRepoMock, nil, nil)

		url, err := obj.Create(ctx, "https://example.com")
		assert.NoError(t, err)
		assert.NotNil(t, url)

		assert.Equal(t, 1, urlRepoMock.InsertCount)
	})

	t.Run("Repository Error", func(t *testing.T) {
		urlRepoMock := &entity.URLRepositoryMock{
			InsertFn: func(ctx context.Context, url *entity.URL) error {
				return assert.AnError
			},
		}

		obj := NewURLService(log, urlRepoMock, nil, nil)

		url, err := obj.Create(ctx, "https://example.com")
		assert.Error(t, err)
		assert.Nil(t, url)

		assert.Equal(t, 1, urlRepoMock.InsertCount)
	})
}

func TestURL_GetURL(t *testing.T) {
	t.Run("Success no cache", func(t *testing.T) {
		expectedURL := &entity.URL{
			ID:       "12345",
			Original: "https://example.com",
		}

		urlRepoMock := &entity.URLRepositoryMock{
			FindByIDFn: func(ctx context.Context, id string) (*entity.URL, error) {
				assert.Equal(t, expectedURL.ID, id)
				return expectedURL, nil
			},
		}

		cacheRepoMock := &entity.CacheRepositoryMock{
			GetFn: func(ctx context.Context, id string) (string, error) {
				assert.Equal(t, expectedURL.ID, id)
				return "", nil
			},
			SetFn: func(ctx context.Context, key, value string, TTL time.Duration) error {
				assert.Equal(t, expectedURL.ID, key)
				assert.Equal(t, expectedURL.Original, value)
				assert.Equal(t, 10*time.Minute, TTL)
				return nil
			},
		}

		publisherMock := &entity.PublisherMock{
			PublishFn: func(ctx context.Context, msg entity.Event) error {
				obj, ok := msg.(*entity.AccessEvent)
				assert.True(t, ok)
				assert.Equal(t, expectedURL.ID, obj.ID)
				return nil
			},
		}

		svc := NewURLService(log, urlRepoMock, cacheRepoMock, publisherMock)

		URL, err := svc.GetURL(ctx, expectedURL.ID)
		assert.NoError(t, err)
		assert.NotNil(t, URL)

		assert.Equal(t, expectedURL.Original, URL)
		assert.Equal(t, 1, urlRepoMock.FindByIDCount)
		assert.Equal(t, 1, cacheRepoMock.GetCount)
		assert.Equal(t, 1, cacheRepoMock.SetCount)
		assert.Equal(t, 1, publisherMock.PublishCount)
	})

	t.Run("Success cache", func(t *testing.T) {
		expectedURL := &entity.URL{
			ID:       "12345",
			Original: "https://example.com",
		}

		urlRepoMock := &entity.URLRepositoryMock{}

		cacheRepoMock := &entity.CacheRepositoryMock{
			GetFn: func(ctx context.Context, id string) (string, error) {
				return expectedURL.Original, nil
			},
		}

		publisherMock := &entity.PublisherMock{
			PublishFn: func(ctx context.Context, msg entity.Event) error {
				obj, ok := msg.(*entity.AccessEvent)
				assert.True(t, ok)
				assert.Equal(t, expectedURL.ID, obj.ID)
				return nil
			},
		}

		svc := NewURLService(log, urlRepoMock, cacheRepoMock, publisherMock)

		URL, err := svc.GetURL(ctx, expectedURL.ID)
		assert.NoError(t, err)
		assert.NotNil(t, URL)

		assert.Equal(t, expectedURL.Original, URL)
		assert.Equal(t, 0, urlRepoMock.FindByIDCount)
		assert.Equal(t, 1, cacheRepoMock.GetCount)
		assert.Equal(t, 0, cacheRepoMock.SetCount)
		assert.Equal(t, 1, publisherMock.PublishCount)
	})

	t.Run("cache/pubisher Error", func(t *testing.T) {
		buf := bytes.Buffer{}
		zlog := zerolog.New(&buf)
		localLog := logger.NewWithLogger(zlog)

		expectedURL := &entity.URL{
			ID:       "12345",
			Original: "https://example.com",
		}

		urlRepoMock := &entity.URLRepositoryMock{
			FindByIDFn: func(ctx context.Context, id string) (*entity.URL, error) {
				assert.Equal(t, expectedURL.ID, id)
				return expectedURL, nil
			},
		}

		cacheRepoMock := &entity.CacheRepositoryMock{
			GetFn: func(ctx context.Context, id string) (string, error) {
				return "", assert.AnError
			},
			SetFn: func(ctx context.Context, key, value string, TTL time.Duration) error {
				assert.Equal(t, expectedURL.ID, key)
				assert.Equal(t, expectedURL.Original, value)
				assert.Equal(t, 10*time.Minute, TTL)
				return assert.AnError
			},
		}

		publisherMock := &entity.PublisherMock{
			PublishFn: func(ctx context.Context, msg entity.Event) error {
				obj, ok := msg.(*entity.AccessEvent)
				assert.True(t, ok)
				assert.Equal(t, expectedURL.ID, obj.ID)
				return assert.AnError
			},
		}

		svc := NewURLService(localLog, urlRepoMock, cacheRepoMock, publisherMock)

		URL, err := svc.GetURL(ctx, expectedURL.ID)
		assert.NoError(t, err)
		assert.NotNil(t, URL)

		assert.Equal(t, expectedURL.Original, URL)
		assert.Equal(t, 1, urlRepoMock.FindByIDCount)
		assert.Equal(t, 1, cacheRepoMock.GetCount)
		assert.Equal(t, 1, cacheRepoMock.SetCount)
		assert.Equal(t, 1, publisherMock.PublishCount)

		logStr := buf.String()
		assert.Contains(t, logStr, "failed to get URL from cache")
		assert.Contains(t, logStr, "failed to cache URL")
		assert.Contains(t, logStr, "failed to publish access event")
	})

	t.Run("Repository Error", func(t *testing.T) {
		urlRepoMock := &entity.URLRepositoryMock{
			FindByIDFn: func(ctx context.Context, id string) (*entity.URL, error) {
				return nil, assert.AnError
			},
		}

		cacheRepoMock := &entity.CacheRepositoryMock{
			GetFn: func(ctx context.Context, id string) (string, error) {
				return "", assert.AnError
			},
			SetFn: func(ctx context.Context, key, value string, TTL time.Duration) error {
				return nil
			},
		}

		publisherMock := &entity.PublisherMock{
			PublishFn: func(ctx context.Context, msg entity.Event) error {
				return nil
			},
		}

		svc := NewURLService(log, urlRepoMock, cacheRepoMock, publisherMock)

		URL, err := svc.GetURL(ctx, "666")
		assert.Error(t, err)
		assert.Empty(t, URL)

		assert.Equal(t, 1, urlRepoMock.FindByIDCount)
		assert.Equal(t, 1, cacheRepoMock.GetCount)
		assert.Equal(t, 0, cacheRepoMock.SetCount)
		assert.Equal(t, 0, publisherMock.PublishCount)
	})

}

func TestURL_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		urlRepoMock := &entity.URLRepositoryMock{
			DeleteByIDFn: func(ctx context.Context, id string) error {
				assert.Equal(t, "12345", id)
				return nil
			},
		}

		cacheRepoMock := &entity.CacheRepositoryMock{
			DeleteFn: func(ctx context.Context, id string) error {
				assert.Equal(t, "12345", id)
				return nil
			},
		}

		svc := NewURLService(log, urlRepoMock, cacheRepoMock, nil)
		err := svc.Delete(ctx, "12345")
		assert.NoError(t, err)

		assert.Equal(t, 1, urlRepoMock.DeleteByIDCount)
		assert.Equal(t, 1, cacheRepoMock.DeleteCount)
	})

	t.Run("Repository Error", func(t *testing.T) {
		urlRepoMock := &entity.URLRepositoryMock{
			DeleteByIDFn: func(ctx context.Context, id string) error {
				return assert.AnError
			},
		}

		cacheRepoMock := &entity.CacheRepositoryMock{
			DeleteFn: func(ctx context.Context, id string) error {
				return nil
			},
		}

		svc := NewURLService(log, urlRepoMock, cacheRepoMock, nil)
		err := svc.Delete(ctx, "12345")
		assert.Error(t, err)

		assert.Equal(t, 1, urlRepoMock.DeleteByIDCount)
		assert.Equal(t, 1, cacheRepoMock.DeleteCount)
	})

	t.Run("Cache Error", func(t *testing.T) {
		buf := bytes.Buffer{}
		zlog := zerolog.New(&buf)
		localLog := logger.NewWithLogger(zlog)

		urlRepoMock := &entity.URLRepositoryMock{
			DeleteByIDFn: func(ctx context.Context, id string) error {
				return nil
			},
		}

		cacheRepoMock := &entity.CacheRepositoryMock{
			DeleteFn: func(ctx context.Context, id string) error {
				return assert.AnError
			},
		}

		svc := NewURLService(localLog, urlRepoMock, cacheRepoMock, nil)
		err := svc.Delete(ctx, "12345")
		assert.Nil(t, err)

		assert.Equal(t, 1, urlRepoMock.DeleteByIDCount)
		assert.Equal(t, 1, cacheRepoMock.DeleteCount)
		assert.Contains(t, buf.String(), "failed to delete URL from cache")
	})
}

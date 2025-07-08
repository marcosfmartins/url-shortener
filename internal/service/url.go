package service

import (
	"context"
	"github.com/marcosfmartins/url-shortener/internal/entity"
	"github.com/marcosfmartins/url-shortener/pkg/id"
	"strings"
	"time"
)

const TTL = 10 * time.Minute

type URL struct {
	db        entity.URLRepository
	cache     entity.CacheRepository
	publisher entity.Publisher
	logger    entity.Logger
}

func NewURLService(logger entity.Logger, db entity.URLRepository, cache entity.CacheRepository, publisher entity.Publisher) *URL {
	return &URL{
		db:        db,
		cache:     cache,
		publisher: publisher,
		logger:    logger,
	}
}

func (s *URL) Create(ctx context.Context, url string) (*entity.URL, error) {
	ID, err := id.GenerateID()
	if err != nil {
		return nil, err
	}

	obj := entity.URL{
		ID:        ID,
		Original:  url,
		CreatedAt: time.Now(),
	}

	err = s.db.Insert(ctx, &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func (s *URL) Get(ctx context.Context, ID string) (*entity.URL, error) {
	return s.db.FindByID(ctx, ID)
}

func (s *URL) GetURL(ctx context.Context, ID string) (string, error) {
	cache, err := s.cache.Get(ctx, ID)
	if err != nil {
		s.logger.Err(err).Error("failed to get URL from cache")
	}

	result := cache

	if strings.TrimSpace(cache) == "" {
		url, err := s.db.FindByID(ctx, ID)
		if err != nil {
			return "", err
		}

		result = url.Original

		err = s.cache.Set(ctx, ID, url.Original, TTL)
		if err != nil {
			s.logger.Err(err).Error("failed to cache URL")
		}
	}

	err = s.publisher.Publish(ctx, &entity.AccessEvent{ID: ID, Timestamp: time.Now()})
	if err != nil {
		s.logger.Err(err).Error("failed to publish access event")
	}

	return result, nil
}

func (s *URL) Delete(ctx context.Context, ID string) error {
	err := s.cache.Delete(ctx, ID)
	if err != nil {
		s.logger.Err(err).Error("failed to delete URL from cache")
	}

	return s.db.DeleteByID(ctx, ID)
}

func (s *URL) Increment(ctx context.Context, urls []entity.URL) error {
	if len(urls) == 0 {
		return nil
	}

	err := s.db.Increment(ctx, urls)
	if err != nil {
		s.logger.Err(err).Error("failed to increment URL hits")
		return err
	}

	return nil
}

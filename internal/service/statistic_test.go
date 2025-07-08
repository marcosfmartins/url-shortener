package service

import (
	"context"
	"github.com/marcosfmartins/url_shortener/internal/entity"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNewStatisticService(t *testing.T) {
	svc := NewStatisticService(nil, nil)
	assert.NotNil(t, svc)
}

func TestStatisticService_Processor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)

	expected := entity.URL{ID: "123", Hits: 1}
	var events []entity.URL

	urlSvc := &entity.URLServiceMock{
		IncrementFn: func(ctx context.Context, obj []entity.URL) error {
			events = append(events, obj...)
			wg.Done()
			return nil
		},
	}

	consumer := &entity.ConsumerMock{
		GetMessageFn: func(ctx context.Context) (entity.Message, error) {
			cancel()
			return &entity.MessageMock{Data: []byte(`{"id":"123"}`)}, nil
		},
		CommitFn: func(ctx context.Context, msg entity.Message) error {
			return nil
		},
	}

	svc := NewStatisticService(urlSvc, consumer)

	flushCtx, flushCancel := context.WithCancel(context.Background())
	defer flushCancel()

	go svc.SchedulerFlush(flushCtx)
	go svc.Processor(ctx)
	wg.Wait()

	assert.Equal(t, 1, urlSvc.IncrementCount)
	assert.Equal(t, 1, consumer.GetMessageCount)
	assert.Equal(t, 1, consumer.CommitCount)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, expected.ID, events[0].ID)
	assert.Equal(t, expected.Hits, events[0].Hits)

}

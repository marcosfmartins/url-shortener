package service

import (
	"context"
	"github.com/marcosfmartins/url-shortener/internal/entity"
	"sync"
	"time"
)

const bufferSize = 1000
const tickerInterval = 1 * time.Second

type StatisticService struct {
	url      entity.URLService
	consumer entity.Consumer
	mutex    sync.Mutex
	buffer   map[string]entity.URL
	messages []entity.Message
}

func NewStatisticService(url entity.URLService, consumer entity.Consumer) *StatisticService {
	return &StatisticService{
		url:      url,
		consumer: consumer,
		buffer:   make(map[string]entity.URL),
		messages: make([]entity.Message, 0, bufferSize),
	}
}

func (ss *StatisticService) SchedulerFlush(ctx context.Context) {
	flusherTicker := time.NewTicker(tickerInterval)
	defer flusherTicker.Stop()

loop:
	for {
		select {
		case <-flusherTicker.C:
			ss.flush()
		case <-ctx.Done():
			break loop
		}
	}
}

func (ss *StatisticService) Processor(ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			ss.process(ctx)
		}
	}
}

func (ss *StatisticService) process(ctx context.Context) {
	msg, err := ss.consumer.GetMessage(ctx)
	if err != nil {
		return
	}

	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	event := entity.NewAccessEvent(msg.Value())

	if obj, ok := ss.buffer[event.ID]; ok {
		obj.Hits++
		obj.LastAccess = &event.Timestamp
		ss.buffer[event.ID] = obj
	} else {
		ss.buffer[event.ID] = entity.URL{
			ID:         event.ID,
			Hits:       1,
			LastAccess: &event.Timestamp,
		}
	}

	ss.messages = append(ss.messages, msg)
}

func (ss *StatisticService) flush() {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	if len(ss.messages) == 0 {
		return
	}

	bulk := make([]entity.URL, 0, len(ss.buffer))
	for _, value := range ss.buffer {
		bulk = append(bulk, value)
	}

	err := ss.url.Increment(context.Background(), bulk)
	if err != nil {
		return
	}

	for _, msg := range ss.messages {
		_ = ss.consumer.Commit(context.Background(), msg)
	}

	ss.buffer = make(map[string]entity.URL)
	ss.messages = make([]entity.Message, 0, bufferSize)
}

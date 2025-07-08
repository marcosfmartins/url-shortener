package entity

import "context"

type PublisherMock struct {
	PublishFn    func(ctx context.Context, msg Event) error
	PublishCount int
	CloseFn      func() error
	CloseCount   int
}

func (pm *PublisherMock) Publish(ctx context.Context, msg Event) error {
	pm.PublishCount++
	return pm.PublishFn(ctx, msg)
}

func (pm *PublisherMock) Close() error {
	pm.CloseCount++

	if pm.CloseFn != nil {
		return pm.CloseFn()
	}

	return nil
}

type ConsumerMock struct {
	GetMessageFn    func(ctx context.Context) (Message, error)
	GetMessageCount int
	CommitFn        func(ctx context.Context, msg Message) error
	CommitCount     int
	CloseFn         func() error
	CloseCount      int
}

func (cm *ConsumerMock) GetMessage(ctx context.Context) (Message, error) {
	cm.GetMessageCount++
	return cm.GetMessageFn(ctx)
}

func (cm *ConsumerMock) Commit(ctx context.Context, msg Message) error {
	cm.CommitCount++
	return cm.CommitFn(ctx, msg)
}

func (cm *ConsumerMock) Close() error {
	cm.CloseCount++

	if cm.CloseFn != nil {
		return cm.CloseFn()
	}

	return nil
}

type MessageMock struct {
	Data []byte
}

func (mm *MessageMock) Value() []byte {
	return mm.Data
}

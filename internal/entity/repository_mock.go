package entity

import (
	"context"
	"time"
)

type URLRepositoryMock struct {
	InsertFn        func(ctx context.Context, url *URL) error
	InsertCount     int
	FindByIDFn      func(ctx context.Context, id string) (*URL, error)
	FindByIDCount   int
	DeleteByIDFn    func(ctx context.Context, id string) error
	DeleteByIDCount int
	IncrementFn     func(ctx context.Context, objs []URL) error
	IncrementCount  int
}

func (urm *URLRepositoryMock) Insert(ctx context.Context, url *URL) error {
	urm.InsertCount++
	return urm.InsertFn(ctx, url)
}

func (urm *URLRepositoryMock) FindByID(ctx context.Context, id string) (*URL, error) {
	urm.FindByIDCount++
	return urm.FindByIDFn(ctx, id)
}

func (urm *URLRepositoryMock) DeleteByID(ctx context.Context, id string) error {
	urm.DeleteByIDCount++
	return urm.DeleteByIDFn(ctx, id)
}

func (urm *URLRepositoryMock) Increment(ctx context.Context, objs []URL) error {
	urm.IncrementCount++
	return urm.IncrementFn(ctx, objs)
}

type CacheRepositoryMock struct {
	GetFn       func(ctx context.Context, id string) (string, error)
	GetCount    int
	SetFn       func(ctx context.Context, key, value string, TTL time.Duration) error
	SetCount    int
	DeleteFn    func(ctx context.Context, id string) error
	DeleteCount int
}

func (crm *CacheRepositoryMock) Get(ctx context.Context, id string) (string, error) {
	crm.GetCount++
	return crm.GetFn(ctx, id)
}

func (crm *CacheRepositoryMock) Set(ctx context.Context, key, value string, TTL time.Duration) error {
	crm.SetCount++
	return crm.SetFn(ctx, key, value, TTL)
}

func (crm *CacheRepositoryMock) Delete(ctx context.Context, id string) error {
	crm.DeleteCount++
	return crm.DeleteFn(ctx, id)
}

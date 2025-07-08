package entity

import "context"

type URLServiceMock struct {
	CreateFn       func(ctx context.Context, url string) (*URL, error)
	CreateCount    int
	GetFn          func(ctx context.Context, ID string) (*URL, error)
	GetCount       int
	GetURLFn       func(ctx context.Context, ID string) (string, error)
	GetURLCount    int
	DeleteFn       func(ctx context.Context, ID string) error
	DeleteCount    int
	IncrementFn    func(ctx context.Context, objs []URL) error
	IncrementCount int
}

func (u *URLServiceMock) Create(ctx context.Context, url string) (*URL, error) {
	u.CreateCount++
	return u.CreateFn(ctx, url)
}

func (u *URLServiceMock) Get(ctx context.Context, ID string) (*URL, error) {
	u.GetCount++
	return u.GetFn(ctx, ID)
}

func (u *URLServiceMock) GetURL(ctx context.Context, ID string) (string, error) {
	u.GetURLCount++
	return u.GetURLFn(ctx, ID)
}

func (u *URLServiceMock) Delete(ctx context.Context, ID string) error {
	u.DeleteCount++
	return u.DeleteFn(ctx, ID)
}

func (u *URLServiceMock) Increment(ctx context.Context, objs []URL) error {
	u.IncrementCount++
	return u.IncrementFn(ctx, objs)
}

package entity

import "context"

type URLService interface {
	Create(ctx context.Context, url string) (*URL, error)
	Get(ctx context.Context, ID string) (*URL, error)
	GetURL(ctx context.Context, ID string) (string, error)
	Delete(ctx context.Context, ID string) error
	Increment(ctx context.Context, objs []URL) error
}

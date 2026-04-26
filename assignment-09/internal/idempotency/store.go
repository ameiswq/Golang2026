package idempotency

import "context"

type CachedResponse struct {
	StatusCode int
	Body string
	Completed bool
}

type Store interface {
	StartProcessing(ctx context.Context, key string) (bool, error)
	Get(ctx context.Context, key string) (*CachedResponse, bool, error)
	Finish(ctx context.Context, key string, statusCode int, body string) error
}

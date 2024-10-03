package cache

import (
	"context"
	"time"
)

type UserCache interface {
	Get(ctx context.Context, key string) (string, error)
	SetWithExpiration(ctx context.Context, key string, data []byte, expiration time.Duration) error
}

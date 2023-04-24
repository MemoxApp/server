package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type (
	Svc interface {
		Get(ctx context.Context, key string, ttl time.Duration, f func() ([]byte, error)) ([]byte, error)
		Del(ctx context.Context, key string)
	}

	defaultCacheSvc struct {
		r *redis.Client
	}
)

func NewCacheSvc(redis *redis.Client) Svc {
	return &defaultCacheSvc{
		r: redis,
	}
}

func (s *defaultCacheSvc) Get(ctx context.Context, key string, ttl time.Duration, f func() ([]byte, error)) ([]byte, error) {
	result := s.r.Get(ctx, key)
	if result.Err() != nil {
		data, err := f()
		if err != nil {
			return data, err
		}

		_ = s.r.Set(ctx, key, data, ttl)
		return data, err
	} else {
		return result.Bytes()
	}
}

func (s *defaultCacheSvc) Del(ctx context.Context, key string) {
	s.r.Del(ctx, key)
}

package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type (
	Svc struct {
		r *redis.Client
	}
)

func NewCacheSvc(redis *redis.Client) *Svc {
	return &Svc{
		r: redis,
	}
}

func (s *Svc) Get(ctx context.Context, key string, ttl time.Duration, f func() ([]byte, error)) ([]byte, error) {
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

func (s *Svc) Del(ctx context.Context, key string) {
	s.r.Del(ctx, key)
}

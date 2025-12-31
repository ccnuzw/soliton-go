package lock

import (
	"context"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

// Locker is the interface for obtaining locks.
type Locker interface {
	Obtain(ctx context.Context, key string, ttl time.Duration) (Lock, error)
}

// Lock represents a held lock.
type Lock interface {
	Release(ctx context.Context) error
}

// RedisLocker implements Locker using Redis.
type RedisLocker struct {
	client *redislock.Client
}

// NewRedisLocker creates a new RedisLocker.
func NewRedisLocker(redisClient redis.UniversalClient) *RedisLocker {
	return &RedisLocker{
		client: redislock.New(redisClient),
	}
}

func (l *RedisLocker) Obtain(ctx context.Context, key string, ttl time.Duration) (Lock, error) {
	// Linear backoff for retry
	backoff := redislock.LinearBackoff(100 * time.Millisecond)
	lock, err := l.client.Obtain(ctx, key, ttl, &redislock.Options{
		RetryStrategy: backoff,
	})
	if err != nil {
		if err == redislock.ErrNotObtained {
			return nil, fmt.Errorf("could not obtain lock for key %s", key)
		}
		return nil, fmt.Errorf("failed to obtain lock: %w", err)
	}
	return &redisLockWrapper{lock: lock}, nil
}

type redisLockWrapper struct {
	lock *redislock.Lock
}

func (l *redisLockWrapper) Release(ctx context.Context) error {
	return l.lock.Release(ctx)
}

package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"time"
)

type IdempotentCache struct {
	client *redis.Client
}

func New(redisConn string) *IdempotentCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &IdempotentCache{rdb}
}

func (c *IdempotentCache) Get(ctx context.Context, key uuid.UUID) (err error) {
	_, err = c.client.Get(ctx, key.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = ErrCacheMiss
			return
		} else {
			return
		}
	}
	return
}

func (c *IdempotentCache) Set(ctx context.Context, key uuid.UUID) (err error) {
	return c.client.Set(
		ctx,
		key.String(),
		"",
		time.Minute*15,
	).Err()
}

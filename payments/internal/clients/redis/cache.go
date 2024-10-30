package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const cacheTTL = time.Minute * 15

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
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get data from redis: %w", err)
	}
	return
}

func (c *IdempotentCache) Set(ctx context.Context, key uuid.UUID) (err error) {
	err = c.client.Set(
		ctx,
		key.String(),
		"",
		cacheTTL,
	).Err()
	if err != nil {
		return fmt.Errorf("failed to set: %w", err)
	}
	return nil
}

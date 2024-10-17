package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NStegura/saga/tgbot/internal/domain"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache реализует storage (корзина)
type Cache struct {
	client *redis.Client
}

func New(redisConn string) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &Cache{rdb}
}

func (c *Cache) keyFormat(userID int64) string {
	return fmt.Sprintf("user:%d:state", userID)
}

func (c *Cache) Get(ctx context.Context, userID int64) (userState domain.UserState, err error) {
	val, err := c.client.Get(ctx, c.keyFormat(userID)).Result()
	if err != nil {
		if err == redis.Nil {
			return domain.NewUserState(userID), nil
		}
		return userState, fmt.Errorf("failed to get user state: %w", err)
	}

	err = json.Unmarshal([]byte(val), &userState)
	if err != nil {
		return userState, fmt.Errorf("failed to unmarshal user state: %w", err)
	}
	return userState, nil
}

func (c *Cache) Set(ctx context.Context, userState domain.UserState) (err error) {
	data, err := json.Marshal(userState)
	if err != nil {
		return fmt.Errorf("failed to save user state: %w", err)
	}
	return c.client.Set(
		ctx,
		c.keyFormat(userState.UserID),
		data,
		time.Minute*60,
	).Err()
}

func (c *Cache) Del(ctx context.Context, userID int64) (err error) {
	return c.client.Del(
		ctx,
		c.keyFormat(userID),
	).Err()
}

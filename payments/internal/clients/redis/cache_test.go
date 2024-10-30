package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIdempotentCache_Get_Success(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	cache := &IdempotentCache{client: mockRedis}

	key := uuid.New()
	mock.ExpectGet(key.String()).SetVal("some value")

	err := cache.Get(context.Background(), key)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIdempotentCache_Get_CacheMiss(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	cache := &IdempotentCache{client: mockRedis}

	key := uuid.New()
	mock.ExpectGet(key.String()).RedisNil()

	err := cache.Get(context.Background(), key)

	assert.ErrorIs(t, err, ErrCacheMiss)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIdempotentCache_Get_InternalError(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	cache := &IdempotentCache{client: mockRedis}

	key := uuid.New()
	mock.ExpectGet(key.String()).SetErr(errors.New("internal redis error"))

	err := cache.Get(context.Background(), key)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get data from redis")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIdempotentCache_Set_Success(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	cache := &IdempotentCache{client: mockRedis}

	key := uuid.New()
	mock.ExpectSet(key.String(), "", time.Minute*15).SetVal("OK")

	err := cache.Set(context.Background(), key)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIdempotentCache_Set_Error(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	cache := &IdempotentCache{client: mockRedis}

	key := uuid.New()
	mock.ExpectSet(key.String(), "", time.Minute*15).SetErr(errors.New("redis set error"))

	err := cache.Set(context.Background(), key)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis set error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

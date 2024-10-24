package redis

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"

	"github.com/stretchr/testify/assert"

	"github.com/NStegura/saga/tgbot/internal/domain"
)

func TestCache_Get_Success(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	cache := &Cache{client: mockRedis}

	userID := int64(12345)
	userState := domain.NewUserState(userID)
	userStateJSON, _ := json.Marshal(userState)

	mock.ExpectGet(cache.keyFormat(userID)).SetVal(string(userStateJSON))

	state, err := cache.Get(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, userState, state)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_Get_CacheMiss(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	cache := &Cache{client: mockRedis}

	userID := int64(12345)
	mock.ExpectGet(cache.keyFormat(userID)).RedisNil()

	state, err := cache.Get(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, domain.NewUserState(userID), state)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_Get_Error(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	cache := &Cache{client: mockRedis}

	userID := int64(12345)
	mock.ExpectGet(cache.keyFormat(userID)).SetErr(errors.New("internal redis error"))

	_, err := cache.Get(context.Background(), userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get user state")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_Set_Success(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	cache := &Cache{client: mockRedis}

	userID := int64(12345)
	userState := domain.NewUserState(userID)
	userStateJSON, _ := json.Marshal(userState)

	mock.ExpectSet(cache.keyFormat(userID), userStateJSON, time.Minute*60).SetVal("OK")

	err := cache.Set(context.Background(), userState)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_Set_Error(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	cache := &Cache{client: mockRedis}

	userID := int64(12345)
	userState := domain.NewUserState(userID)
	userStateJSON, _ := json.Marshal(userState)

	mock.ExpectSet(cache.keyFormat(userID), userStateJSON, time.Minute*60).SetErr(errors.New("redis set error"))

	err := cache.Set(context.Background(), userState)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis set error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_Del_Success(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	cache := &Cache{client: mockRedis}

	userID := int64(12345)
	mock.ExpectDel(cache.keyFormat(userID)).SetVal(1)

	err := cache.Del(context.Background(), userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCache_Del_Error(t *testing.T) {
	mockRedis, mock := redismock.NewClientMock()
	cache := &Cache{client: mockRedis}

	userID := int64(12345)
	mock.ExpectDel(cache.keyFormat(userID)).SetErr(errors.New("redis delete error"))

	err := cache.Del(context.Background(), userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis delete error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

package ordercons

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_Success(t *testing.T) {
	err := os.Setenv("CONSUMER_TOPICS", "topic1,topic2")
	assert.NoError(t, err)
	err = os.Setenv("CONSUMER_BROKERS", "broker1,broker2")
	assert.NoError(t, err)
	err = os.Setenv("DB_DSN", "postgres://user:password@localhost:5432/dbname")
	assert.NoError(t, err)
	err = os.Setenv("REDIS_DSN", "redis://localhost:6379")
	assert.NoError(t, err)
	err = os.Setenv("ORDER_CLI_CONN", "order-service:50051")
	assert.NoError(t, err)
	err = os.Setenv("LOG_LEVEL", "INFO")
	assert.NoError(t, err)
	err = os.Setenv("SHUTDOWN_TIMEOUT", "10s")
	assert.NoError(t, err)

	defer func() {
		_ = os.Unsetenv("CONSUMER_TOPICS")
		_ = os.Unsetenv("CONSUMER_BROKERS")
		_ = os.Unsetenv("DB_DSN")
		_ = os.Unsetenv("REDIS_DSN")
		_ = os.Unsetenv("ORDER_CLI_CONN")
		_ = os.Unsetenv("LOG_LEVEL")
		_ = os.Unsetenv("SHUTDOWN_TIMEOUT")
	}()

	cfg, err := New()
	assert.NoError(t, err)

	assert.Equal(t, []string{"topic1", "topic2"}, cfg.Consumer.Topics)
	assert.Equal(t, []string{"broker1", "broker2"}, cfg.Consumer.Brokers)
	assert.Equal(t, "postgres://user:password@localhost:5432/dbname", cfg.DB.DSN)
	assert.Equal(t, "redis://localhost:6379", cfg.Redis.DSN)
	assert.Equal(t, "order-service:50051", cfg.OrderCli.CONN)
	assert.Equal(t, "INFO", cfg.LogLevel)
	assert.Equal(t, 10*time.Second, cfg.ShutdownTimeout)
}

func TestNewConfig_DefaultValues(t *testing.T) {
	err := os.Setenv("CONSUMER_TOPICS", "topic1,topic2")
	assert.NoError(t, err)
	err = os.Setenv("CONSUMER_BROKERS", "broker1,broker2")
	assert.NoError(t, err)
	err = os.Setenv("DB_DSN", "postgres://user:password@localhost:5432/dbname")
	assert.NoError(t, err)
	err = os.Setenv("REDIS_DSN", "redis://localhost:6379")
	assert.NoError(t, err)
	err = os.Setenv("ORDER_CLI_CONN", "order-service:50051")
	assert.NoError(t, err)

	defer func() {
		_ = os.Unsetenv("CONSUMER_TOPICS")
		_ = os.Unsetenv("CONSUMER_BROKERS")
		_ = os.Unsetenv("DB_DSN")
		_ = os.Unsetenv("REDIS_DSN")
		_ = os.Unsetenv("ORDER_CLI_CONN")
	}()

	cfg, err := New()
	assert.NoError(t, err)

	assert.Equal(t, "DEBUG", cfg.LogLevel)
	assert.Equal(t, 5*time.Second, cfg.ShutdownTimeout)
}

func TestNewConfig_MissingRequiredValues(t *testing.T) {
	_, err := New()

	assert.Error(t, err)
}

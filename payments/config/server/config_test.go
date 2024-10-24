package server

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_Success(t *testing.T) {
	err := os.Setenv("SERVER_GRPC_ADDR", "localhost:50051")
	assert.NoError(t, err)
	err = os.Setenv("CRON_PRODUCER_BROKERS", "broker1,broker2")
	assert.NoError(t, err)
	err = os.Setenv("DB_DSN", "postgres://user:password@localhost:5432/dbname")
	assert.NoError(t, err)
	err = os.Setenv("LOG_LEVEL", "INFO")
	assert.NoError(t, err)

	defer func() {
		_ = os.Unsetenv("SERVER_GRPC_ADDR")
		_ = os.Unsetenv("CRON_PRODUCER_BROKERS")
		_ = os.Unsetenv("DB_DSN")
		_ = os.Unsetenv("LOG_LEVEL")
	}()

	cfg, err := New()
	assert.NoError(t, err)

	assert.Equal(t, "localhost:50051", cfg.Server.GRPCAddr)
	assert.Equal(t, []string{"broker1", "broker2"}, cfg.Cron.Producer.Brokers)
	assert.Equal(t, "postgres://user:password@localhost:5432/dbname", cfg.DB.DSN)
	assert.Equal(t, "INFO", cfg.LogLevel)
	assert.Equal(t, 5*time.Second, cfg.Server.ShutdownTimeout)
	assert.Equal(t, 5*time.Second, cfg.Cron.Frequency)
	assert.Equal(t, 2, cfg.Cron.RateLimit)
	assert.Equal(t, 20, cfg.Cron.EventsLimit)
	assert.Equal(t, 20*time.Second, cfg.Cron.Reserve)
}

func TestNewConfig_DefaultValues(t *testing.T) {
	err := os.Setenv("SERVER_GRPC_ADDR", "localhost:50051")
	assert.NoError(t, err)
	err = os.Setenv("CRON_PRODUCER_BROKERS", "broker1,broker2")
	assert.NoError(t, err)
	err = os.Setenv("DB_DSN", "postgres://user:password@localhost:5432/dbname")
	assert.NoError(t, err)

	defer func() {
		_ = os.Unsetenv("SERVER_GRPC_ADDR")
		_ = os.Unsetenv("CRON_PRODUCER_BROKERS")
		_ = os.Unsetenv("DB_DSN")
	}()

	cfg, err := New()
	assert.NoError(t, err)

	assert.Equal(t, "DEBUG", cfg.LogLevel)
	assert.Equal(t, 5*time.Second, cfg.Server.ShutdownTimeout)
	assert.Equal(t, 5*time.Second, cfg.Cron.Frequency)
	assert.Equal(t, 2, cfg.Cron.RateLimit)
	assert.Equal(t, 20, cfg.Cron.EventsLimit)
	assert.Equal(t, 20*time.Second, cfg.Cron.Reserve)
}

func TestNewConfig_MissingRequiredValues(t *testing.T) {
	_, err := New()

	assert.Error(t, err)
}

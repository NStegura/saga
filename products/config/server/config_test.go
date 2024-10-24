package server

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_Success(t *testing.T) {
	err := os.Setenv("SERVER_GRPC_ADDR", ":8080")
	assert.NoError(t, err)
	err = os.Setenv("DB_DSN", "postgres://user:password@localhost:5432/dbname")
	assert.NoError(t, err)
	err = os.Setenv("PRODUCER_BROKERS", "broker1,broker2")
	assert.NoError(t, err)
	err = os.Setenv("CRON_FREQUENCY", "10s")
	assert.NoError(t, err)
	err = os.Setenv("CRON_RATE_LIMIT", "5")
	assert.NoError(t, err)
	err = os.Setenv("CRON_EVENTS_LIMIT", "50")
	assert.NoError(t, err)
	err = os.Setenv("CRON_RESERVE", "30s")
	assert.NoError(t, err)
	err = os.Setenv("LOG_LEVEL", "INFO")
	assert.NoError(t, err)

	defer func() {
		_ = os.Unsetenv("SERVER_GRPC_ADDR")
		_ = os.Unsetenv("DB_DSN")
		_ = os.Unsetenv("PRODUCER_BROKERS")
		_ = os.Unsetenv("CRON_FREQUENCY")
		_ = os.Unsetenv("CRON_RATE_LIMIT")
		_ = os.Unsetenv("CRON_EVENTS_LIMIT")
		_ = os.Unsetenv("CRON_RESERVE")
		_ = os.Unsetenv("LOG_LEVEL")
	}()

	cfg, err := New()
	assert.NoError(t, err)

	assert.Equal(t, ":8080", cfg.Server.GRPCAddr)
	assert.Equal(t, "postgres://user:password@localhost:5432/dbname", cfg.DB.DSN)
	assert.Equal(t, []string{"broker1", "broker2"}, cfg.Cron.Producer.Brokers)
	assert.Equal(t, 10*time.Second, cfg.Cron.Frequency)
	assert.Equal(t, 5, cfg.Cron.RateLimit)
	assert.Equal(t, 50, cfg.Cron.EventsLimit)
	assert.Equal(t, 30*time.Second, cfg.Cron.Reserve)
	assert.Equal(t, "INFO", cfg.LogLevel)
}

func TestNewConfig_DefaultValues(t *testing.T) {
	err := os.Setenv("SERVER_GRPC_ADDR", ":8080")
	assert.NoError(t, err)
	err = os.Setenv("DB_DSN", "postgres://user:password@localhost:5432/dbname")
	assert.NoError(t, err)
	err = os.Setenv("PRODUCER_BROKERS", "broker1,broker2")
	assert.NoError(t, err)

	defer func() {
		_ = os.Unsetenv("SERVER_GRPC_ADDR")
		_ = os.Unsetenv("DB_DSN")
		_ = os.Unsetenv("PRODUCER_BROKERS")
	}()

	cfg, err := New()
	assert.NoError(t, err)

	assert.Equal(t, 5*time.Second, cfg.Cron.Frequency)
	assert.Equal(t, 2, cfg.Cron.RateLimit)
	assert.Equal(t, 20, cfg.Cron.EventsLimit)
	assert.Equal(t, 20*time.Second, cfg.Cron.Reserve)
	assert.Equal(t, "DEBUG", cfg.LogLevel)
	assert.Equal(t, 5*time.Second, cfg.Server.ShutdownTimeout)
}

func TestNewConfig_MissingRequiredValues(t *testing.T) {
	_, err := New()

	assert.Error(t, err)
}

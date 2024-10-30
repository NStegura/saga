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
	err = os.Setenv("SERVER_SHUTDOWN_TIMEOUT", "10s")
	assert.NoError(t, err)
	err = os.Setenv("DB_DSN", "postgres://user:password@localhost:5432/dbname")
	assert.NoError(t, err)
	err = os.Setenv("CRON_PRODUCER_BROKERS", "broker1,broker2")
	assert.NoError(t, err)
	err = os.Setenv("CRON_FREQUENCY", "5s")
	assert.NoError(t, err)
	err = os.Setenv("CRON_RATE_LIMIT", "3")
	assert.NoError(t, err)
	err = os.Setenv("CRON_EVENTS_LIMIT", "30")
	assert.NoError(t, err)
	err = os.Setenv("CRON_RESERVE", "15s")
	assert.NoError(t, err)

	defer func() {
		_ = os.Unsetenv("SERVER_GRPC_ADDR")
		_ = os.Unsetenv("SERVER_SHUTDOWN_TIMEOUT")
		_ = os.Unsetenv("DB_DSN")
		_ = os.Unsetenv("CRON_PRODUCER_BROKERS")
		_ = os.Unsetenv("CRON_FREQUENCY")
		_ = os.Unsetenv("CRON_RATE_LIMIT")
		_ = os.Unsetenv("CRON_EVENTS_LIMIT")
		_ = os.Unsetenv("CRON_RESERVE")
	}()

	cfg, err := New()
	assert.NoError(t, err)

	assert.Equal(t, "localhost:50051", cfg.Server.GRPCAddr)
	assert.Equal(t, 10*time.Second, cfg.Server.ShutdownTimeout)
	assert.Equal(t, "postgres://user:password@localhost:5432/dbname", cfg.DB.DSN)
	assert.ElementsMatch(t, []string{"broker1", "broker2"}, cfg.Cron.Producer.Brokers)
	assert.Equal(t, 5*time.Second, cfg.Cron.Frequency)
	assert.Equal(t, 3, cfg.Cron.RateLimit)
	assert.Equal(t, 30, cfg.Cron.EventsLimit)
	assert.Equal(t, 15*time.Second, cfg.Cron.Reserve)
}

func TestNewConfig_DefaultValues(t *testing.T) {
	err := os.Setenv("SERVER_GRPC_ADDR", "localhost:50051")
	assert.NoError(t, err)
	err = os.Setenv("DB_DSN", "postgres://user:password@localhost:5432/dbname")
	assert.NoError(t, err)
	err = os.Setenv("CRON_PRODUCER_BROKERS", "broker1,broker2")
	assert.NoError(t, err)

	defer func() {
		_ = os.Unsetenv("SERVER_GRPC_ADDR")
		_ = os.Unsetenv("DB_DSN")
		_ = os.Unsetenv("CRON_PRODUCER_BROKERS")
	}()

	cfg, err := New()
	assert.NoError(t, err)

	assert.Equal(t, "localhost:50051", cfg.Server.GRPCAddr)
	assert.Equal(t, 5*time.Second, cfg.Server.ShutdownTimeout)
	assert.Equal(t, "postgres://user:password@localhost:5432/dbname", cfg.DB.DSN)
	assert.ElementsMatch(t, []string{"broker1", "broker2"}, cfg.Cron.Producer.Brokers)
	assert.Equal(t, 5*time.Second, cfg.Cron.Frequency)
	assert.Equal(t, 2, cfg.Cron.RateLimit)
	assert.Equal(t, 20, cfg.Cron.EventsLimit)
	assert.Equal(t, 20*time.Second, cfg.Cron.Reserve)
}

func TestNewConfig_MissingRequiredValues(t *testing.T) {
	_, err := New()
	assert.Error(t, err)
}

func TestNewConfig_MissingProducerBrokers(t *testing.T) {
	err := os.Setenv("SERVER_GRPC_ADDR", "localhost:50051")
	assert.NoError(t, err)
	err = os.Setenv("DB_DSN", "postgres://user:password@localhost:5432/dbname")
	assert.NoError(t, err)

	defer func() {
		_ = os.Unsetenv("SERVER_GRPC_ADDR")
		_ = os.Unsetenv("DB_DSN")
	}()

	_, err = New()
	assert.Error(t, err)
}

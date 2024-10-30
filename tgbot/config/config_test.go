package config

import (
	"os"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig_Success(t *testing.T) {
	err := os.Setenv("TG_BOT_TOKEN", "test-token")
	assert.NoError(t, err)
	err = os.Setenv("REDIS_DSN", "redis://localhost:6379")
	assert.NoError(t, err)
	err = os.Setenv("ORDER_CLI_CONN", "localhost:8081")
	assert.NoError(t, err)
	err = os.Setenv("PAYMENT_CLI_CONN", "localhost:8083")
	assert.NoError(t, err)
	err = os.Setenv("PRODUCT_CLI_CONN", "localhost:8082")
	assert.NoError(t, err)
	err = os.Setenv("LOG_LEVEL", "INFO")
	assert.NoError(t, err)
	err = os.Setenv("SHUTDOWN_TIMEOUT", "10s")
	assert.NoError(t, err)

	defer func() {
		_ = os.Unsetenv("TG_BOT_TOKEN")
		_ = os.Unsetenv("REDIS_DSN")
		_ = os.Unsetenv("ORDER_CLI_CONN")
		_ = os.Unsetenv("PAYMENT_CLI_CONN")
		_ = os.Unsetenv("PRODUCT_CLI_CONN")
		_ = os.Unsetenv("LOG_LEVEL")
		_ = os.Unsetenv("SHUTDOWN_TIMEOUT")
	}()

	cfg, err := New()
	err = envconfig.Process("", &cfg)
	assert.NoError(t, err)

	assert.Equal(t, "test-token", cfg.TgBot.Token)
	assert.Equal(t, "redis://localhost:6379", cfg.Redis.DSN)
	assert.Equal(t, "localhost:8081", cfg.OrderCli.CONN)
	assert.Equal(t, "localhost:8083", cfg.PaymentCli.CONN)
	assert.Equal(t, "localhost:8082", cfg.ProductCli.CONN)
	assert.Equal(t, "INFO", cfg.LogLevel)
	assert.Equal(t, 10*time.Second, cfg.ShutdownTimeout)
}

func TestNewConfig_DefaultValues(t *testing.T) {
	err := os.Setenv("TG_BOT_TOKEN", "test-token")
	assert.NoError(t, err)

	defer func() {
		_ = os.Unsetenv("TG_BOT_TOKEN")
	}()

	cfg, err := New()
	err = envconfig.Process("", &cfg)
	assert.NoError(t, err)

	assert.Equal(t, "test-token", cfg.TgBot.Token)
	assert.Equal(t, "0.0.0.0:6379", cfg.Redis.DSN)
	assert.Equal(t, "0.0.0.0:8081", cfg.OrderCli.CONN)
	assert.Equal(t, "0.0.0.0:8083", cfg.PaymentCli.CONN)
	assert.Equal(t, "0.0.0.0:8082", cfg.ProductCli.CONN)
	assert.Equal(t, "DEBUG", cfg.LogLevel)
	assert.Equal(t, 5*time.Second, cfg.ShutdownTimeout)
}

func TestNewConfig_MissingRequiredValues(t *testing.T) {
	_, err := New()

	assert.Error(t, err, "должна вернуться ошибка при отсутствии обязательных переменных окружения")
}

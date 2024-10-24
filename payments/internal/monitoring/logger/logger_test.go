package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInit_Success(t *testing.T) {
	logger, err := Init("INFO")
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.Equal(t, logrus.InfoLevel, logger.Level)
}

func TestInit_InvalidLogLevel(t *testing.T) {
	logger, err := Init("INVALID")
	assert.Error(t, err)
	assert.Nil(t, logger)
}

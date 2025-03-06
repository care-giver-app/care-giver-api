package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetLogger(t *testing.T) {
	logger := GetLogger(InfoLevel)
	assert.Equal(t, logger.Level(), zap.InfoLevel)
}

func TestGetLoggerWithEnv(t *testing.T) {
	logger := GetLoggerWithEnv(InfoLevel, "test")
	assert.Equal(t, logger.Level(), zap.InfoLevel)
}

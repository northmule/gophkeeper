package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger_Success(t *testing.T) {
	logger, err := NewLogger("info")
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.SugaredLogger)
}

func TestNewLogger_InvalidLevel(t *testing.T) {
	logger, err := NewLogger("invalid_level")
	assert.Error(t, err)
	assert.Nil(t, logger)
}

func TestNewLogger_DebugLevel(t *testing.T) {
	logger, err := NewLogger("debug")
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.SugaredLogger)
	assert.Equal(t, zapcore.DebugLevel, logger.Level())
}

func TestNewLogger_WarnLevel(t *testing.T) {
	logger, err := NewLogger("warn")
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.SugaredLogger)
	assert.Equal(t, zapcore.WarnLevel, logger.Level())
}

func TestNewLogger_ErrorLevel(t *testing.T) {
	logger, err := NewLogger("error")
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.SugaredLogger)
	assert.Equal(t, zapcore.ErrorLevel, logger.Level())
}

func TestNewLogger_FatalLevel(t *testing.T) {
	logger, err := NewLogger("fatal")
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.SugaredLogger)
	assert.Equal(t, zapcore.FatalLevel, logger.Level())
}

func TestNewLogger_PanicLevel(t *testing.T) {
	logger, err := NewLogger("panic")
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.SugaredLogger)
	assert.Equal(t, zapcore.PanicLevel, logger.Level())
}

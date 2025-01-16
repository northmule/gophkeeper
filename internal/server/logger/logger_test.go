package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	t.Run("Valid log level", func(t *testing.T) {
		logger, err := NewLogger("info")
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("Invalid log level", func(t *testing.T) {
		logger, err := NewLogger("invalid")
		assert.Error(t, err)
		assert.Nil(t, logger)
	})

}

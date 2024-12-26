package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_HandleGetHealth_Success(t *testing.T) {
	logger, _ := logger.NewLogger("info")
	handler := NewHealthHandler(logger)

	req, _ := http.NewRequest("GET", "/health", nil)
	res := httptest.NewRecorder()

	handler.HandleGetHealth(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "ok")
}

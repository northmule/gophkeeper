package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/northmule/gophkeeper/internal/server/logger"
	appMock "github.com/northmule/gophkeeper/internal/server/repository/mock"
	"github.com/stretchr/testify/assert"
)

func TestDefiningAppRoutes(t *testing.T) {

	mockRepository := new(appMock.MockManager)
	mockStorage := new(appMock.MockDBQuery)
	mockSessionStorage := new(appMock.MockSessionManager)
	mockAccessService := new(appMock.MockAccessService)
	mockCryptService := new(appMock.MockCryptService)

	l, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()
	cfg.Value().PathKeys = t.TempDir()

	appRoutes := NewAppRoutes(mockRepository, mockStorage, mockSessionStorage, l, cfg, mockAccessService, mockCryptService)

	jwt := new(jwtauth.JWTAuth)
	mockAccessService.On("FillJWTToken").Return(jwt)
	router := appRoutes.DefiningAppRoutes()

	assert.NotNil(t, router)

	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

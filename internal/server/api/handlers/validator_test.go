package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/stretchr/testify/assert"
)

func TestValidatorHandler_HandleValidation_Success(t *testing.T) {
	l, _ := logger.NewLogger("info")
	reqBody := `{"name": "test", "value": "test value"}`
	req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	next := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

	})

	handler := NewValidatorHandler(&textDataRequest{}, l)

	rr := httptest.NewRecorder()

	handler.HandleValidation(next).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestValidatorHandler_HandleValidation_ValidationErrors(t *testing.T) {
	l, _ := logger.NewLogger("info")

	next := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

	})
	t.Run("textDataRequest", func(t *testing.T) {
		reqBody := `{"name": "", "value": "test value"}`
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		handler := NewValidatorHandler(new(textDataRequest), l)
		rr := httptest.NewRecorder()
		handler.HandleValidation(next).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("registrationRequest", func(t *testing.T) {
		reqBody := `{"login": ""}`
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		handler := NewValidatorHandler(new(registrationRequest), l)
		rr := httptest.NewRecorder()
		handler.HandleValidation(next).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("authenticationRequest", func(t *testing.T) {
		reqBody := `{"name": "", "login": "1"}`
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		handler := NewValidatorHandler(new(authenticationRequest), l)
		rr := httptest.NewRecorder()
		handler.HandleValidation(next).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("cardDataRequest", func(t *testing.T) {
		reqBody := `{"name": ""}`
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		handler := NewValidatorHandler(new(cardDataRequest), l)
		rr := httptest.NewRecorder()
		handler.HandleValidation(next).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("fileDataInitRequest", func(t *testing.T) {
		reqBody := `{"name": ""}`
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		handler := NewValidatorHandler(new(fileDataInitRequest), l)
		rr := httptest.NewRecorder()
		handler.HandleValidation(next).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

}

func TestValidatorHandler_HandleValidation_ValidationSuccess(t *testing.T) {
	l, _ := logger.NewLogger("info")

	next := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

	})
	t.Run("textDataRequest", func(t *testing.T) {
		reqBody := `{"name": "test", "value": "test value"}`
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		handler := NewValidatorHandler(new(textDataRequest), l)
		rr := httptest.NewRecorder()
		handler.HandleValidation(next).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
	t.Run("registrationRequest", func(t *testing.T) {
		reqBody := `{"login": "123456", "password":"1123132", "email":"admin@mail.ru"}`
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		handler := NewValidatorHandler(new(registrationRequest), l)
		rr := httptest.NewRecorder()
		handler.HandleValidation(next).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
	t.Run("authenticationRequest", func(t *testing.T) {
		reqBody := `{"password": "1234567890", "login": "12341215"}`
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		handler := NewValidatorHandler(new(authenticationRequest), l)
		rr := httptest.NewRecorder()
		handler.HandleValidation(next).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
	t.Run("cardDataRequest", func(t *testing.T) {
		reqBody := `{"name": "112213"}`
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		handler := NewValidatorHandler(new(cardDataRequest), l)
		rr := httptest.NewRecorder()
		handler.HandleValidation(next).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
	t.Run("fileDataInitRequest", func(t *testing.T) {
		reqBody := `{"name": "232323", "extension":".123", "file_name":"123"}`
		req, _ := http.NewRequest("POST", "/test", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		handler := NewValidatorHandler(new(fileDataInitRequest), l)
		rr := httptest.NewRecorder()
		handler.HandleValidation(next).ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

}

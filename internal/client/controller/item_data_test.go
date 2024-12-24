package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/stretchr/testify/assert"
)

func TestItemData_Send_Success(t *testing.T) {
	cryptService := NewCryptMock(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !strings.Contains(r.URL.Path, "dataUUID") {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		encryptedData, _ := cryptService.EncryptAES([]byte(`{"Data":{}}`))
		w.WriteHeader(http.StatusOK)
		w.Write(encryptedData)
	}))
	defer server.Close()

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}

	mockConfig := makeMockConfig(server.URL)
	controller := NewItemData(mockConfig, cryptService, log)

	response, err := controller.Send("token", "dataUUID")
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestItemData_Send_Unauthorized(t *testing.T) {
	cryptService := NewCryptMock(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}

	mockConfig := makeMockConfig(server.URL)
	controller := NewItemData(mockConfig, cryptService, log)
	_, err = controller.Send("invalid_token", "dataUUID")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "вы не авторизованы")
}

func TestItemData_Send_BadRequest(t *testing.T) {
	cryptService := NewCryptMock(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}))
	defer server.Close()

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}

	mockConfig := makeMockConfig(server.URL)
	controller := NewItemData(mockConfig, cryptService, log)

	_, err = controller.Send("token", "invalid_dataUUID")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка в запросе")
}

func TestItemData_Send_UnknownError(t *testing.T) {
	cryptService := NewCryptMock(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}

	mockConfig := makeMockConfig(server.URL)
	controller := NewItemData(mockConfig, cryptService, log)

	_, err = controller.Send("token", "dataUUID")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "не известная ошибка")
}

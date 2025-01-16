package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/northmule/gophkeeper/internal/client/logger"
)

func TestGridDataSend(t *testing.T) {
	cryptService := NewCryptMock(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/api/v1/items_list") {
			t.Errorf("Expected path /api/v1/items_list, got %s", r.URL.Path)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer validtoken" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		rawBody, _ := cryptService.EncryptAES([]byte(`{"items":[{"number": "123"}]}`))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(rawBody)

	}))

	defer server.Close()

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}

	mockConfig := makeMockConfig(server.URL)
	controller := NewGridData(mockConfig, cryptService, log)

	t.Run("ok", func(t *testing.T) {
		_, err := controller.Send("validtoken")
		if err != nil {
			t.Errorf("Send failed: %v", err)
		}

	})

	t.Run("no_validtoken", func(t *testing.T) {
		_, err := controller.Send("no_validtoken")
		if err == nil || !strings.Contains(err.Error(), "вы не авторизованы") {
			t.Errorf("Send should have failed with unknown error: %v", err)
		}
	})
}

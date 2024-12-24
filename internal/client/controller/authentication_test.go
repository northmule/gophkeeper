package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/northmule/gophkeeper/internal/client/logger"
)

func TestAuthenticationSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
			return
		}

		if r.URL.Path != "/api/v1/login" {
			t.Errorf("Expected path /api/v1/login, got %s", r.URL.Path)
			return
		}

		var requestData authenticationRequest
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			return
		}

		if requestData.Login == "validuser" && requestData.Password == "validpass" {
			w.Header().Set("Authorization", "Bearer validtoken")
			w.WriteHeader(http.StatusOK)
			return
		}

		if requestData.Login == "unauthorizeduser" && requestData.Password == "unauthorizedpass" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}

	mockConfig := makeMockConfig(server.URL)

	authController := NewAuthentication(mockConfig, log)

	response, err := authController.Send("validuser", "validpass")
	if err != nil {
		t.Errorf("Send failed with valid credentials: %v", err)
	}

	if response.Value != "validtoken" {
		t.Errorf("Expected token 'validtoken', got '%s'", response.Value)
	}

	_, err = authController.Send("unauthorizeduser", "unauthorizedpass")
	if err == nil || !strings.Contains(err.Error(), "не верная пара логин/пароль") {
		t.Errorf("Send should have failed with unauthorized credentials: %v", err)
	}

	_, err = authController.Send("unknownuser", "unknownpass")
	if err == nil || !strings.Contains(err.Error(), "не известная ошибка") {
		t.Errorf("Send should have failed with unknown error: %v", err)
	}

	_, err = authController.Send("", "validpass")
	if err == nil {
		t.Errorf("Send should have failed with empty login")
	}

	_, err = authController.Send("validuser", "")
	if err == nil {
		t.Errorf("Send should have failed with empty password")
	}

}

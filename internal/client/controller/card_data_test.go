package controller

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"github.com/northmule/gophkeeper/internal/common/util"
)

// CryptMock мокк
type CryptMock struct {
	pubKey     *rsa.PublicKey
	privateKey *rsa.PrivateKey
	aesKey     []byte
}

// NewCryptMock мокк
func NewCryptMock(t *testing.T) *CryptMock {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}
	return &CryptMock{
		pubKey:     &key.PublicKey,
		privateKey: key,
		aesKey:     make([]byte, 32),
	}
}

// EncryptRSA мокк
func (c *CryptMock) EncryptRSA(data []byte) ([]byte, error) {
	return util.DataEncryptRSA(data, c.pubKey)
}

// DecryptRSA мокк
func (c *CryptMock) DecryptRSA(data []byte) ([]byte, error) {
	return util.DataDecryptRSA(data, c.privateKey)
}

// EncryptAES мокк
func (c *CryptMock) EncryptAES(data []byte) ([]byte, error) {
	return util.DataEncryptAES(data, c.aesKey)
}

// DecryptAES мокк
func (c *CryptMock) DecryptAES(data []byte) ([]byte, error) {
	return util.DataDecryptAES(data, c.aesKey)
}

func TestCardDataSend(t *testing.T) {
	cryptService := NewCryptMock(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
			return
		}

		if r.URL.Path != "/api/v1/save_card_data" {
			t.Errorf("Expected path /api/v1/save_card_data, got %s", r.URL.Path)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer validtoken" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		bodyBytes, _ := io.ReadAll(r.Body)
		rawBody, _ := cryptService.DecryptAES(bodyBytes)
		buf := bytes.NewBuffer(rawBody)

		var requestData model_data.CardDataRequest
		if err := json.NewDecoder(buf).Decode(&requestData); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			return
		}

		if requestData.CardNumber == "1234567890123456" && requestData.FullNameHolder == "John Doe" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if requestData.CardNumber == "badrequest" {
			w.WriteHeader(http.StatusBadRequest)
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

	cardDataController := NewCardData(mockConfig, cryptService, log)

	t.Run("ok", func(t *testing.T) {
		requestData := &model_data.CardDataRequest{
			CardNumber:     "1234567890123456",
			FullNameHolder: "John Doe",
		}

		response, err := cardDataController.Send("validtoken", requestData)
		if err != nil {
			t.Errorf("Send failed with valid credentials: %v", err)
		}

		if response.Value != "ok" {
			t.Errorf("Expected response 'ok', got '%s'", response.Value)
		}
	})

	t.Run("badrequest", func(t *testing.T) {
		requestData := &model_data.CardDataRequest{
			CardNumber: "badrequest",
		}
		response, err := cardDataController.Send("validtoken", requestData)
		if err == nil || !strings.Contains(err.Error(), "ошибка в запросе") {
			t.Errorf("Send should have failed with unknown error: %v", err)
		}
		_ = response
	})

	t.Run("no_validtoken", func(t *testing.T) {
		requestData := &model_data.CardDataRequest{}
		response, err := cardDataController.Send("no_validtoken", requestData)
		if err == nil || !strings.Contains(err.Error(), "вы не авторизованы") {
			t.Errorf("Send should have failed with unknown error: %v", err)
		}
		_ = response
	})

}

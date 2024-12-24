package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/common/model_data"
)

func makeMockConfig(server string) *config.Config {
	validConfig := `
ServerAddress: "{ServerAddress}"
LogLevel: "info"
FilePath: "/tmp"
PathKeys: "/tmp"
PathPublicKeyServer: "/tmp"
OverwriteKeys: false
`
	validConfig = strings.Replace(validConfig, "{ServerAddress}", server, 1)
	_ = os.WriteFile("client.yaml", []byte(validConfig), 0644)
	defer os.Remove("client.yaml")

	mockConfig := config.NewConfig()
	_ = mockConfig.Init()

	return mockConfig
}

func TestFileDataSend(t *testing.T) {
	cryptService := NewCryptMock(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
			return
		}

		if r.URL.Path != "/api/v1/file_data/init" {
			t.Errorf("Expected path /api/v1/file_data/init, got %s", r.URL.Path)
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

		var requestData model_data.FileDataInitRequest
		if err := json.NewDecoder(buf).Decode(&requestData); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			return
		}

		if requestData.UUID == "1234567890123456" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"upload_path":"http://ok"}`))
			return
		}

		if requestData.UUID == "badrequest" {
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
	controller := NewFileData(mockConfig, cryptService, log)

	t.Run("ok", func(t *testing.T) {
		requestData := &model_data.FileDataInitRequest{
			UUID: "1234567890123456",
		}

		response, err := controller.Send("validtoken", requestData)
		if err != nil {
			t.Errorf("Send failed: %v", err)
		}

		if response.UploadPath != "http://ok" {
			t.Errorf("Expected response 'ok', got '%s'", response.UploadPath)
		}
	})

	t.Run("badrequest", func(t *testing.T) {
		requestData := &model_data.FileDataInitRequest{
			UUID: "badrequest",
		}
		response, err := controller.Send("validtoken", requestData)
		if err == nil || !strings.Contains(err.Error(), "ошибка в запросе") {
			t.Errorf("Send should have failed with unknown error: %v", err)
		}
		_ = response
	})

	t.Run("no_validtoken", func(t *testing.T) {
		requestData := &model_data.FileDataInitRequest{}
		response, err := controller.Send("no_validtoken", requestData)
		if err == nil || !strings.Contains(err.Error(), "вы не авторизованы") {
			t.Errorf("Send failed: %v", err)
		}
		_ = response
	})

}

func TestFileDataUploadFile(t *testing.T) {
	cryptService := NewCryptMock(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/api/v1") {
			t.Errorf("Expected path starting with /api/v1, got %s", r.URL.Path)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer validtoken" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			t.Errorf("Expected Content-Type starting with multipart/form-data, got %s", contentType)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}

	mockConfig := makeMockConfig(server.URL)
	controller := NewFileData(mockConfig, cryptService, log)

	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("Hello")
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	tempFile.Close()

	t.Run("ok", func(t *testing.T) {
		err = controller.UploadFile("validtoken", "/upload", tempFile)
		if err != nil {
			t.Errorf("Send failed: %v", err)
		}

	})

	t.Run("no_validtoken", func(t *testing.T) {
		err = controller.UploadFile("no_validtoken", "/upload", tempFile)
		if err == nil {
			t.Errorf("Send failed: %v", err)
		}
	})

}

func TestFileDataDownLoadFile(t *testing.T) {
	cryptService := NewCryptMock(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/api/v1/file_data/get") {
			t.Errorf("Expected path /api/v1/file_data/get, got %s", r.URL.Path)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer validtoken" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("вы не авторизованы"))
			return
		}

		if strings.Contains(r.URL.Path, "no_valid_file") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ошибка в запросе"))
			return
		}
		if strings.Contains(r.URL.Path, "valid_file_no_ok") {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("bad"))
			return
		}

		if strings.Contains(r.URL.Path, "valid_file_1123") {
			rawBody, _ := cryptService.EncryptAES([]byte("file_data"))
			_, _ = w.Write(rawBody)
			return
		}

		bodyBytes, _ := io.ReadAll(r.Body)
		rawBody, _ := cryptService.DecryptAES(bodyBytes)
		buf := bytes.NewBuffer(rawBody)

		var requestData model_data.FileDataInitRequest
		if err := json.NewDecoder(buf).Decode(&requestData); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			return
		}

		if requestData.UUID == "1234567890123456" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"upload_path":"http://ok"}`))
			return
		}

		if requestData.UUID == "badrequest" {
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
	controller := NewFileData(mockConfig, cryptService, log)

	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("Hello")
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	tempFile.Close()

	t.Run("ok", func(t *testing.T) {
		err = controller.DownLoadFile("validtoken", "file_name", "valid_file_1123")
		if err != nil {
			t.Errorf("Send failed: %v", err)
		}
	})

	t.Run("badrequest", func(t *testing.T) {
		err = controller.DownLoadFile("validtoken", "file_name", "no_valid_file")
		if err == nil || !strings.Contains(err.Error(), "ошибка в запросе") {
			t.Errorf("Send should have failed with unknown error: %v", err)
		}

	})

	t.Run("no_validtoken", func(t *testing.T) {
		err = controller.DownLoadFile("no_validtoken", "file_name", "valid_file_1123")
		if err == nil || !strings.Contains(err.Error(), "вы не авторизованы") {
			t.Errorf("Send failed: %v", err)
		}
	})

	t.Run("no_ok", func(t *testing.T) {
		err = controller.DownLoadFile("validtoken", "file_name", "valid_file_no_ok")
		if err == nil || !strings.Contains(err.Error(), "bad") {
			t.Errorf("Send failed: %v", err)
		}
	})

}

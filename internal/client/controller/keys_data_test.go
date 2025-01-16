package controller

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/keys"
	"github.com/stretchr/testify/assert"
)

func TestKeysData_UploadClientPublicKey_Success(t *testing.T) {
	tempFile, err := os.Create(path.Join("", keys.PublicKeyFileName))
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("public_key_data")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	cryptService := NewCryptMock(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") == "" {
			http.Error(w, "Content-Type missing", http.StatusBadRequest)
			return
		}
		_, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
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
	controller := NewKeysData(mockConfig, cryptService, log)

	err = controller.UploadClientPublicKey("token")
	assert.NoError(t, err)
}

func TestKeysData_UploadClientPublicKey_Unauthorized(t *testing.T) {

	tempFile, err := os.Create(path.Join("", keys.PublicKeyFileName))
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("public_key_data")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	mockConfig := makeMockConfig(server.URL)
	controller := NewKeysData(mockConfig, cryptService, log)

	err = controller.UploadClientPublicKey("invalid_token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "вы не авторизованы")
}

func TestKeysData_UploadClientPublicKey_FileOpenError(t *testing.T) {
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	mockConfig := makeMockConfig("http://uuu.loc")
	controller := NewKeysData(mockConfig, cryptService, log)

	err = controller.UploadClientPublicKey("token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestKeysData_UploadClientPublicKey_MultipartFormError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	mockConfig := makeMockConfig(server.URL)
	controller := NewKeysData(mockConfig, cryptService, log)

	tempFile, err := os.Create(path.Join("", keys.PublicKeyFileName))
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("public_key_data")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	err = controller.UploadClientPublicKey("token")
	assert.Error(t, err)
}

func TestKeysData_DownloadPublicServerKey_Success(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("server_public_key_data"))
	}))
	defer server.Close()

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	mockConfig := makeMockConfig(server.URL)
	controller := NewKeysData(mockConfig, cryptService, log)

	err = controller.DownloadPublicServerKey("token")
	assert.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(keys.PublicKeyFileName))
	assert.NoError(t, err)
	assert.Equal(t, "server_public_key_data", string(data))
}

func TestKeysData_DownloadPublicServerKey_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	mockConfig := makeMockConfig(server.URL)
	controller := NewKeysData(mockConfig, cryptService, log)

	err = controller.DownloadPublicServerKey("invalid_token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unauthorized")
}

func TestKeysData_DownloadPublicServerKey_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}))
	defer server.Close()

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	mockConfig := makeMockConfig(server.URL)
	controller := NewKeysData(mockConfig, cryptService, log)

	err = controller.DownloadPublicServerKey("token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Bad Request")
}

func TestKeysData_DownloadPublicServerKey_UnknownError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	mockConfig := makeMockConfig(server.URL)
	controller := NewKeysData(mockConfig, cryptService, log)

	err = controller.DownloadPublicServerKey("token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Internal Server Error")
}

func TestUploadClientPrivateKey_Success(t *testing.T) {
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/save_client_private_key", r.URL.Path)
		assert.Equal(t, "Bearer test_token", r.Header.Get("Authorization"))

		err := r.ParseMultipartForm(32 << 20)
		assert.NoError(t, err)
		file, _, err := r.FormFile(data_type.FileField)
		assert.NoError(t, err)
		defer file.Close()

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, file)
		assert.NoError(t, err)
		assert.NotEmpty(t, buf.String())

		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	cryptService := NewCryptMock(t)
	mockConfig := makeMockConfig(testServer.URL)
	controller := NewKeysData(mockConfig, cryptService, log)

	err = controller.UploadClientPrivateKey("test_token")
	assert.NoError(t, err)
}

func TestUploadClientPrivateKey_EmptyToken(t *testing.T) {
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	mockConfig := makeMockConfig("")
	controller := NewKeysData(mockConfig, cryptService, log)

	err = controller.UploadClientPrivateKey("")
	assert.Error(t, err)
}

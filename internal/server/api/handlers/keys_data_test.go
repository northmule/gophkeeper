package handlers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/keys"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/northmule/gophkeeper/internal/server/logger"
	appMock "github.com/northmule/gophkeeper/internal/server/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestKeysDataHandler_HandleSaveClientPublicKey_Successful(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockUserRepository := new(appMock.MockUserDataModelRepository)
	mockCryptService := new(appMock.MockCryptService)
	l, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()
	cfg.Value().PathKeys = t.TempDir()

	mockRepository.On("User").Return(mockUserRepository)
	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user123", nil)
	mockUserRepository.On("SetPublicKey", mock.Anything, "publicKey", "user123").Return(nil)

	handler := NewKeysDataHandler(mockAccessService, mockCryptService, mockRepository, cfg, l)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(data_type.FileField, "publicKey.pem")
	io.WriteString(part, "publicKey")
	writer.Close()
	req := httptest.NewRequest("POST", "/keys/public", body)
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler.HandleSaveClientPublicKey(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAccessService.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestKeysDataHandler_HandleDownloadServerPublicKey_Successful(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockUserRepository := new(appMock.MockUserDataModelRepository)
	mockCryptService := new(appMock.MockCryptService)
	l, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()
	cfg.Value().PathKeys = t.TempDir()

	mockUserRepository.On("SetPublicKey", mock.Anything, "publicKey", "user123").Return(nil)

	publicKeyPath := filepath.Join(cfg.Value().PathKeys, keys.PublicKeyFileName)
	os.WriteFile(publicKeyPath, []byte("serverPublicKey"), 0644)

	user := new(models.User)
	user.UUID = "user123"

	mockRepository.On("User").Return(mockUserRepository)
	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user123", nil)
	mockUserRepository.On("FindOneByUUID", mock.Anything, "user123").Return(user, nil)

	handler := NewKeysDataHandler(mockAccessService, mockCryptService, mockRepository, cfg, l)

	req := httptest.NewRequest("GET", "/keys/public", nil)
	rr := httptest.NewRecorder()

	handler.HandleDownloadServerPublicKey(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "serverPublicKey", rr.Body.String())
	mockAccessService.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestKeysDataHandler_HandleDownloadServerPublicKey_InvalidJWTToken(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockCryptService := new(appMock.MockCryptService)
	l, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()
	cfg.Value().PathKeys = t.TempDir()

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("", fmt.Errorf("invalid token"))

	handler := NewKeysDataHandler(mockAccessService, mockCryptService, mockRepository, cfg, l)

	req := httptest.NewRequest("GET", "/keys/public", nil)
	rr := httptest.NewRecorder()

	handler.HandleDownloadServerPublicKey(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAccessService.AssertExpectations(t)
}

func TestKeysDataHandler_HandleDownloadServerPublicKey_NonExistentUser(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockCryptService := new(appMock.MockCryptService)
	mockUserRepository := new(appMock.MockUserDataModelRepository)

	l, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()
	cfg.Value().PathKeys = t.TempDir()

	mockRepository.On("User").Return(mockUserRepository)
	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user123", nil)
	mockUserRepository.On("FindOneByUUID", mock.Anything, "user123").Return(nil, fmt.Errorf("user not found"))

	handler := NewKeysDataHandler(mockAccessService, mockCryptService, mockRepository, cfg, l)

	req := httptest.NewRequest("GET", "/keys/public", nil)
	rr := httptest.NewRecorder()

	handler.HandleDownloadServerPublicKey(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAccessService.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestKeysDataHandler_HandleSaveClientPrivateKey_Successful(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockCryptService := new(appMock.MockCryptService)
	mockUserRepository := new(appMock.MockUserDataModelRepository)

	l, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()
	cfg.Value().PathKeys = t.TempDir()

	mockRepository.On("User").Return(mockUserRepository)
	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user123", nil)
	mockCryptService.On("DecryptRSA", []byte("encryptedPrivateKey")).Return([]byte("privateKey"), nil)
	mockUserRepository.On("SetPrivateClientKey", mock.Anything, "privateKey", "user123").Return(nil)

	handler := NewKeysDataHandler(mockAccessService, mockCryptService, mockRepository, cfg, l)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(data_type.FileField, "privateKey.pem")
	io.WriteString(part, "encryptedPrivateKey")
	writer.Close()
	req := httptest.NewRequest("POST", "/keys/private", body)
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler.HandleSaveClientPrivateKey(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAccessService.AssertExpectations(t)
	mockCryptService.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestKeysDataHandler_HandleSaveClientPrivateKey_InvalidJWTToken(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockCryptService := new(appMock.MockCryptService)

	l, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()
	cfg.Value().PathKeys = t.TempDir()

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("", fmt.Errorf("invalid token"))

	handler := NewKeysDataHandler(mockAccessService, mockCryptService, mockRepository, cfg, l)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(data_type.FileField, "privateKey.pem")
	io.WriteString(part, "encryptedPrivateKey")
	writer.Close()
	req := httptest.NewRequest("POST", "/keys/private", body)
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler.HandleSaveClientPrivateKey(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAccessService.AssertExpectations(t)
}

func TestKeysDataHandler_HandleSaveClientPrivateKey_InvalidFileUpload(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockCryptService := new(appMock.MockCryptService)

	l, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()
	cfg.Value().PathKeys = t.TempDir()

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user123", nil)

	handler := NewKeysDataHandler(mockAccessService, mockCryptService, mockRepository, cfg, l)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()
	req := httptest.NewRequest("POST", "/keys/private", body)
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler.HandleSaveClientPrivateKey(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockAccessService.AssertExpectations(t)
}

func TestKeysDataHandler_HandleSaveClientPrivateKey_DecryptionError(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockCryptService := new(appMock.MockCryptService)

	l, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()
	cfg.Value().PathKeys = t.TempDir()

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user123", nil)
	mockCryptService.On("DecryptRSA", []byte("encryptedPrivateKey")).Return(nil, fmt.Errorf("decryption failed"))

	handler := NewKeysDataHandler(mockAccessService, mockCryptService, mockRepository, cfg, l)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(data_type.FileField, "privateKey.pem")
	io.WriteString(part, "encryptedPrivateKey")
	writer.Close()
	req := httptest.NewRequest("POST", "/keys/private", body)
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler.HandleSaveClientPrivateKey(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockAccessService.AssertExpectations(t)
	mockCryptService.AssertExpectations(t)
}

func TestKeysDataHandler_HandleSaveClientPrivateKey_RepositoryError(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockCryptService := new(appMock.MockCryptService)
	mockUserRepository := new(appMock.MockUserDataModelRepository)

	l, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()
	cfg.Value().PathKeys = t.TempDir()

	mockRepository.On("User").Return(mockUserRepository)
	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user123", nil)
	mockCryptService.On("DecryptRSA", []byte("encryptedPrivateKey")).Return([]byte("privateKey"), nil)
	mockUserRepository.On("SetPrivateClientKey", mock.Anything, "privateKey", "user123").Return(fmt.Errorf("repository error"))

	handler := NewKeysDataHandler(mockAccessService, mockCryptService, mockRepository, cfg, l)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(data_type.FileField, "privateKey.pem")
	io.WriteString(part, "encryptedPrivateKey")
	writer.Close()
	req := httptest.NewRequest("POST", "/keys/private", body)
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler.HandleSaveClientPrivateKey(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockAccessService.AssertExpectations(t)
	mockCryptService.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestKeysDataHandler_HandleSaveClientPublicKey_InvalidFileUpload(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockCryptService := new(appMock.MockCryptService)

	l, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()
	cfg.Value().PathKeys = t.TempDir()

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user123", nil)

	handler := NewKeysDataHandler(mockAccessService, mockCryptService, mockRepository, cfg, l)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// файла нет
	writer.Close()
	req := httptest.NewRequest("POST", "/keys/public", body)
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler.HandleSaveClientPublicKey(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockAccessService.AssertExpectations(t)
}

func TestKeysDataHandler_HandleSaveClientPublicKey_RepositoryError(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockCryptService := new(appMock.MockCryptService)
	mockUserRepository := new(appMock.MockUserDataModelRepository)

	l, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()
	cfg.Value().PathKeys = t.TempDir()

	mockRepository.On("User").Return(mockUserRepository)
	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user123", nil)
	mockUserRepository.On("SetPublicKey", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("repository error"))

	handler := NewKeysDataHandler(mockAccessService, mockCryptService, mockRepository, cfg, l)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(data_type.FileField, "publicKey.pem")
	io.WriteString(part, "publicKey")
	writer.Close()
	req := httptest.NewRequest("POST", "/keys/public", body)
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	handler.HandleSaveClientPublicKey(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockAccessService.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

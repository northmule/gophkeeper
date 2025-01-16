package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/render"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/common/util"
	"github.com/northmule/gophkeeper/internal/server/logger"
	appMock "github.com/northmule/gophkeeper/internal/server/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleDecryptData_SuccessfulDecryption(t *testing.T) {
	mockRepository := new(appMock.MockManager)
	mockAccessService := new(appMock.MockAccessService)
	mockUserRepository := new(appMock.MockUserDataModelRepository)
	logger, _ := logger.NewLogger("info")

	mockRepository.On("User").Return(mockUserRepository)

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)

	privateKey := string(make([]byte, 32))
	user := new(models.User)
	user.UUID = "userUUID"
	user.PrivateClientKey = privateKey

	mockUserRepository.On("FindOneByUUID", mock.Anything, "userUUID").Return(user, nil)
	encryptedData, _ := util.DataEncryptAES([]byte("test_data"), []byte(privateKey))

	req := httptest.NewRequest("POST", "/decrypt", bytes.NewBuffer(encryptedData))
	req.Header.Set("Authorization", "Bearer valid_token")
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	handler := NewDecryptDataHandler(mockAccessService, mockRepository, logger)
	handler.HandleDecryptData(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		render.JSON(res, req, "test_data")
	})).ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestHandleDecryptData_InvalidJWTToken(t *testing.T) {
	mockRepository := new(appMock.MockManager)
	mockAccessService := new(appMock.MockAccessService)
	logger, _ := logger.NewLogger("info")

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("", fmt.Errorf("invalid token"))

	req := httptest.NewRequest("POST", "/decrypt", bytes.NewBuffer([]byte("test_data")))
	req.Header.Set("Authorization", "Bearer invalid_token")
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	handler := NewDecryptDataHandler(mockAccessService, mockRepository, logger)
	handler.HandleDecryptData(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		render.JSON(res, req, "test_data")
	})).ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestHandleDecryptData_NonExistentUserUUID(t *testing.T) {
	mockRepository := new(appMock.MockManager)
	mockAccessService := new(appMock.MockAccessService)
	mockUserRepository := new(appMock.MockUserDataModelRepository)
	logger, _ := logger.NewLogger("info")

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("non_existent_uuid", nil)

	mockRepository.On("User").Return(mockUserRepository)
	mockUserRepository.On("FindOneByUUID", mock.Anything, "non_existent_uuid").Return(nil, fmt.Errorf("user not found"))

	req := httptest.NewRequest("POST", "/decrypt", bytes.NewBuffer([]byte("test_data")))
	req.Header.Set("Authorization", "Bearer valid_token")
	res := httptest.NewRecorder()

	handler := NewDecryptDataHandler(mockAccessService, mockRepository, logger)
	handler.HandleDecryptData(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		render.JSON(res, req, "test_data")
	})).ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestHandleDecryptData_EmptyRequestBody(t *testing.T) {
	mockRepository := new(appMock.MockManager)
	mockAccessService := new(appMock.MockAccessService)
	mockUserRepository := new(appMock.MockUserDataModelRepository)
	logger, _ := logger.NewLogger("info")

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)

	privateKey := string(make([]byte, 32))
	user := new(models.User)
	user.UUID = "userUUID"
	user.PrivateClientKey = privateKey

	mockRepository.On("User").Return(mockUserRepository)
	mockUserRepository.On("FindOneByUUID", mock.Anything, "userUUID").Return(user, nil)

	req := httptest.NewRequest("POST", "/decrypt", bytes.NewBuffer([]byte("")))
	req.Header.Set("Authorization", "Bearer valid_token")
	res := httptest.NewRecorder()

	handler := NewDecryptDataHandler(mockAccessService, mockRepository, logger)
	handler.HandleDecryptData(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		render.JSON(res, req, "test_data")
	})).ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestHandleEncryptData_SuccessfulEncryption(t *testing.T) {
	mockRepository := new(appMock.MockManager)
	mockAccessService := new(appMock.MockAccessService)
	mockUserRepository := new(appMock.MockUserDataModelRepository)
	logger, _ := logger.NewLogger("info")

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)

	privateKey := string(make([]byte, 32))
	user := new(models.User)
	user.UUID = "userUUID"
	user.PrivateClientKey = privateKey

	mockRepository.On("User").Return(mockUserRepository)
	mockUserRepository.On("FindOneByUUID", mock.Anything, "userUUID").Return(user, nil)

	req := httptest.NewRequest("POST", "/encrypt", bytes.NewBuffer([]byte("test_data")))
	req.Header.Set("Authorization", "Bearer valid_token")
	res := httptest.NewRecorder()

	handler := NewDecryptDataHandler(mockAccessService, mockRepository, logger)
	handler.HandleEncryptData(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("test_data"))
	})).ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	encryptedData := res.Body.Bytes()
	decryptedData, err := util.DataDecryptAES(encryptedData, []byte(privateKey))
	if err != nil {
		t.Errorf("Failed to decrypt response body: %v", err)
	}
	assert.Equal(t, "test_data", string(decryptedData))
}

func TestHandleEncryptData_InvalidJWTToken(t *testing.T) {
	mockRepository := new(appMock.MockManager)
	mockAccessService := new(appMock.MockAccessService)
	logger, _ := logger.NewLogger("info")

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("", fmt.Errorf("invalid token"))

	req := httptest.NewRequest("POST", "/encrypt", bytes.NewBuffer([]byte("test_data")))
	req.Header.Set("Authorization", "Bearer invalid_token")
	res := httptest.NewRecorder()

	handler := NewDecryptDataHandler(mockAccessService, mockRepository, logger)
	handler.HandleEncryptData(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("test_data"))
	})).ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestHandleEncryptData_NonExistentUserUUID(t *testing.T) {
	mockRepository := new(appMock.MockManager)
	mockAccessService := new(appMock.MockAccessService)
	mockUserRepository := new(appMock.MockUserDataModelRepository)
	logger, _ := logger.NewLogger("info")

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("non_existent_uuid", nil)

	mockRepository.On("User").Return(mockUserRepository)
	mockUserRepository.On("FindOneByUUID", mock.Anything, "non_existent_uuid").Return(nil, fmt.Errorf("user not found"))

	req := httptest.NewRequest("POST", "/encrypt", bytes.NewBuffer([]byte("test_data")))
	req.Header.Set("Authorization", "Bearer valid_token")
	res := httptest.NewRecorder()

	handler := NewDecryptDataHandler(mockAccessService, mockRepository, logger)
	handler.HandleEncryptData(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("test_data"))
	})).ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

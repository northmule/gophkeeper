package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/api/rctx"
	"github.com/northmule/gophkeeper/internal/server/logger"
	appMock "github.com/northmule/gophkeeper/internal/server/repository/mock"
	"github.com/northmule/gophkeeper/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestHandleRegistration(t *testing.T) {

	t.Run("Successful Registration", func(t *testing.T) {
		mockAccessService := new(appMock.MockAccessService)
		mockRepository := new(appMock.MockManager)
		mockUserRepository := new(appMock.MockUserDataModelRepository)
		mockSessionManager := new(appMock.MockSessionManager)
		mockTxDBQuery := new(appMock.MockTxDBQuery)
		mockQuery := new(appMock.MockDBQuery)

		mockRepository.On("User").Return(mockUserRepository)
		l, _ := logger.NewLogger("info")

		handler := NewRegistrationHandler(mockRepository, mockSessionManager, mockAccessService, l)

		mockQuery.On("Begin").Return(mockTxDBQuery, nil)
		transaction, _ := storage.NewTransaction(mockQuery)

		mockUserRepository.On("FindOneByLogin", mock.Anything, "testuser").Return(nil, nil)
		mockAccessService.On("PasswordHash", "testpassword").Return("hashedpassword", nil)
		mockUserRepository.On("TxCreateNewUser", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil)

		reqBody := `{"login": "testuser", "password": "testpassword", "email": "test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
		req = req.WithContext(context.WithValue(req.Context(), rctx.TransactionCtxKey, transaction))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.HandleRegistration(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Login Already Exists", func(t *testing.T) {
		mockAccessService := new(appMock.MockAccessService)
		mockRepository := new(appMock.MockManager)
		mockUserRepository := new(appMock.MockUserDataModelRepository)
		mockSessionManager := new(appMock.MockSessionManager)

		mockRepository.On("User").Return(mockUserRepository)
		l, _ := logger.NewLogger("info")

		handler := NewRegistrationHandler(mockRepository, mockSessionManager, mockAccessService, l)
		mockUserRepository.On("FindOneByLogin", mock.Anything, "testuser").Return(&models.User{Login: "testuser"}, nil)

		reqBody := `{"login": "testuser", "password": "testpassword", "email": "test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.HandleRegistration(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)
	})

}

func TestHandleRegistration_TxCreateNewUserError(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockUserRepository := new(appMock.MockUserDataModelRepository)
	mockSessionManager := new(appMock.MockSessionManager)
	mockTxDBQuery := new(appMock.MockTxDBQuery)
	mockQuery := new(appMock.MockDBQuery)
	mockQuery.On("Begin").Return(mockTxDBQuery, nil)
	transaction, _ := storage.NewTransaction(mockQuery)

	mockRepository.On("User").Return(mockUserRepository)
	l, _ := logger.NewLogger("info")

	handler := NewRegistrationHandler(mockRepository, mockSessionManager, mockAccessService, l)

	mockUserRepository.On("FindOneByLogin", mock.Anything, mock.Anything).Return(nil, nil)
	mockUserRepository.On("TxCreateNewUser", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), errors.New("database error"))
	mockAccessService.On("PasswordHash", mock.Anything).Return("hashedpassword", nil)

	reqBody := `{"login": "testuser", "password": "testpassword", "email": "test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := context.WithValue(req.Context(), rctx.TransactionCtxKey, transaction)
	req = req.WithContext(ctx)

	res := httptest.NewRecorder()

	handler.HandleRegistration(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	mockUserRepository.AssertExpectations(t)
	mockAccessService.AssertExpectations(t)

}

func TestHandleRegistration_FindOneByLoginError(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockUserRepository := new(appMock.MockUserDataModelRepository)
	mockSessionManager := new(appMock.MockSessionManager)

	mockRepository.On("User").Return(mockUserRepository)
	l, _ := logger.NewLogger("info")

	handler := NewRegistrationHandler(mockRepository, mockSessionManager, mockAccessService, l)

	mockUserRepository.On("FindOneByLogin", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

	reqBody := `{"login": "testuser", "password": "testpassword", "email": "test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	handler.HandleRegistration(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	mockUserRepository.AssertExpectations(t)
}

func TestHandleAuthentication(t *testing.T) {

	t.Run("Invalid request body", func(t *testing.T) {
		mockAccessService := new(appMock.MockAccessService)
		mockRepository := new(appMock.MockManager)
		mockUserRepository := new(appMock.MockUserDataModelRepository)
		l, _ := logger.NewLogger("info")

		handler := NewRegistrationHandler(mockRepository, nil, mockAccessService, l)

		mockRepository.On("User").Return(mockUserRepository)
		user := new(models.User)
		user.Login = "login"
		mockUserRepository.On("FindOneByLogin", mock.Anything, mock.Anything).Return(user, nil)
		mockAccessService.On("PasswordHash", mock.Anything).Return("hashedpassword", nil)
		reqBody := `{"login": "existinguser"}`
		req := httptest.NewRequest(http.MethodPost, "/authenticate", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.HandleAuthentication(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("User not found", func(t *testing.T) {
		mockAccessService := new(appMock.MockAccessService)
		mockRepository := new(appMock.MockManager)
		mockUserRepository := new(appMock.MockUserDataModelRepository)
		l, _ := logger.NewLogger("info")

		handler := NewRegistrationHandler(mockRepository, nil, mockAccessService, l)

		mockRepository.On("User").Return(mockUserRepository)
		mockAccessService.On("PasswordHash", mock.Anything).Return("hashedpassword", nil)
		mockUserRepository.On("FindOneByLogin", mock.Anything, "nonexistentuser").Return(nil, nil)
		reqBody := `{"login": "nonexistentuser", "password": "password"}`
		req := httptest.NewRequest(http.MethodPost, "/authenticate", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.HandleAuthentication(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("Successful authentication", func(t *testing.T) {
		mockAccessService := new(appMock.MockAccessService)
		mockRepository := new(appMock.MockManager)
		mockUserRepository := new(appMock.MockUserDataModelRepository)
		l, _ := logger.NewLogger("info")

		handler := NewRegistrationHandler(mockRepository, nil, mockAccessService, l)

		mockRepository.On("User").Return(mockUserRepository)
		user := new(models.User)
		user.Login = "login"

		mockAccessService.On("PasswordHash", mock.Anything).Return("hashedpassword", nil)
		mockUserRepository.On("FindOneByLogin", mock.Anything, "existinguser").Return(&models.User{Login: "existinguser", Password: "hashedpassword"}, nil)

		reqBody := `{"login": "existinguser", "password": "password"}`
		req := httptest.NewRequest(http.MethodPost, "/authenticate", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.HandleAuthentication(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("FindOneByLogin error", func(t *testing.T) {
		mockAccessService := new(appMock.MockAccessService)
		mockRepository := new(appMock.MockManager)
		mockUserRepository := new(appMock.MockUserDataModelRepository)
		l, _ := logger.NewLogger("info")

		handler := NewRegistrationHandler(mockRepository, nil, mockAccessService, l)

		mockRepository.On("User").Return(mockUserRepository)
		mockUserRepository.On("FindOneByLogin", mock.Anything, "existinguser").Return(nil, errors.New("database error"))

		reqBody := `{"login": "existinguser", "password": "password"}`
		req := httptest.NewRequest(http.MethodPost, "/authenticate", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler.HandleAuthentication(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

}

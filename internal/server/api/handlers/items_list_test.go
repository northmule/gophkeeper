package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/logger"
	appMock "github.com/northmule/gophkeeper/internal/server/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestItemsListHandler_HandleItemsList_SuccessfulRetrievalWithSpecifiedOffsetAndLimit(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockLogger, _ := logger.NewLogger("info")

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user123", nil)
	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockOwnerRepo.On("AllOwnerData", mock.Anything, "user123", 10, 50).Return([]models.OwnerData{
		{DataUUID: "item1", UserUUID: "user_uid", DataType: "Type 1", DataTypeName: "type", DataName: "name"},
	}, nil)

	handler := NewItemsListHandler(mockAccessService, mockRepository, mockLogger)

	req, _ := http.NewRequest("GET", "/items?offset=10&limit=50", nil)
	rr := httptest.NewRecorder()

	handler.HandleItemsList(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAccessService.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestItemsListHandler_HandleItemsList_RetrievalWithMaximumLimit(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockLogger, _ := logger.NewLogger("info")

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user123", nil)
	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockOwnerRepo.On("AllOwnerData", mock.Anything, "user123", 0, 200).Return([]models.OwnerData{
		{DataUUID: "item1", UserUUID: "user_uid", DataType: "Type 1", DataTypeName: "type", DataName: "name"},
	}, nil)

	handler := NewItemsListHandler(mockAccessService, mockRepository, mockLogger)

	req, _ := http.NewRequest("GET", "/items?limit=200", nil)
	rr := httptest.NewRecorder()

	handler.HandleItemsList(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAccessService.AssertExpectations(t)
	mockRepository.AssertExpectations(t)
}

func TestItemsListHandler_HandleItemsList_InvalidJWTToken(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockLogger, _ := logger.NewLogger("info")

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("", fmt.Errorf("invalid token"))

	handler := NewItemsListHandler(mockAccessService, mockRepository, mockLogger)

	req, _ := http.NewRequest("GET", "/items", nil)
	rr := httptest.NewRecorder()

	handler.HandleItemsList(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAccessService.AssertExpectations(t)
}

func TestItemsListHandler_HandleItemsList_InvalidOffset(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockLogger, _ := logger.NewLogger("info")

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user123", nil)
	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockOwnerRepo.On("AllOwnerData", mock.Anything, "user123", 0, mock.Anything).Return([]models.OwnerData{
		{DataUUID: "item1", UserUUID: "user_uid", DataType: "Type 1", DataTypeName: "type", DataName: "name"},
	}, nil)

	handler := NewItemsListHandler(mockAccessService, mockRepository, mockLogger)

	req, _ := http.NewRequest("GET", "/items?offset=abc", nil)
	rr := httptest.NewRecorder()

	handler.HandleItemsList(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAccessService.AssertExpectations(t)
}

func TestItemsListHandler_HandleItemsList_InvalidLimit(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockLogger, _ := logger.NewLogger("info")

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user123", nil)
	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockOwnerRepo.On("AllOwnerData", mock.Anything, "user123", mock.Anything, 0).Return([]models.OwnerData{}, nil)

	handler := NewItemsListHandler(mockAccessService, mockRepository, mockLogger)

	req, _ := http.NewRequest("GET", "/items?limit=abc", nil)
	rr := httptest.NewRecorder()

	handler.HandleItemsList(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAccessService.AssertExpectations(t)
}

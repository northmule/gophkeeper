package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/logger"
	appMock "github.com/northmule/gophkeeper/internal/server/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleSave_NewTextData(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockCardDataRepo := new(appMock.MockCardDataModelRepository)
	mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
	mockTextDataRepo := new(appMock.MockTextDataModelRepository)
	mockRepository := new(appMock.MockManager)
	l, _ := logger.NewLogger("info")

	mockRepository.On("CardData").Return(mockCardDataRepo)
	mockRepository.On("MetaData").Return(mockMetaDataRepo)
	mockRepository.On("TextData").Return(mockTextDataRepo)
	mockRepository.On("Owner").Return(mockOwnerRepo)

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user-uuid", nil)
	mockTextDataRepo.On("Add", mock.Anything, mock.Anything).Return(int64(1), nil)
	mockOwnerRepo.On("Add", mock.Anything, mock.Anything).Return(int64(1), nil)
	mockMetaDataRepo.On("Add", mock.Anything, mock.Anything).Return(int64(1), nil)

	handler := NewTextDataHandler(mockAccessService, mockRepository, l)

	reqBody, _ := json.Marshal(model_data.TextDataRequest{
		Name:  "Test Text",
		Value: "This is a test text",
		Meta: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	})

	req, _ := http.NewRequest("POST", "/textdata", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleSave(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAccessService.AssertExpectations(t)
	mockMetaDataRepo.AssertExpectations(t)
}

func TestHandleSave_UpdateTextData(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockCardDataRepo := new(appMock.MockCardDataModelRepository)
	mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
	mockTextDataRepo := new(appMock.MockTextDataModelRepository)
	mockRepository := new(appMock.MockManager)
	l, _ := logger.NewLogger("info")

	mockRepository.On("CardData").Return(mockCardDataRepo)
	mockRepository.On("MetaData").Return(mockMetaDataRepo)
	mockRepository.On("TextData").Return(mockTextDataRepo)
	mockRepository.On("Owner").Return(mockOwnerRepo)

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user-uuid", nil)
	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUIDAndDataType", mock.Anything, "user-uuid", "data-uuid", data_type.TextType).Return(&models.Owner{DataUUID: "data-uuid"}, nil)
	mockMetaDataRepo.On("ReplaceMetaByDataUUID", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	textData := new(models.TextData)
	textData.UUID = "data-uuid"
	mockTextDataRepo.On("FindOneByUUID", mock.Anything, "data-uuid").Return(textData, nil)
	mockTextDataRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	handler := NewTextDataHandler(mockAccessService, mockRepository, l)

	reqBody, _ := json.Marshal(model_data.TextDataRequest{
		UUID:  "data-uuid",
		Name:  "Updated Test Text",
		Value: "This is an updated test text",
		Meta: map[string]string{
			"key1": "updated-value1",
			"key2": "updated-value2",
		},
	})

	req, _ := http.NewRequest("POST", "/textdata", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleSave(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAccessService.AssertExpectations(t)
	mockTextDataRepo.AssertExpectations(t)
}

func TestHandleSave_InvalidJSON(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockLogger, _ := logger.NewLogger("info")

	handler := NewTextDataHandler(mockAccessService, mockRepository, mockLogger)

	reqBody := "invalid json"
	req, _ := http.NewRequest("POST", "/textdata", nil)
	req.Body = http.NoBody
	_ = json.NewDecoder(req.Body).Decode(&reqBody)

	rr := httptest.NewRecorder()

	handler.HandleSave(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAccessService.AssertExpectations(t)
}

func TestHandleSave_InvalidJWTToken(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockRepository := new(appMock.MockManager)
	mockLogger, _ := logger.NewLogger("info")

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("", fmt.Errorf("invalid token"))

	handler := NewTextDataHandler(mockAccessService, mockRepository, mockLogger)

	reqBody, _ := json.Marshal(model_data.TextDataRequest{
		Name:  "Test Text",
		Value: "This is a test text",
		Meta: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	})

	req, _ := http.NewRequest("POST", "/textdata", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandleSave(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockAccessService.AssertExpectations(t)
}

func TestHandleSave_NonExistentDataUUIDForUpdate(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockCardDataRepo := new(appMock.MockCardDataModelRepository)
	mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
	mockTextDataRepo := new(appMock.MockTextDataModelRepository)
	mockRepository := new(appMock.MockManager)
	l, _ := logger.NewLogger("info")

	mockRepository.On("CardData").Return(mockCardDataRepo)
	mockRepository.On("MetaData").Return(mockMetaDataRepo)
	mockRepository.On("TextData").Return(mockTextDataRepo)
	mockRepository.On("Owner").Return(mockOwnerRepo)

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user-uuid", nil)
	mockTextDataRepo.On("TextData").Return(mockTextDataRepo)
	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUIDAndDataType", mock.Anything, "user-uuid", "non-existent-uuid", data_type.TextType).Return(nil, nil)

	handler := NewTextDataHandler(mockAccessService, mockRepository, l)

	reqBody, _ := json.Marshal(model_data.TextDataRequest{
		UUID:  "non-existent-uuid",
		Name:  "Updated Test Text",
		Value: "This is an updated test text",
		Meta: map[string]string{
			"key1": "updated-value1",
			"key2": "updated-value2",
		},
	})

	req, _ := http.NewRequest("POST", "/textdata", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleSave(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockAccessService.AssertExpectations(t)
	mockOwnerRepo.AssertExpectations(t)
}

func TestHandleSave_RepositoryErrorOnReplaceMeta(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockCardDataRepo := new(appMock.MockCardDataModelRepository)
	mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
	mockTextDataRepo := new(appMock.MockTextDataModelRepository)
	mockRepository := new(appMock.MockManager)
	l, _ := logger.NewLogger("info")

	mockRepository.On("CardData").Return(mockCardDataRepo)
	mockRepository.On("MetaData").Return(mockMetaDataRepo)
	mockRepository.On("TextData").Return(mockTextDataRepo)
	mockRepository.On("Owner").Return(mockOwnerRepo)

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("user-uuid", nil)
	mockTextDataRepo.On("Add", mock.Anything, mock.Anything).Return(int64(1), nil)
	mockOwnerRepo.On("Add", mock.Anything, mock.Anything).Return(int64(1), nil)
	mockMetaDataRepo.On("ReplaceMetaByDataUUID", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("repository error"))
	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUIDAndDataType", mock.Anything, mock.Anything, mock.Anything, data_type.TextType).Return(&models.Owner{DataUUID: "data-uuid"}, nil)
	mockTextDataRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	textData := new(models.TextData)
	textData.UUID = "data-uuid"
	mockTextDataRepo.On("FindOneByUUID", mock.Anything, mock.Anything).Return(textData, nil)
	handler := NewTextDataHandler(mockAccessService, mockRepository, l)

	reqBody, _ := json.Marshal(model_data.TextDataRequest{
		UUID:  "uuid",
		Name:  "Test Text",
		Value: "This is a test text",
	})

	req, _ := http.NewRequest("POST", "/textdata", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler.HandleSave(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockAccessService.AssertExpectations(t)
	mockMetaDataRepo.AssertExpectations(t)
}

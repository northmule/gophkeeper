package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/northmule/gophkeeper/internal/server/logger"
	appMock "github.com/northmule/gophkeeper/internal/server/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestItemDataHandler_HandleItem_SuccessfulCardData(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockCardDataRepo := new(appMock.MockCardDataModelRepository)
	mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
	mockRepository := new(appMock.MockManager)
	logger, _ := logger.NewLogger("info")

	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockRepository.On("CardData").Return(mockCardDataRepo)
	mockRepository.On("MetaData").Return(mockMetaDataRepo)

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)

	dataUUID := uuid.NewString()

	owner := &models.Owner{
		UserUUID: "userUUID",
		DataType: data_type.CardType,
		DataUUID: dataUUID,
	}

	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUID", mock.Anything, "userUUID", dataUUID).Return(owner, nil)

	cardData := &models.CardData{
		Name: "Test Card",
		Value: models.CardDataValueV1{
			CardNumber:           "1234567890123456",
			NameBank:             "Test Bank",
			CurrentAccountNumber: "1234567890",
			FullNameHolder:       "John Doe",
			PhoneHolder:          "1234567890",
			SecurityCode:         "123",
			ValidityPeriod:       time.Now(),
		},
	}
	cardData.UUID = dataUUID

	mockCardDataRepo.On("FindOneByUUID", mock.Anything, dataUUID).Return(cardData, nil)

	metaData := []models.MetaData{
		{
			MetaName: "test",
			MetaValue: models.MetaDataValue{
				Value: "value",
			},
			DataUUID: dataUUID,
		},
	}

	mockMetaDataRepo.On("FindOneByUUID", mock.Anything, dataUUID).Return(metaData, nil)

	handler := NewItemDataHandler(mockAccessService, mockRepository, logger)

	req := httptest.NewRequest("GET", "/items/{uuid}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("uuid", dataUUID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	res := httptest.NewRecorder()

	handler.HandleItem(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), `"is_card":true`)
	assert.Contains(t, res.Body.String(), `1234567890123456`)
}

func TestItemDataHandler_HandleItem_SuccessfulTextData(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockTextDataRepo := new(appMock.MockTextDataModelRepository)
	mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
	mockRepository := new(appMock.MockManager)
	logger, _ := logger.NewLogger("info")

	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockRepository.On("TextData").Return(mockTextDataRepo)
	mockRepository.On("MetaData").Return(mockMetaDataRepo)

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)

	dataUUID := uuid.NewString()

	owner := &models.Owner{
		UserUUID: "userUUID",
		DataType: data_type.TextType,
		DataUUID: dataUUID,
	}

	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUID", mock.Anything, "userUUID", dataUUID).Return(owner, nil)

	textData := &models.TextData{
		Name:  "Test Text",
		Value: "This is a test text data.",
	}
	textData.UUID = dataUUID

	mockTextDataRepo.On("FindOneByUUID", mock.Anything, dataUUID).Return(textData, nil)

	metaData := []models.MetaData{
		{
			MetaName: "test",
			MetaValue: models.MetaDataValue{
				Value: "value",
			},
			DataUUID: dataUUID,
		},
	}

	mockMetaDataRepo.On("FindOneByUUID", mock.Anything, dataUUID).Return(metaData, nil)

	handler := NewItemDataHandler(mockAccessService, mockRepository, logger)

	req := httptest.NewRequest("GET", "/items/{uuid}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("uuid", dataUUID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	res := httptest.NewRecorder()

	handler.HandleItem(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), `"is_text":true`)
	assert.Contains(t, res.Body.String(), `"name":"Test Text"`)
	assert.Contains(t, res.Body.String(), `"value":"This is a test text data."`)
	assert.Contains(t, res.Body.String(), `"test":"value"`)
}

func TestItemDataHandler_HandleItem_SuccessfulFileData(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockFileDataRepo := new(appMock.MockFileDataModelRepository)
	mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
	mockRepository := new(appMock.MockManager)
	logger, _ := logger.NewLogger("info")

	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockRepository.On("FileData").Return(mockFileDataRepo)
	mockRepository.On("MetaData").Return(mockMetaDataRepo)

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)

	dataUUID := uuid.NewString()

	owner := &models.Owner{
		UserUUID: "userUUID",
		DataType: data_type.BinaryType,
		DataUUID: dataUUID,
	}

	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUID", mock.Anything, "userUUID", dataUUID).Return(owner, nil)

	fileData := &models.FileData{
		Name:      "Test File",
		FileName:  "test.pdf",
		Size:      1024,
		Extension: ".pdf",
		MimeType:  "application/pdf",
	}
	fileData.UUID = dataUUID

	mockFileDataRepo.On("FindOneByUUID", mock.Anything, dataUUID).Return(fileData, nil)

	metaData := []models.MetaData{
		{
			MetaName: "test",
			MetaValue: models.MetaDataValue{
				Value: "value",
			},
			DataUUID: dataUUID,
		},
	}

	mockMetaDataRepo.On("FindOneByUUID", mock.Anything, dataUUID).Return(metaData, nil)

	handler := NewItemDataHandler(mockAccessService, mockRepository, logger)

	req := httptest.NewRequest("GET", "/items/{uuid}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("uuid", dataUUID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	res := httptest.NewRecorder()

	handler.HandleItem(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), `"is_file":true`)
	assert.Contains(t, res.Body.String(), `"name":"Test File"`)
	assert.Contains(t, res.Body.String(), `"file_name":"test.pdf"`)
	assert.Contains(t, res.Body.String(), `"test":"value"`)
}

func TestItemDataHandler_HandleItem_EmptyUUID(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockFileDataRepo := new(appMock.MockFileDataModelRepository)
	mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
	mockRepository := new(appMock.MockManager)
	logger, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()

	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockRepository.On("FileData").Return(mockFileDataRepo)
	mockRepository.On("MetaData").Return(mockMetaDataRepo)

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)

	handler := NewItemDataHandler(mockAccessService, mockRepository, logger)

	req := httptest.NewRequest("GET", "/items/{uuid}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("uuid", "")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	res := httptest.NewRecorder()

	handler.HandleItem(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestItemDataHandler_HandleItem_NonExistentUserUUID(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockFileDataRepo := new(appMock.MockFileDataModelRepository)
	mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
	mockRepository := new(appMock.MockManager)
	logger, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()

	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockRepository.On("FileData").Return(mockFileDataRepo)
	mockRepository.On("MetaData").Return(mockMetaDataRepo)

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("", fmt.Errorf("user not found"))

	handler := NewItemDataHandler(mockAccessService, mockRepository, logger)

	req := httptest.NewRequest("GET", "/items/{uuid}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("uuid", "non-existent-uuid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	res := httptest.NewRecorder()

	handler.HandleItem(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestItemDataHandler_HandleItem_NonExistentDataUUID(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockFileDataRepo := new(appMock.MockFileDataModelRepository)
	mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
	mockRepository := new(appMock.MockManager)
	logger, _ := logger.NewLogger("info")

	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockRepository.On("FileData").Return(mockFileDataRepo)
	mockRepository.On("MetaData").Return(mockMetaDataRepo)

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)

	dataUUID := uuid.NewString()

	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUID", mock.Anything, "userUUID", dataUUID).Return(nil, nil)

	handler := NewItemDataHandler(mockAccessService, mockRepository, logger)

	req := httptest.NewRequest("GET", "/items/{uuid}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("uuid", dataUUID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	res := httptest.NewRecorder()

	handler.HandleItem(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

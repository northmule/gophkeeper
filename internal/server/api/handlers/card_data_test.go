package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/logger"
	appMock "github.com/northmule/gophkeeper/internal/server/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleSave_SuccessfulCreation(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockCardDataRepo := new(appMock.MockCardDataModelRepository)
	mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
	mockRepository := new(appMock.MockManager)

	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockRepository.On("CardData").Return(mockCardDataRepo)
	mockRepository.On("MetaData").Return(mockMetaDataRepo)

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)
	cardData := new(models.CardData)
	cardData.Name = "Test Card"
	cardData.UUID = uuid.NewString()
	cardData.ObjectType = data_type.CardType
	cardData.Value = models.CardDataValueV1{
		CardNumber:           "1234567890123456",
		ValidityPeriod:       time.Now().AddDate(1, 0, 0),
		SecurityCode:         "123",
		FullNameHolder:       "John Doe",
		NameBank:             "Test Bank",
		PhoneHolder:          "1234567890",
		CurrentAccountNumber: "12345678901234567890",
	}

	owner := &models.Owner{
		UserUUID: "userUUID",
		DataType: data_type.CardType,
		DataUUID: cardData.UUID,
	}
	metaData := &models.MetaData{
		MetaName: "test",
		MetaValue: models.MetaDataValue{
			Value: "value",
		},
		DataUUID: cardData.UUID,
	}

	_ = owner
	_ = metaData

	mockCardDataRepo.On("Add", mock.Anything, mock.Anything).Return(int64(1), nil)
	mockOwnerRepo.On("Add", mock.Anything, mock.Anything).Return(int64(1), nil)
	mockMetaDataRepo.On("Add", mock.Anything, mock.Anything).Return(int64(1), nil)

	logger, _ := logger.NewLogger("info")

	requestData := new(cardDataRequest)
	requestData.Name = "Test Card"
	requestData.CardNumber = "1234567890123456"
	requestData.ValidityPeriod = time.Now().AddDate(1, 0, 0).Format(time.RFC3339)
	requestData.SecurityCode = "123"
	requestData.FullNameHolder = "John Doe"
	requestData.NameBank = "Test Bank"
	requestData.PhoneHolder = "1234567890"
	requestData.CurrentAccountNumber = "12345678901234567890"
	requestData.Meta = map[string]string{"test": "value"}
	reqBody, _ := json.Marshal(requestData)

	req, _ := http.NewRequest("POST", "/card", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json") // обязательно для bind структуры
	res := httptest.NewRecorder()

	handler := NewCardDataHandler(mockAccessService, mockRepository, logger)
	handler.HandleSave(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

//func TestHandleSave_SuccessfulUpdate(t *testing.T) {
//	mockAccessService := new(MockAccessService)
//	mockRepository := new(MockRepository)
//	mockOwnerRepo := new(MockOwnerRepository)
//	mockCardDataRepo := new(MockCardDataRepository)
//	mockMetaDataRepo := new(MockMetaDataRepository)
//
//	mockRepository.On("Owner").Return(mockOwnerRepo)
//	mockRepository.On("CardData").Return(mockCardDataRepo)
//	mockRepository.On("MetaData").Return(mockMetaDataRepo)
//
//	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)
//
//	dataUUID := uuid.NewString()
//	cardData := &models.CardData{
//		Name:       "Test Card",
//		UUID:       dataUUID,
//		ObjectType: data_type.CardType,
//		Value: models.CardDataValue{
//			CardNumber:           "1234567890123456",
//			ValidityPeriod:       time.Now().AddDate(1, 0, 0),
//			SecurityCode:         "123",
//			FullNameHolder:       "John Doe",
//			NameBank:             "Test Bank",
//			PhoneHolder:          "1234567890",
//			CurrentAccountNumber: "12345678901234567890",
//		},
//	}
//
//	mockCardDataRepo.On("FindOneByUUID", mock.Anything, dataUUID).Return(cardData, nil)
//
//	owner := &models.Owner{
//		UserUUID: "userUUID",
//		DataType: data_type.CardType,
//		DataUUID: dataUUID,
//	}
//
//	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUIDAndDataType", mock.Anything, "userUUID", dataUUID, data_type.CardType).Return(owner, nil)
//
//	mockCardDataRepo.On("Update", mock.Anything, cardData).Return(nil)
//
//	newValidityPeriod := time.Now().AddDate(2, 0, 0)
//	newMetaData := []models.MetaData{
//		{
//			MetaName: "test",
//			MetaValue: models.MetaValue{
//				Value: "value",
//			},
//			DataUUID: dataUUID,
//		},
//	}
//
//	mockMetaDataRepo.On("ReplaceMetaByDataUUID", mock.Anything, dataUUID, newMetaData).Return(nil)
//
//	handler := NewCardDataHandler(mockAccessService, mockRepository, logger.NewLogger())
//
//	reqBody, _ := json.Marshal(model_data.CardDataRequest{
//		UUID:                 dataUUID,
//		Name:                 "Updated Test Card",
//		CardNumber:           "6543210987654321",
//		ValidityPeriod:       newValidityPeriod.Format(time.RFC3339),
//		SecurityCode:         "321",
//		FullNameHolder:       "Jane Doe",
//		NameBank:             "Updated Test Bank",
//		PhoneHolder:          "0987654321",
//		CurrentAccountNumber: "09876543210987654321",
//		Meta:                 map[string]string{"test": "value"},
//	})
//
//	req, _ := http.NewRequest("POST", "/card", bytes.NewBuffer(reqBody))
//	res := httptest.NewRecorder()
//
//	handler.HandleSave(res, req)
//
//	assert.Equal(t, http.StatusOK, res.Code)
//}
//
//func TestHandleSave_EmptyRequest(t *testing.T) {
//	mockAccessService := new(MockAccessService)
//	mockRepository := new(MockRepository)
//	mockOwnerRepo := new(MockOwnerRepository)
//	mockCardDataRepo := new(MockCardDataRepository)
//	mockMetaDataRepo := new(MockMetaDataRepository)
//
//	mockRepository.On("Owner").Return(mockOwnerRepo)
//	mockRepository.On("CardData").Return(mockCardDataRepo)
//	mockRepository.On("MetaData").Return(mockMetaDataRepo)
//
//	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)
//
//	handler := NewCardDataHandler(mockAccessService, mockRepository, logger.NewLogger())
//
//	reqBody, _ := json.Marshal(model_data.CardDataRequest{})
//
//	req, _ := http.NewRequest("POST", "/card", bytes.NewBuffer(reqBody))
//	res := httptest.NewRecorder()
//
//	handler.HandleSave(res, req)
//
//	assert.Equal(t, http.StatusBadRequest, res.Code)
//}
//
//func TestHandleSave_InvalidValidityPeriod(t *testing.T) {
//	mockAccessService := new(MockAccessService)
//	mockRepository := new(MockRepository)
//	mockOwnerRepo := new(MockOwnerRepository)
//	mockCardDataRepo := new(MockCardDataRepository)
//	mockMetaDataRepo := new(MockMetaDataRepository)
//
//	mockRepository.On("Owner").Return(mockOwnerRepo)
//	mockRepository.On("CardData").Return(mockCardDataRepo)
//	mockRepository.On("MetaData").Return(mockMetaDataRepo)
//
//	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)
//
//	handler := NewCardDataHandler(mockAccessService, mockRepository, logger.NewLogger())
//
//	reqBody, _ := json.Marshal(model_data.CardDataRequest{
//		Name:                 "Test Card",
//		CardNumber:           "1234567890123456",
//		ValidityPeriod:       "invalid-date", // Invalid date format
//		SecurityCode:         "123",
//		FullNameHolder:       "John Doe",
//		NameBank:             "Test Bank",
//		PhoneHolder:          "1234567890",
//		CurrentAccountNumber: "12345678901234567890",
//		Meta:                 map[string]string{"test": "value"},
//	})
//
//	req, _ := http.NewRequest("POST", "/card", bytes.NewBuffer(reqBody))
//	res := httptest.NewRecorder()
//
//	handler.HandleSave(res, req)
//
//	assert.Equal(t, http.StatusBadRequest, res.Code)
//}
//
//func TestHandleSave_NonExistentUserUUID(t *testing.T) {
//	mockAccessService := new(MockAccessService)
//	mockRepository := new(MockRepository)
//	mockOwnerRepo := new(MockOwnerRepository)
//	mockCardDataRepo := new(MockCardDataRepository)
//	mockMetaDataRepo := new(MockMetaDataRepository)
//
//	mockRepository.On("Owner").Return(mockOwnerRepo)
//	mockRepository.On("CardData").Return(mockCardDataRepo)
//	mockRepository.On("MetaData").Return(mockMetaDataRepo)
//
//	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("", fmt.Errorf("user not found"))
//
//	handler := NewCardDataHandler(mockAccessService, mockRepository, logger.NewLogger())
//
//	reqBody, _ := json.Marshal(model_data.CardDataRequest{
//		Name:                 "Test Card",
//		CardNumber:           "1234567890123456",
//		ValidityPeriod:       time.Now().AddDate(1, 0, 0).Format(time.RFC3339),
//		SecurityCode:         "123",
//		FullNameHolder:       "John Doe",
//		NameBank:             "Test Bank",
//		PhoneHolder:          "1234567890",
//		CurrentAccountNumber: "12345678901234567890",
//		Meta:                 map[string]string{"test": "value"},
//	})
//
//	req, _ := http.NewRequest("POST", "/card", bytes.NewBuffer(reqBody))
//	res := httptest.NewRecorder()
//
//	handler.HandleSave(res, req)
//
//	assert.Equal(t, http.StatusBadRequest, res.Code)
//}
//
//func TestHandleSave_NonExistentDataUUID(t *testing.T) {
//	mockAccessService := new(MockAccessService)
//	mockRepository := new(MockRepository)
//	mockOwnerRepo := new(MockOwnerRepository)
//	mockCardDataRepo := new(MockCardDataRepository)
//	mockMetaDataRepo := new(MockMetaDataRepository)
//
//	mockRepository.On("Owner").Return(mockOwnerRepo)
//	mockRepository.On("CardData").Return(mockCardDataRepo)
//	mockRepository.On("MetaData").Return(mockMetaDataRepo)
//
//	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)
//
//	dataUUID := uuid.NewString()
//
//	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUIDAndDataType", mock.Anything, "userUUID", dataUUID, data_type.CardType).Return(nil, nil)
//
//	handler := NewCardDataHandler(mockAccessService, mockRepository, logger.NewLogger())
//
//	reqBody, _ := json.Marshal(model_data.CardDataRequest{
//		UUID:                 dataUUID,
//		Name:                 "Test Card",
//		CardNumber:           "1234567890123456",
//		ValidityPeriod:       time.Now().AddDate(1, 0, 0).Format(time.RFC3339),
//		SecurityCode:         "123",
//		FullNameHolder:       "John Doe",
//		NameBank:             "Test Bank",
//		PhoneHolder:          "1234567890",
//		CurrentAccountNumber: "12345678901234567890",
//		Meta:                 map[string]string{"test": "value"},
//	})
//
//	req, _ := http.NewRequest("POST", "/card", bytes.NewBuffer(reqBody))
//	res := httptest.NewRecorder()
//
//	handler.HandleSave(res, req)
//
//	assert.Equal(t, http.StatusNotFound, res.Code)
//}
//
//func TestHandleSave_InvalidJSON(t *testing.T) {
//	mockAccessService := new(MockAccessService)
//	mockRepository := new(MockRepository)
//	mockOwnerRepo := new(MockOwnerRepository)
//	mockCardDataRepo := new(MockCardDataRepository)
//	mockMetaDataRepo := new(MockMetaDataRepository)
//
//	mockRepository.On("Owner").Return(mockOwnerRepo)
//	mockRepository.On("CardData").Return(mockCardDataRepo)
//	mockRepository.On("MetaData").Return(mockMetaDataRepo)
//
//	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("userUUID", nil)
//
//	handler := NewCardDataHandler(mockAccessService, mockRepository, logger.NewLogger())
//
//	reqBody := []byte(`{ "invalid": "json" }`)
//
//	req, _ := http.NewRequest("POST", "/card", bytes.NewBuffer(reqBody))
//	res := httptest.NewRecorder()
//
//	handler.HandleSave(res, req)
//
//	assert.Equal(t, http.StatusBadRequest, res.Code)
//}

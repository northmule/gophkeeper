package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/northmule/gophkeeper/internal/server/logger"
	appMock "github.com/northmule/gophkeeper/internal/server/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestFileDataHandleInit_SuccessfulCreation(t *testing.T) {
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

	mockFileDataRepo.On("Add", mock.Anything, mock.Anything).Return(int64(1), nil)
	mockOwnerRepo.On("Add", mock.Anything, mock.Anything).Return(int64(1), nil)
	mockMetaDataRepo.On("Add", mock.Anything, mock.Anything).Return(int64(1), nil)

	requestData := new(fileDataInitRequest)
	requestData.Name = "Test File"
	requestData.FileName = "test.pdf"
	requestData.Size = 1024
	requestData.Extension = ".pdf"
	requestData.MimeType = "application/pdf"
	requestData.Meta = map[string]string{"test": "value"}
	reqBody, _ := json.Marshal(requestData)

	req, _ := http.NewRequest("POST", "/file_data/init", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	handler := NewFileDataHandler(mockAccessService, mockRepository, cfg, logger)
	handler.HandleInit(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestFileDataHandleInit_SuccessfulUpdate(t *testing.T) {
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

	fileData := new(models.FileData)
	fileData.Name = "Test File"
	fileData.UUID = dataUUID
	fileData.MimeType = "application/pdf"
	fileData.FileName = "test.pdf"
	fileData.Size = 1024
	fileData.Extension = ".pdf"
	fileData.PathTmp = os.TempDir() + "/" + fileData.UUID
	fileData.Path = os.TempDir() + "/load_" + fileData.UUID
	fileData.Storage = "local://"
	fileData.Uploaded = false

	mockFileDataRepo.On("FindOneByUUID", mock.Anything, dataUUID).Return(fileData, nil)

	owner := &models.Owner{
		UserUUID: "userUUID",
		DataType: data_type.BinaryType,
		DataUUID: dataUUID,
	}

	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUIDAndDataType", mock.Anything, "userUUID", dataUUID, data_type.BinaryType).Return(owner, nil)

	mockFileDataRepo.On("Update", mock.Anything, fileData).Return(nil)

	newMetaData := []models.MetaData{
		{
			MetaName: "test",
			MetaValue: models.MetaDataValue{
				Value: "value",
			},
			DataUUID: dataUUID,
		},
	}

	mockMetaDataRepo.On("ReplaceMetaByDataUUID", mock.Anything, dataUUID, newMetaData).Return(nil)

	handler := NewFileDataHandler(mockAccessService, mockRepository, &config.Config{}, logger)

	reqBody, _ := json.Marshal(model_data.FileDataInitRequest{
		UUID:      dataUUID,
		Name:      "Updated Test File",
		FileName:  "updated_test.pdf",
		Size:      2048,
		Extension: ".pdf",
		MimeType:  "application/pdf",
		Meta:      map[string]string{"test": "value"},
	})

	req, _ := http.NewRequest("POST", "/file_data/init", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	handler.HandleInit(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestFileDataHandleInit_EmptyRequest(t *testing.T) {
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

	handler := NewFileDataHandler(mockAccessService, mockRepository, &config.Config{}, logger)

	reqBody, _ := json.Marshal(model_data.FileDataInitRequest{})

	req, _ := http.NewRequest("POST", "/file_data/init", bytes.NewBuffer(reqBody))
	res := httptest.NewRecorder()

	handler.HandleInit(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestFileDataHandleInit_NonExistentUserUUID(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockFileDataRepo := new(appMock.MockFileDataModelRepository)
	mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
	mockRepository := new(appMock.MockManager)
	logger, _ := logger.NewLogger("info")

	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockRepository.On("FileData").Return(mockFileDataRepo)
	mockRepository.On("MetaData").Return(mockMetaDataRepo)

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("", fmt.Errorf("user not found"))

	handler := NewFileDataHandler(mockAccessService, mockRepository, &config.Config{}, logger)

	reqBody, _ := json.Marshal(model_data.FileDataInitRequest{
		Name:      "Test File",
		FileName:  "test.pdf",
		Size:      1024,
		Extension: ".pdf",
		MimeType:  "application/pdf",
		Meta:      map[string]string{"test": "value"},
	})

	req, _ := http.NewRequest("POST", "/file_data/init", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	handler.HandleInit(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestFileDataHandleInit_NonExistentDataUUID(t *testing.T) {
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

	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUIDAndDataType", mock.Anything, "userUUID", dataUUID, data_type.BinaryType).Return(nil, nil)

	handler := NewFileDataHandler(mockAccessService, mockRepository, &config.Config{}, logger)

	reqBody, _ := json.Marshal(model_data.FileDataInitRequest{
		UUID:      dataUUID,
		Name:      "Test File",
		FileName:  "test.pdf",
		Size:      1024,
		Extension: ".pdf",
		MimeType:  "application/pdf",
		Meta:      map[string]string{"test": "value"},
	})

	req, _ := http.NewRequest("POST", "/file_data/init", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	handler.HandleInit(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestFileData_HandleAction(t *testing.T) {

	t.Run("SuccessfulFileHandling", func(t *testing.T) {

		mockAccessService := new(appMock.MockAccessService)
		mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
		mockRepository := new(appMock.MockManager)
		mockFileDataRepo := new(appMock.MockFileDataModelRepository)
		mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
		logger, _ := logger.NewLogger("info")
		cfg := config.NewConfig()
		_ = cfg.Init()

		mockRepository.On("Owner").Return(mockOwnerRepo)
		mockRepository.On("FileData").Return(mockFileDataRepo)
		mockRepository.On("MetaData").Return(mockMetaDataRepo)

		handler := &FileDataHandler{
			accessService: mockAccessService,
			manager:       mockRepository,
			log:           logger,
		}

		fileUUID := uuid.NewString()

		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("file_uuid", fileUUID)
		ctx.URLParams.Add("part", "valid-part")

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile(data_type.FileField, "test.pdf")
		_, _ = part.Write([]byte("test file content"))
		_ = writer.Close()

		req := httptest.NewRequest("POST", "/files/{file_uuid}/{part}", body)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req.Header.Set("Content-Type", writer.FormDataContentType())

		owner := new(models.Owner)
		owner.UserUUID = "123"
		owner.UserUUID = "321"
		dataUUID := fileUUID
		fileData := new(models.FileData)
		fileData.Name = "Test File"
		fileData.UUID = dataUUID
		fileData.MimeType = "application/pdf"
		fileData.FileName = "test.pdf"
		fileData.Size = 1024
		fileData.Extension = ".pdf"
		fileData.PathTmp = os.TempDir() + "/" + fileData.UUID
		fileData.Path = os.TempDir() + "/load_" + fileData.UUID
		fileData.Storage = "local://"
		fileData.Uploaded = false
		mockFileDataRepo.On("FindOneByUUID", mock.Anything, dataUUID).Return(fileData, nil)
		mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("valid-user-uuid", nil)
		mockOwnerRepo.On("FindOneByUserUUIDAndDataUUIDAndDataType", mock.Anything, "valid-user-uuid", fileUUID, data_type.BinaryType).Return(owner, nil)
		mockFileDataRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

		res := httptest.NewRecorder()
		handler.HandleAction(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		mockAccessService.AssertExpectations(t)
		mockOwnerRepo.AssertExpectations(t)
	})

	t.Run("UserUUIDNotFound", func(t *testing.T) {

		mockAccessService := new(appMock.MockAccessService)
		mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
		mockRepository := new(appMock.MockManager)
		mockFileDataRepo := new(appMock.MockFileDataModelRepository)
		mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
		logger, _ := logger.NewLogger("info")
		cfg := config.NewConfig()
		_ = cfg.Init()

		mockRepository.On("Owner").Return(mockOwnerRepo)
		mockRepository.On("FileData").Return(mockFileDataRepo)
		mockRepository.On("MetaData").Return(mockMetaDataRepo)

		handler := &FileDataHandler{
			accessService: mockAccessService,
			manager:       mockRepository,
			log:           logger,
		}

		req := httptest.NewRequest("POST", "/files/{file_uuid}/{part}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("file_uuid", "valid-file-uuid")
		ctx.URLParams.Add("part", "valid-part")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("11212", nil)
		mockOwnerRepo.On("FindOneByUserUUIDAndDataUUIDAndDataType", mock.Anything, mock.Anything, mock.Anything, data_type.BinaryType).Return(nil, nil)
		res := httptest.NewRecorder()
		handler.HandleAction(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		mockAccessService.AssertExpectations(t)
	})

	t.Run("DataUUIDNotFound", func(t *testing.T) {
		mockAccessService := new(appMock.MockAccessService)
		mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
		mockRepository := new(appMock.MockManager)
		mockFileDataRepo := new(appMock.MockFileDataModelRepository)
		mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
		logger, _ := logger.NewLogger("info")
		cfg := config.NewConfig()
		_ = cfg.Init()

		mockRepository.On("Owner").Return(mockOwnerRepo)
		mockRepository.On("FileData").Return(mockFileDataRepo)
		mockRepository.On("MetaData").Return(mockMetaDataRepo)

		handler := &FileDataHandler{
			accessService: mockAccessService,
			manager:       mockRepository,
			log:           logger,
		}

		req := httptest.NewRequest("POST", "/files/{file_uuid}/{part}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("file_uuid", "valid-file-uuid")
		ctx.URLParams.Add("part", "valid-part")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("valid-user-uuid", nil)
		mockOwnerRepo.On("FindOneByUserUUIDAndDataUUIDAndDataType", mock.Anything, "valid-user-uuid", "valid-file-uuid", data_type.BinaryType).Return(nil, nil)

		res := httptest.NewRecorder()
		handler.HandleAction(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		mockAccessService.AssertExpectations(t)
		mockOwnerRepo.AssertExpectations(t)
	})

	t.Run("InvalidJWTToken", func(t *testing.T) {
		mockAccessService := new(appMock.MockAccessService)
		mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
		mockRepository := new(appMock.MockManager)
		mockFileDataRepo := new(appMock.MockFileDataModelRepository)
		mockMetaDataRepo := new(appMock.MockMetaDataModelRepository)
		logger, _ := logger.NewLogger("info")
		cfg := config.NewConfig()
		_ = cfg.Init()

		mockRepository.On("Owner").Return(mockOwnerRepo)
		mockRepository.On("FileData").Return(mockFileDataRepo)
		mockRepository.On("MetaData").Return(mockMetaDataRepo)

		handler := &FileDataHandler{
			accessService: mockAccessService,
			manager:       mockRepository,
			log:           logger,
		}

		req := httptest.NewRequest("POST", "/files/{file_uuid}/{part}", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("file_uuid", "valid-file-uuid")
		ctx.URLParams.Add("part", "valid-part")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("", fmt.Errorf("invalid token"))

		res := httptest.NewRecorder()
		handler.HandleAction(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		mockAccessService.AssertExpectations(t)
	})
}

func TestFileDataHandleGetAction_SuccessfulFileDownload(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockFileDataRepo := new(appMock.MockFileDataModelRepository)
	mockRepository := new(appMock.MockManager)
	logger, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()

	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockRepository.On("FileData").Return(mockFileDataRepo)

	handler := &FileDataHandler{
		accessService: mockAccessService,
		manager:       mockRepository,
		log:           logger,
	}

	fileUUID := uuid.NewString()

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("file_uuid", fileUUID)
	ctx.URLParams.Add("part", "0")

	req := httptest.NewRequest("GET", "/files/{file_uuid}/{part}", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	fileName := uuid.NewString()

	owner := new(models.Owner)
	owner.UserUUID = "valid-user-uuid"
	dataUUID := fileUUID
	fileData := new(models.FileData)
	fileData.Name = "Test File"
	fileData.UUID = dataUUID
	fileData.MimeType = "application/pdf"
	fileData.FileName = fileName + "_test.pdf"
	fileData.Size = 1024
	fileData.Extension = ".pdf"
	fileData.Path = os.TempDir()
	fileData.Storage = "local://"
	fileData.Uploaded = true

	mockFileDataRepo.On("FindOneByUUID", mock.Anything, dataUUID).Return(fileData, nil)
	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("valid-user-uuid", nil)
	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUIDAndDataType", mock.Anything, "valid-user-uuid", fileUUID, data_type.BinaryType).Return(owner, nil)

	tempFile, err := os.Create(path.Join(fileData.Path, fileData.FileName))
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	testData := []byte("test file content")
	if _, err := tempFile.Write(testData); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	res := httptest.NewRecorder()
	handler.HandleGetAction(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, testData, res.Body.Bytes())
	mockAccessService.AssertExpectations(t)
	mockOwnerRepo.AssertExpectations(t)
	mockFileDataRepo.AssertExpectations(t)
}

func TestFileDataHandleGetAction_UserUUIDNotFound(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockFileDataRepo := new(appMock.MockFileDataModelRepository)
	mockRepository := new(appMock.MockManager)
	logger, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()

	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockRepository.On("FileData").Return(mockFileDataRepo)

	handler := &FileDataHandler{
		accessService: mockAccessService,
		manager:       mockRepository,
		log:           logger,
	}

	req := httptest.NewRequest("GET", "/files/{file_uuid}/{part}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("file_uuid", "valid-file-uuid")
	ctx.URLParams.Add("part", "0")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("11212", nil)
	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUIDAndDataType", mock.Anything, mock.Anything, mock.Anything, data_type.BinaryType).Return(nil, nil)

	res := httptest.NewRecorder()
	handler.HandleGetAction(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
	mockAccessService.AssertExpectations(t)
	mockOwnerRepo.AssertExpectations(t)
}

func TestFileDataHandleGetAction_DataUUIDNotFound(t *testing.T) {
	mockAccessService := new(appMock.MockAccessService)
	mockOwnerRepo := new(appMock.MockOwnerDataModelRepository)
	mockFileDataRepo := new(appMock.MockFileDataModelRepository)
	mockRepository := new(appMock.MockManager)
	logger, _ := logger.NewLogger("info")
	cfg := config.NewConfig()
	_ = cfg.Init()

	mockRepository.On("Owner").Return(mockOwnerRepo)
	mockRepository.On("FileData").Return(mockFileDataRepo)

	handler := &FileDataHandler{
		accessService: mockAccessService,
		manager:       mockRepository,
		log:           logger,
	}

	req := httptest.NewRequest("GET", "/files/{file_uuid}/{part}", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("file_uuid", "valid-file-uuid")
	ctx.URLParams.Add("part", "0")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	mockAccessService.On("GetUserUUIDByJWTToken", mock.Anything).Return("valid-user-uuid", nil)
	mockOwnerRepo.On("FindOneByUserUUIDAndDataUUIDAndDataType", mock.Anything, "valid-user-uuid", "valid-file-uuid", data_type.BinaryType).Return(nil, nil)

	res := httptest.NewRecorder()
	handler.HandleGetAction(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
	mockAccessService.AssertExpectations(t)
	mockOwnerRepo.AssertExpectations(t)
}

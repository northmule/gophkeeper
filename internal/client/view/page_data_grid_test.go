package view

import (
	"path"
	"strings"
	"testing"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/controller"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/service"
	"github.com/northmule/gophkeeper/internal/client/storage"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockManagerController is a mock implementation of ManagerController.
type MockManagerController struct {
	mock.Mock
}

func (m *MockManagerController) Authentication() *controller.Authentication {
	args := m.Called()
	return args.Get(0).(*controller.Authentication)
}

func (m *MockManagerController) CardData() *controller.CardData {
	args := m.Called()
	return args.Get(0).(*controller.CardData)
}

func (m *MockManagerController) TextData() controller.TextDataController {
	args := m.Called()
	return args.Get(0).(*mockTextData)
}

func (m *MockManagerController) FileData() controller.FileDataController {
	args := m.Called()
	return args.Get(0).(*mockFileData)
}

func (m *MockManagerController) GridData() controller.GridDataController {
	args := m.Called()
	return args.Get(0).(controller.GridDataController)
}

func (m *MockManagerController) ItemData() controller.ItemDataController {
	args := m.Called()
	return args.Get(0).(controller.ItemDataController)
}

func (m *MockManagerController) KeysData() *controller.KeysData {
	args := m.Called()
	return args.Get(0).(*controller.KeysData)
}

func (m *MockManagerController) Registration() controller.RegistrationController {
	args := m.Called()
	return args.Get(0).(*mockRegistration)
}

// MockGridDataController is a mock implementation of GridDataController.
type MockGridDataController struct {
	mock.Mock
}

func (m *MockGridDataController) Send(token string) (*controller.GridDataResponse, error) {
	args := m.Called(token)
	return args.Get(0).(*controller.GridDataResponse), args.Error(1)
}

// MockItemDataController is a mock implementation of ItemDataController.
type MockItemDataController struct {
	mock.Mock
}

func (m *MockItemDataController) Send(token string, dataUUID string) (*model_data.DataByUUIDResponse, error) {
	args := m.Called(token, dataUUID)
	return args.Get(0).(*model_data.DataByUUIDResponse), args.Error(1)
}

func TestPageDataGrid_Init(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	cryptService, _ := service.NewCrypt(mockCfg)
	manager, _ := controller.NewManager(mockCfg, cryptService, log)
	memoryStorage := storage.NewMemoryStorage()
	mainPage := newPageIndex(manager, memoryStorage, log)

	pa := newPageAction(mainPage)
	pdg := newPageDataGrid(mainPage, pa)
	assert.NotNil(t, pdg)
}

func TestNewPageDataGrid(t *testing.T) {
	mockManagerController := new(MockManagerController)
	mockGridDataController := new(MockGridDataController)
	mockStorage := storage.NewMemoryStorage()
	log, _ := logger.NewLogger("info")
	mockPageIndex := newPageIndex(mockManagerController, mockStorage, log)
	mockPageAction := newPageAction(mockPageIndex)

	mockStorage.SetToken("token")
	mockManagerController.On("GridData").Return(mockGridDataController)

	result := new(controller.GridDataResponse)
	result.Items = []model_data.ItemDataResponse{
		{Number: "1", Type: "Card", Name: "Card1", UUID: "uuid1"},
		{Number: "2", Type: "Text", Name: "Text1", UUID: "uuid2"},
	}
	mockGridDataController.On("Send", "token").Return(result, nil)

	pdg := newPageDataGrid(mockPageIndex, mockPageAction)
	assert.NotNil(t, pdg)
	assert.NotNil(t, pdg.table)
	assert.Len(t, pdg.table.Rows(), 2)
}

func TestPageDataGrid_View(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	cryptService, _ := service.NewCrypt(mockCfg)
	manager, _ := controller.NewManager(mockCfg, cryptService, log)
	memoryStorage := storage.NewMemoryStorage()
	mainPage := newPageIndex(manager, memoryStorage, log)
	pa := newPageAction(mainPage)

	page := newPageDataGrid(mainPage, pa)
	result := page.View()
	assert.True(t, strings.Contains(result, "Все данные"))
	assert.True(t, strings.Contains(result, "вверх/вниз: для переключения"))
	assert.True(t, strings.Contains(result, "enter: просмотреть данные"))
	assert.True(t, strings.Contains(result, "ctrl+c: вернуться"))
}

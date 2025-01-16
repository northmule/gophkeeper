package view

import (
	"path"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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

func (m *MockManagerController) Authentication() controller.AuthenticationDataController {
	args := m.Called()
	return args.Get(0).(controller.AuthenticationDataController)
}

func (m *MockManagerController) CardData() controller.CardDataController {
	args := m.Called()
	return args.Get(0).(controller.CardDataController)
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

func (m *MockManagerController) KeysData() controller.KeyDataController {
	args := m.Called()
	return args.Get(0).(controller.KeyDataController)
}

func (m *MockManagerController) Registration() controller.RegistrationController {
	args := m.Called()
	return args.Get(0).(*mockRegistration)
}

// MockGridDataController mock
type MockGridDataController struct {
	mock.Mock
}

func (m *MockGridDataController) Send(token string) (*controller.GridDataResponse, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*controller.GridDataResponse), args.Error(1)
}

// MockItemDataController mock
type MockItemDataController struct {
	mock.Mock
}

func (m *MockItemDataController) Send(token string, dataUUID string) (*model_data.DataByUUIDResponse, error) {
	args := m.Called(token, dataUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model_data.DataByUUIDResponse), args.Error(1)
}

// MockAuthenticationDataController mock
type MockAuthenticationDataController struct {
	mock.Mock
}

func (m *MockAuthenticationDataController) Send(login string, password string) (*controller.AuthenticationResponse, error) {
	args := m.Called(login, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*controller.AuthenticationResponse), args.Error(1)
}

// MockKeyDataController mock
type MockKeyDataController struct {
	mock.Mock
}

func (m *MockKeyDataController) UploadClientPublicKey(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockKeyDataController) DownloadPublicServerKey(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockKeyDataController) UploadClientPrivateKey(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

// MockCardDataController mock
type MockCardDataController struct {
	mock.Mock
}

func (m *MockCardDataController) Send(token string, requestData *model_data.CardDataRequest) (*controller.CardDataResponse, error) {
	args := m.Called(token, requestData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*controller.CardDataResponse), args.Error(1)
}

func TestPageDataGrid_Init(t *testing.T) {
	mockCfg, _ := config.NewConfig()
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
	mockCfg, _ := config.NewConfig()
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

func TestPageDataGrid_Update(t *testing.T) {
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	msg := tea.KeyMsg{Type: tea.KeyEnter}

	t.Run("enter IsFile", func(t *testing.T) {
		mockManagerController := new(MockManagerController)
		mockItemData := new(MockItemDataController)
		mockGridData := new(MockGridDataController)
		mockManagerController.On("ItemData").Return(mockItemData)
		mockManagerController.On("GridData").Return(mockGridData)

		responseData := new(controller.GridDataResponse)
		responseData.Items = []model_data.ItemDataResponse{
			{
				Number: "1",
				Type:   "type1",
				Name:   "name1",
				UUID:   "111",
			},
			{
				Number: "2",
				Type:   "type2",
				Name:   "name2",
				UUID:   "22222",
			},
		}
		mockGridData.On("Send", mock.Anything).Return(responseData, nil)

		mainPage := newPageIndex(mockManagerController, memoryStorage, log)
		actionPage := newPageAction(mainPage)
		mockItemData.On("Send", mock.Anything, mock.Anything).Return(&model_data.DataByUUIDResponse{IsFile: true}, nil)
		pa := newPageDataGrid(mainPage, actionPage)
		m, _ := pa.Update(msg)
		assert.NotNil(t, m)

	})

	t.Run("enter IsText", func(t *testing.T) {
		mockManagerController := new(MockManagerController)
		mockItemData := new(MockItemDataController)
		mockGridData := new(MockGridDataController)
		mockManagerController.On("ItemData").Return(mockItemData)
		mockManagerController.On("GridData").Return(mockGridData)

		responseData := new(controller.GridDataResponse)
		responseData.Items = []model_data.ItemDataResponse{
			{
				Number: "1",
				Type:   "type1",
				Name:   "name1",
				UUID:   "111",
			},
			{
				Number: "2",
				Type:   "type2",
				Name:   "name2",
				UUID:   "22222",
			},
		}
		mockGridData.On("Send", mock.Anything).Return(responseData, nil)

		mainPage := newPageIndex(mockManagerController, memoryStorage, log)
		actionPage := newPageAction(mainPage)
		mockItemData.On("Send", mock.Anything, mock.Anything).Return(&model_data.DataByUUIDResponse{IsText: true}, nil)
		pa := newPageDataGrid(mainPage, actionPage)
		m, _ := pa.Update(msg)
		assert.NotNil(t, m)

	})

	t.Run("enter IsCard", func(t *testing.T) {
		mockManagerController := new(MockManagerController)
		mockItemData := new(MockItemDataController)
		mockGridData := new(MockGridDataController)
		mockManagerController.On("ItemData").Return(mockItemData)
		mockManagerController.On("GridData").Return(mockGridData)

		responseData := new(controller.GridDataResponse)
		responseData.Items = []model_data.ItemDataResponse{
			{
				Number: "1",
				Type:   "type1",
				Name:   "name1",
				UUID:   "111",
			},
			{
				Number: "2",
				Type:   "type2",
				Name:   "name2",
				UUID:   "22222",
			},
		}
		mockGridData.On("Send", mock.Anything).Return(responseData, nil)

		mainPage := newPageIndex(mockManagerController, memoryStorage, log)
		actionPage := newPageAction(mainPage)
		mockItemData.On("Send", mock.Anything, mock.Anything).Return(&model_data.DataByUUIDResponse{IsCard: true}, nil)
		pa := newPageDataGrid(mainPage, actionPage)
		m, _ := pa.Update(msg)
		assert.NotNil(t, m)

	})

	t.Run("enter no type", func(t *testing.T) {
		mockManagerController := new(MockManagerController)
		mockItemData := new(MockItemDataController)
		mockGridData := new(MockGridDataController)
		mockManagerController.On("ItemData").Return(mockItemData)
		mockManagerController.On("GridData").Return(mockGridData)

		responseData := new(controller.GridDataResponse)
		responseData.Items = []model_data.ItemDataResponse{
			{
				Number: "1",
				Type:   "type1",
				Name:   "name1",
				UUID:   "111",
			},
			{
				Number: "2",
				Type:   "type2",
				Name:   "name2",
				UUID:   "22222",
			},
		}
		mockGridData.On("Send", mock.Anything).Return(responseData, nil)

		mainPage := newPageIndex(mockManagerController, memoryStorage, log)
		actionPage := newPageAction(mainPage)
		mockItemData.On("Send", mock.Anything, mock.Anything).Return(&model_data.DataByUUIDResponse{}, nil)
		pa := newPageDataGrid(mainPage, actionPage)
		m, _ := pa.Update(msg)
		assert.NotNil(t, m)

	})
}

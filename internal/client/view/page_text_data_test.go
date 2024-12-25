package view

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/controller"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/storage"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTextData struct {
	mock.Mock
}

func (m *mockTextData) Send(token string, requestData *model_data.TextDataRequest) (*controller.TextDataResponse, error) {
	args := m.Called(token, requestData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*controller.TextDataResponse), args.Error(1)
}

func TestNewPageTextData(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = "testpath"
	mockCfg.Value().PathKeys = "testpath"
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	mainPage := newPageIndex(manager, memoryStorage, log)

	page := newPageTextData(mainPage)
	assert.NotNil(t, page)
	assert.NotNil(t, page.name)
	assert.NotNil(t, page.text)
	assert.NotNil(t, page.meta1)
	assert.NotNil(t, page.meta2)
}

func TestPageTextData_SetEditableData(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = "testpath"
	mockCfg.Value().PathKeys = "testpath"
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	mainPage := newPageIndex(manager, memoryStorage, log)

	page := newPageTextData(mainPage)
	data := &model_data.TextDataRequest{
		UUID:  "test-uuid",
		Name:  "test-name",
		Value: "test-value",
		Meta: map[string]string{
			data_type.MetaNameNote:    "test-note",
			data_type.MetaNameWebSite: "test-website",
		},
	}
	page.SetEditableData(data)
	assert.Equal(t, "test-uuid", page.uuid)
	assert.Equal(t, "test-name", page.name.Value())
	assert.Equal(t, "test-value", page.text.Value())
	assert.Equal(t, "test-note", page.meta1.Value())
	assert.Equal(t, "test-website", page.meta2.Value())
	assert.True(t, page.isEditable)

}

func TestPageTextData_SetPageGrid(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = "testpath"
	mockCfg.Value().PathKeys = "testpath"
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	mainPage := newPageIndex(manager, memoryStorage, log)

	page := newPageTextData(mainPage)
	gridPage := &pageDataGrid{}
	page.SetPageGrid(gridPage)
	assert.Equal(t, gridPage, page.gridPage)
}

func TestPageTextData_Init(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = "testpath"
	mockCfg.Value().PathKeys = "testpath"
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	mainPage := newPageIndex(manager, memoryStorage, log)

	page := newPageTextData(mainPage)
	cmd := page.Init()
	assert.NotNil(t, cmd)
}

func TestPageTextData_Update(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = "testpath"
	mockCfg.Value().PathKeys = "testpath"
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	mainPage := newPageIndex(manager, memoryStorage, log)

	textDataCtrl := new(mockTextData)
	manager.On("TextData").Return(textDataCtrl)

	page := newPageTextData(mainPage)

	page.Choice = 0
	msg := tea.KeyMsg{Type: tea.KeyDown}
	model, cmd := page.Update(msg)
	assert.Equal(t, 1, model.(*pageTextData).Choice)
	assert.Nil(t, cmd)

	page.Choice = 1
	msg = tea.KeyMsg{Type: tea.KeyUp}
	model, cmd = page.Update(msg)
	assert.Equal(t, 0, model.(*pageTextData).Choice)
	assert.Nil(t, cmd)

	page.Choice = 0
	msg = tea.KeyMsg{Type: tea.KeyTab}
	model, cmd = page.Update(msg)
	assert.Equal(t, 1, model.(*pageTextData).Choice)
	assert.Nil(t, cmd)

	page.Choice = 4
	page.uuid = "test-uuid"
	page.name.SetValue("test-name")
	page.text.SetValue("test-value")
	page.meta1.SetValue("test-note")
	page.meta2.SetValue("test-website")
	response := new(controller.TextDataResponse)
	textDataCtrl.On("Send", "test-token", &model_data.TextDataRequest{
		UUID:  "test-uuid",
		Name:  "test-name",
		Value: "test-value",
		Meta: map[string]string{
			data_type.MetaNameNote:    "test-note",
			data_type.MetaNameWebSite: "test-website",
		},
	}).Return(response, nil)

	memoryStorage.SetToken("test-token")
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = page.Update(msg)
	assert.Equal(t, "Данные сохранены", page.responseMessage)
	assert.IsType(t, &pageAction{}, model)

	page.Choice = 4
	page.uuid = "test-uuid"
	page.name.SetValue("test-name")
	page.text.SetValue("test-value")
	page.meta1.SetValue("test-note")
	page.meta2.SetValue("test-website")
	textDataCtrl.On("Send", "test-token2", &model_data.TextDataRequest{
		UUID:  "test-uuid",
		Name:  "test-name",
		Value: "test-value",
		Meta: map[string]string{
			data_type.MetaNameNote:    "test-note",
			data_type.MetaNameWebSite: "test-website",
		},
	}).Return(nil, errors.New("send failed"))
	memoryStorage.SetToken("test-token2")
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = page.Update(msg)
	assert.Equal(t, "send failed", page.responseMessage)
	assert.IsType(t, &pageTextData{}, model)

	page.Choice = 5
	page.isEditable = true
	gridPage := &pageDataGrid{}
	page.SetPageGrid(gridPage)
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = page.Update(msg)
	assert.Equal(t, gridPage, model)
	assert.Nil(t, cmd)

	page.Choice = 5
	page.isEditable = false
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = page.Update(msg)
	assert.IsType(t, &pageAction{}, model)
	assert.Nil(t, cmd)

	page.Choice = 0
	msg = tea.KeyMsg{Type: tea.KeyLeft}
	model, cmd = page.Update(msg)
	assert.Equal(t, 0, page.Choice)
	assert.NotNil(t, cmd)

	page.Choice = 5
	msg = tea.KeyMsg{Type: tea.KeyDown}
	model, cmd = page.Update(msg)
	assert.Equal(t, 5, page.Choice)
	assert.Nil(t, cmd)

	page.Choice = 0
	msg = tea.KeyMsg{Type: tea.KeyUp}
	model, cmd = page.Update(msg)
	assert.Equal(t, 0, page.Choice)
	assert.Nil(t, cmd)
}

func TestPageTextData_View(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = "testpath"
	mockCfg.Value().PathKeys = "testpath"
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	mainPage := newPageIndex(manager, memoryStorage, log)

	page := newPageTextData(mainPage)

	page.Choice = 0
	page.name.SetValue("test-name")
	page.text.SetValue("test-value")
	page.meta1.SetValue("test-note")
	page.meta2.SetValue("test-website")
	page.responseMessage = ""
	view := page.View()
	assert.Contains(t, view, "Произвольные текстовые данные")
	assert.Contains(t, view, "test-name")
	assert.Contains(t, view, "test-value")
	assert.Contains(t, view, "test-note")
	assert.Contains(t, view, "test-website")
	assert.Contains(t, view, "вверх/вниз: для переключения")
	assert.Contains(t, view, "enter: начать ввод значения")
	assert.Contains(t, view, "Отправить")
	assert.Contains(t, view, "Вернуться")
}

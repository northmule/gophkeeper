package view

import (
	"os"
	"path"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
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

type mockFileData struct {
	mock.Mock
}

func (m *mockFileData) Send(token string, requestData *model_data.FileDataInitRequest) (*controller.FileDataResponse, error) {
	args := m.Called(token, requestData)
	return args.Get(0).(*controller.FileDataResponse), args.Error(1)
}

func (m *mockFileData) UploadFile(token, uploadPath string, file *os.File) error {
	args := m.Called(token, uploadPath, file)
	return args.Error(0)
}

func (m *mockFileData) DownLoadFile(token, fileName, uuid string) error {
	args := m.Called(token, fileName, uuid)
	return args.Error(0)
}

func TestNewPageFileData(t *testing.T) {
	mainPage := &pageIndex{}
	page := newPageFileData(mainPage)
	assert.NotNil(t, page)
	assert.Equal(t, mainPage, page.mainPage)
	assert.Equal(t, 0, page.Choice)
	assert.Equal(t, false, page.Chosen)
	assert.Equal(t, "", page.responseMessage)
	assert.Nil(t, page.err)
	assert.Equal(t, "", page.selectedFile)
	assert.Equal(t, "", page.uuid)
	assert.Equal(t, "", page.mimeType)
	assert.Equal(t, "", page.extension)
	assert.Equal(t, "", page.fileName)
	assert.Equal(t, int64(0), page.size)
	assert.Equal(t, "Название данных", page.name.Placeholder)
	assert.Equal(t, "Для выбора файла нажмите enter", page.filePath.Placeholder)
	assert.Equal(t, data_type.TranslateDataType(data_type.MetaNameNote), page.meta1.Placeholder)
	assert.Equal(t, data_type.TranslateDataType(data_type.MetaNameWebSite), page.meta2.Placeholder)
}

func TestSetEditableData(t *testing.T) {
	page := &pageFileData{}
	data := &model_data.FileDataInitRequest{
		UUID:     "test-uuid",
		Name:     "test-name",
		FileName: "test-file-name",
		Meta: map[string]string{
			data_type.MetaNameNote:    "test-note",
			data_type.MetaNameWebSite: "test-website",
		},
	}
	page.SetEditableData(data)
	assert.Equal(t, "test-uuid", page.uuid)
	assert.Equal(t, "test-name", page.name.Value())
	assert.Equal(t, "test-file-name", page.filePath.Value())
	assert.Equal(t, "test-file-name", page.fileName)
	assert.Equal(t, "test-note", page.meta1.Value())
	assert.Equal(t, "test-website", page.meta2.Value())
	assert.Equal(t, true, page.isEditable)
}

func TestSetPageGrid(t *testing.T) {
	page := &pageFileData{}
	gridPage := &pageDataGrid{}
	page.SetPageGrid(gridPage)
	assert.Equal(t, gridPage, page.gridPage)
}

func TestInit(t *testing.T) {
	page := &pageFileData{
		selectedFile: "test-file-path",
	}
	page.Init()
	assert.Equal(t, "test-file-path", page.filePath.Value())
}

func TestUpdate(t *testing.T) {

	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	fileDataCtrl := new(mockFileData)
	manager := new(MockManagerController)
	manager.On("FileData").Return(fileDataCtrl)

	mainPage := newPageIndex(manager, memoryStorage, log)
	memoryStorage.SetToken("test-token")

	page := newPageFileData(mainPage)
	page.name.SetValue("test-name")
	page.filePath.SetValue("test-file-path")
	page.meta1.SetValue("test-note")
	page.meta2.SetValue("test-website")

	msg := tea.KeyMsg{Type: tea.KeyDown}
	_, cmd := page.Update(msg)
	assert.Equal(t, 1, page.Choice)
	assert.Nil(t, cmd)

	page.Choice = 4

	page.mainPage.managerController = manager
	page.selectedFile = "tmp_file"
	tempFile, err := os.Create(path.Join(page.selectedFile))
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	fileDataCtrl.On("Send", "test-token", mock.Anything).Return(&controller.FileDataResponse{UploadPath: "test-upload-path"}, nil)
	fileDataCtrl.On("UploadFile", "test-token", "test-upload-path", mock.Anything).Return(nil)

	msg = tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd = page.Update(msg)
	assert.Equal(t, "Файл загружен", page.responseMessage)
	assert.NotNil(t, cmd)

	page.Choice = 5
	page.uuid = "test-uuid"
	page.fileName = "test-file-name"
	fileDataCtrl.On("DownLoadFile", "test-token", "test-file-name", "test-uuid").Return(nil)

	msg = tea.KeyMsg{Type: tea.KeyCtrlD}
	_, cmd = page.Update(msg)
	assert.Equal(t, "Файл получен", page.responseMessage)
	assert.NotNil(t, cmd)
}

func TestView(t *testing.T) {
	page := &pageFileData{
		name:     textinput.New(),
		filePath: textinput.New(),
		meta1:    textinput.New(),
		meta2:    textinput.New(),
	}
	page.name.SetValue("test-name")
	page.filePath.SetValue("test-file-path")
	page.meta1.SetValue("test-note")
	page.meta2.SetValue("test-website")

	view := page.View()
	assert.Contains(t, view, "test-name")
	assert.Contains(t, view, "test-file-path")
	assert.Contains(t, view, "test-note")
	assert.Contains(t, view, "test-website")
}

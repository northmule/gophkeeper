package view

import (
	"errors"
	"path"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/controller"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/service"
	"github.com/northmule/gophkeeper/internal/client/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPageAuthentication_Init(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	cryptService, _ := service.NewCrypt(mockCfg)
	manager, _ := controller.NewManager(mockCfg, cryptService, log)
	memoryStorage := storage.NewMemoryStorage()
	mainPage := newPageIndex(manager, memoryStorage, log)

	page := newPageAuthentication(mainPage)
	cmd := page.Init()
	assert.NotNil(t, cmd)
}

func TestPageAuthentication_Update(t *testing.T) {
	tests := []struct {
		name     string
		msg      tea.Msg
		expected pageAction
	}{
		{"down key", tea.KeyMsg{Type: tea.KeyDown}, pageAction{Choice: 1}},
		{"up key", tea.KeyMsg{Type: tea.KeyUp}, pageAction{Choice: 0}},
		{"enter key on choice 0", tea.KeyMsg{Type: tea.KeyEnter}, pageAction{Choice: 0}},
		{"invalid key", tea.KeyMsg{Type: tea.KeyCtrlC}, pageAction{Choice: 0}},
	}

	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	cryptService, _ := service.NewCrypt(mockCfg)
	manager, _ := controller.NewManager(mockCfg, cryptService, log)
	memoryStorage := storage.NewMemoryStorage()
	mainPage := newPageIndex(manager, memoryStorage, log)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := newPageAuthentication(mainPage)
			_, _ = page.Update(tt.msg)
			assert.Equal(t, tt.expected.Choice, page.Choice)
		})
	}

	// прочие кейсы
	t.Run("choice 2 positive", func(t *testing.T) {

		mockManagerController := new(MockManagerController)
		mockAuthentication := new(MockAuthenticationDataController)

		mockKeyData := new(MockKeyDataController)
		mockManagerController.On("Authentication").Return(mockAuthentication)
		mockManagerController.On("KeysData").Return(mockKeyData)

		mockKeyData.On("UploadClientPublicKey", mock.Anything).Return(nil)
		mockKeyData.On("DownloadPublicServerKey", mock.Anything).Return(nil)
		mockKeyData.On("UploadClientPrivateKey", mock.Anything).Return(nil)

		mockAuthentication.On("Send", mock.Anything, mock.Anything).Return(&controller.AuthenticationResponse{Value: "ok"}, nil)

		mainPage := newPageIndex(mockManagerController, memoryStorage, log)
		pa := pageAuthentication{Choice: 2, mainPage: mainPage}
		pa.login.SetValue("login")
		pa.login.SetValue("password")
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.NotNil(t, m)
	})

	t.Run("choice 2 negative 0", func(t *testing.T) {

		mockManagerController := new(MockManagerController)
		mockAuthentication := new(MockAuthenticationDataController)
		mockManagerController.On("Authentication").Return(mockAuthentication)

		mockAuthentication.On("Send", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

		mainPage := newPageIndex(mockManagerController, memoryStorage, log)
		pa := pageAuthentication{Choice: 2, mainPage: mainPage}
		pa.login.SetValue("login")
		pa.login.SetValue("password")
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.NotEmpty(t, pa.responseMessage)
		assert.NotNil(t, m)
	})

	t.Run("choice 2 negative 1", func(t *testing.T) {

		mockManagerController := new(MockManagerController)
		mockAuthentication := new(MockAuthenticationDataController)

		mockKeyData := new(MockKeyDataController)
		mockManagerController.On("Authentication").Return(mockAuthentication)
		mockManagerController.On("KeysData").Return(mockKeyData)

		mockKeyData.On("UploadClientPublicKey", mock.Anything).Return(errors.New("error"))
		mockKeyData.On("DownloadPublicServerKey", mock.Anything).Return(nil)
		mockKeyData.On("UploadClientPrivateKey", mock.Anything).Return(nil)

		mockAuthentication.On("Send", mock.Anything, mock.Anything).Return(&controller.AuthenticationResponse{Value: "ok"}, nil)

		mainPage := newPageIndex(mockManagerController, memoryStorage, log)
		pa := pageAuthentication{Choice: 2, mainPage: mainPage}
		pa.login.SetValue("login")
		pa.login.SetValue("password")
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.NotEmpty(t, pa.responseMessage)
		assert.NotNil(t, m)
	})

	t.Run("choice 2 negative 2", func(t *testing.T) {

		mockManagerController := new(MockManagerController)
		mockAuthentication := new(MockAuthenticationDataController)

		mockKeyData := new(MockKeyDataController)
		mockManagerController.On("Authentication").Return(mockAuthentication)
		mockManagerController.On("KeysData").Return(mockKeyData)

		mockKeyData.On("UploadClientPublicKey", mock.Anything).Return(nil)
		mockKeyData.On("DownloadPublicServerKey", mock.Anything).Return(errors.New("error"))
		mockKeyData.On("UploadClientPrivateKey", mock.Anything).Return(nil)

		mockAuthentication.On("Send", mock.Anything, mock.Anything).Return(&controller.AuthenticationResponse{Value: "ok"}, nil)

		mainPage := newPageIndex(mockManagerController, memoryStorage, log)
		pa := pageAuthentication{Choice: 2, mainPage: mainPage}
		pa.login.SetValue("login")
		pa.login.SetValue("password")
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.NotEmpty(t, pa.responseMessage)
		assert.NotNil(t, m)
	})

	t.Run("choice 2 negative 3", func(t *testing.T) {

		mockManagerController := new(MockManagerController)
		mockAuthentication := new(MockAuthenticationDataController)

		mockKeyData := new(MockKeyDataController)
		mockManagerController.On("Authentication").Return(mockAuthentication)
		mockManagerController.On("KeysData").Return(mockKeyData)

		mockKeyData.On("UploadClientPublicKey", mock.Anything).Return(nil)
		mockKeyData.On("DownloadPublicServerKey", mock.Anything).Return(nil)
		mockKeyData.On("UploadClientPrivateKey", mock.Anything).Return(errors.New("error"))

		mockAuthentication.On("Send", mock.Anything, mock.Anything).Return(&controller.AuthenticationResponse{Value: "ok"}, nil)

		mainPage := newPageIndex(mockManagerController, memoryStorage, log)
		pa := pageAuthentication{Choice: 2, mainPage: mainPage}
		pa.login.SetValue("login")
		pa.login.SetValue("password")
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.NotEmpty(t, pa.responseMessage)
		assert.NotNil(t, m)
	})

	t.Run("choice 3", func(t *testing.T) {
		pa := pageAuthentication{Choice: 3, mainPage: mainPage}
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.NotNil(t, m)
	})

	t.Run("choice 10", func(t *testing.T) {
		pa := pageAuthentication{Choice: 10, mainPage: mainPage}
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.NotNil(t, m)
	})

}

func TestPageAuthentication_View(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	cryptService, _ := service.NewCrypt(mockCfg)
	manager, _ := controller.NewManager(mockCfg, cryptService, log)
	memoryStorage := storage.NewMemoryStorage()
	mainPage := newPageIndex(manager, memoryStorage, log)

	page := newPageAuthentication(mainPage)
	result := page.View()
	assert.True(t, strings.Contains(result, "Авторизация"))
	assert.True(t, strings.Contains(result, "вверх/вниз: для переключения"))
	assert.True(t, strings.Contains(result, "enter: начать ввод значения"))
	assert.True(t, strings.Contains(result, "Отправить"))
	assert.True(t, strings.Contains(result, "Вернуться"))

}

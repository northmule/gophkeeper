package view

import (
	"errors"
	"path"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/controller"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRegistration struct {
	mock.Mock
}

func (m *mockRegistration) Send(login, password, email string) (*controller.RegistrationResponse, error) {
	args := m.Called(login, password, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*controller.RegistrationResponse), args.Error(1)
}

func TestPageRegistration_Init(t *testing.T) {
	mockCfg, _ := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	mainPage := newPageIndex(manager, memoryStorage, log)

	pg := newPageRegistration(mainPage)
	cmd := pg.Init()
	assert.NotNil(t, cmd)
}

func TestPageRegistration_Update(t *testing.T) {
	mockCfg, _ := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	mainPage := newPageIndex(manager, memoryStorage, log)

	registrationCtrl := new(mockRegistration)
	manager.On("Registration").Return(registrationCtrl)

	pg := newPageRegistration(mainPage)

	msg := tea.KeyMsg{Type: tea.KeyDown}
	model, cmd := pg.Update(msg)
	assert.Equal(t, 1, model.(*pageRegistration).Choice)
	assert.Nil(t, cmd)

	msg = tea.KeyMsg{Type: tea.KeyUp}
	model, cmd = pg.Update(msg)
	assert.Equal(t, 0, model.(*pageRegistration).Choice)
	assert.Nil(t, cmd)

	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = pg.Update(msg)
	assert.Equal(t, 0, model.(*pageRegistration).Choice)

	// Test with "enter" and Choice = 1
	pg.Choice = 1
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = pg.Update(msg)
	assert.Equal(t, 1, model.(*pageRegistration).Choice)

	pg.Choice = 2
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = pg.Update(msg)
	assert.Equal(t, 2, model.(*pageRegistration).Choice)

	pg.Choice = 3
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = pg.Update(msg)
	assert.Equal(t, 3, model.(*pageRegistration).Choice)

	pg.Choice = 4
	pg.password.SetValue("password1")
	pg.password2.SetValue("password2")
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = pg.Update(msg)
	assert.Equal(t, 4, model.(*pageRegistration).Choice)
	assert.NotNil(t, cmd)
	assert.Equal(t, "Пароли не совпадают", model.(*pageRegistration).responseMessage)

	pg.password.SetValue("")
	pg.password2.SetValue("")
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = pg.Update(msg)
	assert.Equal(t, 4, model.(*pageRegistration).Choice)
	assert.NotNil(t, cmd)
	assert.Equal(t, "Заполните все поля", model.(*pageRegistration).responseMessage)

	pg.login.SetValue("testuser")
	pg.password.SetValue("password1")
	pg.password2.SetValue("password1")
	pg.email.SetValue("test@example.com")

	response := new(controller.RegistrationResponse)
	response.Value = "ok"
	registrationCtrl.On("Send", "testuser", "password1", "test@example.com").Return(response, nil)
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = pg.Update(msg)
	assert.Equal(t, mainPage, model)
	assert.NotNil(t, cmd)
	assert.Equal(t, "Вы зарегистрированы", pg.responseMessage)

	pg.login.SetValue("testuser")
	pg.password.SetValue("password2")
	pg.password2.SetValue("password2")
	pg.email.SetValue("test2@example.com")
	registrationCtrl.On("Send", "testuser", "password2", "test2@example.com").Return(nil, errors.New("registration failed"))
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = pg.Update(msg)
	assert.Equal(t, 4, pg.Choice)
	assert.NotNil(t, cmd)
	assert.Equal(t, "registration failed", pg.responseMessage)

	pg.Choice = 5
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = pg.Update(msg)
	assert.Equal(t, mainPage, model)
	assert.Nil(t, cmd)

	pg.Choice = 0
	msg = tea.KeyMsg{Type: tea.KeyTab}
	model, cmd = pg.Update(msg)
	assert.Equal(t, 1, pg.Choice)
	assert.Nil(t, cmd)
}

func TestPageRegistration_View(t *testing.T) {
	mockCfg, _ := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	mainPage := newPageIndex(manager, memoryStorage, log)

	pg := newPageRegistration(mainPage)

	pg.Choice = 0
	pg.login.SetValue("testuser")
	pg.password.SetValue("password1")
	pg.password2.SetValue("password2")
	pg.email.SetValue("test@example.com")
	pg.responseMessage = ""
	view := pg.View()
	assert.Contains(t, view, "Регистрация")
	assert.Contains(t, view, "testuser")
	assert.Contains(t, view, "password1")
	assert.Contains(t, view, "password2")
	assert.Contains(t, view, "test@example.com")
	assert.Contains(t, view, "вверх/вниз: для переключения")
	assert.Contains(t, view, "enter: начать ввод значения")
}

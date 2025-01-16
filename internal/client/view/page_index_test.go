package view

import (
	"path"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/storage"
	"github.com/stretchr/testify/assert"
)

func TestPageIndex_Init(t *testing.T) {
	mockCfg, _ := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)

	pi := newPageIndex(manager, memoryStorage, log)
	cmd := pi.Init()
	assert.Nil(t, cmd)
}

func TestPageIndex_Update(t *testing.T) {
	mockCfg, _ := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	pi := newPageIndex(manager, memoryStorage, log)

	pi.Choice = 0
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd := pi.Update(msg)
	assert.NotNil(t, model)
	assert.Nil(t, cmd)

	pi.Choice = 1
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = pi.Update(msg)
	assert.NotNil(t, model)
	assert.Nil(t, cmd)

	pi.Choice = 2
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd = pi.Update(msg)
	assert.NotNil(t, model)
	assert.Nil(t, cmd)

	pi.Choice = 1
	msg = tea.KeyMsg{Type: tea.KeyUp}
	model, cmd = pi.Update(msg)
	assert.Equal(t, 0, model.(*pageIndex).Choice)
	assert.Nil(t, cmd)

	pi.Choice = 1
	msg = tea.KeyMsg{Type: tea.KeyDown}
	model, cmd = pi.Update(msg)
	assert.Equal(t, 2, model.(*pageIndex).Choice)
	assert.Nil(t, cmd)

	pi.Choice = 0
	msg = tea.KeyMsg{Type: tea.KeyCtrlC}
	model, cmd = pi.Update(msg)
	assert.True(t, model.(*pageIndex).Quitting)

	pi.Choice = 0
	msg = tea.KeyMsg{Type: tea.KeyTab}
	model, cmd = pi.Update(msg)
	assert.Equal(t, 0, model.(*pageIndex).Choice)
	assert.Nil(t, cmd)
}

func TestPageIndex_View(t *testing.T) {

	mockCfg, _ := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	pi := newPageIndex(manager, memoryStorage, log)

	pi.Choice = 0
	view := pi.View()
	assert.Contains(t, view, "Выберите действия")
	assert.Contains(t, view, "Авторизация")
	assert.Contains(t, view, "Регистрация")
	assert.Contains(t, view, "Справка")
	assert.Contains(t, view, "вверх/вниз: для переключения")
	assert.Contains(t, view, "enter: выбрать")
	assert.Contains(t, view, "q, esc: quit")
	assert.Contains(t, view, "Авторизация")
	assert.Contains(t, view, "Регистрация")
	assert.Contains(t, view, "Справка")

}

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

func TestPageHelp_Init(t *testing.T) {
	mockCfg, _ := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	mainPage := newPageIndex(manager, memoryStorage, log)

	ph := newPageHelp(mainPage)
	cmd := ph.Init()
	assert.Nil(t, cmd)
}

func TestPageHelp_Update(t *testing.T) {
	mockCfg, _ := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	mainPage := newPageIndex(manager, memoryStorage, log)

	ph := newPageHelp(mainPage)

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd := ph.Update(msg)
	assert.Equal(t, mainPage, model)
	assert.Nil(t, cmd)

	msg = tea.KeyMsg{Type: tea.KeyTab}
	model, cmd = ph.Update(msg)
	assert.Equal(t, ph, model)
	assert.Nil(t, cmd)
}

func TestPageHelp_View(t *testing.T) {
	mockCfg, _ := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	mainPage := newPageIndex(manager, memoryStorage, log)

	ph := newPageHelp(mainPage)

	ph.Choice = 0
	view := ph.View()
	assert.Contains(t, view, "Менеджер паролей GophKeeper")
	assert.Contains(t, view, "YP: 32 go  (2024)")
	assert.Contains(t, view, "Вернуться")

	ph.Choice = 1
	view = ph.View()
	assert.Contains(t, view, "Менеджер паролей GophKeeper")
	assert.Contains(t, view, "YP: 32 go  (2024)")
	assert.Contains(t, view, "Вернуться")
}

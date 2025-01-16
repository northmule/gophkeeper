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
	"github.com/stretchr/testify/assert"
)

func TestPageAction_Init(t *testing.T) {
	mockCfg, _ := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	cryptService, _ := service.NewCrypt(mockCfg)
	manager, _ := controller.NewManager(mockCfg, cryptService, log)
	memoryStorage := storage.NewMemoryStorage()
	mainPage := newPageIndex(manager, memoryStorage, log)

	page := newPageAction(mainPage)
	cmd := page.Init()
	assert.Nil(t, cmd)
}

func TestPageAction_Update(t *testing.T) {
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

	mockCfg, _ := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	cryptService, _ := service.NewCrypt(mockCfg)
	manager, _ := controller.NewManager(mockCfg, cryptService, log)
	memoryStorage := storage.NewMemoryStorage()
	mainPage := newPageIndex(manager, memoryStorage, log)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := newPageAction(mainPage)
			_, _ = page.Update(tt.msg)
			assert.Equal(t, tt.expected.Choice, page.Choice)
		})
	}
	// прочие кейсы
	t.Run("choice 1", func(t *testing.T) {
		pa := pageAction{Choice: 1, mainPage: mainPage}
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.NotNil(t, m)
	})
	t.Run("choice 2", func(t *testing.T) {
		pa := pageAction{Choice: 2, mainPage: mainPage}
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.NotNil(t, m)
	})
	t.Run("choice 3", func(t *testing.T) {
		pa := pageAction{Choice: 3, mainPage: mainPage}
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.NotNil(t, m)
	})
	t.Run("choice 4", func(t *testing.T) {
		pa := pageAction{Choice: 4, mainPage: mainPage}
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.NotNil(t, m)
	})
}

func TestPageAction_View(t *testing.T) {
	mockCfg, _ := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	cryptService, _ := service.NewCrypt(mockCfg)
	manager, _ := controller.NewManager(mockCfg, cryptService, log)
	memoryStorage := storage.NewMemoryStorage()
	mainPage := newPageIndex(manager, memoryStorage, log)

	page := newPageAction(mainPage)
	result := page.View()
	assert.True(t, strings.Contains(result, "Доступные действия"))
	assert.True(t, strings.Contains(result, "Добавить данные банковских карт"))
	assert.True(t, strings.Contains(result, "Добавить произвольные текстовые данные"))
	assert.True(t, strings.Contains(result, "Добавить бинарные данные"))
	assert.True(t, strings.Contains(result, "Показать мои данные"))
	assert.True(t, strings.Contains(result, "Выйти"))
	assert.True(t, strings.Contains(result, "вверх/вниз: для переключения • enter: выбрать"))

}

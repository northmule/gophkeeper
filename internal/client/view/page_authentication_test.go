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

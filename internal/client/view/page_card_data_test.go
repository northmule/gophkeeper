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

func TestPageCardData_Init(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	cryptService, _ := service.NewCrypt(mockCfg)
	manager, _ := controller.NewManager(mockCfg, cryptService, log)
	memoryStorage := storage.NewMemoryStorage()
	mainPage := newPageIndex(manager, memoryStorage, log)
	page := newPageCardData(mainPage)
	cmd := page.Init()
	assert.NotNil(t, cmd)
}

func TestPageCardData_Update(t *testing.T) {
	tests := []struct {
		name     string
		msg      tea.Msg
		expected pageCardData
	}{
		{"down key", tea.KeyMsg{Type: tea.KeyDown}, pageCardData{Choice: 1}},
		{"up key", tea.KeyMsg{Type: tea.KeyUp}, pageCardData{Choice: 0}},
		{"invalid key", tea.KeyMsg{Type: tea.KeyCtrlC}, pageCardData{Choice: 0}},
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
			page := newPageCardData(mainPage)
			_, _ = page.Update(tt.msg)
			assert.Equal(t, tt.expected.Choice, page.Choice)
		})
	}
}

func TestPageCardData_View(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	cryptService, _ := service.NewCrypt(mockCfg)
	manager, _ := controller.NewManager(mockCfg, cryptService, log)
	memoryStorage := storage.NewMemoryStorage()
	mainPage := newPageIndex(manager, memoryStorage, log)

	page := newPageCardData(mainPage)
	result := page.View()
	assert.True(t, strings.Contains(result, "Данные банковских карт"))
	assert.True(t, strings.Contains(result, "enter: начать ввод значения"))
	assert.True(t, strings.Contains(result, "Отправить"))
	assert.True(t, strings.Contains(result, "Вернуться"))
}

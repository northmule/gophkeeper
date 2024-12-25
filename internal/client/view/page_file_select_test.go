package view

import (
	"errors"
	"path"
	"testing"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/storage"
	"github.com/stretchr/testify/assert"
)

func TestPageFileSelect_Init(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)

	mainPage := newPageIndex(manager, memoryStorage, log)
	memoryStorage.SetToken("test-token")
	pfd := newPageFileData(mainPage)
	pfs := newPageFileSelect(pfd)
	cmd := pfs.Init()
	assert.NotNil(t, cmd)
}

func TestPageFileSelect_View(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	memoryStorage := storage.NewMemoryStorage()
	manager := new(MockManagerController)
	mainPage := newPageIndex(manager, memoryStorage, log)
	memoryStorage.SetToken("test-token")

	pfd := newPageFileData(mainPage)
	pfs := newPageFileSelect(pfd)

	view := pfs.View()
	assert.Contains(t, view, "Выберите файл")
	assert.Contains(t, view, "Сейчас выбран: ")

	pfs.err = errors.New("test error")
	view = pfs.View()
	assert.Contains(t, view, "test error")

	pfs.selectedFile = "test.go"
	pfs.err = nil
	view = pfs.View()
	assert.Contains(t, view, "Сейчас выбран: test.go")
	assert.Contains(t, view, "Нажмите ещё enter что бы отправить")
}

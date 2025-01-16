package controller

import (
	"testing"

	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewManager_Success(t *testing.T) {
	mockConfig := makeMockConfig("")
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	manager, err := NewManager(mockConfig, cryptService, log)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.NotNil(t, manager.authentication)
	assert.NotNil(t, manager.cardData)
	assert.NotNil(t, manager.textData)
	assert.NotNil(t, manager.fileData)
	assert.NotNil(t, manager.gridData)
	assert.NotNil(t, manager.itemData)
	assert.NotNil(t, manager.keysData)
	assert.NotNil(t, manager.registration)
}

func TestManager_Authentication(t *testing.T) {
	mockConfig := makeMockConfig("")
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	manager, err := NewManager(mockConfig, cryptService, log)
	assert.NoError(t, err)

	auth := manager.Authentication()
	assert.NotNil(t, auth)
}

func TestManager_CardData(t *testing.T) {
	mockConfig := makeMockConfig("")
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	manager, err := NewManager(mockConfig, cryptService, log)
	assert.NoError(t, err)

	cardData := manager.CardData()
	assert.NotNil(t, cardData)
}

func TestManager_TextData(t *testing.T) {
	mockConfig := makeMockConfig("")
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	manager, err := NewManager(mockConfig, cryptService, log)
	assert.NoError(t, err)

	textData := manager.TextData()
	assert.NotNil(t, textData)
}

func TestManager_FileData(t *testing.T) {
	mockConfig := makeMockConfig("")
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	manager, err := NewManager(mockConfig, cryptService, log)
	assert.NoError(t, err)

	data := manager.FileData()
	assert.NotNil(t, data)
}
func TestManager_ItemData(t *testing.T) {
	mockConfig := makeMockConfig("")
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	manager, err := NewManager(mockConfig, cryptService, log)
	assert.NoError(t, err)

	data := manager.ItemData()
	assert.NotNil(t, data)
}

func TestManager_KeysData(t *testing.T) {
	mockConfig := makeMockConfig("")
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	manager, err := NewManager(mockConfig, cryptService, log)
	assert.NoError(t, err)

	data := manager.KeysData()
	assert.NotNil(t, data)
}
func TestManager_Registration(t *testing.T) {
	mockConfig := makeMockConfig("")
	log, err := logger.NewLogger("info")
	if err != nil {
		t.Errorf(err.Error())
	}
	cryptService := NewCryptMock(t)
	manager, err := NewManager(mockConfig, cryptService, log)
	assert.NoError(t, err)

	data := manager.Registration()
	assert.NotNil(t, data)
}

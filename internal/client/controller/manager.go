package controller

import (
	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/service"
)

// Manager менеджер контроллеров
type Manager struct {
	logger *logger.Logger

	// Контроллеры
	authentication *Authentication
	cardData       *CardData
	textData       *TextData
	fileData       *FileData
	gridData       *GridData
	itemData       *ItemData
	keysData       *KeysData

	cfg *config.Config
}

// NewManager конструктор
func NewManager(cfg *config.Config, logger *logger.Logger) (*Manager, error) {

	cryptService, err := service.NewCrypt(cfg)
	if err != nil {
		return nil, err
	}

	return &Manager{
		logger:         logger,
		authentication: NewAuthentication(cfg, logger),
		cardData:       NewCardData(cfg, cryptService, logger),
		textData:       NewTextData(cfg, cryptService, logger),
		fileData:       NewFileData(cfg, cryptService, logger),
		gridData:       NewGridData(cfg, logger),
		itemData:       NewItemData(cfg, logger),
		keysData:       NewKeysData(cfg, cryptService, logger),
	}, nil
}

// ManagerController интерфейс для передачи в модели
type ManagerController interface {
	Authentication() *Authentication
	CardData() *CardData
	TextData() *TextData
	FileData() *FileData
	GridData() *GridData
	ItemData() *ItemData
	KeysData() *KeysData
}

// Authentication контроллер
func (manager *Manager) Authentication() *Authentication {
	return manager.authentication
}

// CardData контроллер
func (manager *Manager) CardData() *CardData {
	return manager.cardData
}

// TextData контроллер
func (manager *Manager) TextData() *TextData {
	return manager.textData
}

// FileData контроллер
func (manager *Manager) FileData() *FileData {
	return manager.fileData
}

// GridData контроллер
func (manager *Manager) GridData() *GridData {
	return manager.gridData
}

// ItemData контроллер
func (manager *Manager) ItemData() *ItemData {
	return manager.itemData
}

// KeysData контроллер
func (manager *Manager) KeysData() *KeysData {
	return manager.keysData
}

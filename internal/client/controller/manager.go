package controller

import (
	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
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

	cfg *config.Config
}

// NewManager конструктор
func NewManager(cfg *config.Config, logger *logger.Logger) *Manager {
	return &Manager{
		logger:         logger,
		authentication: NewAuthentication(cfg, logger),
		cardData:       NewCardData(cfg, logger),
		textData:       NewTextData(cfg, logger),
		fileData:       NewFileData(cfg, logger),
		gridData:       NewGridData(cfg, logger),
	}
}

// ManagerController интерфейс для передачи в модели
type ManagerController interface {
	Authentication() *Authentication
	CardData() *CardData
	TextData() *TextData
	FileData() *FileData
	GridData() *GridData
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

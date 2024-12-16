package controller

import (
	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
)

// Manager менеджер контроллеров
type Manager struct {
	logger *logger.Logger

	authentication *Authentication
	cardData       *CardData
	textData       *TextData

	cfg *config.Config
}

// NewManager конструктор
func NewManager(cfg *config.Config, logger *logger.Logger) *Manager {
	return &Manager{
		logger:         logger,
		authentication: NewAuthentication(cfg, logger),
		cardData:       NewCardData(cfg, logger),
		textData:       NewTextData(cfg, logger),
	}
}

// ManagerController интерфейс для передачи в модели
type ManagerController interface {
	Authentication() *Authentication
	CardData() *CardData
	TextData() *TextData
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
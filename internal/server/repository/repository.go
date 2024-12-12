package repository

import (
	"github.com/northmule/gophkeeper/internal/server/storage"
)

// Manager общий менеджер репозитариев для проекта
type Manager struct {
	user     *UserRepository
	cardData *CardDataRepository
	owner    *OwnerRepository
	metaData *MetaDataRepository
	textData *TextDataRepository
}

// NewManager конструктор
func NewManager(db storage.DBQuery) (Repository, error) {
	userRepository, err := NewUserRepository(db)
	if err != nil {
		return nil, err
	}
	cardDataRepository, err := NewCardDataRepository(db)
	if err != nil {
		return nil, err
	}
	ownerRepository, err := NewOwnerRepository(db)
	if err != nil {
		return nil, err
	}
	metaDataRepository, err := NewMetaDataRepository(db)
	if err != nil {
		return nil, err
	}
	textDataRepository, err := NewTextDataRepository(db)
	if err != nil {
		return nil, err
	}
	instance := &Manager{
		user:     userRepository,
		cardData: cardDataRepository,
		owner:    ownerRepository,
		metaData: metaDataRepository,
		textData: textDataRepository,
	}

	return instance, nil
}

type Repository interface {
	User() *UserRepository
	CardData() *CardDataRepository
	Owner() *OwnerRepository
	MetaData() *MetaDataRepository
	TextData() *TextDataRepository
}

// User репозитацрий
func (m *Manager) User() *UserRepository {
	return m.user
}

// CardData репозитацрий
func (m *Manager) CardData() *CardDataRepository {
	return m.cardData
}

// Owner репозитацрий
func (m *Manager) Owner() *OwnerRepository {
	return m.owner
}

// MetaData репозитацрий
func (m *Manager) MetaData() *MetaDataRepository {
	return m.metaData
}

// TextData репозитацрий
func (m *Manager) TextData() *TextDataRepository {
	return m.textData
}

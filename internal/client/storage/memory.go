package storage

import (
	"sync"

	"github.com/northmule/gophkeeper/internal/common/models"
)

// MemoryStorage хранилище данных на время запуска
type MemoryStorage struct {
	token string // текущий токен авторизации

	// данные синхронизации, ключами явлюятся uuid этих данных
	cardDataList map[string]models.CardData
	metaDataList []models.MetaData
	textDataList map[string]models.TextData
	fileDataList map[string]models.FileData

	mx sync.RWMutex
}

// Storage интерфейс локального хранилища
type Storage interface {
	SetToken(token string)
	Token() string
	ResetToken()
	AddCardDataList(data models.CardData) error
	AddMetaDataList(data models.MetaData) error
	AddTextData(data models.TextData) error
	AddFileData(data models.FileData) error
}

// NewMemoryStorage конструктор
func NewMemoryStorage() *MemoryStorage {
	instance := new(MemoryStorage)

	instance.cardDataList = make(map[string]models.CardData)
	instance.metaDataList = make([]models.MetaData, 0)
	instance.textDataList = make(map[string]models.TextData)
	instance.fileDataList = make(map[string]models.FileData)

	return instance
}

// SetToken добавить токен
func (s *MemoryStorage) SetToken(token string) {
	s.token = token
}

// Token значение токена
func (s *MemoryStorage) Token() string {
	return s.token
}

// ResetToken сбросить токен
func (s *MemoryStorage) ResetToken() {
	s.token = ""
}

// AddCardDataList добавляет или заменяет данные
func (s *MemoryStorage) AddCardDataList(data models.CardData) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.cardDataList[data.UUID] = data

	return nil
}

// AddMetaDataList добавляет или заменяет данные
func (s *MemoryStorage) AddMetaDataList(data models.MetaData) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	var isExist bool
	for k, item := range s.metaDataList {
		if item.DataUUID == data.DataUUID {
			s.metaDataList[k] = data
			isExist = true
		}
	}
	if !isExist {
		s.metaDataList = append(s.metaDataList, data)
	}

	return nil
}

// AddTextData добавляет или заменяет данные
func (s *MemoryStorage) AddTextData(data models.TextData) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.textDataList[data.UUID] = data

	return nil
}

// AddFileData добавляет или заменяет данные
func (s *MemoryStorage) AddFileData(data models.FileData) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.fileDataList[data.UUID] = data

	return nil
}

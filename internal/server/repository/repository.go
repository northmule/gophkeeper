package repository

import (
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/storage"
	"golang.org/x/net/context"
)

// Manager общий менеджер репозитариев для проекта
type Manager struct {
	user     *UserRepository
	cardData *CardDataRepository
	owner    *OwnerRepository
	metaData *MetaDataRepository
	textData *TextDataRepository
	fileData *FileDataRepository
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
	fileDataRepository, err := NewFileDataRepository(db)
	if err != nil {
		return nil, err
	}
	instance := &Manager{
		user:     userRepository,
		cardData: cardDataRepository,
		owner:    ownerRepository,
		metaData: metaDataRepository,
		textData: textDataRepository,
		fileData: fileDataRepository,
	}

	return instance, nil
}

// Repository общий интерфейс репозитария
type Repository interface {
	User() UserDataModelRepository
	CardData() CardDataModelRepository
	Owner() OwnerDataModelRepository
	MetaData() MetaDataModelRepository
	TextData() TextDataModelRepository
	FileData() FileDataModelRepository
}

// UserDataModelRepository интерфейс запросов
type UserDataModelRepository interface {
	FindOneByLogin(ctx context.Context, login string) (*models.User, error)
	FindOneByUUID(ctx context.Context, uuid string) (*models.User, error)
	CreateNewUser(ctx context.Context, user models.User) (int64, error)
	TxCreateNewUser(ctx context.Context, tx storage.TxDBQuery, user models.User) (int64, error)
	SetPublicKey(ctx context.Context, data string, userUUID string) error
	SetPrivateClientKey(ctx context.Context, data string, userUUID string) error
}

// CardDataModelRepository интерфейс запросов
type CardDataModelRepository interface {
	FindOneByUUID(ctx context.Context, uuid string) (*models.CardData, error)
	Add(ctx context.Context, data *models.CardData) (int64, error)
	Update(ctx context.Context, data *models.CardData) error
}

// OwnerDataModelRepository интерфейс запросов
type OwnerDataModelRepository interface {
	FindOneByUserUUIDAndDataUUIDAndDataType(ctx context.Context, userUuid string, dataUuid string, dataType string) (*models.Owner, error)
	FindOneByUserUUIDAndDataUUID(ctx context.Context, userUuid string, dataUuid string) (*models.Owner, error)
	Add(ctx context.Context, data *models.Owner) (int64, error)
	AllOwnerData(ctx context.Context, userUUID string, offset int, limit int) ([]models.OwnerData, error)
}

// MetaDataModelRepository интерфейс запросов
type MetaDataModelRepository interface {
	FindOneByUUID(ctx context.Context, uuid string) ([]models.MetaData, error)
	Add(ctx context.Context, data *models.MetaData) (int64, error)
	ReplaceMetaByDataUUID(ctx context.Context, dataUUID string, metaDataList []models.MetaData) error
}

// TextDataModelRepository интерфейс запросов
type TextDataModelRepository interface {
	FindOneByUUID(ctx context.Context, uuid string) (*models.TextData, error)
	Add(ctx context.Context, data *models.TextData) (int64, error)
	Update(ctx context.Context, data *models.TextData) error
}

// FileDataModelRepository интерфейс запросов
type FileDataModelRepository interface {
	FindOneByUUID(ctx context.Context, uuid string) (*models.FileData, error)
	Add(ctx context.Context, data *models.FileData) (int64, error)
	Update(ctx context.Context, data *models.FileData) error
}

// User репозитацрий
func (m *Manager) User() UserDataModelRepository {
	return m.user
}

// CardData репозитацрий
func (m *Manager) CardData() CardDataModelRepository {
	return m.cardData
}

// Owner репозитацрий
func (m *Manager) Owner() OwnerDataModelRepository {
	return m.owner
}

// MetaData репозитацрий
func (m *Manager) MetaData() MetaDataModelRepository {
	return m.metaData
}

// TextData репозитацрий
func (m *Manager) TextData() TextDataModelRepository {
	return m.textData
}

// FileData репозитацрий
func (m *Manager) FileData() FileDataModelRepository {
	return m.fileData
}

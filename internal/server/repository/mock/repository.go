package mock

import (
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/repository"
	"github.com/northmule/gophkeeper/internal/server/storage"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

// MockUserDataModelRepository is a mock implementation of UserDataModelRepository
type MockUserDataModelRepository struct {
	mock.Mock
}

func (m *MockUserDataModelRepository) FindOneByLogin(ctx context.Context, login string) (*models.User, error) {
	args := m.Called(ctx, login)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserDataModelRepository) FindOneByUUID(ctx context.Context, uuid string) (*models.User, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserDataModelRepository) CreateNewUser(ctx context.Context, user models.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserDataModelRepository) TxCreateNewUser(ctx context.Context, tx storage.TxDBQuery, user models.User) (int64, error) {
	args := m.Called(ctx, tx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserDataModelRepository) SetPublicKey(ctx context.Context, data string, userUUID string) error {
	args := m.Called(ctx, data, userUUID)
	return args.Error(0)
}

func (m *MockUserDataModelRepository) SetPrivateClientKey(ctx context.Context, data string, userUUID string) error {
	args := m.Called(ctx, data, userUUID)
	return args.Error(0)
}

// MockCardDataModelRepository is a mock implementation of CardDataModelRepository
type MockCardDataModelRepository struct {
	mock.Mock
}

func (m *MockCardDataModelRepository) FindOneByUUID(ctx context.Context, uuid string) (*models.CardData, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CardData), args.Error(1)
}

func (m *MockCardDataModelRepository) Add(ctx context.Context, data *models.CardData) (int64, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCardDataModelRepository) Update(ctx context.Context, data *models.CardData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

// MockOwnerDataModelRepository is a mock implementation of OwnerDataModelRepository
type MockOwnerDataModelRepository struct {
	mock.Mock
}

func (m *MockOwnerDataModelRepository) FindOneByUserUUIDAndDataUUIDAndDataType(ctx context.Context, userUuid string, dataUuid string, dataType string) (*models.Owner, error) {
	args := m.Called(ctx, userUuid, dataUuid, dataType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Owner), args.Error(1)
}

func (m *MockOwnerDataModelRepository) FindOneByUserUUIDAndDataUUID(ctx context.Context, userUuid string, dataUuid string) (*models.Owner, error) {
	args := m.Called(ctx, userUuid, dataUuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Owner), args.Error(1)
}

func (m *MockOwnerDataModelRepository) Add(ctx context.Context, data *models.Owner) (int64, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockOwnerDataModelRepository) AllOwnerData(ctx context.Context, userUUID string, offset int, limit int) ([]models.OwnerData, error) {
	args := m.Called(ctx, userUUID, offset, limit)
	return args.Get(0).([]models.OwnerData), args.Error(1)
}

// MockMetaDataModelRepository is a mock implementation of MetaDataModelRepository
type MockMetaDataModelRepository struct {
	mock.Mock
}

func (m *MockMetaDataModelRepository) FindOneByUUID(ctx context.Context, uuid string) ([]models.MetaData, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).([]models.MetaData), args.Error(1)
}

func (m *MockMetaDataModelRepository) Add(ctx context.Context, data *models.MetaData) (int64, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMetaDataModelRepository) ReplaceMetaByDataUUID(ctx context.Context, dataUUID string, metaDataList []models.MetaData) error {
	args := m.Called(ctx, dataUUID, metaDataList)
	return args.Error(0)
}

// MockTextDataModelRepository is a mock implementation of TextDataModelRepository
type MockTextDataModelRepository struct {
	mock.Mock
}

func (m *MockTextDataModelRepository) FindOneByUUID(ctx context.Context, uuid string) (*models.TextData, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TextData), args.Error(1)
}

func (m *MockTextDataModelRepository) Add(ctx context.Context, data *models.TextData) (int64, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTextDataModelRepository) Update(ctx context.Context, data *models.TextData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

// MockFileDataModelRepository is a mock implementation of FileDataModelRepository
type MockFileDataModelRepository struct {
	mock.Mock
}

func (m *MockFileDataModelRepository) FindOneByUUID(ctx context.Context, uuid string) (*models.FileData, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FileData), args.Error(1)
}

func (m *MockFileDataModelRepository) Add(ctx context.Context, data *models.FileData) (int64, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockFileDataModelRepository) Update(ctx context.Context, data *models.FileData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

// MockManager is a mock implementation of Repository
type MockManager struct {
	mock.Mock
	user     *MockUserDataModelRepository
	cardData *MockCardDataModelRepository
	owner    *MockOwnerDataModelRepository
	metaData *MockMetaDataModelRepository
	textData *MockTextDataModelRepository
	fileData *MockFileDataModelRepository
}

func NewMockManager() *MockManager {
	instance := new(MockManager)
	instance.user = new(MockUserDataModelRepository)
	instance.cardData = new(MockCardDataModelRepository)
	instance.owner = new(MockOwnerDataModelRepository)
	instance.metaData = new(MockMetaDataModelRepository)
	instance.textData = new(MockTextDataModelRepository)
	instance.fileData = new(MockFileDataModelRepository)

	return instance

}

func (m *MockManager) User() repository.UserDataModelRepository {
	args := m.Called()
	return args.Get(0).(repository.UserDataModelRepository)
}

func (m *MockManager) CardData() repository.CardDataModelRepository {
	args := m.Called()
	return args.Get(0).(repository.CardDataModelRepository)
}

func (m *MockManager) Owner() repository.OwnerDataModelRepository {
	args := m.Called()
	return args.Get(0).(repository.OwnerDataModelRepository)
}

func (m *MockManager) MetaData() repository.MetaDataModelRepository {
	args := m.Called()
	return args.Get(0).(repository.MetaDataModelRepository)
}

func (m *MockManager) TextData() repository.TextDataModelRepository {
	args := m.Called()
	return args.Get(0).(repository.TextDataModelRepository)
}

func (m *MockManager) FileData() repository.FileDataModelRepository {
	args := m.Called()
	return args.Get(0).(repository.FileDataModelRepository)
}

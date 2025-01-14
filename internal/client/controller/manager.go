package controller

import (
	"os"

	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/service"
	"github.com/northmule/gophkeeper/internal/common/model_data"
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
	registration   *Registration

	cfg *config.Config
}

// NewManager конструктор
func NewManager(cfg *config.Config, cryptService service.Cryptographer, logger *logger.Logger) (*Manager, error) {

	return &Manager{
		logger:         logger,
		authentication: NewAuthentication(cfg, logger),
		cardData:       NewCardData(cfg, cryptService, logger),
		textData:       NewTextData(cfg, cryptService, logger),
		fileData:       NewFileData(cfg, cryptService, logger),
		gridData:       NewGridData(cfg, cryptService, logger),
		itemData:       NewItemData(cfg, cryptService, logger),
		keysData:       NewKeysData(cfg, cryptService, logger),
		registration:   NewRegistration(cfg, logger),
	}, nil
}

// ItemDataController контроллер
type ItemDataController interface {
	Send(token string, dataUUID string) (*model_data.DataByUUIDResponse, error)
}

// GridDataController контроллер
type GridDataController interface {
	Send(token string) (*GridDataResponse, error)
}

// FileDataController контроллер
type FileDataController interface {
	Send(token string, requestData *model_data.FileDataInitRequest) (*FileDataResponse, error)
	UploadFile(token string, url string, file *os.File) error
	DownLoadFile(token string, fileName string, dataUUID string) error
}

// RegistrationController контроллер
type RegistrationController interface {
	Send(login string, password string, email string) (*RegistrationResponse, error)
}

// TextDataController контроллер
type TextDataController interface {
	Send(token string, requestData *model_data.TextDataRequest) (*TextDataResponse, error)
}

// AuthenticationDataController контроллер
type AuthenticationDataController interface {
	Send(login string, password string) (*AuthenticationResponse, error)
}

// KeyDataController контроллер
type KeyDataController interface {
	UploadClientPublicKey(token string) error
	DownloadPublicServerKey(token string) error
	UploadClientPrivateKey(token string) error
}

type CardDataController interface {
	Send(token string, requestData *model_data.CardDataRequest) (*CardDataResponse, error)
}

// Authentication контроллер
func (manager *Manager) Authentication() AuthenticationDataController {
	return manager.authentication
}

// CardData контроллер
func (manager *Manager) CardData() CardDataController {
	return manager.cardData
}

// TextData контроллер
func (manager *Manager) TextData() TextDataController {
	return manager.textData
}

// FileData контроллер
func (manager *Manager) FileData() FileDataController {
	return manager.fileData
}

// GridData контроллер
func (manager *Manager) GridData() GridDataController {
	return manager.gridData
}

// ItemData контроллер
func (manager *Manager) ItemData() ItemDataController {
	return manager.itemData
}

// KeysData контроллер
func (manager *Manager) KeysData() KeyDataController {
	return manager.keysData
}

// Registration контроллер
func (manager *Manager) Registration() RegistrationController {
	return manager.registration
}

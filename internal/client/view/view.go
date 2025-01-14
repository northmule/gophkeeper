package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/gophkeeper/internal/client/controller"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/common/models"
	"golang.org/x/net/context"
)

// ClientView вьюха клиента
type ClientView struct {
	log           *logger.Logger
	manager       ManagerController
	memoryStorage Storage
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

// ManagerController интерфейс для передачи в модели
type ManagerController interface {
	Authentication() controller.AuthenticationDataController
	CardData() controller.CardDataController
	TextData() controller.TextDataController
	FileData() controller.FileDataController
	GridData() controller.GridDataController
	ItemData() controller.ItemDataController
	KeysData() controller.KeyDataController
	Registration() controller.RegistrationController
}

// NewClientView конструктор
func NewClientView(manager ManagerController, storage Storage, log *logger.Logger) *ClientView {
	instance := &ClientView{
		log:           log,
		manager:       manager,
		memoryStorage: storage,
	}

	return instance
}

// InitMain подготовка консольных форм
func (v *ClientView) InitMain(ctx context.Context) error {
	var err error

	p := tea.NewProgram(newPageIndex(v.manager, v.memoryStorage, v.log), tea.WithContext(ctx))
	if _, err = p.Run(); err != nil {
		return err
	}

	return nil
}

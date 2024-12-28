package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/controller"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/service"
	"github.com/northmule/gophkeeper/internal/client/storage"
	"golang.org/x/net/context"
)

// ClientView вьюха клиента
type ClientView struct {
	log *logger.Logger
	cfg *config.Config
}

// NewClientView конструктор
func NewClientView(cfg *config.Config, log *logger.Logger) *ClientView {
	instance := new(ClientView)
	instance.cfg = cfg
	instance.log = log
	return instance
}

// InitMain подготовка консольных форм
func (v *ClientView) InitMain(ctx context.Context) error {

	cryptService, err := service.NewCrypt(v.cfg)
	if err != nil {
		return err
	}
	manager, err := controller.NewManager(v.cfg, cryptService, v.log)
	if err != nil {
		return err
	}
	memoryStorage := storage.NewMemoryStorage()

	p := tea.NewProgram(newPageIndex(manager, memoryStorage, v.log), tea.WithContext(ctx))
	if _, err = p.Run(); err != nil {
		return err
	}

	return nil
}

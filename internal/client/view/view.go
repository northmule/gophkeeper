package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/controller"
	"github.com/northmule/gophkeeper/internal/client/logger"
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

	manager, err := controller.NewManager(v.cfg, v.log)
	if err != nil {
		return err
	}
	memoryStorage := storage.NewMemoryStorage()

	p := tea.NewProgram(newPageIndex(manager, memoryStorage, v.log))
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

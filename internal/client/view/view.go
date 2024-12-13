package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"golang.org/x/net/context"
)

// ClientView вьюха клиента
type ClientView struct {
	log *logger.Logger
}

// NewClientView конструктор
func NewClientView(log *logger.Logger) *ClientView {
	instance := new(ClientView)
	instance.log = log
	return instance
}

func (v *ClientView) InitMain(ctx context.Context) error {

	m := menu{
		options: []menuItem{
			menuItem{
				text:    "new check-in",
				onPress: func() tea.Msg { return toggleCasingMsg{} },
			},
			menuItem{
				text:    "view check-ins",
				onPress: func() tea.Msg { return toggleCasingMsg{} },
			},
		},
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil

}

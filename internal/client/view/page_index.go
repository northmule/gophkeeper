package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/gophkeeper/internal/client/controller"
	"github.com/northmule/gophkeeper/internal/client/logger"
)

// Экран сразу после запуска клиента
type pageIndex struct {
	Choice int
	Chosen bool

	Quitting          bool
	managerController controller.ManagerController
	log               *logger.Logger

	storage Storage
}

func newPageIndex(managerController controller.ManagerController, storage Storage, log *logger.Logger) *pageIndex {
	return &pageIndex{
		log:               log,
		storage:           storage,
		managerController: managerController,
	}
}

// Init Действия при инициализации (загрузка данных и т.д)
func (m *pageIndex) Init() tea.Cmd {
	return nil
}

// Update обновление
func (m *pageIndex) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
		if k == "down" {
			m.Choice++
			if m.Choice > 2 {
				m.Choice = 2
			}
		}
		if k == "up" {
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		}

		if k == "enter" {
			if m.Choice == 0 {
				return newPageAuthentication(m), nil
			}
			if m.Choice == 1 {
				return newPageRegistration(m), nil
			}
			if m.Choice == 2 {
				return newPageHelp(m), nil
			}

		}
	}

	return m, nil
}

// View вид модели( в том числе при старте)
func (m *pageIndex) View() string {
	c := m.Choice

	title := renderTitle("Выберите действия")
	tpl := "%s\n\n"
	tpl += subtleStyle.Render("вверх/вниз: для переключения") + dotStyle +
		subtleStyle.Render("enter: выбрать") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n",
		renderCheckbox("Авторизация", c == 0),
		renderCheckbox("Регистрация", c == 1),
		renderCheckbox("Справка", c == 2),
	)

	s := fmt.Sprintf(tpl, choices)
	return mainStyle.Render(title + "\n" + s + "\n\n")
}

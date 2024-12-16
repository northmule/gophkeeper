package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Экран справки
type pageHelp struct {
	Choice   int
	mainPage *pageIndex
}

func newPageHelp(main *pageIndex) pageHelp {
	return pageHelp{
		mainPage: main,
	}
}

func (m pageHelp) Init() tea.Cmd {
	return nil
}

func (m pageHelp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "enter" {
			return m.mainPage, nil
		}
	}
	return m, nil
}

// View вид модели
func (m pageHelp) View() string {
	c := m.Choice
	title := renderTitle("Менеджер паролей GophKeeper")
	tpl := bodyStyle.Render("GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю надёжно и безопасно хранить логины,\n пароли, бинарные данные и прочую приватную информацию.\n")
	tpl += subtleStyle.Render("YP: 32 go  (2024)")

	choices := fmt.Sprintf(
		"\n%s\n\n",
		renderCheckbox("Вернуться", c == 0),
	)

	s := fmt.Sprint(tpl, choices)

	return mainStyle.Render(title + "\n" + s + "\n\n")
}

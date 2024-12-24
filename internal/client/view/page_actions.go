package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Экран выбора действий (доступные действия по добавлению данных)
type pageAction struct {
	Choice   int
	Chosen   bool
	mainPage *pageIndex
}

func newPageAction(mainPage *pageIndex) pageAction {
	return pageAction{
		mainPage: mainPage,
	}
}

// Init инициализация модели
func (m pageAction) Init() tea.Cmd {
	return nil
}

// Update изменение модели
func (m pageAction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "down" || k == "tab" {
			m.Choice++
			if m.Choice > 4 {
				m.Choice = 4
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
				return newPageCardData(m.mainPage), nil
			}
			if m.Choice == 1 {
				return newPageTextData(m.mainPage), nil
			}
			if m.Choice == 2 {
				p := newPageFileData(m.mainPage)
				return p, p.Init()
			}
			if m.Choice == 3 {
				return newPageDataGrid(m.mainPage, &m), nil
			}

			// выход
			if m.Choice == 4 {
				m.mainPage.storage.ResetToken()
				return m.mainPage, nil
			}

		}
	}

	return m, nil
}

// View вид модели( в том числе при старте)
func (m pageAction) View() string {
	c := m.Choice

	title := renderTitle("Доступные действия")
	tpl := "%s\n\n"
	tpl += subtleStyle.Render("вверх/вниз: для переключения") + dotStyle +
		subtleStyle.Render("enter: выбрать")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n\n%s\n",
		renderCheckbox("Добавить данные банковских карт", c == 0),
		renderCheckbox("Добавить произвольные текстовые данные", c == 1),
		renderCheckbox("Добавить бинарные данные", c == 2),
		renderCheckbox("Показать мои данные", c == 3),
		renderCheckbox("Выйти", c == 4),
	)

	s := fmt.Sprintf(tpl, choices)
	return mainStyle.Render(title + "\n" + s + "\n\n")
}

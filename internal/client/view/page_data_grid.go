package view

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// Экран с списком данных пользоателя
type pageDataGrid struct {
	mainPage   *pageIndex
	actionPage *pageAction
	table      table.Model
}

func newPageDataGrid(mainPage *pageIndex, actionPage *pageAction) pageDataGrid {
	m := pageDataGrid{}
	m.mainPage = mainPage
	m.actionPage = actionPage

	columns := []table.Column{
		{Title: "№", Width: 4},
		{Title: "Тип", Width: 30},
		{Title: "Название", Width: 60},
		{Title: "UUID", Width: 40},
	}

	rowsData, err := m.mainPage.managerController.GridData().Send(m.mainPage.storage.Token())
	if err != nil {
		tea.Println(err)
		return m
	}
	var rows []table.Row
	for _, item := range rowsData.Items {
		rows = append(rows, table.Row{item.Number, item.Type, item.Name, item.UUID})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m.table = t

	return m
}

func (m pageDataGrid) Init() tea.Cmd { return nil }

func (m pageDataGrid) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m.actionPage, nil
		case "enter":
			dataUUID := m.table.SelectedRow()[3]

			itemResponse, err := m.mainPage.managerController.ItemData().Send(m.mainPage.storage.Token(), dataUUID)
			if err != nil {
				return m, tea.Batch(
					tea.Printf("Произошла ошибка: %s!", err),
				)
			}
			if itemResponse.IsCard {
				return newPageCardData(m.mainPage).SetEditableData(&itemResponse.CardData).SetPageGrid(&m), nil
			}

			if itemResponse.IsText {
				return newPageTextData(m.mainPage).SetEditableData(&itemResponse.TextData).SetPageGrid(&m), nil
			}

			if itemResponse.IsFile {
				return newPageFileData(m.mainPage).SetEditableData(&itemResponse.FileData).SetPageGrid(&m), nil
			}

			return m, tea.Batch(
				tea.Printf("Выбраны данные %s!", dataUUID),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m pageDataGrid) View() string {
	title := renderTitle("Все данные")
	tpl := "%s\n\n"
	tpl += subtleStyle.Render("вверх/вниз: для переключения") + dotStyle +
		subtleStyle.Render("enter: просмотреть данные") + dotStyle +
		subtleStyle.Render("ctrl+c: вернуться") + dotStyle

	s := fmt.Sprintf(tpl, baseStyle.Render(m.table.View()))
	return mainStyle.Render(title + "\n" + s + "\n\n")
}

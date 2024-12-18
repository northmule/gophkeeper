package view

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/model_data"
)

// Ввод/редактирование текстовых данных
type pageTextData struct {
	Choice          int
	Chosen          bool
	mainPage        *pageIndex
	responseMessage string

	// идентификатор редактирования
	uuid string
	// поля
	name  textinput.Model
	text  textarea.Model
	meta1 textinput.Model
	meta2 textinput.Model
}

func newPageTextData(mainPage *pageIndex) pageTextData {

	name := textinput.New()
	name.Placeholder = "Название данных"
	name.Focus()
	name.CharLimit = 100
	name.Width = 100

	text := textarea.New()
	text.Placeholder = "Текст"
	text.CharLimit = 1000
	text.MaxHeight = 100

	meta1 := textinput.New()
	meta1.Placeholder = data_type.TranslateDataType(data_type.MetaNameNote)
	meta1.CharLimit = 100
	meta1.Width = 100

	meta2 := textinput.New()
	meta2.Placeholder = data_type.TranslateDataType(data_type.MetaNameWebSite)
	meta2.CharLimit = 100
	meta2.Width = 100

	m := pageTextData{}
	m.mainPage = mainPage

	m.name = name
	m.text = text
	m.meta1 = meta1
	m.meta2 = meta2

	return m
}

func (m pageTextData) Init() tea.Cmd {
	return textinput.Blink
}

func (m pageTextData) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "down" || k == "tab" {
			m.Choice++
			if m.Choice > 5 {
				m.Choice = 5
			}
		}
		if k == "up" {
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		}
		if k == "enter" {
			if m.Choice == 4 {
				requestData := new(model_data.TextDataRequest)

				requestData.UUID = m.uuid
				requestData.Name = m.name.Value()
				requestData.Value = m.text.Value()
				requestData.Meta = make(map[string]string)
				requestData.Meta[data_type.MetaNameNote] = m.meta1.Value()
				requestData.Meta[data_type.MetaNameWebSite] = m.meta2.Value()

				_, err := m.mainPage.managerController.TextData().Send(m.mainPage.storage.Token(), requestData)
				if err != nil {
					m.responseMessage = err.Error()
					return m, nil
				}
				// Данные отправлены
				m.responseMessage = "Данные сохранены"

				return newPageAction(m.mainPage), nil
			}

			if m.Choice == 5 {
				return newPageAction(m.mainPage), nil
			}
		}
	}

	if m.Choice == 0 {
		m.name, cmd = m.name.Update(msg)
		m.name.Focus()
		return m, cmd
	}
	if m.Choice == 1 {
		m.text, cmd = m.text.Update(msg)
		m.text.Focus()
		return m, cmd
	}

	if m.Choice == 2 {
		m.meta1, cmd = m.meta1.Update(msg)
		m.meta1.Focus()
		return m, cmd
	}

	if m.Choice == 3 {
		m.meta2, cmd = m.meta2.Update(msg)
		m.meta2.Focus()
		return m, cmd
	}

	return m, nil
}

func (m pageTextData) View() string {

	c := m.Choice

	title := renderTitle("Произвольные текстовые данные")

	tpl := "%s\n\n"
	tpl += subtleStyle.Render("вверх/вниз: для переключения") + dotStyle +
		subtleStyle.Render("enter: начать ввод значения") + dotStyle +
		responseTextStyle.Render("\n"+m.responseMessage) + dotStyle

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n\n%s\n",
		renderCheckbox(m.name.View(), c == 0),
		renderCheckbox(m.text.View(), c == 1),
		renderCheckbox(m.meta1.View(), c == 2),
		renderCheckbox(m.meta2.View(), c == 3),
		renderCheckbox("Отправить", c == 4),
		renderCheckbox("Вернуться", c == 5),
	)

	s := fmt.Sprintf(tpl, choices)
	return mainStyle.Render(title + "\n" + s + "\n\n")

}

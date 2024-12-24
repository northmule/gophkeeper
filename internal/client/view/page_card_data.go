package view

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/model_data"
)

// Воод/редактирование данных карт
type pageCardData struct {
	Choice          int
	Chosen          bool
	mainPage        *pageIndex
	gridPage        *pageDataGrid
	responseMessage string

	// идентификатор редактирования
	uuid string
	// поля
	name                 textinput.Model
	cardNumber           textinput.Model
	validityPeriod       textinput.Model
	securityCode         textinput.Model
	fullNameHolder       textinput.Model
	nameBank             textinput.Model
	phoneHolder          textinput.Model
	currentAccountNumber textinput.Model
	meta1                textinput.Model
	meta2                textinput.Model

	isEditable bool
}

func newPageCardData(mainPage *pageIndex) *pageCardData {

	name := textinput.New()
	name.Placeholder = "Название данных"
	name.Focus()
	name.CharLimit = 100
	name.Width = 100

	cardNumber := textinput.New()
	cardNumber.Placeholder = "Номер карты"
	cardNumber.CharLimit = 100
	cardNumber.Width = 100

	validityPeriod := textinput.New()
	validityPeriod.Placeholder = "Срок действия"
	validityPeriod.CharLimit = 100
	validityPeriod.Width = 100

	securityCode := textinput.New()
	securityCode.Placeholder = "Защитный код"
	securityCode.CharLimit = 100
	securityCode.Width = 100

	fullNameHolder := textinput.New()
	fullNameHolder.Placeholder = "ФИО держателя"
	fullNameHolder.CharLimit = 100
	fullNameHolder.Width = 100

	nameBank := textinput.New()
	nameBank.Placeholder = "Название банка"
	nameBank.CharLimit = 100
	nameBank.Width = 100

	phoneHolder := textinput.New()
	phoneHolder.Placeholder = "Телефон держателя"
	phoneHolder.CharLimit = 100
	phoneHolder.Width = 100

	currentAccountNumber := textinput.New()
	currentAccountNumber.Placeholder = "Номер счёта"
	currentAccountNumber.CharLimit = 100
	currentAccountNumber.Width = 100

	meta1 := textinput.New()
	meta1.Placeholder = data_type.TranslateDataType(data_type.MetaNameNote)
	meta1.CharLimit = 100
	meta1.Width = 100

	meta2 := textinput.New()
	meta2.Placeholder = data_type.TranslateDataType(data_type.MetaNameWebSite)
	meta2.CharLimit = 100
	meta2.Width = 100

	m := &pageCardData{}
	m.mainPage = mainPage

	m.name = name
	m.cardNumber = cardNumber
	m.validityPeriod = validityPeriod
	m.securityCode = securityCode
	m.fullNameHolder = fullNameHolder
	m.nameBank = nameBank
	m.phoneHolder = phoneHolder
	m.currentAccountNumber = currentAccountNumber
	m.meta1 = meta1
	m.meta2 = meta2

	return m
}

// SetEditableData значения для редактирования
func (m *pageCardData) SetEditableData(data *model_data.CardDataRequest) *pageCardData {
	m.uuid = data.UUID
	m.name.SetValue(data.Name)
	m.cardNumber.SetValue(data.CardNumber)
	m.validityPeriod.SetValue(data.ValidityPeriod)
	m.securityCode.SetValue(data.SecurityCode)
	m.fullNameHolder.SetValue(data.FullNameHolder)
	m.nameBank.SetValue(data.NameBank)
	m.phoneHolder.SetValue(data.PhoneHolder)
	m.currentAccountNumber.SetValue(data.CurrentAccountNumber)

	if v, ok := data.Meta[data_type.MetaNameNote]; ok {
		m.meta1.SetValue(v)
	}
	if v, ok := data.Meta[data_type.MetaNameWebSite]; ok {
		m.meta2.SetValue(v)
	}

	m.isEditable = true

	return m
}

func (m *pageCardData) SetPageGrid(page *pageDataGrid) *pageCardData {
	m.gridPage = page

	return m
}

func (m *pageCardData) Init() tea.Cmd {
	return textinput.Blink
}

func (m *pageCardData) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "down" || k == "tab" {
			m.Choice++
			if m.Choice > 11 {
				m.Choice = 11
			}
		}
		if k == "up" {
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		}
		if k == "enter" {
			if m.Choice == 10 {
				requestData := new(model_data.CardDataRequest)

				requestData.UUID = m.uuid
				requestData.Name = m.name.Value()
				requestData.CardNumber = m.cardNumber.Value()
				requestData.ValidityPeriod = m.validityPeriod.Value()
				requestData.SecurityCode = m.securityCode.Value()
				requestData.FullNameHolder = m.fullNameHolder.Value()
				requestData.NameBank = m.nameBank.Value()
				requestData.PhoneHolder = m.phoneHolder.Value()
				requestData.CurrentAccountNumber = m.currentAccountNumber.Value()
				requestData.Meta = make(map[string]string)
				requestData.Meta[data_type.MetaNameNote] = m.meta1.Value()
				requestData.Meta[data_type.MetaNameWebSite] = m.meta2.Value()

				_, err := m.mainPage.managerController.CardData().Send(m.mainPage.storage.Token(), requestData)
				if err != nil {
					m.responseMessage = err.Error()
					return m, nil
				}
				// Данные отправлены
				m.responseMessage = "Данные сохранены"

				return newPageAction(m.mainPage), nil
			}

			if m.Choice == 11 {
				if m.isEditable {
					return m.gridPage, nil
				}
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
		m.cardNumber, cmd = m.cardNumber.Update(msg)
		m.cardNumber.Focus()
		return m, cmd
	}
	if m.Choice == 2 {
		m.validityPeriod, cmd = m.validityPeriod.Update(msg)
		m.validityPeriod.Focus()
		return m, cmd
	}
	if m.Choice == 3 {
		m.securityCode, cmd = m.securityCode.Update(msg)
		m.securityCode.Focus()
		return m, cmd
	}
	if m.Choice == 4 {
		m.fullNameHolder, cmd = m.fullNameHolder.Update(msg)
		m.fullNameHolder.Focus()
		return m, cmd
	}
	if m.Choice == 5 {
		m.nameBank, cmd = m.nameBank.Update(msg)
		m.nameBank.Focus()
		return m, cmd
	}
	if m.Choice == 6 {
		m.phoneHolder, cmd = m.phoneHolder.Update(msg)
		m.phoneHolder.Focus()
		return m, cmd
	}
	if m.Choice == 7 {
		m.currentAccountNumber, cmd = m.currentAccountNumber.Update(msg)
		m.currentAccountNumber.Focus()
		return m, cmd
	}
	if m.Choice == 8 {
		m.meta1, cmd = m.meta1.Update(msg)
		m.meta1.Focus()
		return m, cmd
	}

	if m.Choice == 9 {
		m.meta2, cmd = m.meta2.Update(msg)
		m.meta2.Focus()
		return m, cmd
	}

	return m, nil
}

func (m *pageCardData) View() string {

	c := m.Choice

	title := renderTitle("Данные банковских карт")

	tpl := "%s\n\n"
	tpl += subtleStyle.Render("вверх/вниз: для переключения") + dotStyle +
		subtleStyle.Render("enter: начать ввод значения") + dotStyle +
		responseTextStyle.Render("\n"+m.responseMessage) + dotStyle

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n\n%s\n",
		renderCheckbox(m.name.View(), c == 0),
		renderCheckbox(m.cardNumber.View(), c == 1),
		renderCheckbox(m.validityPeriod.View(), c == 2),
		renderCheckbox(m.securityCode.View(), c == 3),
		renderCheckbox(m.fullNameHolder.View(), c == 4),
		renderCheckbox(m.nameBank.View(), c == 5),
		renderCheckbox(m.phoneHolder.View(), c == 6),
		renderCheckbox(m.currentAccountNumber.View(), c == 7),
		renderCheckbox(m.meta1.View(), c == 8),
		renderCheckbox(m.meta2.View(), c == 9),
		renderCheckbox("Отправить", c == 10),
		renderCheckbox("Вернуться", c == 11),
	)

	s := fmt.Sprintf(tpl, choices)
	return mainStyle.Render(title + "\n" + s + "\n\n")

}

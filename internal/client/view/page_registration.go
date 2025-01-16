package view

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Экран регистрации
type pageRegistration struct {
	login           textinput.Model
	password        textinput.Model
	password2       textinput.Model
	email           textinput.Model
	err             error
	Choice          int
	mainPage        *pageIndex
	responseMessage string
}

// инициализация модели
func newPageRegistration(main *pageIndex) *pageRegistration {
	login := textinput.New()
	login.Placeholder = "Введите логин"
	login.Focus()
	login.CharLimit = 50
	login.Width = 20

	password := textinput.New()
	password.Placeholder = "Придумайте пароль"
	password.CharLimit = 50
	password.Width = 20

	password2 := textinput.New()
	password2.Placeholder = "Повторите пароль"
	password2.CharLimit = 50
	password2.Width = 20

	email := textinput.New()
	email.Placeholder = "Email"
	email.CharLimit = 50
	email.Width = 20

	m := &pageRegistration{}
	m.login = login
	m.password = password
	m.password2 = password2
	m.email = email

	m.mainPage = main

	return m
}

func (m *pageRegistration) Init() tea.Cmd {
	return textinput.Blink
}

func (m *pageRegistration) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
				if m.password.Value() != m.password2.Value() {
					m.responseMessage = "Пароли не совпадают"
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}
				if m.login.Value() == "" || m.password.Value() == "" || m.email.Value() == "" {
					m.responseMessage = "Заполните все поля"
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}
				_, err := m.mainPage.managerController.Registration().Send(m.login.Value(), m.password.Value(), m.email.Value())
				if err != nil {
					m.responseMessage = err.Error()
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}

				m.responseMessage = "Вы зарегистрированы"
				return m.mainPage, tea.Batch(cmd, clearErrorAfter(3*time.Second))
			}

			if m.Choice == 5 {
				return m.mainPage, nil
			}
		}
	}

	if m.Choice == 0 {
		m.login, cmd = m.login.Update(msg)
		m.login.Focus()
		return m, cmd
	}
	if m.Choice == 1 {
		m.password, cmd = m.password.Update(msg)
		m.password.Focus()
		return m, cmd
	}
	if m.Choice == 2 {
		m.password2, cmd = m.password2.Update(msg)
		m.password2.Focus()
		return m, cmd
	}
	if m.Choice == 3 {
		m.email, cmd = m.email.Update(msg)
		m.email.Focus()
		return m, cmd
	}

	return m, nil
}

func (m *pageRegistration) View() string {

	c := m.Choice

	title := renderTitle("Регистрация")

	tpl := "%s\n\n"
	tpl += subtleStyle.Render("вверх/вниз: для переключения") + dotStyle +
		subtleStyle.Render("enter: начать ввод значения") + dotStyle +
		responseTextStyle.Render("\n"+m.responseMessage) + dotStyle

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n\n%s\n",
		renderCheckbox(m.login.View(), c == 0),
		renderCheckbox(m.password.View(), c == 1),
		renderCheckbox(m.password2.View(), c == 2),
		renderCheckbox(m.email.View(), c == 3),
		renderCheckbox("Отправить", c == 4),
		renderCheckbox("Вернуться", c == 5),
	)

	s := fmt.Sprintf(tpl, choices)
	return mainStyle.Render(title + "\n" + s + "\n\n")

}

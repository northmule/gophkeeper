package view

import (
	"fmt"
	"time"
)
import tea "github.com/charmbracelet/bubbletea"
import "github.com/charmbracelet/bubbles/textinput"

// Экран авторизации
type pageAuthentication struct {
	login           textinput.Model
	password        textinput.Model
	err             error
	Choice          int
	mainPage        *pageIndex
	responseMessage string
}

// инициализация модели
func newPageAuthentication(main *pageIndex) pageAuthentication {
	login := textinput.New()
	login.Placeholder = "Введите логин"
	login.Focus()
	login.CharLimit = 50
	login.Width = 20

	password := textinput.New()
	password.Placeholder = "Введите пароль"
	password.CharLimit = 50
	password.Width = 20

	m := pageAuthentication{}
	m.login = login
	m.password = password

	m.mainPage = main

	return m
}

func (m pageAuthentication) Init() tea.Cmd {
	return textinput.Blink
}

func (m pageAuthentication) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "down" || k == "tab" {
			m.Choice++
			if m.Choice > 3 {
				m.Choice = 3
			}
		}
		if k == "up" {
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		}
		if k == "enter" {
			if m.Choice == 2 {
				r, err := m.mainPage.managerController.Authentication().Send(m.login.Value(), m.password.Value())
				if err != nil {
					m.responseMessage = err.Error()
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}

				// Сохранём токен
				m.mainPage.storage.SetToken(r.Value)
				// Отправляем публичный ключ клиента на сервер
				err = m.mainPage.managerController.KeysData().UploadClientPublicKey(m.mainPage.storage.Token())
				if err != nil {
					m.responseMessage = err.Error()
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}
				// Забираем публичный ключ с сервера
				err = m.mainPage.managerController.KeysData().DownloadPublicServerKey(m.mainPage.storage.Token())
				if err != nil {
					m.responseMessage = err.Error()
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}
				// Отправка приватного ключа (ключ отправляется зашифрованным публичным серверным)
				err = m.mainPage.managerController.KeysData().UploadClientPrivateKey(m.mainPage.storage.Token())
				if err != nil {
					m.responseMessage = err.Error()
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}
				// Авторизация успешна, отображаем следующее меню
				m.responseMessage = "Вы авторизованы"
				return newPageAction(m.mainPage), tea.Batch(cmd, clearErrorAfter(3*time.Second))
			}

			if m.Choice == 3 {
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

	return m, nil
}

func (m pageAuthentication) View() string {

	c := m.Choice

	title := renderTitle("Авторизация")

	tpl := "%s\n\n"
	tpl += subtleStyle.Render("вверх/вниз: для переключения") + dotStyle +
		subtleStyle.Render("enter: начать ввод значения") + dotStyle +
		responseTextStyle.Render("\n"+m.responseMessage) + dotStyle

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n\n%s\n",
		renderCheckbox(m.login.View(), c == 0),
		renderCheckbox(m.password.View(), c == 1),
		renderCheckbox("Отправить", c == 2),
		renderCheckbox("Вернуться", c == 3),
	)

	s := fmt.Sprintf(tpl, choices)
	return mainStyle.Render(title + "\n" + s + "\n\n")

}

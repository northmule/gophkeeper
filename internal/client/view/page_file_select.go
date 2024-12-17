package view

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

// Выбор файла
type pageFileSelect struct {
	pageFileData    *pageFileData
	responseMessage string
	selectedFile    string
	err             error
	filepicker      filepicker.Model
}

func newpageFileSelect(pageFileData *pageFileData) pageFileSelect {

	m := pageFileSelect{}
	m.pageFileData = pageFileData
	fp := filepicker.New()
	// Файлы разрешённые к загрузке
	fp.AllowedTypes = []string{".go", ".txt", ".md", ".png", ".pdf", ".jpg", ".jpeg", ".doc", ".docx", ".xml", ".zip", ".rar", ".mkv", ".mp3", ".mp4"}
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.ShowSize = false
	// Высота терминала выбора файлов
	fp.Height = 10
	fp.ShowPermissions = true
	fp.AutoHeight = true
	fp.Cursor = "+"
	fp.DirAllowed = true

	style := filepicker.DefaultStyles()
	style.Selected = checkboxStyle
	style.Cursor = checkboxStyle
	fp.Styles = style
	m.filepicker = fp

	return m
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

type clearFieldMsg struct {
}

func clearFieldAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearFieldMsg{}
	})
}

func (m pageFileSelect) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m pageFileSelect) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m.pageFileData, nil
		case "enter":

		}

	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// выбран не разрешёный файл
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = errors.New(path + " выбранный файл не разрешён для отправки.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
	}

	// выбрали файл, заполнили путь
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.selectedFile = path
		m.pageFileData.selectedFile = path
		fileData := m.pageFileData
		return fileData, fileData.Init()
	}

	return m, cmd
}

func (m pageFileSelect) View() string {

	var s strings.Builder
	s.WriteString(renderTitle("Выберите файл"))
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else {
		s.WriteString("Сейчас выбран: " + checkboxStyle.Render(m.selectedFile) + "\n")
		if m.selectedFile != "" {
			s.WriteString("\n" + subtleStyle.Render("Нажмите ещё enter что бы отправить"))
		}
	}

	s.WriteString("\n" + m.filepicker.View() + "\n\n")

	s.WriteString(subtleStyle.Render("стрелки: для перемещения по каталогам\n"))
	s.WriteString(subtleStyle.Render("enter: для выбора файла\n"))
	s.WriteString(subtleStyle.Render("ctrl+c: ернуться назад\n"))

	return s.String()
}

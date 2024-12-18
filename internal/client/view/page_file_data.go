package view

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabriel-vasile/mimetype"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/model_data"
)

// Ввод/редактирование данных о файле
type pageFileData struct {
	Choice          int
	Chosen          bool
	mainPage        *pageIndex
	gridPage        *pageDataGrid
	responseMessage string
	err             error
	selectedFile    string

	// идентификатор редактирования
	uuid string
	// поля на основание файла
	mimeType  string
	extension string
	fileName  string
	size      int64
	// поля для ввода
	name     textinput.Model
	filePath textinput.Model
	meta1    textinput.Model
	meta2    textinput.Model

	isEditable bool
}

func newPageFileData(mainPage *pageIndex) pageFileData {

	name := textinput.New()
	name.Placeholder = "Название данных"
	name.Focus()
	name.CharLimit = 100
	name.Width = 30

	filePath := textinput.New()
	filePath.Placeholder = "Для выбора файла нажмите enter"
	filePath.CharLimit = 100
	filePath.Width = 35
	filePath.SetValue("")

	meta1 := textinput.New()
	meta1.Placeholder = data_type.TranslateDataType(data_type.MetaNameNote)
	meta1.CharLimit = 100
	meta1.Width = 35

	meta2 := textinput.New()
	meta2.Placeholder = data_type.TranslateDataType(data_type.MetaNameWebSite)
	meta2.CharLimit = 100
	meta2.Width = 35

	m := pageFileData{}
	m.mainPage = mainPage

	m.name = name
	m.filePath = filePath
	m.meta1 = meta1
	m.meta2 = meta2

	return m
}

// SetEditableData значения для редактирования
func (m pageFileData) SetEditableData(data *model_data.FileDataInitRequest) pageFileData {
	m.uuid = data.UUID
	m.name.SetValue(data.Name)
	m.filePath.SetValue(data.FileName)
	m.fileName = data.FileName

	if v, ok := data.Meta[data_type.MetaNameNote]; ok {
		m.meta1.SetValue(v)
	}
	if v, ok := data.Meta[data_type.MetaNameWebSite]; ok {
		m.meta2.SetValue(v)
	}

	m.isEditable = true

	return m
}

func (m pageFileData) SetPageGrid(page *pageDataGrid) pageFileData {
	m.gridPage = page

	return m
}

func (m pageFileData) Init() tea.Cmd {
	m.filePath.SetValue(m.selectedFile)
	return textinput.Blink
}

func (m pageFileData) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.(type) {
	case clearErrorMsg:
		m.err = nil
		m.responseMessage = ""
	case clearFieldMsg:
		m.name.SetValue("")
		m.filePath.SetValue("")
		m.selectedFile = ""
		m.meta1.SetValue("")
		m.meta2.SetValue("")
	}

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
			if m.Choice == 1 {
				pageFileSelected := newpageFileSelect(&m)
				return pageFileSelected, pageFileSelected.Init()
			}
			if m.Choice == 4 {
				requestData := new(model_data.FileDataInitRequest)

				requestData.UUID = m.uuid
				requestData.Name = m.name.Value()
				// Информация о файле
				if m.selectedFile == "" {
					m.responseMessage = "Файл не выбран, выберите файл"
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}

				file, _ := os.Open(m.selectedFile)
				defer file.Close()
				fileInfo, _ := file.Stat()
				requestData.FileName = fileInfo.Name()
				requestData.Size = fileInfo.Size()

				mtype, _ := mimetype.DetectFile(m.selectedFile)
				requestData.MimeType = mtype.String()
				requestData.Extension = mtype.Extension()

				requestData.Meta = make(map[string]string)
				requestData.Meta[data_type.MetaNameNote] = m.meta1.Value()
				requestData.Meta[data_type.MetaNameWebSite] = m.meta2.Value()

				fresponse, err := m.mainPage.managerController.FileData().Send(m.mainPage.storage.Token(), requestData)
				if err != nil {
					m.responseMessage = err.Error()
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}

				// Отправка самого файла
				err = m.mainPage.managerController.FileData().UploadFile(m.mainPage.storage.Token(), fresponse.UploadPath, file)
				if err != nil {
					m.responseMessage = err.Error()
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}
				m.responseMessage = "Файл загружен"
				return m, tea.Batch(cmd, clearErrorAfter(3*time.Second), clearFieldAfter(1*time.Second))
			}

			if m.Choice == 5 {
				if m.isEditable {
					return m.gridPage, nil
				}
				return newPageAction(m.mainPage), nil
			}
		}

		if k == "ctrl+d" && m.uuid != "" { // скачать файл при редактирование в папку указанную в конфиге
			err := m.mainPage.managerController.FileData().DownLoadFile(m.mainPage.storage.Token(), m.fileName, m.uuid)
			if err != nil {
				m.responseMessage = err.Error()
			}
			m.responseMessage = "Файл получен"
			return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
		}
	}

	// m.filePath.Update(msg)

	if m.Choice == 0 {
		m.name, cmd = m.name.Update(msg)
		m.name.Focus()
		return m, cmd
	}
	if m.Choice == 1 {
		m.filePath, cmd = m.filePath.Update(msg)
		m.filePath.Focus()
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

func (m pageFileData) View() string {

	c := m.Choice
	if m.selectedFile != "" {
		m.filePath.SetValue(m.selectedFile)
	}

	title := renderTitle("Бинарные данные")

	tpl := "%s\n\n"
	tpl += subtleStyle.Render("вверх/вниз: для переключения") + dotStyle +
		subtleStyle.Render("enter: начать ввод значения") + dotStyle +
		responseTextStyle.Render("\n"+m.responseMessage) + dotStyle

	if m.isEditable {
		tpl += subtleStyle.Render("ctrl+d: чтобы скачать файл в папку")
	}

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n\n%s\n",
		renderCheckbox(m.name.View(), c == 0),
		renderCheckbox(m.filePath.View(), c == 1)+subtleStyle.Render(" # файл для отправки"),
		renderCheckbox(m.meta1.View(), c == 2),
		renderCheckbox(m.meta2.View(), c == 3),
		renderCheckbox("Отправить", c == 4),
		renderCheckbox("Вернуться", c == 5),
	)

	s := fmt.Sprintf(tpl, choices)
	return mainStyle.Render(title + "\n" + s + "\n\n")

}

package view

import (
	"errors"
	"path"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/gophkeeper/internal/client/config"
	"github.com/northmule/gophkeeper/internal/client/controller"
	"github.com/northmule/gophkeeper/internal/client/logger"
	"github.com/northmule/gophkeeper/internal/client/service"
	"github.com/northmule/gophkeeper/internal/client/storage"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPageCardData_Init(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	cryptService, _ := service.NewCrypt(mockCfg)
	manager, _ := controller.NewManager(mockCfg, cryptService, log)
	memoryStorage := storage.NewMemoryStorage()
	mainPage := newPageIndex(manager, memoryStorage, log)
	page := newPageCardData(mainPage)
	cmd := page.Init()
	assert.NotNil(t, cmd)
}

func TestPageCardData_Update(t *testing.T) {
	tests := []struct {
		name     string
		msg      tea.Msg
		expected pageCardData
	}{
		{"down key", tea.KeyMsg{Type: tea.KeyDown}, pageCardData{Choice: 1}},
		{"up key", tea.KeyMsg{Type: tea.KeyUp}, pageCardData{Choice: 0}},
		{"invalid key", tea.KeyMsg{Type: tea.KeyCtrlC}, pageCardData{Choice: 0}},
	}
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	cryptService, _ := service.NewCrypt(mockCfg)
	manager, _ := controller.NewManager(mockCfg, cryptService, log)
	memoryStorage := storage.NewMemoryStorage()
	mainPage := newPageIndex(manager, memoryStorage, log)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := newPageCardData(mainPage)
			_, _ = page.Update(tt.msg)
			assert.Equal(t, tt.expected.Choice, page.Choice)
		})
	}

	// прочие кейсы
	t.Run("choice 10", func(t *testing.T) {
		mockManagerController := new(MockManagerController)
		mockCardData := new(MockCardDataController)
		mockManagerController.On("CardData").Return(mockCardData)

		mockCardData.On("Send", mock.Anything, mock.Anything).Return(&controller.CardDataResponse{Value: "ok"}, nil)

		mainPage := newPageIndex(mockManagerController, memoryStorage, log)
		pa := pageCardData{Choice: 10, mainPage: mainPage}
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.Contains(t, "Данные сохранены", pa.responseMessage)
		assert.NotNil(t, m)
	})

	t.Run("choice 10 error", func(t *testing.T) {
		mockManagerController := new(MockManagerController)
		mockCardData := new(MockCardDataController)
		mockManagerController.On("CardData").Return(mockCardData)

		mockCardData.On("Send", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

		mainPage := newPageIndex(mockManagerController, memoryStorage, log)
		pa := pageCardData{Choice: 10, mainPage: mainPage}
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.NotEmpty(t, pa.responseMessage)
		assert.NotNil(t, m)
	})

	t.Run("choice 11", func(t *testing.T) {
		mockManagerController := new(MockManagerController)
		mockCardData := new(MockCardDataController)
		mockManagerController.On("CardData").Return(mockCardData)

		mockCardData.On("Send", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

		mainPage := newPageIndex(mockManagerController, memoryStorage, log)
		pa := pageCardData{Choice: 11, mainPage: mainPage}
		pa.isEditable = false
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		m, _ := pa.Update(msg)
		assert.NotNil(t, m)

		pa.isEditable = true
		msg = tea.KeyMsg{Type: tea.KeyEnter}
		m, _ = pa.Update(msg)
		assert.Nil(t, m)
	})

	t.Run("choice 2-9", func(t *testing.T) {
		mockManagerController := new(MockManagerController)
		mockCardData := new(MockCardDataController)
		mockManagerController.On("CardData").Return(mockCardData)

		mainPage := newPageIndex(mockManagerController, memoryStorage, log)

		data := &model_data.CardDataRequest{
			Name:                 "My Card",
			UUID:                 "123e4567-e89b-12d3-a456-426614174000",
			CardNumber:           "4111111111111111",
			ValidityPeriod:       "2025-12-31T23:59:59Z07:00",
			SecurityCode:         "123",
			FullNameHolder:       "John Doe",
			NameBank:             "Example Bank",
			PhoneHolder:          "+1234567890",
			CurrentAccountNumber: "12345678901234567890",
			Meta: map[string]string{
				data_type.MetaNameNote:    "value1",
				data_type.MetaNameWebSite: "value2",
			},
		}

		pa := newPageCardData(mainPage)
		pa.SetEditableData(data)
		msg := tea.KeyMsg{Type: tea.KeyEnter}

		num := 2
		for {
			if num == 9 {
				break
			}
			pa.Choice = num
			m, _ := pa.Update(msg)
			assert.NotNil(t, m)
			num++
		}

	})
}

func TestPageCardData_View(t *testing.T) {
	mockCfg := config.NewConfig()
	mockCfg.Value().PathPublicKeyServer = path.Join("testpath")
	mockCfg.Value().PathKeys = path.Join("testpath")
	log, _ := logger.NewLogger("info")
	cryptService, _ := service.NewCrypt(mockCfg)
	manager, _ := controller.NewManager(mockCfg, cryptService, log)
	memoryStorage := storage.NewMemoryStorage()
	mainPage := newPageIndex(manager, memoryStorage, log)

	page := newPageCardData(mainPage)
	result := page.View()
	assert.True(t, strings.Contains(result, "Данные банковских карт"))
	assert.True(t, strings.Contains(result, "enter: начать ввод значения"))
	assert.True(t, strings.Contains(result, "Отправить"))
	assert.True(t, strings.Contains(result, "Вернуться"))
}

func TestPageCardData_SetEditableData(t *testing.T) {
	data := &model_data.CardDataRequest{
		Name:                 "My Card",
		UUID:                 "123e4567-e89b-12d3-a456-426614174000",
		CardNumber:           "4111111111111111",
		ValidityPeriod:       "2025-12-31T23:59:59Z07:00",
		SecurityCode:         "123",
		FullNameHolder:       "John Doe",
		NameBank:             "Example Bank",
		PhoneHolder:          "+1234567890",
		CurrentAccountNumber: "12345678901234567890",
		Meta: map[string]string{
			data_type.MetaNameNote:    "value1",
			data_type.MetaNameWebSite: "value2",
		},
	}
	pa := pageCardData{}
	pa.SetEditableData(data)

	assert.True(t, pa.isEditable)
	assert.NotEmpty(t, pa.name.Value())
	assert.NotEmpty(t, pa.fullNameHolder.Value())

}

package model_data

// Общие данные для запросов между клиентом и сервером

// CardDataRequest данные для запросов (клиент и сервер)
type CardDataRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"` // короткое название
	UUID string `json:"uuid" validate:"omitempty,uuid"`         // uuid данных, заполняется при редактирование

	CardNumber           string `json:"card_number" validate:"min=13,max=19"`                          // номер карты
	ValidityPeriod       string `json:"validity_period" validate:"datetime=2006-01-02T15:04:05Z07:00"` // срок действия
	SecurityCode         string `json:"security_code" validate:"max=10"`                               // код безопасности
	FullNameHolder       string `json:"full_name_holder" validate:"min=3,max=100"`                     // ФИО держателя
	NameBank             string `json:"name_bank" validate:"max=100"`                                  // название банка
	PhoneHolder          string `json:"phone_holder" validate:"min=7,max=12"`                          // телефон держателя
	CurrentAccountNumber string `json:"current_account_number" validate:"max=20"`                      // номер расчётного счета

	Meta map[string]string `json:"meta" validate:"max=5,dive,keys,min=3,max=20,endkeys"` // мета данные (имя поля - значение)

}

// TextDataRequest данные для запросов (клиент и сервер)
type TextDataRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"` // короткое название
	UUID string `json:"uuid" validate:"omitempty,uuid"`         // uuid данных, заполняется при редактирование

	Value string `json:"value"` // Текстовые данные

	Meta map[string]string `json:"meta" validate:"max=5,dive,keys,min=3,max=20,endkeys"` // мета данные (имя поля - значение)
}

// FileDataInitRequest Запрос инициализации загрузки файла (основная информация о файле) (клиент и сервер)
type FileDataInitRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"` // короткое название
	UUID     string `json:"uuid" validate:"omitempty,uuid"`         // uuid данных, заполняется при редактирование
	MimeType string `json:"mime_type"`                              // тип файла

	Extension string `json:"extension" validate:"required,min=1,max=10"`  // расширение файла
	FileName  string `json:"file_name" validate:"required,min=3,max=100"` // оригинальное имя файла
	Size      int64  `json:"size"`                                        // размер файла в байтах

	Meta map[string]string `json:"meta" validate:"max=5,dive,keys,min=3,max=20,endkeys"` // мета данные (имя поля - значение)
}

// ItemDataResponse одна еденица данных пользователя
type ItemDataResponse struct {
	Number     string `json:"number"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	UpdateDate string `json:"update_date"`
	UUID       string `json:"uuid"`
}

// ListDataItemsResponse список данных пользователя
type ListDataItemsResponse struct {
	Items []ItemDataResponse `json:"items"`
}

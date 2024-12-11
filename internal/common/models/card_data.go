package models

import (
	"time"
)

// CardData данные бансковских карт
type CardData struct {
	ID         int64           `json:"id"`
	Name       string          `json:"name"`        // короткое название
	ObjectType string          `json:"object_type"` // тип данных из value
	Value      CardDataValueV1 `json:"value"`       // jsonb postgress (тип зависит от ObjectType)
	UUID       string          `json:"uuid"`        // uuid этих данных
}

// CardDataValueV1 значение для Value
type CardDataValueV1 struct {
	CardNumber           string    `json:"card_number"`            // номер карты
	ValidityPeriod       time.Time `json:"validity_period"`        // срок действия
	SecurityCode         string    `json:"security_code"`          // код безопасности
	FullNameHolder       string    `json:"full_name_holder"`       // ФИО держателя
	NameBank             string    `json:"name_bank"`              // название банка
	PhoneHolder          string    `json:"phone_holder"`           // телефон держателя
	CurrentAccountNumber string    `json:"current_account_number"` // номер расчётного счета
}

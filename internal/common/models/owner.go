package models

// Owner Пользователь, владелец данных
type Owner struct {
	ID       int64  `json:"id"`
	UserUUID string `json:"user_uuid"` // uuid пользователя
	DataType string `json:"data_type"` // тип данных @see data_type.go
	DataUUID string `json:"data_uuid"` // uuid данных
}

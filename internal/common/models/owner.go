package models

// Owner Пользователь, владелец данных
type Owner struct {
	ID       int64  `json:"id"`
	UserUUID string `json:"user_uuid"` // uuid пользователя
	DataType string `json:"data_type"` // тип данных @see data_type.go
	DataUUID string `json:"data_uuid"` // uuid данных
}

// OwnerData данные пользователя
type OwnerData struct {
	UserUUID     string `json:"user_uuid"`
	DataUUID     string `json:"data_uuid"`
	DataType     string `json:"data_type"`
	DataTypeName string `json:"data_type_name"`
	DataName     string `json:"data_name"`
}

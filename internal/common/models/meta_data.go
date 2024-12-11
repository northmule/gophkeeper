package models

type MetaData struct {
	ID        int64         `json:"id"`
	MetaName  string        `json:"meta_name"`  // имя св-ва (поля)
	MetaValue MetaDataValue `json:"meta_value"` // значения из св-ва (поля)
	DataUUID  string        `json:"data_uuid"`  // uuid связанных данных
}

type MetaDataValue struct {
	Value string `json:"value"`
}

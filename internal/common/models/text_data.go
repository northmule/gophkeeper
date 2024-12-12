package models

// TextData текст
type TextData struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`  // короткое название
	Value string `json:"value"` // текст
	UUID  string `json:"uuid"`  // uuid этих данных
}

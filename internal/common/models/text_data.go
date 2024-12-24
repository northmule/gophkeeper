package models

// TextData текст
type TextData struct {
	Common
	Name  string `json:"name"`  // короткое название
	Value string `json:"value"` // текст

}

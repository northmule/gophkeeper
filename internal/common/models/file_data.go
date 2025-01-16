package models

// FileData данные файла
type FileData struct {
	Common
	Name      string `json:"name"`      // короткое название
	MimeType  string `json:"mime_type"` // тип файла
	Path      string `json:"path"`      // путь до файла на сервере storage
	PathTmp   string `json:"path_tmp"`  // путь до временной папки с файлом
	Extension string `json:"extension"` // расширение файла
	FileName  string `json:"file_name"` // оригинальное имя файла
	Size      int64  `json:"size"`      // размер файла
	Storage   string `json:"storage"`   // имя сервера где находится файла
	Uploaded  bool   `json:"uploaded"`  // полностью загружен
}

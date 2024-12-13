package handlers

import (
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	"github.com/northmule/gophkeeper/internal/server/services/access"
)

// FileDataHandler обработка запросо на сохранение файлов
type FileDataHandler struct {
	log            *logger.Logger
	accessService  *access.Access
	manager        repository.Repository
	expectedAction map[string]bool
}

// NewFileDataHandler конструктор
func NewFileDataHandler(accessService *access.Access, manager repository.Repository, log *logger.Logger) *FileDataHandler {
	expectedAction := make(map[string]bool)

	expectedAction["load"] = true
	expectedAction["finish"] = true

	return &FileDataHandler{
		accessService:  accessService,
		manager:        manager,
		log:            log,
		expectedAction: expectedAction,
	}
}

// Запрос инициализации загрузки файла (основная информация о файле)
type fileDataInitRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"` // короткое название
	UUID     string `json:"uuid" validate:"omitempty,uuid"`         // uuid данных, заполняется при редактирование
	MimeType string `json:"mime_type"`                              // тип файла

	Extension string `json:"extension" validate:"required,min=1,max=10"`  // расширение файла
	FileName  string `json:"file_name" validate:"required,min=3,max=100"` // оригинальное имя файла
	Size      int64  `json:"size"`                                        // размер файла в байтах

	Meta map[string]string `json:"meta" validate:"max=5,dive,keys,min=3,max=20,endkeys"` // мета данные (имя поля - значение)
}

// Ответ на данные инициализации
type fileDataInitResponse struct {
	UploadPath string `json:"upload_path"` // адрес для зарузки файла post-ом
}

// Bind декодирует json в структуру
func (rr *fileDataInitRequest) Bind(r *http.Request) error {
	return nil
}

// Render рисует json ответ в структуре
func (hr fileDataInitResponse) Render(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// HandleInit инициализация загрузки
func (h *FileDataHandler) HandleInit(res http.ResponseWriter, req *http.Request) {
	var (
		err      error
		userUUID string
		owner    *models.Owner
		fileData *models.FileData
		dataUUID string
	)

	request := new(fileDataInitRequest)
	if err = render.Bind(req, request); err != nil {
		h.log.Info(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}
	userUUID, err = h.accessService.GetUserUUIDByJWTToken(req.Context())
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}

	if request.UUID != "" { // редактирование
		dataUUID = request.UUID
		// владелец данных
		owner, err = h.manager.Owner().FindOneByUserUUIDAndDataUUIDAndDataType(req.Context(), userUUID, dataUUID, data_type.BinaryType)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrBadRequest)
			return
		}
		if owner == nil { // нет данных этого пользователя
			h.log.Infof("owner not found: data_uuid: %s, user_uuid: %s, data_type: %s", dataUUID, userUUID, data_type.BinaryType)
			_ = render.Render(res, req, ErrNotFound)
			return
		}
		fileData, err = h.manager.FileData().FindOneByUUID(req.Context(), dataUUID)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrBadRequest)
			return
		}
		if fileData == nil {
			h.log.Infof("file data not found: uuid %s", owner.DataUUID)
			_ = render.Render(res, req, ErrNotFound)
			return
		}
		// основные данные
		fileData.Name = request.Name
		fileData.FileName = request.FileName
		fileData.Size = request.Size
		fileData.Extension = request.Extension
		fileData.MimeType = request.MimeType

		err = h.manager.FileData().Update(req.Context(), fileData)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}

		// мета поля
		var newMeta []models.MetaData
		if len(request.Meta) > 0 {
			for key, value := range request.Meta {
				metaData := models.MetaData{}
				metaData.MetaName = key
				metaData.MetaValue.Value = value
				metaData.DataUUID = dataUUID
				newMeta = append(newMeta, metaData)
			}
		}
		// перезапись мета
		err = h.manager.MetaData().ReplaceMetaByDataUUID(req.Context(), dataUUID, newMeta)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}
	}

	if request.UUID == "" { // новые данные
		dataUUID = uuid.NewString()

		fileData = new(models.FileData)
		fileData.Name = request.Name
		fileData.MimeType = request.MimeType
		fileData.FileName = request.FileName
		fileData.Size = request.Size
		fileData.Extension = request.Extension

		fileData.UUID = dataUUID
		fileData.PathTmp = os.TempDir() + "/" + dataUUID
		fileData.Path = os.TempDir() + "/load_" + dataUUID // todo в конфиг
		fileData.Storage = "local://"                      // todo в настройки (сейчас не используется)
		fileData.Uploaded = false

		_, err = h.manager.FileData().Add(req.Context(), fileData)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}

		// владелец данных
		owner = new(models.Owner)
		owner.UserUUID = userUUID
		owner.DataType = data_type.BinaryType
		owner.DataUUID = dataUUID

		_, err = h.manager.Owner().Add(req.Context(), owner)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}

		// мета поля
		if len(request.Meta) > 0 {
			for key, value := range request.Meta {
				metaData := new(models.MetaData)
				metaData.MetaName = key
				metaData.MetaValue.Value = value
				metaData.DataUUID = dataUUID
				_, err = h.manager.MetaData().Add(req.Context(), metaData)
				if err != nil {
					h.log.Error(err)
					_ = render.Render(res, req, ErrInternalServerError)
					return
				}
			}
		}
	}

	initResponse := fileDataInitResponse{UploadPath: "/file_data/load/" + dataUUID + "/0"}
	err = render.Render(res, req, initResponse)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
	}

}

// HandleAction создание/обновление данных
func (h *FileDataHandler) HandleAction(res http.ResponseWriter, req *http.Request) {
	var (
		err        error
		userUUID   string
		pathAction string
		dataUUID   string
		pathPart   string
		owner      *models.Owner
		fileData   *models.FileData
	)

	pathAction = chi.URLParam(req, "action")
	dataUUID = chi.URLParam(req, "file_uuid")
	pathPart = chi.URLParam(req, "part")

	if flag, ok := h.expectedAction[pathAction]; !ok || !flag {
		h.log.Info("Expected action: '%v', actual: '%s'", h.expectedAction, pathAction)
		_ = render.Render(res, req, ErrNotFound)
		return
	}

	userUUID, err = h.accessService.GetUserUUIDByJWTToken(req.Context())
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}

	// владелец данных
	owner, err = h.manager.Owner().FindOneByUserUUIDAndDataUUIDAndDataType(req.Context(), userUUID, dataUUID, data_type.BinaryType)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}
	if owner == nil { // нет данных этого пользователя
		h.log.Infof("owner not found: data_uuid: %s, user_uuid: %s, data_type: %s", dataUUID, userUUID, data_type.BinaryType)
		_ = render.Render(res, req, ErrNotFound)
		return
	}

	fileData, err = h.manager.FileData().FindOneByUUID(req.Context(), dataUUID)
	// Всё сразу todo по частям и с part
	buffer := make([]byte, req.ContentLength)
	_, err = io.ReadFull(req.Body, buffer)
	if err != nil && err != io.EOF {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}

	filename := fileData.Path + "/" + fileData.FileName
	err = os.Mkdir(fileData.Path+"/", 0777)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	defer f.Close()
	if _, err = f.Write(buffer); err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	// Файл загруен
	fileData.Uploaded = true
	err = h.manager.FileData().Update(req.Context(), fileData)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	_ = pathPart
}

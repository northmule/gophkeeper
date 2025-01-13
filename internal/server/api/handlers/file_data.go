package handlers

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/config"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"golang.org/x/net/context"
)

// FileDataHandler обработка запросо на сохранение файлов
type FileDataHandler struct {
	log             *logger.Logger
	userFinderByJWT UserFinderByJWT
	fileDataCRUD    FileDataCRUD
	ownerCRUD       OwnerCRUD
	metaDataCRUD    MetaDataCRUD
	cfg             *config.Config
}

// NewFileDataHandler конструктор
func NewFileDataHandler(userFinderByJWT UserFinderByJWT, fileDataCRUD FileDataCRUD, ownerCRUD OwnerCRUD, metaDataCRUD MetaDataCRUD, cfg *config.Config, log *logger.Logger) *FileDataHandler {

	return &FileDataHandler{
		userFinderByJWT: userFinderByJWT,
		log:             log,
		fileDataCRUD:    fileDataCRUD,
		ownerCRUD:       ownerCRUD,
		metaDataCRUD:    metaDataCRUD,
		cfg:             cfg,
	}
}

// FileDataCRUD операции над данными
type FileDataCRUD interface {
	FindOneByUUID(ctx context.Context, uuid string) (*models.FileData, error)
	Add(ctx context.Context, data *models.FileData) (int64, error)
	Update(ctx context.Context, data *models.FileData) error
}

// Запрос инициализации загрузки файла (основная информация о файле)
type fileDataInitRequest struct {
	model_data.FileDataInitRequest
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
	userUUID, err = h.userFinderByJWT.GetUserUUIDByJWTToken(req.Context())
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}

	if request.UUID != "" { // редактирование
		dataUUID = request.UUID
		// владелец данных
		owner, err = h.ownerCRUD.FindOneByUserUUIDAndDataUUIDAndDataType(req.Context(), userUUID, dataUUID, data_type.BinaryType)
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
		fileData, err = h.fileDataCRUD.FindOneByUUID(req.Context(), dataUUID)
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

		err = h.fileDataCRUD.Update(req.Context(), fileData)
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
		err = h.metaDataCRUD.ReplaceMetaByDataUUID(req.Context(), dataUUID, newMeta)
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
		fileBaseDir := os.TempDir()
		if h.cfg.Value().PathFileStorage != "" {
			fileBaseDir = h.cfg.Value().PathFileStorage
		}
		fileData.Path = fileBaseDir + "/load_" + dataUUID
		fileData.Storage = "local://" // todo в настройки (сейчас не используется)
		fileData.Uploaded = false

		_, err = h.fileDataCRUD.Add(req.Context(), fileData)
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

		_, err = h.ownerCRUD.Add(req.Context(), owner)
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
				_, err = h.metaDataCRUD.Add(req.Context(), metaData)
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

// HandleAction принимает содержимое файла отправляемое клиентом
func (h *FileDataHandler) HandleAction(res http.ResponseWriter, req *http.Request) {
	var (
		err         error
		userUUID    string
		dataUUID    string
		pathPart    string
		owner       *models.Owner
		errResponse *ErrResponse
	)

	dataUUID = chi.URLParam(req, "file_uuid")
	pathPart = chi.URLParam(req, "part")

	userUUID, err = h.userFinderByJWT.GetUserUUIDByJWTToken(req.Context())
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}

	// владелец данных
	owner, err = h.ownerCRUD.FindOneByUserUUIDAndDataUUIDAndDataType(req.Context(), userUUID, dataUUID, data_type.BinaryType)
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

	errResponse = h.loadFile(req, dataUUID)

	if errResponse != nil {
		h.log.Error(err)
		_ = render.Render(res, req, errResponse)
		return
	}
	_ = pathPart
}

// HandleGetAction отдаёт клиенту файл
func (h *FileDataHandler) HandleGetAction(res http.ResponseWriter, req *http.Request) {
	var (
		err      error
		userUUID string

		dataUUID    string
		pathPart    string
		owner       *models.Owner
		errResponse *ErrResponse
	)

	dataUUID = chi.URLParam(req, "file_uuid")
	pathPart = chi.URLParam(req, "part")

	userUUID, err = h.userFinderByJWT.GetUserUUIDByJWTToken(req.Context())
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}

	// владелец данных
	owner, err = h.ownerCRUD.FindOneByUserUUIDAndDataUUIDAndDataType(req.Context(), userUUID, dataUUID, data_type.BinaryType)
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

	errResponse = h.downLoadFile(res, req, dataUUID)

	if errResponse != nil {
		h.log.Error(err)
		_ = render.Render(res, req, errResponse)
		return
	}
	_ = pathPart
}

// downLoadFile отдача файла клиенту по запросу
func (h *FileDataHandler) downLoadFile(res http.ResponseWriter, req *http.Request, dataUUID string) *ErrResponse {
	var (
		err      error
		fileData *models.FileData
		file     *os.File
	)
	fileData, err = h.fileDataCRUD.FindOneByUUID(req.Context(), dataUUID)
	if err != nil {
		h.log.Error(err)
		return ErrInternalServerError
	}

	buff, err := os.ReadFile(path.Join(fileData.Path, fileData.FileName))
	if err != nil {
		h.log.Error(err)
		return ErrInternalServerError
	}
	defer file.Close()

	_, err = res.Write(buff)
	if err != nil {
		h.log.Error(err)
		return ErrInternalServerError
	}

	return nil
}

// loadFile Загрузка файла от клиента по запросу
func (h *FileDataHandler) loadFile(req *http.Request, dataUUID string) *ErrResponse {
	var (
		err      error
		fileData *models.FileData
	)
	err = req.ParseMultipartForm(4096)
	if err != nil {
		h.log.Error(err)
		return ErrBadRequest
	}
	requestFile, _, err := req.FormFile(data_type.FileField)
	if err != nil {
		h.log.Error(err)
		return ErrBadRequest
	}
	defer requestFile.Close()

	fileData, err = h.fileDataCRUD.FindOneByUUID(req.Context(), dataUUID)

	// Всё сразу todo по частям и с part
	filename := fileData.Path + "/" + fileData.FileName
	err = os.Mkdir(fileData.Path+"/", 0777)
	if err != nil {
		h.log.Error(err)
		return ErrInternalServerError
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		h.log.Error(err)
		return ErrInternalServerError
	}
	defer f.Close()
	_, err = io.Copy(f, requestFile)
	if err != nil {
		h.log.Error(err)
		return ErrInternalServerError
	}
	// Файл загружен
	fileData.Uploaded = true
	err = h.fileDataCRUD.Update(req.Context(), fileData)
	if err != nil {
		h.log.Error(err)
		return ErrInternalServerError
	}

	return nil
}

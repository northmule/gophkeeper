package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	"github.com/northmule/gophkeeper/internal/server/services/access"
)

// TextDataHandler обрабатывает текстовые данные
type TextDataHandler struct {
	log           *logger.Logger
	accessService *access.Access
	manager       repository.Repository
}

// NewTextDataHandler конструктор
func NewTextDataHandler(accessService *access.Access, manager repository.Repository, log *logger.Logger) *TextDataHandler {
	return &TextDataHandler{
		accessService: accessService,
		manager:       manager,
		log:           log,
	}
}

type textDataRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"` // короткое название
	UUID string `json:"uuid" validate:"omitempty,uuid"`         // uuid данных, заполняется при редактирование

	Value string `json:"value"` // Текстовые данные

	Meta map[string]string `json:"meta" validate:"max=5,dive,keys,min=3,max=20,endkeys"` // мета данные (имя поля - значение)
}

// Bind декодирует json в структуру
func (rr *textDataRequest) Bind(r *http.Request) error {
	return nil
}

// HandleSave создание/обновление данных
func (h *TextDataHandler) HandleSave(res http.ResponseWriter, req *http.Request) {
	var (
		err      error
		userUUID string
		owner    *models.Owner
		textData *models.TextData
	)

	request := new(textDataRequest)
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
		dataUUID := request.UUID
		// владелец данных
		owner, err = h.manager.Owner().FindOneByUserUUIDAndDataUUIDAndDataType(req.Context(), userUUID, dataUUID, data_type.TextType)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrBadRequest)
			return
		}
		if owner == nil { // нет данных этого пользователя
			h.log.Infof("owner not found: data_uuid: %s, user_uuid: %s, data_type: %s", dataUUID, userUUID, data_type.TextType)
			_ = render.Render(res, req, ErrNotFound)
			return
		}
		textData, err = h.manager.TextData().FindOneByUUID(req.Context(), dataUUID)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrBadRequest)
			return
		}
		if textData == nil {
			h.log.Infof("text data not found: uuid %s", owner.DataUUID)
			_ = render.Render(res, req, ErrNotFound)
			return
		}
		// основные данные
		textData.Name = request.Name
		textData.Value = request.Value

		err = h.manager.TextData().Update(req.Context(), textData)
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
		dataUUID := uuid.NewString()

		textData = new(models.TextData)
		textData.Name = request.Name
		textData.Value = request.Value
		textData.UUID = dataUUID

		_, err = h.manager.TextData().Add(req.Context(), textData)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}

		// владелец данных
		owner = new(models.Owner)
		owner.UserUUID = userUUID
		owner.DataType = data_type.TextType
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
}

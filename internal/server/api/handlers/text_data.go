package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"golang.org/x/net/context"
)

// TextDataHandler обрабатывает текстовые данные
type TextDataHandler struct {
	log             *logger.Logger
	userFinderByJWT UserFinderByJWT
	ownerCRUD       OwnerCRUD
	metaDataCRUD    MetaDataCRUD
	textDataCRUD    TextDataCRUD
}

// NewTextDataHandler конструктор
func NewTextDataHandler(userFinderByJWT UserFinderByJWT, ownerCRUD OwnerCRUD, metaDataCRUD MetaDataCRUD, textDataCRUD TextDataCRUD, log *logger.Logger) *TextDataHandler {
	return &TextDataHandler{
		userFinderByJWT: userFinderByJWT,
		metaDataCRUD:    metaDataCRUD,
		ownerCRUD:       ownerCRUD,
		textDataCRUD:    textDataCRUD,
		log:             log,
	}
}

// TextDataCRUD операции над данными
type TextDataCRUD interface {
	FindOneByUUID(ctx context.Context, uuid string) (*models.TextData, error)
	Add(ctx context.Context, data *models.TextData) (int64, error)
	Update(ctx context.Context, data *models.TextData) error
}

type textDataRequest struct {
	model_data.TextDataRequest
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
	userUUID, err = h.userFinderByJWT.GetUserUUIDByJWTToken(req.Context())
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}

	if request.UUID != "" { // редактирование
		dataUUID := request.UUID
		// владелец данных
		owner, err = h.ownerCRUD.FindOneByUserUUIDAndDataUUIDAndDataType(req.Context(), userUUID, dataUUID, data_type.TextType)
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
		textData, err = h.textDataCRUD.FindOneByUUID(req.Context(), dataUUID)
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

		err = h.textDataCRUD.Update(req.Context(), textData)
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
		dataUUID := uuid.NewString()

		textData = new(models.TextData)
		textData.Name = request.Name
		textData.Value = request.Value
		textData.UUID = dataUUID

		_, err = h.textDataCRUD.Add(req.Context(), textData)
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
}

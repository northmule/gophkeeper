package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	"github.com/northmule/gophkeeper/internal/server/services/access"
)

// ItemDataHandler обрабатывает запрос данных по uuid
type ItemDataHandler struct {
	log           *logger.Logger
	accessService access.AccessService
	manager       repository.Repository
}

// NewItemDataHandler конструктор
func NewItemDataHandler(accessService access.AccessService, manager repository.Repository, log *logger.Logger) *ItemDataHandler {
	return &ItemDataHandler{
		accessService: accessService,
		manager:       manager,
		log:           log,
	}
}

type dataByUUIDResponse struct {
	model_data.DataByUUIDResponse
}

// Render рисует структуру в json
func (hr dataByUUIDResponse) Render(res http.ResponseWriter, req *http.Request) error {
	return nil
}

// HandleItem обработк запроса
func (h *ItemDataHandler) HandleItem(res http.ResponseWriter, req *http.Request) {
	var (
		dataUUID string
		userUUID string
		err      error
		owner    *models.Owner
	)
	dataUUID = chi.URLParam(req, "uuid")

	if dataUUID == "" {
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

	owner, err = h.manager.Owner().FindOneByUserUUIDAndDataUUID(req.Context(), userUUID, dataUUID)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}
	if owner == nil { // нет данных этого пользователя
		h.log.Infof("owner not found: data_uuid: %s, user_uuid: %s", dataUUID, userUUID)
		_ = render.Render(res, req, ErrNotFound)
		return
	}

	var (
		cardData     *models.CardData
		textData     *models.TextData
		fileData     *models.FileData
		metaData     []models.MetaData
		dataResponse *dataByUUIDResponse
	)
	metaData, err = h.manager.MetaData().FindOneByUUID(req.Context(), owner.DataUUID)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	dataResponse = new(dataByUUIDResponse)
	// Данные карт
	if owner.DataType == data_type.CardType {
		cardData, err = h.manager.CardData().FindOneByUUID(req.Context(), owner.DataUUID)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}

		dataResponse.IsCard = true
		dataResponse.CardData.UUID = cardData.UUID
		dataResponse.CardData.Name = cardData.Name
		dataResponse.CardData.CardNumber = cardData.Value.CardNumber
		dataResponse.CardData.NameBank = cardData.Value.NameBank
		dataResponse.CardData.CurrentAccountNumber = cardData.Value.CurrentAccountNumber
		dataResponse.CardData.FullNameHolder = cardData.Value.FullNameHolder
		dataResponse.CardData.PhoneHolder = cardData.Value.PhoneHolder
		dataResponse.CardData.SecurityCode = cardData.Value.SecurityCode
		dataResponse.CardData.ValidityPeriod = cardData.Value.ValidityPeriod.Format(time.RFC3339)

		if len(metaData) > 0 {
			existMeta := make(map[string]string)
			for _, value := range metaData {
				existMeta[value.MetaName] = value.MetaValue.Value
			}
			dataResponse.CardData.Meta = existMeta
		}

	}
	// Текстовые данные
	if owner.DataType == data_type.TextType {
		textData, err = h.manager.TextData().FindOneByUUID(req.Context(), owner.DataUUID)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}
		dataResponse.IsText = true
		dataResponse.TextData.Name = textData.Name
		dataResponse.TextData.UUID = textData.UUID
		dataResponse.TextData.Value = textData.Value

		if len(metaData) > 0 {
			existMeta := make(map[string]string)
			for _, value := range metaData {
				existMeta[value.MetaName] = value.MetaValue.Value
			}
			dataResponse.TextData.Meta = existMeta
		}
	}
	// Бинарные данные
	if owner.DataType == data_type.BinaryType {
		fileData, err = h.manager.FileData().FindOneByUUID(req.Context(), owner.DataUUID)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}
		dataResponse.IsFile = true
		dataResponse.FileData.Name = fileData.Name
		dataResponse.FileData.UUID = fileData.UUID
		dataResponse.FileData.FileName = fileData.FileName
		dataResponse.FileData.Size = fileData.Size
		dataResponse.FileData.Extension = fileData.Extension
		dataResponse.FileData.MimeType = fileData.MimeType

		if len(metaData) > 0 {
			existMeta := make(map[string]string)
			for _, value := range metaData {
				existMeta[value.MetaName] = value.MetaValue.Value
			}
			dataResponse.FileData.Meta = existMeta
		}
	}

	err = render.Render(res, req, dataResponse)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
	}
}

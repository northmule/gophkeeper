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
)

// ItemDataHandler обрабатывает запрос данных по uuid
type ItemDataHandler struct {
	log             *logger.Logger
	userFinderByJWT UserFinderByJWT
	ownerCRUD       OwnerCRUD
	cardDataCRUD    CardDataCRUD
	metaDataCRUD    MetaDataCRUD
	fileDataCRUD    FileDataCRUD
	textDataCRUD    TextDataCRUD
}

// NewItemDataHandler конструктор
func NewItemDataHandler(userFinderByJWT UserFinderByJWT, cardDataCRUD CardDataCRUD, metaDataCRUD MetaDataCRUD, fileDataCRUD FileDataCRUD, textDataCRUD TextDataCRUD, ownerCRUD OwnerCRUD, log *logger.Logger) *ItemDataHandler {
	return &ItemDataHandler{
		userFinderByJWT: userFinderByJWT,
		cardDataCRUD:    cardDataCRUD,
		metaDataCRUD:    metaDataCRUD,
		fileDataCRUD:    fileDataCRUD,
		textDataCRUD:    textDataCRUD,
		ownerCRUD:       ownerCRUD,
		log:             log,
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

	userUUID, err = h.userFinderByJWT.GetUserUUIDByJWTToken(req.Context())
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}

	owner, err = h.ownerCRUD.FindOneByUserUUIDAndDataUUID(req.Context(), userUUID, dataUUID)
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
	metaData, err = h.metaDataCRUD.FindOneByUUID(req.Context(), owner.DataUUID)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	dataResponse = new(dataByUUIDResponse)

	switch owner.DataType {

	case data_type.CardType: // Данные карт
		cardData, err = h.cardDataCRUD.FindOneByUUID(req.Context(), owner.DataUUID)
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

	case data_type.TextType: // Текстовые данные
		textData, err = h.textDataCRUD.FindOneByUUID(req.Context(), owner.DataUUID)
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
	case data_type.BinaryType: // Бинарные данные
		fileData, err = h.fileDataCRUD.FindOneByUUID(req.Context(), owner.DataUUID)
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

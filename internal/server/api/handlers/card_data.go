package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/model_data"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	"github.com/northmule/gophkeeper/internal/server/services/access"
)

// CardDataHandler обрабатывает данные карт
type CardDataHandler struct {
	log           *logger.Logger
	accessService *access.Access
	manager       repository.Repository
}

// NewCardDataHandler конструктор
func NewCardDataHandler(accessService *access.Access, manager repository.Repository, log *logger.Logger) *CardDataHandler {
	return &CardDataHandler{
		accessService: accessService,
		manager:       manager,
		log:           log,
	}
}

type cardDataRequest struct {
	model_data.CardDataRequest
}

// Bind декодирует json в структуру
func (rr *cardDataRequest) Bind(r *http.Request) error {
	return nil
}

// HandleSave создание/обновление данных карты
func (h *CardDataHandler) HandleSave(res http.ResponseWriter, req *http.Request) {

	var (
		err      error
		userUUID string
		owner    *models.Owner
		cardData *models.CardData
	)

	request := new(cardDataRequest)
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
		owner, err = h.manager.Owner().FindOneByUserUUIDAndDataUUIDAndDataType(req.Context(), userUUID, dataUUID, data_type.CardType)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrBadRequest)
			return
		}
		if owner == nil { // нет данных этого пользователя
			h.log.Infof("owner not found: data_uuid: %s, user_uuid: %s, data_type: %s", dataUUID, userUUID, data_type.CardType)
			_ = render.Render(res, req, ErrNotFound)
			return
		}
		cardData, err = h.manager.CardData().FindOneByUUID(req.Context(), dataUUID)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrBadRequest)
			return
		}
		if cardData == nil {
			h.log.Infof("card data not found: uuid %s", owner.DataUUID)
			_ = render.Render(res, req, ErrNotFound)
			return
		}
		// основные данные
		cardData.Name = request.Name
		validityPeriod, _ := time.Parse(time.RFC3339, request.ValidityPeriod)
		cardData.Value.CardNumber = request.CardNumber
		cardData.Value.ValidityPeriod = validityPeriod
		cardData.Value.SecurityCode = request.SecurityCode
		cardData.Value.FullNameHolder = request.FullNameHolder
		cardData.Value.NameBank = request.NameBank
		cardData.Value.PhoneHolder = request.PhoneHolder
		cardData.Value.CurrentAccountNumber = request.CurrentAccountNumber

		err = h.manager.CardData().Update(req.Context(), cardData)
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
		// основные данные
		dataUUID := uuid.NewString()
		cardData = new(models.CardData)
		cardData.Name = request.Name
		cardData.UUID = dataUUID
		cardData.ObjectType = data_type.CardType

		validityPeriod, _ := time.Parse(time.RFC3339, request.ValidityPeriod)

		cardData.Value.CardNumber = request.CardNumber
		cardData.Value.ValidityPeriod = validityPeriod
		cardData.Value.SecurityCode = request.SecurityCode
		cardData.Value.FullNameHolder = request.FullNameHolder
		cardData.Value.NameBank = request.NameBank
		cardData.Value.PhoneHolder = request.PhoneHolder
		cardData.Value.CurrentAccountNumber = request.CurrentAccountNumber

		_, err = h.manager.CardData().Add(req.Context(), cardData)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}
		// владелец данных
		owner = new(models.Owner)
		owner.UserUUID = userUUID
		owner.DataType = data_type.CardType
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

	// общие проверки

	if cardData == nil {
		h.log.Error("card data empty")
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}
	if owner == nil {
		h.log.Error("owner empty")
		_ = render.Render(res, req, ErrInternalServerError)
		return
	}

	//OK
}

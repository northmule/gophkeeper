package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/northmule/gophkeeper/internal/common/data_type"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	"github.com/northmule/gophkeeper/internal/server/services/access"
)

// ItemDataHandler обрабатывает запрос данных по uuid
type ItemDataHandler struct {
	log           *logger.Logger
	accessService *access.Access
	manager       repository.Repository
}

// NewItemDataHandler конструктор
func NewItemDataHandler(accessService *access.Access, manager repository.Repository, log *logger.Logger) *ItemDataHandler {
	return &ItemDataHandler{
		accessService: accessService,
		manager:       manager,
		log:           log,
	}
}

func (h *ItemDataHandler) HandleItem(res http.ResponseWriter, req *http.Request) {
	var (
		dataUUID string
		userUUID string
		err      error
		owner    *models.Owner
	)
	dataUUID = chi.URLParam(req, "uuid")

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
	// Данные карт
	if owner.DataType == data_type.CardType {

	}
	// Текстовые данные
	if owner.DataType == data_type.TextType {

	}
	// Бинарные данные
	if owner.DataType == data_type.BinaryType {

	}
}

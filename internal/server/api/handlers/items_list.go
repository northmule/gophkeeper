package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	"github.com/northmule/gophkeeper/internal/server/services/access"
)

type ItemsListHandler struct {
	log           *logger.Logger
	accessService *access.Access
	manager       repository.Repository
}

func NewItemsListHandler(accessService *access.Access, manager repository.Repository, log *logger.Logger) *ItemsListHandler {
	return &ItemsListHandler{
		accessService: accessService,
		manager:       manager,
		log:           log,
	}
}

func (ih *ItemsListHandler) HandleItemsList(res http.ResponseWriter, req *http.Request) {

	userUUID, err := ih.accessService.GetUserUUIDByJWTToken(req.Context())
	if err != nil {
		ih.log.Error(err)
		_ = render.Render(res, req, ErrBadRequest)
		return
	}
	_ = userUUID

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(nil)
}

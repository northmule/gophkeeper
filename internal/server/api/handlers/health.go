package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/northmule/gophkeeper/internal/server/logger"
)

type HealthHandler struct {
	log *logger.Logger
}

func NewHealthHandler(log *logger.Logger) *HealthHandler {
	return &HealthHandler{log: log}
}

type healthResponse struct {
	OK bool `json:"ok"`
}

func (hr healthResponse) Render(res http.ResponseWriter, req *http.Request) error {
	return nil
}

func (h *HealthHandler) HandleGetHealth(res http.ResponseWriter, req *http.Request) {
	health := healthResponse{OK: true}
	err := render.Render(res, req, health)
	if err != nil {
		h.log.Error(err)
		_ = render.Render(res, req, ErrInternalServerError)
	}
}

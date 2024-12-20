package handlers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/common/util"
	"github.com/northmule/gophkeeper/internal/server/logger"
	"github.com/northmule/gophkeeper/internal/server/repository"
	"github.com/northmule/gophkeeper/internal/server/services/access"
)

// DecryptDataHandler Расшифровывает входящий запрос
type DecryptDataHandler struct {
	log           *logger.Logger
	accessService *access.Access
	manager       repository.Repository
}

// NewDecryptDataHandler конструктор
func NewDecryptDataHandler(accessService *access.Access, manager repository.Repository, log *logger.Logger) *DecryptDataHandler {
	return &DecryptDataHandler{
		log:           log,
		accessService: accessService,
		manager:       manager,
	}
}

// HandleDecryptData расшифрует входящий запрос
func (h *DecryptDataHandler) HandleDecryptData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var (
			err              error
			userUUID         string
			user             *models.User
			bodyBytesDecrypt []byte
		)

		userUUID, err = h.accessService.GetUserUUIDByJWTToken(req.Context())
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrBadRequest)
			return
		}

		user, err = h.manager.User().FindOneByUUID(req.Context(), userUUID)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}

		// копия body
		bodyBytes, _ := io.ReadAll(req.Body)
		// Расшифровываем тело запроса
		bodyBytesDecrypt, err = util.DataDecryptAES(bodyBytes, []byte(user.PrivateClientKey))
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}
		// Восстанавливаем body после чтения
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytesDecrypt))

		next.ServeHTTP(res, req)
	})
}

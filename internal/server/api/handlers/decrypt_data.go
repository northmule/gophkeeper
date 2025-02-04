package handlers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/northmule/gophkeeper/internal/common/models"
	"github.com/northmule/gophkeeper/internal/common/util"
	"github.com/northmule/gophkeeper/internal/server/logger"
)

// DecryptDataHandler Расшифровывает входящий запрос
type DecryptDataHandler struct {
	log             *logger.Logger
	userFinderByJWT UserFinderByJWT
	userFinder      UserFinder
}

// NewDecryptDataHandler конструктор
func NewDecryptDataHandler(userFinderByJWT UserFinderByJWT, userFinder UserFinder, log *logger.Logger) *DecryptDataHandler {
	return &DecryptDataHandler{
		log:             log,
		userFinderByJWT: userFinderByJWT,
		userFinder:      userFinder,
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

		userUUID, err = h.userFinderByJWT.GetUserUUIDByJWTToken(req.Context())
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrBadRequest)
			return
		}

		user, err = h.userFinder.FindOneByUUID(req.Context(), userUUID)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}

		// копия body
		bodyBytes, _ := io.ReadAll(req.Body)
		if len(bodyBytes) == 0 {
			h.log.Error(err)
			_ = render.Render(res, req, ErrBadRequest)
			return
		}
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

type MixedResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (m *MixedResponseWriter) Write(p []byte) (int, error) {
	return m.buf.Write(p)
}

// HandleEncryptData зашифрует исходящие данные
func (h *DecryptDataHandler) HandleEncryptData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var (
			err              error
			userUUID         string
			user             *models.User
			bodyBytesEncrypt []byte
		)

		mixedResponseWriter := &MixedResponseWriter{
			ResponseWriter: res,
			buf:            &bytes.Buffer{},
		}

		next.ServeHTTP(mixedResponseWriter, req)

		userUUID, err = h.userFinderByJWT.GetUserUUIDByJWTToken(req.Context())
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrBadRequest)
			return
		}

		user, err = h.userFinder.FindOneByUUID(req.Context(), userUUID)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}
		// Шифруем ответ
		bodyBytesEncrypt, err = util.DataEncryptAES(mixedResponseWriter.buf.Bytes(), []byte(user.PrivateClientKey))
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}
		_, err = mixedResponseWriter.ResponseWriter.Write(bodyBytesEncrypt)
		if err != nil {
			h.log.Error(err)
			_ = render.Render(res, req, ErrInternalServerError)
			return
		}
	})
}

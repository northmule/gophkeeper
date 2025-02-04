package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/northmule/gophkeeper/internal/server/logger"
)

// ValidatorHandler Валидатор
type ValidatorHandler struct {
	log           *logger.Logger
	requestStruct any
}

// NewValidatorHandler конструктор
func NewValidatorHandler(requestStruct any, log *logger.Logger) *ValidatorHandler {
	return &ValidatorHandler{
		log:           log,
		requestStruct: requestStruct,
	}
}

// HandleValidation валидирует запрос
func (v *ValidatorHandler) HandleValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var err error

		st := v.requestStruct
		// копия body
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		validate := validator.New(validator.WithRequiredStructEnabled()) // doc: https://pkg.go.dev/github.com/go-playground/validator/v10
		// Список типов для валидации из текущих хэндлеров
		switch requestType := st.(type) {
		case *registrationRequest:

			err = render.Bind(req, requestType)
			err = errors.Join(err, validate.Struct(requestType))
		case *authenticationRequest:

			err = render.Bind(req, requestType)
			err = errors.Join(err, validate.Struct(requestType))
		case *cardDataRequest:

			err = render.Bind(req, requestType)
			err = errors.Join(err, validate.Struct(requestType))
		case *textDataRequest:

			err = render.Bind(req, requestType)
			err = errors.Join(err, validate.Struct(requestType))
		case *fileDataInitRequest:

			err = render.Bind(req, requestType)
			err = errors.Join(err, validate.Struct(requestType))

			// Пропускаем не известные
		default:
			v.log.Info(fmt.Sprintf("skip validation for %v", &v.requestStruct))
			next.ServeHTTP(res, req)
			return
		}
		// Восстанавливаем body после чтения в bind
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err != nil {
			var verr validator.ValidationErrors
			if errors.As(err, &verr) {
				v.log.Info(err)
				_ = render.Render(res, req, ErrValidation(verr))
			} else {
				v.log.Error(err)
				_ = render.Render(res, req, ErrBadRequest)
			}
			return
		}

		next.ServeHTTP(res, req)
	})
}

package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

var (
	ErrNotFound            = &ErrResponse{HTTPStatusCode: http.StatusNotFound, StatusText: "Resource not found."}
	ErrBadRequest          = &ErrResponse{HTTPStatusCode: http.StatusBadRequest, StatusText: "Bad request"}
	ErrInternalServerError = &ErrResponse{HTTPStatusCode: http.StatusInternalServerError, StatusText: "Internal Server Error"}
	ErrUnauthorized        = &ErrResponse{HTTPStatusCode: http.StatusUnauthorized, StatusText: "Authentication failed"}
)

func ErrConflict(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusConflict,
		StatusText:     "Duplicate Key",
		ErrorText:      err.Error(),
	}
}

func ErrValidation(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Validation error",
		ErrorText:      err.Error(),
	}
}

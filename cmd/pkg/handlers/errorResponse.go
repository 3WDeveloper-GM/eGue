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

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func NewErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	errorResponse := &ErrResponse{
		Err:            err,
		HTTPStatusCode: statusCode,
		ErrorText:      err.Error(),
	}

	render.Status(r, statusCode)
	render.JSON(w, r, errorResponse)
}

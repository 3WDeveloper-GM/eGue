package handlers

import (
	"log"
	"net/http"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/jsonio"
)

type errResponse struct {
	io jsonio.JsonIORW
}

func NewErrResponse(io jsonio.JsonIORW) *errResponse {
	return &errResponse{io: io}
}

func (errR *errResponse) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	envelope := map[string]interface{}{
		"error": message,
	}

	err := errR.io.WriteJSON(w, status, envelope, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}
}

func (errR *errResponse) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)

	message := "the server encountered a problem and could not process your request."
	errR.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (errR *errResponse) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errR.errorResponse(w, r, http.StatusBadRequest, err)
}

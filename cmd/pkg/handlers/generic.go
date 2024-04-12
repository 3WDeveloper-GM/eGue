package handlers

import (
	"net/http"
)

// HealthCheckHandler returns a response whether or not the backend is live.
func (sh *SearchHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "available",
	}

	err := sh.io.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		responser := NewErrResponse(sh.io)
		responser.serverErrorResponse(w, r, err)
		return
	}
}

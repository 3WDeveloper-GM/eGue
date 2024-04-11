package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

func (sh *SearchHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "available",
	}

	render.JSON(w, r, response)
}

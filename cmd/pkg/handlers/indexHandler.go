package handlers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/render"
)

type IndexHandler struct {
	indexer DBIndex
}

func NewIndexHandler(indexer DBIndex) *IndexHandler {
	return &IndexHandler{indexer: indexer}
}

func (ih *IndexHandler) IndexMails(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	root, err := os.Getwd()
	if err != nil {
		NewErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	err = ih.indexer.Index(ctx, root)
	if err != nil {
		NewErrorResponse(w, r, http.StatusInternalServerError, err)
	}

	response := map[string]interface{}{
		"sent": true,
	}

	render.JSON(w, r, response)

}

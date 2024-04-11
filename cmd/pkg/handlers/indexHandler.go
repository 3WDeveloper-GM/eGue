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

// IndexMails: the handler gets the necessary parameters in order to run the
// Index() method that is part of the DBIndex interface. When the payload is
// sent, it sends a message to the client signaling that the mails are indexed.
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

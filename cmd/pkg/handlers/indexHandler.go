package handlers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/jsonio"
)

type IndexHandler struct {
	indexer DBIndex
	io      jsonio.JsonIORW
}

func NewIndexHandler(indexer DBIndex, io jsonio.JsonIORW) *IndexHandler {
	return &IndexHandler{indexer: indexer, io: io}
}

// IndexMails runs the Index() method that is specified by the DBIndex() interface
// this method is for demoing purposes, the best way to index emails is still running
// the indexer binary generated in the indexer_standalone part of the crawler package.
func (ih *IndexHandler) IndexMails(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	root, err := os.Getwd()
	if err != nil {
		responser := NewErrResponse(ih.io)
		responser.serverErrorResponse(w, r, err)
		return
	}

	err = ih.indexer.Index(ctx, root)
	if err != nil {
		responser := NewErrResponse(ih.io)
		responser.serverErrorResponse(w, r, err)
	}

	response := map[string]interface{}{
		"sent": true,
	}

	err = ih.io.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		responser := NewErrResponse(ih.io)
		responser.serverErrorResponse(w, r, err)
	}

}

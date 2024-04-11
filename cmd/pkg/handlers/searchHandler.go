package handlers

import (
	"errors"
	"net/http"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/adapter/zs"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/domain"
	Validator "github.com/3WDeveloper-GM/pipeline/cmd/pkg/domain/validator"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/jsonio"
	"github.com/go-chi/render"
)

type SearchHandler struct {
	searcher DBSearch
	io       jsonio.JsonIORW
}

func NewSearchHandler(searcher DBSearch, io jsonio.JsonIORW) *SearchHandler {
	return &SearchHandler{
		searcher: searcher,
		io:       io,
	}
}

func (sh *SearchHandler) SearchMails(w http.ResponseWriter, r *http.Request) {
	var input zs.SearchRequest

	// read the json from the user
	err := sh.io.ReadJSON(w, r, &input)
	if err != nil {
		NewErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	//log.Println(input)

	// validate the json instance
	v := Validator.NewValidator()
	if !domain.ValidateInput(v, &input) {
		errs := errors.Join(v.Errors()...)
		NewErrorResponse(w, r, http.StatusBadRequest, errs)
		return
	}

	//log.Println("validated!")

	sh.searcher.SetInput(input)

	// Get the mails with the search method
	returnedMails, err := sh.searcher.Search(domain.Index)
	if err != nil {
		NewErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	res := map[string]interface{}{
		"length": len(returnedMails),
		"result": returnedMails,
	}

	render.JSON(w, r, res)
}

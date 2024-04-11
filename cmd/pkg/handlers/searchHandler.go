package handlers

import (
	"errors"
	"net/http"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/adapter/zs"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/domain"
	Validator "github.com/3WDeveloper-GM/pipeline/cmd/pkg/domain/validator"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/jsonio"
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

// SearchMails validates and reads the query input that comes from the client
// then runs the Search() method implemented by the DBSearch Interface. Then it
// returns a JSON response back to the client.
func (sh *SearchHandler) SearchMails(w http.ResponseWriter, r *http.Request) {
	var input zs.SearchRequest

	// read the json from the user
	err := sh.io.ReadJSON(w, r, &input)
	if err != nil {
		responser := NewErrResponse(sh.io)
		responser.serverErrorResponse(w, r, err)
		return
	}

	//log.Println(input)

	// validate the json instance
	v := Validator.NewValidator()
	if !domain.ValidateInput(v, &input) {
		errs := errors.Join(v.Errors()...)
		responser := NewErrResponse(sh.io)
		responser.badRequestResponse(w, r, errs)
		return
	}

	//log.Println("validated!")

	sh.searcher.SetInput(input)

	// Get the mails with the search method
	returnedMails, err := sh.searcher.Search(domain.Index)
	if err != nil {
		responser := NewErrResponse(sh.io)
		responser.serverErrorResponse(w, r, err)
		return
	}

	res := map[string]interface{}{
		"length": len(returnedMails),
		"result": returnedMails,
	}

	err = sh.io.WriteJSON(w, http.StatusOK, res, nil)
	if err != nil {
		responser := NewErrResponse(sh.io)
		responser.serverErrorResponse(w, r, err)
	}
}

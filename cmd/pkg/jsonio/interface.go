package jsonio

import "net/http"

type JsonIORW interface {
	ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error
}

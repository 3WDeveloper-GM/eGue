package jsonio

import "net/http"

type JsonIORW interface {
	ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error
	WriteJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error
}

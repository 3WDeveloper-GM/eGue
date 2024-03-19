package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Envelope map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {

	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	json = append(json, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {

		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError
		var invalidunmarshalErr *json.InvalidUnmarshalError

		switch {

		case errors.As(err, &syntaxErr):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxErr.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeErr):
			if unmarshalTypeErr.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeErr.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeErr.Offset)

		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")

		case errors.As(err, &invalidunmarshalErr):
			panic(err)

		default:
			return err
		}
	}
	return nil
}

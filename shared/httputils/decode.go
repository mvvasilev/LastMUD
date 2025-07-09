package httputils

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeBody[T any](reader io.Reader, body T) (err error) {
	err = json.NewDecoder(reader).Decode(body)
	return
}

func EncodeBody[T any](writer io.Writer, body T) (err error) {
	return json.NewEncoder(writer).Encode(body)
}

func WriteUnhandledError(w http.ResponseWriter, err error) {
	http.Error(w, NewHTTPError(http.StatusInternalServerError, err.Error()).Error(), http.StatusInternalServerError)
}

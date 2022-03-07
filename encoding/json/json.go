package json

import (
	"net/http"

	"github.com/abemedia/go-don"
	"github.com/goccy/go-json"
)

func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).DecodeContext(r.Context(), v)
}

func encodeJSON(w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func init() {
	don.RegisterDecoder("application/json", decodeJSON)
	don.RegisterEncoder("application/json", encodeJSON)
}

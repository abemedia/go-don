package form

import (
	"net/http"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/decoder"
)

var MemoryLimit int64 = 1 << 20 // 1MB

var dec = decoder.NewDecoder("form")

func decodeForm(r *http.Request, v interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	return dec.Decode(decoder.MapGetter(r.Form), v)
}

func decodeMultipartForm(r *http.Request, v interface{}) error {
	if err := r.ParseMultipartForm(MemoryLimit); err != nil {
		return err
	}

	return dec.Decode(decoder.MapGetter(r.Form), v)
}

func init() {
	don.RegisterDecoder("application/x-www-form-urlencoded", decodeForm)
	don.RegisterDecoder("multipart/form-data", decodeMultipartForm)
}
